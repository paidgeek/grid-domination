package griddomination

import (
	"google.golang.org/appengine/datastore"
	"errors"
	"golang.org/x/net/context"
	"github.com/pquerna/ffjson/ffjson"
	"google.golang.org/appengine"
)

func GetPlayer(ctx context.Context, id string) *Player {
	if len(id) == 0 {
		return nil
	}

	key := datastore.NewKey(ctx, "Player", id, 0, nil)
	player := &Player{}

	if err := datastore.Get(ctx, key, player); err != nil {
		return nil
	}

	player.Id = id

	return player
}

func PutPlayer(ctx context.Context, player *Player) error {
	if player == nil {
		return errors.New("player is nil")
	}

	key := datastore.NewKey(ctx, "Player", player.Id, 0, nil)
	_, err := datastore.Put(ctx, key, player)

	return err
}


func GetChunk(ctx context.Context, id string) *Chunk {
	if len(id) == 0 {
		return nil
	}

	key := datastore.NewKey(ctx, "Chunk", id, 0, nil)
	chunk := &Chunk{}

	if err := datastore.Get(ctx, key, chunk); err != nil {
		return nil
	}

	chunk.Id = id
	chunk.Cells = make(map[string]Cell)
	err := ffjson.Unmarshal(chunk.CellsBinary, &chunk.Cells)

	if err != nil {
		return nil
	}
	
	return chunk
}

func GetChunks(ctx context.Context, ids []string) []*Chunk {
	chunks := make([]*Chunk, len(ids))
	keys := make([]*datastore.Key, len(ids))

	for i, id := range ids {
		keys[i] = datastore.NewKey(ctx, "Chunk", id, 0, nil)
	}

	if err := datastore.GetMulti(ctx, keys, chunks); err != nil {
		if multiErr, ok := err.(appengine.MultiError); ok {
			for i, err := range multiErr {
				chunk := chunks[i]

				if err != nil {
					if chunk == nil {
						chunk = &Chunk{}
					} else {

					}

					chunk.Id = ids[i]
					chunk.Cells = make(map[string]Cell)
					chunks[i] = chunk
				} else {
					chunk.Id = ids[i]
					chunk.Cells = make(map[string]Cell)
					ffjson.Unmarshal(chunk.CellsBinary, &chunk.Cells)
				}
			}
		}
	}

	return chunks
}

func PutChunk(ctx context.Context, chunk *Chunk) error {
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

	return err
}
