package convert

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Convert starts the conversion of in to out and returns a progress channel
// and an error channel. Both channels are closed when conversion finishes.
// Cancelling ctx kills the ffmpeg process and removes any partial output.
func Convert(ctx context.Context, in, out string) (<-chan Progress, <-chan error) {
	progCh := make(chan Progress, 8)
	errCh := make(chan error, 1)

	go func() {
		defer close(progCh)
		defer close(errCh)

		err := runConvert(ctx, in, out, progCh)
		if err != nil {
			if ctx.Err() != nil {
				os.Remove(out)
				errCh <- ctx.Err()
				return
			}
			errCh <- err
			return
		}
		errCh <- nil
	}()

	return progCh, errCh
}

func runConvert(ctx context.Context, in, out string, progCh chan<- Progress) error {
	r, err := probe(in)
	if err != nil {
		return fmt.Errorf("probe %q: %w", in, err)
	}

	durationSec := 0.0
	if r.Format.Duration != "" {
		durationSec, _ = strconv.ParseFloat(r.Format.Duration, 64)
	}

	var args []string
	if needsReencode(r) {
		args = reencodeArgs(in, out)
	} else {
		args = remuxArgs(in, out)
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return err
	}

	parseProgress(stdout, durationSec, progCh)

	if err := cmd.Wait(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		stderr := strings.TrimSpace(stderrBuf.String())
		if stderr != "" {
			return fmt.Errorf("%w: %s", err, stderr)
		}
		return err
	}
	return nil
}
