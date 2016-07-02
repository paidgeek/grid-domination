package main

import (
	"errors"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
)

type Player struct {
	Id           string `datastore:"-" json:"id"`
	SessionToken string `datastore:",noindex" json:"session_token"`
}

var PlayerDatabase *PlayerDatastoreDatabase = &PlayerDatastoreDatabase{}

type PlayerDatastoreDatabase struct{}

func (db PlayerDatastoreDatabase) GetPlayer(ctx context.Context, id string) *Player {
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

func (db PlayerDatastoreDatabase) PutPlayer(ctx context.Context, player *Player) error {
	if player == nil {
		return errors.New("player is nil")
	}

	key := datastore.NewKey(ctx, "Player", player.Id, 0, nil)
	_, err := datastore.Put(ctx, key, player)

	return err
}
