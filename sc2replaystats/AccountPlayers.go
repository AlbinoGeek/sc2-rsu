package sc2replaystats

// AccountPlayer represents the JSON format of a Player owned by a server Account
type AccountPlayer struct {
	ID      uint   `json:"players_id"`
	Default uint   `json:"default"`
	Player  Player `json:"player"`
}
