package convert

import "testing"

func TestIsSupportedInput(t *testing.T) {
	cases := []struct {
		name string
		path string
		want bool
	}{
		{name: "lowercase mov", path: "a.mov", want: true},
		{name: "uppercase MOV", path: "a.MOV", want: true},
		{name: "mixed case Mov", path: "a.Mov", want: true},
		{name: "lowercase qt", path: "a.qt", want: true},
		{name: "uppercase QT", path: "a.QT", want: true},
		{name: "mixed case Qt", path: "a.Qt", want: true},
		{name: "dotted base name", path: "my.vacation.video.mov", want: true},
		{name: "path with dirs", path: "/tmp/dir/clip.qt", want: true},
		{name: "empty path", path: "", want: false},
		{name: "no extension", path: "a", want: false},
		{name: "mp4", path: "a.mp4", want: false},
		{name: "txt", path: "a.txt", want: false},
		{name: "mkv", path: "a.mkv", want: false},
		{name: "m4v", path: "a.m4v", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsSupportedInput(tc.path); got != tc.want {
				t.Errorf("IsSupportedInput(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}
