package sc2replaystats

import (
	"fmt"
	"time"
)

// Player represents the JSON format that the server stores Toon/Characters in
type Player struct {
	ID                  uint      `json:"players_id"`
	Name                string    `json:"players_name"`
	BattleNetURL        string    `json:"battle_net_url"`
	BattleTagName       string    `json:"battle_tag_name"`
	BattleTagID         uint      `json:"battle_tag_id"`
	CharacterID         uint      `json:"character_link_id"`
	LegacyBattleTagName string    `json:"legacy_battle_tag_name"`
	LegacyLinkID        uint      `json:"legacy_link_id"`
	LegacyLinkRealm     uint      `json:"legacy_link_realm"`
	ReplaysURL          string    `json:"players_replays_url"`
	Updated             time.Time `json:"updated_at"`
}

// BattleTag returns the assembled friends tag of a given Toon/Character
func (p Player) BattleTag() string {
	return fmt.Sprintf("%s#%d", p.BattleTagName, p.BattleTagID)
}
