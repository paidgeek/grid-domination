package griddomination

type Player struct {
	Id string `datastore:"-" json:"id"`
	SessionToken string `datastore:",noindex" json:"session_token"`
}
