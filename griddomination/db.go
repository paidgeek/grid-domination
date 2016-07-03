package griddomination

import (
	"google.golang.org/appengine/datastore"
	"errors"
	"golang.org/x/net/context"
	"github.com/pquerna/ffjson/ffjson"
	"google.golang.org/appengine/memcache"
	"time"
)

var MemcachePlayerExpiration time.Duration = 5 * time.Minute
var MemcacheChunkExpiration time.Duration = 5 * time.Minute

func getPlayer(ctx context.Context, id string) *Player {
	if len(id) == 0 {
		return nil
	}

	player := &Player{}

	if _, err := memcache.Gob.Get(ctx, "player." + id, player); err == nil {
		return player
	}

	key := datastore.NewKey(ctx, "Player", id, 0, nil)

	if err := datastore.Get(ctx, key, player); err != nil {
		return nil
	}

	player.Id = id

	memcache.Gob.Set(ctx, &memcache.Item{
		Key:"player." + id,
		Object:&player,
		Expiration:MemcachePlayerExpiration,
	})

	return player
}

func putPlayer(ctx context.Context, player *Player) error {
	if player == nil {
		return errors.New("player is nil")
	}

	key := datastore.NewKey(ctx, "Player", player.Id, 0, nil)
	_, err := datastore.Put(ctx, key, player)

	if err == nil {
		memcache.Gob.Set(ctx, &memcache.Item{
			Key:"player." + player.Id,
			Object:&player,
			Expiration:MemcachePlayerExpiration,
		})
	}

	return err
}

func getChunk(ctx context.Context, id string) *Chunk {
	if len(id) == 0 {
		return nil
	}

	chunk := &Chunk{}

	if _, err := memcache.Gob.Get(ctx, "chunk." + id, chunk); err == nil {
		return chunk
	}

	key := datastore.NewKey(ctx, "Chunk", id, 0, nil)

	if err := datastore.Get(ctx, key, chunk); err == nil {
		err := ffjson.Unmarshal(chunk.CellsBinary, &chunk.Cells)

		if err != nil {
			return nil
		}
	}

	chunk.Id = id
	chunk.Cells = make(map[string]Cell)

	memcache.Gob.Set(ctx, &memcache.Item{
		Key:"chunk." + id,
		Object:chunk,
		Expiration:MemcacheChunkExpiration,
	})

	return chunk
}

func getChunks(ctx context.Context, ids []string) []*Chunk {
	chunks := make([]*Chunk, len(ids))
	var keys []*datastore.Key

	for i, id := range ids {
		chunk := &Chunk{}

		if _, err := memcache.Gob.Get(ctx, "chunk." + id, chunk); err == nil {
			chunks[i] = chunk
		} else {
			keys = append(keys, datastore.NewKey(ctx, "Chunk", id, 0, nil))
		}
	}

	if len(keys) == 0 {
		return chunks
	}

	dbChunks := make([]*Chunk, len(keys))
	datastore.GetMulti(ctx, keys, dbChunks)

	for i, chunk := range chunks {
		if chunk == nil {
			id := ids[i]

			for j, key := range keys {
				if key.StringID() == id {
					chunk = dbChunks[j]

					break
				}
			}
		}

		if chunk == nil {
			chunk = &Chunk{}
		}

		chunk.Id = ids[i]
		chunk.Cells = make(map[string]Cell)

		if len(chunk.CellsBinary) != 0 {
			ffjson.Unmarshal(chunk.CellsBinary, &chunk.Cells)
			chunk.Update()
		}

		memcache.Gob.Set(ctx, &memcache.Item{
			Key:"chunk." + chunk.Id,
			Object:chunk,
			Expiration:MemcacheChunkExpiration,
		})

		chunks[i] = chunk
	}

	return chunks
}

func putChunk(ctx context.Context, chunk *Chunk) error {
	if chunk == nil {
		return errors.New("chunk is nil")
	}

	var err error
	chunk.CellsBinary, err = ffjson.Marshal(chunk.Cells)

	if err != nil {
		return err
	}

	key := datastore.NewKey(ctx, "Chunk", chunk.Id, 0, nil)
	_, err = datastore.Put(ctx, key, chunk)

	if err == nil {
		memcache.Gob.Set(ctx, &memcache.Item{
			Key:"chunk." + chunk.Id,
			Object:chunk,
			Expiration:MemcacheChunkExpiration,
		})
	}

	return err
}
