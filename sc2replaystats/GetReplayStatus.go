package sc2replaystats

import (
	"fmt"
	"net/http"
)

// GetReplayStatus tries to retrieve the replayID associated with a given
// replayQueueID -- returning an empty string if it's still processing
func (client *Client) GetReplayStatus(replayQueueID string) (replayID string, err error) {
	res, err := client.request(http.MethodGet, fmt.Sprintf("replay/status/%s", replayQueueID), "", nil)

	// return the replay_id if present
	if rid, ok := res["replay_id"]; ok {
		replayID = rid
	}

	return
}
