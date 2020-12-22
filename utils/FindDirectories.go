package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FindDirectoriesBySuffix searches the given root recursively and returns
// all paths where the last directory component has the supplied suffix,
// optionally ignoring errors "access denied" and "permission denied"
func FindDirectoriesBySuffix(root, suffix string, ignoreAccessErrors bool) (results []string, err error) {
	results = make([]string, 0)

	err = filepath.Walk(root, func(path string, info os.FileInfo, wErr error) error {
		if wErr != nil {
			if ignoreAccessErrors && (strings.Contains(wErr.Error(), "access denied") ||
				strings.Contains(wErr.Error(), "permission denied")) {
				return nil
			}

			return wErr
		}

		if info.IsDir() && strings.HasSuffix(path, suffix) {
			results = append(results, path)
		}

		return nil
	})

	return
}
