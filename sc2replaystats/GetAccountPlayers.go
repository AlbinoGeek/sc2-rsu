package sc2replaystats

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// GetAccountPlayers returns an slice of AccountPlayers from the server,
// showing all those Accounts/Toons associated with the given API key
func (client *Client) GetAccountPlayers() (players []AccountPlayer, err error) {
	result, err := client.requestBytes(http.MethodGet, "account/players", "", nil)

	if err != nil {
		return nil, fmt.Errorf("GetAccountPlayers: %v", err)
	}

	players = make([]AccountPlayer, 0)
	err = jsoniter.Unmarshal(result, &players)

	return
}
