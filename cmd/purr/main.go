package main

import (
	"fmt"
	"os"
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
	inputs, fast, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, arg := range inputs {
		if _, err := os.Stat(arg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	}
	return tui.New(inputs, fast)
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
		fmt.Fprintln(os.Stderr, "no .mov or .qt files in current directory")
		os.Exit(1)
	}
	return tui.NewApp(picker, false)
}

func run(m tea.Model) {
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func parseArgs(argv []string) (inputs []string, fast bool, err error) {
	for _, arg := range argv {
		switch {
		case arg == "--fast" || arg == "-f":
			fast = true
		case strings.HasPrefix(arg, "-"):
			return nil, false, fmt.Errorf("unknown flag %q", arg)
		case !convert.IsSupportedInput(arg):
			return nil, false, fmt.Errorf("%q is not a .mov or .qt file", arg)
		default:
			inputs = append(inputs, arg)
		}
	}
	return inputs, fast, nil
}
