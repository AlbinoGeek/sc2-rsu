package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// StripPathParts takes a given filepath and removes the specified number of
// path parts; if stripRight is positive, they are removed from the right
// (end), if negative, they are removed from the left (beginning).
func StripPathParts(path string, stripRight int) string {
	abs := filepath.IsAbs(path)

	// "normalize" paths by removing trailing separators
	if path[len(path)-1] == filepath.Separator {
		path = path[:len(path)-1]
	}

	parts := strings.Split(path, string(filepath.Separator))
	if stripRight < 0 {
		// stripping more parts than we have
		if len(parts) <= -stripRight {
			return ""
		}

		return strings.Join(parts[-stripRight:], string(filepath.Separator))
	}

	if stripRight > 0 {
		// stripping more parts than we have
		if len(parts) <= stripRight {
			if abs {
				return "/"
			}
			return ""
		}

		path = strings.Join(parts[:len(parts)-stripRight], string(filepath.Separator))

		// restore absolute-ness for re-assembled path
		if abs && (len(path) == 0 || path[0] != filepath.Separator) {
			return fmt.Sprintf("%c%s", filepath.Separator, path)
		}
	}

	return path
}
