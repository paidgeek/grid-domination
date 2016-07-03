package griddomination

import (
	"net/http"
	"google.golang.org/appengine"
	"github.com/gorilla/context"
	"strings"
)

func authenticator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 {
			responseError(w, "no Authorization header found", http.StatusUnauthorized)
			return
		}

		auth := strings.Split(authHeader, ".")

		if len(auth) != 2 {
			responseError(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}

		playerId := auth[0]
		sessionToken := auth[1]

		player := GetPlayer(ctx, playerId)

		if player == nil || player.SessionToken != sessionToken {
			responseError(w, "invalid session token", http.StatusUnauthorized)
			return
		}

		player.Id = playerId

		context.Set(r, "player", player)

		h.ServeHTTP(w, r)
	})
}
