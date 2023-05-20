package restfulserver

import (
	"net/http"
	"whale/pkg/matcher"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/letjoy-club/mida-tool/proxy"
)

type Resp struct {
	Error string `json:"error"`
}

func Mount(r chi.Router) {
	r.Post("/start-matching", func(w http.ResponseWriter, r *http.Request) {
		matcher := matcher.Matcher{}
		err := matcher.Match(r.Context())
		if err != nil {
			render.JSON(w, r, Resp{Error: err.Error()})
		} else {
			render.JSON(w, r, Resp{})
		}
	})
	r.Handle("/ws", proxy.ProxyHandler())
}
