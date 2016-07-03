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
	"strings"
	"math/rand"
	"fmt"
)

func claimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := gorillaContext.Get(r, "ctx").(context.Context)
	player := gorillaContext.Get(r, "player").(*Player)

	chunkId := vars["chunk_id"]
	chunkX, chunkY, chunkErr := locationFromId(chunkId)
	cellIdStr := vars["cell_id"]
	cellId, cellErr := strconv.ParseInt(vars["cell_id"], 10, 64)

	if cellErr != nil || chunkErr != nil || cellId < 0 || cellId >= 64 {
		responseError(w, "invalid ids", http.StatusBadRequest)
		return
	}

	// claim
	var chunk *Chunk

	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		chunk = getChunk(ctx, chunkId)
		hasChanged := false
		isSuccess := false

		chunk.X = chunkX
		chunk.Y = chunkY
		chunk.Update()

		if cell, ok := chunk.Cells[cellIdStr]; ok {
			if canSteal(&cell, player) && canClaim(ctx, player, chunk, cellId) {
				// steal
				cell.PlayerId = player.Id
				cell.ClaimedAt = time.Now().UTC()
				cell.IsStealing = cell.IsOwned
				cell.IsOwned = false

				chunk.Cells[cellIdStr] = cell
				hasChanged = true
				isSuccess = true
			}
		} else {
			if canClaim(ctx, player, chunk, cellId) {
				// first time claim
				chunk.Cells[cellIdStr] = Cell{
					PlayerId:player.Id,
					ClaimedAt:time.Now().UTC(),
				}
				hasChanged = true
				isSuccess = true
			}
		}

		if hasChanged {
			if err := putChunk(ctx, chunk); err != nil {
				return err
			}
		}

		if isSuccess {
			player.Score++
			player.Pixels += int64(rand.Intn(5))

			return nil
		} else {
			return errors.New("cannot claim")
		}
	}, &datastore.TransactionOptions{XG:true})

	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJson(w, ClaimMessage{
		Chunk: chunk,
		Player:player.ToPrivatePlayer(),
	})
}

func getChunksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := gorillaContext.Get(r, "ctx").(context.Context)
	player := gorillaContext.Get(r, "player").(*Player)
	coords := strings.Split(vars["chunk_ids"], ",")

	chunks := getChunks(ctx, coords)

	for _, chunk := range chunks {
		chunk.Update()
	}

	responseJson(w, GetChunksMessage{
		Chunks: chunks,
		Player: player.ToPrivatePlayer(),
	})
}

func canSteal(cell *Cell, player *Player) bool {
	return cell.IsOwned && cell.PlayerId != player.Id
}

func canClaim(ctx context.Context, player *Player, chunk *Chunk, cellId int64) bool {
	cx := cellId % 8
	cy := cellId / 8

	if checkCell(ctx, cx - 1, cy, chunk, player.Id) ||
	checkCell(ctx, cx + 1, cy, chunk, player.Id) ||
	checkCell(ctx, cx, cy - 1, chunk, player.Id) ||
	checkCell(ctx, cx, cy + 1, chunk, player.Id) {
		return true
	}

	return false
}

func checkCell(ctx context.Context, x int64, y int64, chunk *Chunk, playerId string) bool {
	changeChunk := false
	chunkX := chunk.X
	chunkY := chunk.Y

	if x < 0 {
		x = 7
		chunkX--
		changeChunk = true
	} else if x >= 8 {
		x = 0
		chunkX++
		changeChunk = true
	} else if y < 0 {
		y = 7
		chunkY--
		changeChunk = true
	} else if y >= 8 {
		y = 0
		chunkY++
		changeChunk = true
	}

	if changeChunk {
		chunk = getChunk(ctx, fmt.Sprintf("%v.%v", chunkX, chunkY))
		chunk.Update()
	}

	if cell, ok := chunk.Cells[fmt.Sprint(y * 8 + x)]; ok {
		if cell.PlayerId == playerId {
			return true
		}
	}

	return false
}
