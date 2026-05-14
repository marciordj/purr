package convert

import (
	"strings"
	"testing"
	"time"
)

func TestParseProgress_partial(t *testing.T) {
	input := `frame=10
fps=30.0
out_time_ms=500000
speed=1.5x
progress=continue
frame=20
fps=30.0
out_time_ms=1000000
speed=2.0x
progress=continue
`
	ch := make(chan Progress, 10)
	go func() {
		parseProgress(strings.NewReader(input), 4.0, ch)
		close(ch)
	}()

	var snaps []Progress
	for p := range ch {
		snaps = append(snaps, p)
	}

	if len(snaps) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snaps))
	}

	// 0.5s / 4s = 12.5%
	if got := snaps[0].Percent; got < 12 || got > 13 {
		t.Errorf("snapshot 0: expected ~12.5%%, got %.2f", got)
	}
	if snaps[0].Speed != 1.5 {
		t.Errorf("snapshot 0: expected speed 1.5, got %f", snaps[0].Speed)
	}
	// ETA: (4 - 0.5) / 1.5 ≈ 2.33s
	if got := snaps[0].ETA; got < 2*time.Second || got > 3*time.Second {
		t.Errorf("snapshot 0: expected ETA ~2.3s, got %v", got)
	}

	// 1.0s / 4s = 25%
	if got := snaps[1].Percent; got < 24 || got > 26 {
		t.Errorf("snapshot 1: expected ~25%%, got %.2f", got)
	}
	if snaps[1].Speed != 2.0 {
		t.Errorf("snapshot 1: expected speed 2.0, got %f", snaps[1].Speed)
	}
}

func TestParseProgress_end(t *testing.T) {
	input := `out_time_ms=4000000
speed=4.0x
progress=end
`
	ch := make(chan Progress, 10)
	go func() {
		parseProgress(strings.NewReader(input), 4.0, ch)
		close(ch)
	}()

	var snaps []Progress
	for p := range ch {
		snaps = append(snaps, p)
	}

	if len(snaps) == 0 {
		t.Fatal("expected at least one snapshot")
	}
	if last := snaps[len(snaps)-1]; last.Percent != 100 {
		t.Errorf("expected 100%% at end, got %.2f", last.Percent)
	}
}

func TestParseProgress_zeroDuration(t *testing.T) {
	input := `out_time_ms=500000
speed=1.0x
progress=continue
progress=end
`
	ch := make(chan Progress, 10)
	go func() {
		parseProgress(strings.NewReader(input), 0, ch)
		close(ch)
	}()

	var count int
	for range ch {
		count++
	}
	// progress=continue skipped (zeroDuration), progress=end emits one snapshot
	if count != 1 {
		t.Errorf("expected 1 snapshot (end only), got %d", count)
	}
}
