package sc2replaystats

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// GetReplayStatus tries to retrieve the replayID associated with a given
// replayQueueID -- returning an empty string if it's still processing
func GetReplayStatus(apikey string, replayQueueID string) (replayID string, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(
		"%s/replay/status/%s", APIRoot, replayQueueID,
	), nil)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	req.Header.Set("Authorization", apikey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("sc2replaystats API returned error: %v", resp.Status)
	}

	response := make(map[string]string)
	if err = jsoniter.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	// return the replay_id if present
	if rid, ok := response["replay_id"]; ok {
		return rid, nil
	}

	return "", nil
}
