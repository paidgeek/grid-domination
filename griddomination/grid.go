package griddomination

import (
	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
	"errors"
	"time"
	"google.golang.org/appengine"
)

func claimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := appengine.NewContext(r)
	player := gorillaContext.Get(r, "player").(*Player)

	chunkLocation := LocationFromId(vars["chunk_id"])
	cellIdStr := vars["cell_id"]
	cellId, err := strconv.ParseInt(vars["cell_id"], 10, 64)

	if err != nil || cellId < 0 || cellId >= 64 {
		responseError(w, "invalid ids", http.StatusBadRequest)
		return
	}

	// claim
	var chunk *Chunk

	err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		chunkId := chunkLocation.ToId()
		chunk = GetChunk(ctx, chunkId)
		hasChanged := false
		isSuccess := false

		if chunk == nil {
			chunk = &Chunk{}
			chunk.Id = chunkId
			chunk.Cells = make(map[string]Cell)
		} else {
			hasChanged = updateChunk(chunk)
		}

		if cell, ok := chunk.Cells[cellIdStr]; ok {
			if canClaim(player, chunk, &cell) {
				cell.PlayerId = player.Id
				cell.ClaimedAt = time.Now().UTC()
				cell.IsStealing = cell.IsOwned
				cell.IsOwned = false

				chunk.Cells[cellIdStr] = cell
				hasChanged = true
				isSuccess = true
			}
		} else {
			chunk.Cells[cellIdStr] = Cell{
				PlayerId:player.Id,
				ClaimedAt:time.Now().UTC(),
			}
			hasChanged = true
			isSuccess = true
		}

		if hasChanged {
			if err := PutChunk(ctx, chunk); err != nil {
				return err
			}
		}

		if isSuccess {
			return nil
		} else {
			return errors.New("")
		}
	}, nil)

	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJson(w, chunk)
}

func canClaim(player *Player, chunk *Chunk, cell *Cell) bool {
	return true
}

func updateChunk(chunk *Chunk) bool {
	hasChanged := false
	now := time.Now().UTC()

	for i := 0; i < 64; i++ {
		id := strconv.Itoa(i)

		if cell, ok := chunk.Cells[id]; ok {
			diffMinutes := now.Sub(cell.ClaimedAt).Minutes()

			if diffMinutes >= 0.09 {
				cell.IsOwned = true
				cell.IsStealing = false
				chunk.Cells[id] = cell
				hasChanged = true
			}
		}
	}

	return hasChanged
}
