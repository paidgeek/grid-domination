package griddomination

import (
	"errors"
	"fmt"
	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func claimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := gorillaContext.Get(r, "ctx").(context.Context)
	player := gorillaContext.Get(r, "player").(*Player)
	isTake := false

	if take, err := strconv.ParseBool(r.URL.Query().Get("take")); err == nil {
		isTake = take
	}

	chunkId := vars["chunk_id"]
	cellId, cellErr := strconv.ParseInt(vars["cell_id"], 10, 64)
	if cellErr != nil || cellId < 0 || cellId >= 64 {
		responseError(w, "invalid ids", http.StatusBadRequest)
		return
	}

	chunk := getChunk(ctx, chunkId)
	if chunk == nil {
		responseError(w, "invalid chunk", http.StatusBadRequest)
		return
	}

	err := claim(ctx, cellId, chunk, player, isTake)

	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJson(w, ClaimMessage{
		Chunk:  chunk,
		Player: player.ToPrivatePlayer(),
	})
}

func takeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := gorillaContext.Get(r, "ctx").(context.Context)
	player := gorillaContext.Get(r, "player").(*Player)

	chunkId := vars["chunk_id"]
	cellIdStr := vars["cell_id"]
	cellId, cellErr := strconv.ParseInt(vars["cell_id"], 10, 64)

	if cellErr != nil || cellId < 0 || cellId >= 64 {
		responseError(w, "invalid ids", http.StatusBadRequest)
		return
	}

	var chunk *Chunk

	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		chunk = getChunk(ctx, chunkId)

		if chunk == nil {
			return errors.New("invalid chunk")
		}

		chunk.Update()

		var cell *Cell

		if c, ok := chunk.Cells[cellIdStr]; ok {
			cell = c
		} else {
			return errors.New("invalid cell")
		}

		if cell.IsOwned || cell.PlayerId != player.Id {
			return errors.New("cannot take")
		}

		cost := cell.GetTakeCost()

		if player.Pixels < cost {
			return errors.New("not enough pixels")
		}

		cell.IsStealing = false
		cell.IsOwned = true
		chunk.Cells[cellIdStr] = cell

		if err := putChunk(ctx, chunk); err != nil {
			return err
		}

		player.Score++
		player.Pixels -= cost
		player.Reward()

		if err := putPlayer(ctx, player); err != nil {
			return err
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJson(w, ClaimMessage{
		Chunk:  chunk,
		Player: player.ToPrivatePlayer(),
	})
}

func getChunksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := gorillaContext.Get(r, "ctx").(context.Context)
	player := gorillaContext.Get(r, "player").(*Player)
	coords := strings.Split(vars["chunk_ids"], ",")

	chunks := getChunks(ctx, coords)

	for _, chunk := range chunks {
		if chunk != nil {
			chunk.Update()
		}
	}

	responseJson(w, GetChunksMessage{
		Chunks: chunks,
		Player: player.ToPrivatePlayer(),
	})
}

func claim(ctx context.Context, cellId int64, chunk *Chunk, player *Player, isTake bool) error {
	cellIdStr := strconv.FormatInt(cellId, 10)

	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		isSuccess := false
		isSteal := false
		chunk.Update()

		cell := chunk.Cells[cellIdStr]

		if player.Score == 0 {
			isSuccess = cell == nil
		} else {
			if cell == nil && hasNeighbours(ctx, player, chunk, cellId) {
				isSuccess = true
			} else if cell != nil && cell.IsOwned && cell.PlayerId != player.Id && hasNeighbours(ctx, player, chunk, cellId) {
				isSuccess = true
				isSteal = true
			}
		}

		if isSuccess {
			if cell == nil {
				cell = &Cell{}
			}

			cell.SetClaimDurationForPlayer(player)
			cost := cell.GetTakeCost()

			if isTake {
				if player.Pixels < cost {
					return errors.New("not enough pixels")
				}

				player.Pixels -= cost
				cell.IsOwned = true
			}

			cell.PlayerId = player.Id
			cell.ClaimedAt = time.Now().UTC()
			cell.IsStealing = isSteal

			chunk.Cells[cellIdStr] = cell

			if isSteal {
				otherPlayer := getPlayer(ctx, cell.PlayerId)

				if otherPlayer == nil {
					return errors.New("other player is nil")
				}

				otherPlayer.Score--

				if err := putPlayer(ctx, otherPlayer); err != nil {
					return err
				}
			}

			player.Score++
			player.Reward()

			if err := putPlayer(ctx, player); err != nil {
				return err
			}

			if err := putChunk(ctx, chunk); err != nil {
				return err
			}

			return nil
		}

		return errors.New("cannot claim")
	}, &datastore.TransactionOptions{XG: true})

	return err
}

func hasNeighbours(ctx context.Context, player *Player, chunk *Chunk, cellId int64) bool {
	cx := cellId % 8
	cy := cellId / 8

	if checkCell(ctx, cx-1, cy, chunk, player.Id) ||
		checkCell(ctx, cx+1, cy, chunk, player.Id) ||
		checkCell(ctx, cx, cy-1, chunk, player.Id) ||
		checkCell(ctx, cx, cy+1, chunk, player.Id) {
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

		if chunk != nil {
			chunk.Update()
		}
	}

	if cell, ok := chunk.Cells[strconv.FormatInt(y*8+x, 10)]; ok {
		if cell.PlayerId == playerId {
			return true
		}
	}

	return false
}
