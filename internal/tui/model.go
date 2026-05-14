package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marciord/purr/internal/convert"
)

type status int

const (
	statusQueued status = iota
	statusRunning
	statusDone
	statusSkipped
	statusFailed
	statusCancelled
)

type job struct {
	input    string
	output   string
	status   status
	progress convert.Progress
	err      error
}

// Model is the Bubble Tea model for the purr TUI.
type Model struct {
	jobs        []job
	current     int
	cancel      context.CancelFunc
	bar         progress.Model
	confirming  bool
	abortQueue  bool
	progCh      <-chan convert.Progress
	errCh       <-chan error
	width       int
}

// New creates a Model with one job per input path.
func New(inputs []string) Model {
	jobs := make([]job, len(inputs))
	for i, in := range inputs {
		ext := filepath.Ext(in)
		out := strings.TrimSuffix(in, ext) + ".mp4"
		jobs[i] = job{input: in, output: out}
	}
	return Model{
		jobs:    jobs,
		current: -1,
		bar:     progress.New(progress.WithDefaultGradient(), progress.WithWidth(40)),
		width:   80,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return startMsg{} }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		barWidth := msg.Width - 6
		if barWidth < 10 {
			barWidth = 10
		}
		m.bar = progress.New(progress.WithDefaultGradient(), progress.WithWidth(barWidth))
		return m, nil

	case startMsg:
		return m.startNext()

	case progressMsg:
		if msg.index == m.current {
			m.jobs[m.current].progress = msg.p
		}
		return m, waitForProgress(m.progCh, m.current)

	case jobDoneMsg:
		if msg.index >= 0 && msg.index < len(m.jobs) {
			switch {
			case errors.Is(msg.err, context.Canceled):
				m.jobs[msg.index].status = statusCancelled
			case msg.err != nil:
				m.jobs[msg.index].status = statusFailed
				m.jobs[msg.index].err = msg.err
			default:
				m.jobs[msg.index].status = statusDone
			}
		}
		if m.abortQueue {
			return m, tea.Quit
		}
		return m.startNext()

	case tea.KeyMsg:
		if m.confirming {
			switch msg.String() {
			case "y", "Y":
				m.confirming = false
				m.abortQueue = true
				if m.cancel != nil {
					m.cancel()
				}
			default:
				m.confirming = false
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c":
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		case "q":
			if m.current < 0 || m.jobs[m.current].status != statusRunning {
				return m, tea.Quit
			}
		case "esc":
			if m.current >= 0 && m.jobs[m.current].status == statusRunning {
				m.confirming = true
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("purr — .mov → .mp4"))
	b.WriteString("\n\n")

	for _, j := range m.jobs {
		name := filepath.Base(j.input)
		switch j.status {
		case statusQueued:
			b.WriteString(styleQueued.Render("  ○  " + name))
		case statusRunning:
			b.WriteString(styleRunning.Render("  ▸  " + name))
			if j.progress.Speed > 0 {
				b.WriteString(styleRunning.Render(fmt.Sprintf("  %.1fx", j.progress.Speed)))
			}
			b.WriteString("\n")
			b.WriteString("     " + m.bar.ViewAs(j.progress.Percent/100))
			if j.progress.ETA > 0 {
				b.WriteString("\n     " + styleQueued.Render("ETA "+j.progress.ETA.Round(time.Second).String()))
			}
		case statusDone:
			b.WriteString(styleDone.Render("  ✓  " + name))
		case statusSkipped:
			b.WriteString(styleSkipped.Render("  –  "+name) +
				styleQueued.Render(" (output exists, skipped)"))
		case statusFailed:
			b.WriteString(styleFailed.Render("  ✗  " + name))
			if j.err != nil {
				b.WriteString("\n     " + styleFailed.Render(j.err.Error()))
			}
		case statusCancelled:
			b.WriteString(styleCancelled.Render("  ◌  " + name + " (cancelled)"))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if m.confirming {
		b.WriteString(styleConfirm.Render("  Abort remaining jobs? [y/N] "))
	} else {
		b.WriteString(styleFooter.Render("  esc cancel job  ctrl+c abort all  q quit"))
	}
	b.WriteString("\n")

	return b.String()
}

func (m Model) startNext() (Model, tea.Cmd) {
	for i := range m.jobs {
		if m.jobs[i].status != statusQueued {
			continue
		}
		if _, err := os.Stat(m.jobs[i].output); err == nil {
			m.jobs[i].status = statusSkipped
			continue
		}
		m.current = i
		m.jobs[i].status = statusRunning
		ctx, cancel := context.WithCancel(context.Background())
		m.cancel = cancel
		progCh, errCh := convert.Convert(ctx, m.jobs[i].input, m.jobs[i].output)
		m.progCh = progCh
		m.errCh = errCh
		return m, tea.Batch(
			waitForProgress(progCh, i),
			waitForError(errCh, i),
		)
	}
	return m, tea.Quit
}

func waitForProgress(ch <-chan convert.Progress, index int) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			return nil
		}
		return progressMsg{index: index, p: p}
	}
}

func waitForError(ch <-chan error, index int) tea.Cmd {
	return func() tea.Msg {
		return jobDoneMsg{index: index, err: <-ch}
	}
}
