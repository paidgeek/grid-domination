package griddomination

import (
	"time"
	"strconv"
	"math"
	"math/rand"
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
	PlayerId      string `json:"player_id"`
	ClaimedAt     time.Time `json:"claimed_at"`
	ClaimDuration time.Duration `json:"claim_duration"`
	IsOwned       bool `json:"is_owned"`
	IsStealing    bool `json:"is_stealing"`
}

type Chunk struct {
	Id          string `datastore:"-" json:"id"`
	CellsBinary []byte `datastore:",noindex" json:"-"`
	Cells       map[string]*Cell `datastore:"-" json:"cells"`
	X           int64 `datastore:"-" json:"-"`
	Y           int64 `datastore:"-" json:"-"`
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

func (chunk *Chunk) Update() bool {
	hasChanged := false
	now := time.Now().UTC()

	for i := 0; i < 64; i++ {
		id := strconv.Itoa(i)

		if cell, ok := chunk.Cells[id]; ok {
			if now.Sub(cell.ClaimedAt) >= cell.ClaimDuration {
				cell.IsOwned = true
				cell.IsStealing = false
				chunk.Cells[id] = cell
				hasChanged = true
			}
		}
	}

	return hasChanged
}

func (player *Player) Reward() {
	player.Pixels += rand.Int63n(5) + 1
}

func (cell *Cell) SetClaimDurationForPlayer(player *Player) {
	cell.ClaimDuration = time.Minute
}

func (cell *Cell) GetTakeCost() int64 {
	diff := time.Now().UTC().Sub(cell.ClaimedAt.Add(cell.ClaimDuration))
	m := diff.Minutes()

	if m < 0 {
		m = 0
	}

	return int64(math.Ceil(math.Sqrt(m) + 5.0))
}
