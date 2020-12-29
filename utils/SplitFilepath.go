package utils

import (
	"path/filepath"
	"strings"
)

// SplitFilepath takes a relative or absolute path to a file, and returns the
// separated path, filename, and file extension. If an extension is present,
// it will start with a dot character, otherwise it will be an empty string.
func SplitFilepath(name string) (fpath, fname, ext string) {
	fpath = filepath.Dir(name)
	ext = filepath.Ext(name)

	// Trim path
	fname = strings.TrimPrefix(name, fpath)

	// Trim path sep
	fname = strings.TrimPrefix(fname, string(filepath.Separator))

	// Trim extension
	fname = strings.TrimSuffix(fname, ext)

	return
}
