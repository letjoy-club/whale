package restfulserver

import (
	"encoding/json"
	"net/http"
	"whale/pkg/matcher"
	"whale/pkg/models"
	"whale/pkg/modelutil"

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
	r.Post("/notify-new-motion-offer", func(w http.ResponseWriter, r *http.Request) {
		param := models.NotifyNewMotionOfferMessageParam{}
		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			render.JSON(w, r, Resp{Error: err.Error()})
			return
		}
		err := modelutil.NotifyNewMotionOffer(r.Context(), param.Begin, param.End)
		if err != nil {
			render.JSON(w, r, Resp{Error: err.Error()})
		} else {
			render.JSON(w, r, Resp{})
		}
	})
	r.Post("/clear-out-date-offer", func(w http.ResponseWriter, r *http.Request) {
		if err := modelutil.ClearOutDateMotionOffer(r.Context()); err != nil {
			render.JSON(w, r, Resp{Error: err.Error()})
		} else {
			render.JSON(w, r, Resp{})
		}
	})
	r.Post("/refresh-duration-constraint", func(w http.ResponseWriter, r *http.Request) {
		if err := modelutil.RefreshDurationConstraint(r.Context()); err != nil {
			render.JSON(w, r, Resp{Error: err.Error()})
		} else {
			render.JSON(w, r, Resp{})
		}
	})
	r.Handle("/ws", proxy.ProxyHandler())
}
