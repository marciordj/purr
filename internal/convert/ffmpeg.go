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

func reencodeArgs(in, out string) []string {
	return []string{
		"-v", "error",
		"-nostdin",
		"-i", in,
		"-c:v", "libx264",
		"-crf", "18",
		"-preset", "medium",
		"-c:a", "aac",
		"-movflags", "+faststart",
		"-progress", "pipe:1",
		"-nostats",
		"-y",
		out,
	}
}
