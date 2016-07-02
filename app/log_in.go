package griddomination

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func logInHandler(req *http.Request, r render.Render, params martini.Params) {
	accessToken := params["access_token"]
	ctx := appengine.NewContext(req)
	client := urlfetch.Client(ctx)

	grapResponse, err := client.Get(GraphAccessTokenUrl + accessToken)
	defer grapResponse.Body.Close()
	if err != nil {
		r.JSON(401, map[string]interface{}{"error": err.Error()})
		return
	}

	body, err := ioutil.ReadAll(grapResponse.Body)
	if err != nil {
		r.JSON(401, map[string]interface{}{"error": err.Error()})
		return
	}

	graphBody := make(map[string]interface{})
	json.Unmarshal(body, &graphBody)
	graphData := graphBody["data"].(map[string]interface{})

	if graphData["app_id"] != FacebookAppId {
		r.JSON(401, map[string]interface{}{"error": "invalid app"})
		return
	}

	userId := graphData["user_id"].(string)
	playerKey := datastore.NewKey(ctx, "Player", userId, 0, nil)

	player := &Player{}
	err = datastore.Get(ctx, playerKey, player)

	if err != nil {
		player = &Player{
			SessionToken: generateSessionToken(),
		}
	} else {
		player.SessionToken = generateSessionToken()
	}

	_, err = datastore.Put(ctx, playerKey, player)

	if err != nil {
		r.JSON(401, map[string]interface{}{"error": err.Error()})
		return
	}

	r.JSON(200, player)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1 << letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func generateSessionToken() string {
	n := 64
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
