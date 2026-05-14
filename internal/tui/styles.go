package tui

import "github.com/charmbracelet/lipgloss"

var (
	styleHeader   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	styleQueued   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleRunning  = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	styleDone     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	styleSkipped  = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	styleFailed   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	styleCancelled = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	styleConfirm  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
	styleFooter   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)
