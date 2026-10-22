package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPicker(t *testing.T) {
	cases := []struct {
		name       string
		entryNames []string // trailing "/" means directory
		want       []string
		wantEmpty  bool
	}{
		{
			name:       "mixed mov and qt files",
			entryNames: []string{"b.mov", "a.MOV", "c.qt", "d.QT", "e.Qt", "f.mov"},
			want:       []string{"a.MOV", "b.mov", "c.qt", "d.QT", "e.Qt", "f.mov"},
		},
		{
			name:       "directories excluded",
			entryNames: []string{"movie.mov/", "sample.qt/", "real.mov"},
			want:       []string{"real.mov"},
		},
		{
			name:       "invalid extensions excluded",
			entryNames: []string{"a.txt", "b.mp4", "c.mkv", "noext", "d.mov"},
			want:       []string{"d.mov"},
		},
		{
			name:       "no supported files",
			entryNames: []string{"a.txt", "b.mp4", "notes.qt-"},
			wantEmpty:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, entry := range tc.entryNames {
				if strings.HasSuffix(entry, "/") {
					if err := os.Mkdir(filepath.Join(dir, strings.TrimSuffix(entry, "/")), 0o755); err != nil {
						t.Fatal(err)
					}
					continue
				}
				if err := os.WriteFile(filepath.Join(dir, entry), []byte("x"), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			picker, err := NewPicker(dir)
			if err != nil {
				t.Fatalf("NewPicker: %v", err)
			}

			if got := picker.files; !equalStrings(got, tc.want) {
				t.Errorf("files = %v, want %v", got, tc.want)
			}
			if picker.Empty() != tc.wantEmpty {
				t.Errorf("Empty() = %v, want %v", picker.Empty(), tc.wantEmpty)
			}
		})
	}
}

func TestNewPickerMissingDir(t *testing.T) {
	_, err := NewPicker(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
}

// equalStrings compares contents, treating nil and empty slices as equal.
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
