package sc2replaystats

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kataras/golog"
)

// GetReplayStatus tries to retrieve the replayID associated with a given
// replayQueueID -- returning an empty string if it's still processing
func (client *Client) GetReplayStatus(replayQueueID string) (replayID string, err error) {
	res, err := client.requestMap(http.MethodGet, fmt.Sprintf("replay/status/%s", replayQueueID), "", nil)

	if err != nil {
		return "", fmt.Errorf("GetReplayStatus: %v", err)
	}

	// return the replay_id if present
	if rid, ok := res["replay_id"]; ok {
		replayID = rid
	}

	if e, ok := res["error"]; ok {
		if strings.HasPrefix(e, "Duplicate Replay") {
			parts := strings.Split(e, ": ")
			if len(parts) == 2 {
				replayID = parts[1]
				return
			}
		}

		golog.Debugf("sc2replaystats failure: %v", e)

		err = fmt.Errorf("replay processing failed")
	}

	return
}
