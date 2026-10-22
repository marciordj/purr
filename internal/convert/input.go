package convert

import (
	"path/filepath"
	"strings"
)

// IsSupportedInput reports whether path names a supported input file:
// .mov or .qt (both QuickTime containers), case-insensitive.
func IsSupportedInput(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mov" || ext == ".qt"
}
