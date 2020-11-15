package sc2utils

import (
	"fmt"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

// FindReplaysRoot recursively searches a given scanRoot directory, returning
// all paths that seem like they could hold StarCraft II replays, organized
// by sub-directories containing "<accountID>/<toonID>/Replays/<gameType>"
func FindReplaysRoot(scanRoot string) (replayRoots []string, err error) {
	paths, err := utils.FindDirectoriesBySuffix(scanRoot, multiplayerSuffix, true)
	if err != nil {
		return nil, fmt.Errorf("FindDirectory error: %v", err)
	}

	i := 0
	uniq := make(map[string]struct{})
	for _, p := range paths {
		// strip "accountID/toonID/Replays/Multiplayer" suffix
		p = utils.StripPathParts(p, 4)
		if _, duplicate := uniq[p]; !duplicate && p != "/" {
			uniq[p] = struct{}{}
			paths[i] = p
			i++
		}
	}

	return paths[:i], nil
}
