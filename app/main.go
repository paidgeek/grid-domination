package griddomination

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func init() {
	m := martini.Classic()

	m.Use(render.Renderer())
	m.Get("/", func() string {
		return "Hello, world!"
	})
	m.Post("/log_in/:access_token", logInHandler)
	m.Post("/grid/(?P<chunk_id>[0-9]+\\.[0-9]+)/(?P<cell_id>[0-9]+)", claimHandler)

	http.Handle("/", m)
}
