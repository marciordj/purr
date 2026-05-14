package tui

import tea "github.com/charmbracelet/bubbletea"

type phase int

const (
	phasePicker phase = iota
	phaseRunner
)

// App is the top-level model when purr is run without arguments.
// It sequences the Picker into the runner Model.
type App struct {
	phase  phase
	picker Picker
	runner Model
}

// NewApp creates an App starting at the picker screen.
func NewApp(picker Picker) App {
	return App{picker: picker}
}

func (a App) Init() tea.Cmd {
	return a.picker.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch a.phase {
	case phasePicker:
		m, cmd := a.picker.Update(msg)
		a.picker = m.(Picker)
		if a.picker.done {
			a.runner = New(a.picker.Selected())
			a.phase = phaseRunner
			return a, a.runner.Init()
		}
		return a, cmd

	case phaseRunner:
		m, cmd := a.runner.Update(msg)
		a.runner = m.(Model)
		return a, cmd
	}
	return a, nil
}

func (a App) View() string {
	if a.phase == phasePicker {
		return a.picker.View()
	}
	return a.runner.View()
}
