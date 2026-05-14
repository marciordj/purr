package convert

import (
	"errors"
	"os/exec"
)

func EnsureDeps() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return errors.New("ffmpeg not found in PATH — install with: brew install ffmpeg")
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return errors.New("ffprobe not found in PATH — install with: brew install ffmpeg")
	}
	return nil
}

func remuxArgs(in, out string) []string {
	return []string{
		"-v", "error",
		"-nostdin",
		"-i", in,
		"-c", "copy",
		"-movflags", "+faststart",
		"-progress", "pipe:1",
		"-nostats",
		"-y",
		out,
	}
}

func compressArgs(in, out string, audioCompatible bool) []string {
	args := []string{
		"-v", "error",
		"-nostdin",
		"-i", in,
		"-c:v", "hevc_videotoolbox",
		"-q:v", "65",
		"-tag:v", "hvc1",
		"-allow_sw", "1",
	}
	if audioCompatible {
		args = append(args, "-c:a", "copy")
	} else {
		args = append(args, "-c:a", "aac", "-b:a", "192k")
	}
	return append(args,
		"-movflags", "+faststart",
		"-progress", "pipe:1",
		"-nostats",
		"-y",
		out,
	)
}

func chooseArgs(in, out string, r *probeResult, fast bool) []string {
	if fast && !needsReencode(r) {
		return remuxArgs(in, out)
	}
	return compressArgs(in, out, audioIsAAC(r))
}
