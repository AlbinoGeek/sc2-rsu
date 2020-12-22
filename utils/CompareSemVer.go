package utils

import (
	"strconv"
	"strings"
)

// CompareSemVer returns 1 if the semantic version `newVer` is greater than
// `oldVer` and -1 otherwise. A semantic version is one that follows the
// format "1.2.3". If either oldVer or newVer startss with the letter "v",
// it will be stripped before comparison, along with any hyphenenated (-)
// suffix. (e.g: "v0.1-alpha" becomes "0.1")
func CompareSemVer(oldVer, newVer string) int {
	// any version is better than no version
	if newVer == "" {
		return -1
	}

	if oldVer == "" {
		return 1
	}

	// strip "v" prefixes
	if newVer[0] == 'v' {
		newVer = newVer[1:]
	}

	if oldVer[0] == 'v' {
		oldVer = oldVer[1:]
	}

	nparts := strings.Split(strings.Split(newVer, "-")[0], ".")
	oparts := strings.Split(strings.Split(oldVer, "-")[0], ".")

	// compare each part
	for i := range nparts {
		if len(oparts)-1 < i {
			// newVer is longer
			return 1
		}

		// parts should be numeric, otherwise they are skipped
		n, err1 := strconv.Atoi(nparts[i])
		o, err2 := strconv.Atoi(oparts[i])

		if err1 == nil && err2 == nil && n > o {
			return 1
		}
	}

	return -1
}
