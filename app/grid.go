package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
)

func claimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chunkId := vars["chunk_id"]
	cellId := vars["cell_id"]

	player := context.Get(r, "player").(*Player)

	responseJson(w, map[string]interface{}{"data": chunkId + ", " + cellId, "id":player.Id})
}
