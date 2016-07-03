package griddomination

import (
	"net/http"
	"google.golang.org/appengine"
	"github.com/gorilla/context"
	"time"
)

func authenticator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)

		playerId := r.Header.Get("X-Player-Id")
		sessionToken := r.Header.Get("X-Session-Token")

		if len(playerId) == 0 || len(sessionToken) == 0 {
			responseError(w, "invalid auth data", http.StatusUnauthorized)
			return
		}

		player := getPlayer(ctx, playerId)

		if player == nil || player.SessionToken != sessionToken {
			responseError(w, "invalid session token", http.StatusUnauthorized)
			return
		}

		if player.LastActionAt.IsZero() {
			player.LastActionAt = time.Now().UTC()
		}

		if (time.Now().UTC().Sub(player.LastActionAt) >= MaxAwayDuration) {
			responseError(w, "auto logged out", http.StatusUnauthorized)

			return
		}

		player.Id = playerId
		player.LastActionAt = time.Now().UTC()

		context.Set(r, "player", player)
		context.Set(r, "ctx", ctx)

		h.ServeHTTP(w, r)

		putPlayer(ctx, player)
	})
}
