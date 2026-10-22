package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewOutputNaming(t *testing.T) {
	inputs := []string{"a.mov", "b.MOV", "c.qt", "d.QT", "e.Qt", "/tmp/vids/my clip.qt"}
	wantOuts := []string{"a.mp4", "b.mp4", "c.mp4", "d.mp4", "e.mp4", "/tmp/vids/my clip.mp4"}

	m := New(inputs, false)
	if len(m.jobs) != len(inputs) {
		t.Fatalf("job count = %d, want %d", len(m.jobs), len(inputs))
	}
	for i, j := range m.jobs {
		if j.input != inputs[i] {
			t.Errorf("job %d input = %q, want %q", i, j.input, inputs[i])
		}
		if j.output != wantOuts[i] {
			t.Errorf("job %d output = %q, want %q", i, j.output, wantOuts[i])
		}
		if j.status != statusQueued {
			t.Errorf("job %d status = %v, want statusQueued", i, j.status)
		}
	}
}

func TestStartNextSkipsExistingOutput(t *testing.T) {
	dir := t.TempDir()
	names := []string{"a.mov", "b.QT", "c.Qt"}
	for _, in := range names {
		if err := os.WriteFile(filepath.Join(dir, in), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		out := strings.TrimSuffix(in, filepath.Ext(in)) + ".mp4"
		if err := os.WriteFile(filepath.Join(dir, out), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	inputs := make([]string, len(names))
	for i, n := range names {
		inputs[i] = filepath.Join(dir, n)
	}

	m := New(inputs, false)
	m2, cmd := m.startNext()

	if len(m2.jobs) != len(inputs) {
		t.Fatalf("job count = %d, want %d", len(m2.jobs), len(inputs))
	}
	if m2.current != -1 {
		t.Errorf("current = %d, want -1 (nothing should start running)", m2.current)
	}
	for i, j := range m2.jobs {
		if j.status != statusSkipped {
			t.Errorf("job %d (%s) status = %v, want statusSkipped", i, j.input, j.status)
		}
	}
	if cmd == nil {
		t.Error("expected a stop command for a queue with nothing left to do")
	}
}
