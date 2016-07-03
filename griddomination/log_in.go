package griddomination

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"time"
)

func logInHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accessToken := vars["access_token"]
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	grapResponse, err := client.Get(GraphAccessTokenUrl + accessToken)
	defer grapResponse.Body.Close()
	if err != nil {
		responseError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(grapResponse.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	graphBody := make(map[string]interface{})
	json.Unmarshal(body, &graphBody)
	graphData := graphBody["data"].(map[string]interface{})

	if graphData["app_id"] != FacebookAppId {
		responseError(w, "invalid app", http.StatusUnauthorized)
		return
	}

	userId := graphData["user_id"].(string)
	player := getPlayer(ctx, userId)

	if player == nil {
		player = &Player{
			Id:           userId,
			SessionToken: generateSessionToken(),
		}
	} else {
		player.Id = userId
		player.SessionToken = generateSessionToken()
	}

	if player.SessionToken == "" {
		responseError(w, "", http.StatusUnauthorized)
		return
	}

	player.LastActionAt = time.Now().UTC()

	if err = putPlayer(ctx, player); err != nil {
		responseError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	responseJson(w, player.ToPrivatePlayer())
}
