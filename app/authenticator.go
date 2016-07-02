package griddomination

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"github.com/gorilla/context"
	"strings"
)

func authenticator(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := appengine.NewContext(r)
			authHeader := r.Header.Get("Authorization")

			if len(authHeader) == 0 {
				responseError(w, "", http.StatusUnauthorized)
				return
			}

			auth := strings.Split(authHeader, ".")

			if len(auth) != 2 {
				responseError(w, "", http.StatusUnauthorized)
				return
			}

			playerId := auth[0]
			sessionToken := auth[1]
			playerKey := datastore.NewKey(ctx, "Player", playerId, 0, nil)

			var player Player
			err := datastore.Get(ctx, playerKey, &player)

			if err != nil || player.SessionToken != sessionToken {
				responseError(w, "", http.StatusUnauthorized)
				return
			}

			player.Id = playerId

			context.Set(r, "player", player)
			context.Set(r, "ctx", ctx)

			h.ServeHTTP(w, r)
		})
}
