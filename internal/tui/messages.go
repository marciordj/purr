package tui

import "github.com/marciord/purr/internal/convert"

type startMsg struct{}

type progressMsg struct {
	index int
	p     convert.Progress
}

type jobDoneMsg struct {
	index int
	err   error
}
