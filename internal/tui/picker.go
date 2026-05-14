package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Picker is the file-selection screen shown when purr is run without arguments.
type Picker struct {
	files    []string
	paths    []string
	cursor   int
	selected map[int]bool
	done     bool
}

type movFile struct{ name, path string }

// NewPicker scans dir for .mov files (case-insensitive) and returns a ready Picker.
func NewPicker(dir string) (Picker, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Picker{}, err
	}

	var found []movFile
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".mov") {
			found = append(found, movFile{e.Name(), filepath.Join(dir, e.Name())})
		}
	}
	sort.Slice(found, func(i, j int) bool {
		return strings.ToLower(found[i].name) < strings.ToLower(found[j].name)
	})

	files := make([]string, len(found))
	paths := make([]string, len(found))
	for i, f := range found {
		files[i] = f.name
		paths[i] = f.path
	}

	return Picker{
		files:    files,
		paths:    paths,
		selected: make(map[int]bool),
	}, nil
}

// Empty reports whether no .mov files were found.
func (p Picker) Empty() bool { return len(p.files) == 0 }

// Selected returns the absolute paths of all checked files, in list order.
func (p Picker) Selected() []string {
	var s []string
	for i, path := range p.paths {
		if p.selected[i] {
			s = append(s, path)
		}
	}
	return s
}

func (p Picker) Init() tea.Cmd { return nil }

func (p Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.files)-1 {
				p.cursor++
			}
		case " ":
			p.selected[p.cursor] = !p.selected[p.cursor]
		case "enter":
			if len(p.Selected()) > 0 {
				p.done = true
			}
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p Picker) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("purr — select .mov files"))
	b.WriteString("\n\n")

	for i, name := range p.files {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		check := "[ ]"
		if p.selected[i] {
			check = "[x]"
		}
		line := cursor + check + " " + name
		if i == p.cursor {
			b.WriteString(styleRunning.Render(line))
		} else {
			b.WriteString(styleQueued.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleFooter.Render("  ↑/↓ move   space toggle   enter convert   q quit"))
	b.WriteString("\n")

	return b.String()
}
