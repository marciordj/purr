package main

import (
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	cases := []struct {
		name       string
		argv       []string
		wantInputs []string
		wantFast   bool
		wantErr    string
	}{
		{
			name:       "single file",
			argv:       []string{"a.mov"},
			wantInputs: []string{"a.mov"},
			wantFast:   false,
		},
		{
			name:       "--fast flag",
			argv:       []string{"--fast", "a.mov", "b.mov"},
			wantInputs: []string{"a.mov", "b.mov"},
			wantFast:   true,
		},
		{
			name:       "-f shorthand",
			argv:       []string{"-f", "a.mov"},
			wantInputs: []string{"a.mov"},
			wantFast:   true,
		},
		{
			name:    "non-.mov file",
			argv:    []string{"a.txt"},
			wantErr: "not a .mov",
		},
		{
			name:    "unknown flag",
			argv:    []string{"--unknown"},
			wantErr: "unknown flag",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inputs, fast, err := parseArgs(tc.argv)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tc.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fast != tc.wantFast {
				t.Errorf("fast = %v, want %v", fast, tc.wantFast)
			}
			if len(inputs) != len(tc.wantInputs) {
				t.Fatalf("inputs = %v, want %v", inputs, tc.wantInputs)
			}
			for i := range inputs {
				if inputs[i] != tc.wantInputs[i] {
					t.Errorf("inputs[%d] = %q, want %q", i, inputs[i], tc.wantInputs[i])
				}
			}
		})
	}
}
