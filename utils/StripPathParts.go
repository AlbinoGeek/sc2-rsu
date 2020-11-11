package utils

import (
	"path/filepath"
	"strings"
)

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
