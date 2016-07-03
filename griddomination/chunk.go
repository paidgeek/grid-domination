package griddomination

import (
	"time"
)

type Cell struct {
	PlayerId   string `json:"player_id"`
	ClaimedAt  time.Time `json:"claimed_at"`
	IsOwned    bool `json:"is_owned"`
	IsStealing bool `json:"is_stealing"`
}

type Chunk struct {
	Id    string `datastore:"-" json:"id"`
	CellsBinary []byte `datastore:",noindex" json:"-"`
	Cells map[string]Cell `datastore:"-" json:"cells"`
}
