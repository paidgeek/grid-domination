package griddomination

import (
	"time"
)

type Player struct {
	Id           string `datastore:"-" json:"id"`
	SessionToken string `datastore:",noindex" json:"-"`
	Score        int64 `json:"score"`
	Pixels       int64 `datastore:",noindex" json:"-"`
	LastActionAt time.Time `datastore:",noindex" json:"last_action_at"`
}

type PrivatePlayer struct {
	Id           string `json:"id"`
	SessionToken string `json:"session_token"`
	Score        int64 `json:"score"`
	Pixels       int64 `json:"pixels"`
	LastActionAt time.Time `json:"last_action_at"`
}

type Cell struct {
	PlayerId   string `json:"player_id"`
	ClaimedAt  time.Time `json:"claimed_at"`
	IsOwned    bool `json:"is_owned"`
	IsStealing bool `json:"is_stealing"`
}

type Chunk struct {
	Id          string `datastore:"-" json:"id"`
	CellsBinary []byte `datastore:",noindex" json:"-"`
	Cells       map[string]Cell `datastore:"-" json:"cells"`
}

func (player *Player) ToPrivatePlayer() *PrivatePlayer {
	return &PrivatePlayer{
		Id:player.Id,
		SessionToken:player.SessionToken,
		Score:player.Score,
		Pixels:player.Pixels,
		LastActionAt:player.LastActionAt,
	}
}

type ClaimMessage struct {
	Chunk  *Chunk `json:"chunk"`
	Player *PrivatePlayer `json:"player"`
}

type GetChunksMessage struct {
	Chunks []*Chunk `json:"chunks"`
	Player *PrivatePlayer `json:"player"`
}
