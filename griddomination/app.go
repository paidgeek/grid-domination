package griddomination

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/v1/log_in/{access_token}", logInHandler).
	Methods("POST")

	r.Handle("/v1/grid/{chunk_id:-?[0-9]+\\.-?[0-9]+}/{cell_id:[0-9]+}", authenticator(http.HandlerFunc(claimHandler))).
	Methods("POST")

	r.Handle("/v1/grid/{chunk_ids:(-?[0-9]+\\.-?[0-9]+)(,-?[0-9]+\\.-?[0-9]+)*}", authenticator(http.HandlerFunc(getChunksHandler))).
	Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
}
