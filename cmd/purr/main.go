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

func main() {
	if err := convert.EnsureDeps(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if len(os.Args) >= 2 {
		run(directMode())
		return
	}

	run(pickerMode())
}

func directMode() tea.Model {
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
	return tui.New(inputs)
}

func pickerMode() tea.Model {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	picker, err := tui.NewPicker(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if picker.Empty() {
		fmt.Fprintln(os.Stderr, "no .mov files in current directory")
		os.Exit(1)
	}
	return tui.NewApp(picker)
}

func run(m tea.Model) {
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
