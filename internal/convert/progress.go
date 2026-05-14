package convert

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

// Progress is a snapshot of conversion state emitted by Convert.
type Progress struct {
	Percent float64
	Speed   float64       // 1.2 = 1.2× realtime
	ETA     time.Duration
}

// parseProgress reads ffmpeg -progress output and sends Progress snapshots to ch.
// durationSec is the total input duration obtained from ffprobe.
func parseProgress(r io.Reader, durationSec float64, ch chan<- Progress) {
	scanner := bufio.NewScanner(r)
	var outTimeSec, speed float64

	for scanner.Scan() {
		k, v, ok := strings.Cut(scanner.Text(), "=")
		if !ok {
			continue
		}
		switch k {
		case "out_time_ms": // microseconds despite the name
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				outTimeSec = n / 1_000_000
			}
		case "speed":
			v = strings.TrimSuffix(strings.TrimSpace(v), "x")
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				speed = n
			}
		case "progress":
			if v == "end" {
				ch <- Progress{Percent: 100, Speed: speed}
				return
			}
			if durationSec <= 0 {
				continue
			}
			pct := (outTimeSec / durationSec) * 100
			if pct > 100 {
				pct = 100
			}
			var eta time.Duration
			if speed > 0 && durationSec > outTimeSec {
				remaining := (durationSec - outTimeSec) / speed
				eta = time.Duration(remaining * float64(time.Second))
			}
			ch <- Progress{Percent: pct, Speed: speed, ETA: eta}
		}
	}
}
