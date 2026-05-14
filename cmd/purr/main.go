package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marciord/purr/internal/convert"
	"github.com/marciord/purr/internal/tui"
)

const usage = `purr — .mov to .mp4 converter

Usage:
  purr <file.mov> [file.mov ...]

Requirements:
  ffmpeg and ffprobe must be available in PATH.`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	if err := convert.EnsureDeps(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	var inputs []string
	for _, arg := range os.Args[1:] {
		if !strings.EqualFold(filepath.Ext(arg), ".mov") {
			fmt.Fprintf(os.Stderr, "error: %q is not a .mov file\n", arg)
			os.Exit(1)
		}
		if _, err := os.Stat(arg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		inputs = append(inputs, arg)
	}

	if _, err := tea.NewProgram(tui.New(inputs)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
