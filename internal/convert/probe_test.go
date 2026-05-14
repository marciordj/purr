package convert

import "testing"

func TestAudioIsAAC(t *testing.T) {
	stream := func(codec, kind string) struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
	} {
		return struct {
			CodecName string `json:"codec_name"`
			CodecType string `json:"codec_type"`
		}{CodecName: codec, CodecType: kind}
	}

	cases := []struct {
		name   string
		r      *probeResult
		wantOK bool
	}{
		{"aac audio", &probeResult{Streams: []struct {
			CodecName string `json:"codec_name"`
			CodecType string `json:"codec_type"`
		}{stream("h264", "video"), stream("aac", "audio")}}, true},
		{"pcm audio", &probeResult{Streams: []struct {
			CodecName string `json:"codec_name"`
			CodecType string `json:"codec_type"`
		}{stream("h264", "video"), stream("pcm_s16le", "audio")}}, false},
		{"no audio", &probeResult{Streams: []struct {
			CodecName string `json:"codec_name"`
			CodecType string `json:"codec_type"`
		}{stream("h264", "video")}}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := audioIsAAC(tc.r); got != tc.wantOK {
				t.Errorf("audioIsAAC() = %v, want %v", got, tc.wantOK)
			}
		})
	}
}
