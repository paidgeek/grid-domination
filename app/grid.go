package griddomination

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func claimHandler(req *http.Request, r render.Render, params martini.Params) {
	chunkId := params["chunk_id"]
	cellId := params["cell_id"]

	r.Text(200, chunkId + ", " + cellId)
}
