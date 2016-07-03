package griddomination

type Player struct {
	Id           string `datastore:"-" json:"id"`
	SessionToken string `datastore:",noindex" json:"session_token"`
	Score        int64 `json:"score"`
	Pixels       int64 `datastore:",noindex" json:"pixels"`
}
