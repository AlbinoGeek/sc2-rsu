package utils

import (
	"strconv"
	"strings"
)

// CompareSemVer returns 1 if the semantic version `new` is greater than `old`
// and -1 otherwise. A semantic version is one that follows the format "1.2.3".
// If either old or new startss with the letter "v", it will be stripped before
// comparison, along with any hyphenenated (-) suffix. (e.g: "v0.1-alpha")
func CompareSemVer(old, new string) int {
	// any version is better than no version
	if new == "" {
		return -1
	}

	if old == "" {
		return 1
	}

	// strip "v" prefixes
	if new[0] == 'v' {
		new = new[1:]
	}

	if old[0] == 'v' {
		old = old[1:]
	}

	nparts := strings.Split(strings.Split(new, "-")[0], ".")
	oparts := strings.Split(strings.Split(old, "-")[0], ".")

	// compare each part
	for i := range nparts {
		if len(oparts)-1 < i {
			// new version is longer
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
