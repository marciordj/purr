package convert

import (
	"encoding/json"
	"os/exec"
)

type probeResult struct {
	Streams []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

var mp4CompatibleVideo = map[string]bool{
	"h264": true, "hevc": true, "h265": true,
	"vp9": true, "av1": true, "mpeg4": true,
}

var mp4CompatibleAudio = map[string]bool{
	"aac": true, "mp3": true, "opus": true,
	"ac3": true, "eac3": true, "flac": true,
}

func probe(path string) (*probeResult, error) {
	out, err := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-show_format",
		path,
	).Output()
	if err != nil {
		return nil, err
	}
	var r probeResult
	if err := json.Unmarshal(out, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func audioIsAAC(r *probeResult) bool {
	for _, s := range r.Streams {
		if s.CodecType == "audio" {
			return s.CodecName == "aac"
		}
	}
	return false
}

func needsReencode(r *probeResult) bool {
	for _, s := range r.Streams {
		switch s.CodecType {
		case "video":
			if !mp4CompatibleVideo[s.CodecName] {
				return true
			}
		case "audio":
			if !mp4CompatibleAudio[s.CodecName] {
				return true
			}
		}
	}
	return false
}
