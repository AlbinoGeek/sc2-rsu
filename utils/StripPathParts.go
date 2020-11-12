package utils

import (
	"path/filepath"
	"strings"
)

// StripPathParts takes a given filepath and removes the specified number of
// path parts; if stripRight is positive, they are removed from the right
// (end), if negative, they are removed from the left (beginning).
func StripPathParts(path string, stripRight int) string {
	if stripRight < 0 {
		if parts := strings.SplitAfter(path, string(filepath.Separator)); len(parts) >= -stripRight {
			path = strings.Join(parts[len(parts)-1+stripRight:], "")
		}
	} else if stripRight > 0 {
		if parts := strings.SplitAfter(path, string(filepath.Separator)); len(parts) >= stripRight {
			path = strings.Join(parts[:len(parts)-stripRight], "")
		}
	}

	return path
}
