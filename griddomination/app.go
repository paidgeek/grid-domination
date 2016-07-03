package griddomination

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	r.HandleFunc("/log_in/{access_token}", logInHandler).
	Methods("POST")
	r.Handle("/grid/{chunk_id:[0-9]+\\.[0-9]+}/{cell_id:[0-9]+}", authenticator(http.HandlerFunc(claimHandler)))

	http.Handle("/", r)
}
