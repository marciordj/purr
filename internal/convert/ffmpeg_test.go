package convert

import (
	"slices"
	"testing"
)

func TestCompressArgs_AudioCompatible(t *testing.T) {
	args := compressArgs("in.mov", "out.mp4", true)
	for _, want := range []string{"-c:v", "hevc_videotoolbox", "-q:v", "65", "-tag:v", "hvc1", "-c:a", "copy", "-movflags", "+faststart"} {
		if !slices.Contains(args, want) {
			t.Errorf("args missing %q; got %v", want, args)
		}
	}
	if slices.Contains(args, "aac") {
		t.Errorf("args should not contain aac when audio is compatible; got %v", args)
	}
}

func TestCompressArgs_AudioReencode(t *testing.T) {
	args := compressArgs("in.mov", "out.mp4", false)
	for _, want := range []string{"-c:a", "aac", "-b:a", "192k"} {
		if !slices.Contains(args, want) {
			t.Errorf("args missing %q; got %v", want, args)
		}
	}
	if slices.Contains(args, "copy") {
		t.Errorf("args should not contain 'copy' when audio needs re-encode; got %v", args)
	}
}

func TestChooseArgs_DefaultAlwaysCompress(t *testing.T) {
	h264 := &probeResult{Streams: []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
	}{{CodecName: "h264", CodecType: "video"}, {CodecName: "aac", CodecType: "audio"}}}

	args := chooseArgs("in.mov", "out.mp4", h264, false)
	if !slices.Contains(args, "hevc_videotoolbox") {
		t.Errorf("default mode must use hevc_videotoolbox; got %v", args)
	}
	// must not be a full remux ("-c copy" copies all streams including video)
	for i, a := range args {
		if a == "-c" && i+1 < len(args) && args[i+1] == "copy" {
			t.Errorf("default mode must not use full stream copy (-c copy); got %v", args)
		}
	}
}

func TestChooseArgs_FastRemuxH264(t *testing.T) {
	h264 := &probeResult{Streams: []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
	}{{CodecName: "h264", CodecType: "video"}, {CodecName: "aac", CodecType: "audio"}}}

	args := chooseArgs("in.mov", "out.mp4", h264, true)
	if !slices.Contains(args, "copy") {
		t.Errorf("--fast with compatible codec must use copy (remux); got %v", args)
	}
}

func TestChooseArgs_FastFallbackProRes(t *testing.T) {
	prores := &probeResult{Streams: []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
	}{{CodecName: "prores", CodecType: "video"}, {CodecName: "aac", CodecType: "audio"}}}

	args := chooseArgs("in.mov", "out.mp4", prores, true)
	if !slices.Contains(args, "hevc_videotoolbox") {
		t.Errorf("--fast with incompatible codec must fall back to compress; got %v", args)
	}
}
