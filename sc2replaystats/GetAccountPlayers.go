package sc2replaystats

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// GetAccountPlayers returns an slice of AccountPlayers from the server,
// showing all those Accounts/Toons associated with the given API key
func (client *Client) GetAccountPlayers() (players []AccountPlayer, err error) {
	result, err := client.requestBytes(http.MethodGet, "account/players", "", nil)

	players = make([]AccountPlayer, 0)
	err = jsoniter.Unmarshal(result, &players)

	return
}
