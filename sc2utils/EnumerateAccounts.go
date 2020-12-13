package sc2utils

import (
	"fmt"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

var multiplayerSuffix = "ultiplayer"

// EnumerateAccounts searches a given replaysRoot and returns a slice
// containing the accounts and toons which could be found, in the format
// "AccountID/ToonID" -- which is the same as their replay folder path.
func EnumerateAccounts(replaysRoot string) (accountIDs []string, err error) {
	paths, err := utils.FindDirectoriesBySuffix(replaysRoot, multiplayerSuffix, true)
	if err != nil {
		return nil, fmt.Errorf("FindDirectory error: %v", err)
	}

	i := 0
	uniq := make(map[string]struct{})

	for _, p := range paths {
		// strip "/Replays/Multiplayer" suffix
		p = utils.StripPathParts(p, 2)

		if _, duplicate := uniq[p]; !duplicate {
			uniq[p] = struct{}{}

			// strip replaysRoot prefix from resulting path
			paths[i] = p[len(replaysRoot):]
			i++
		}
	}

	return paths[:i], nil
}
