package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func FindDirectoriesBySuffix(root string, suffix string, ignoreAccessErrors bool) (results []string, err error) {
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
