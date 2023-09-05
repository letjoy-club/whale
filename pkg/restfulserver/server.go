package restfulserver

import (
	"net/http"
	"strconv"
	"whale/pkg/matcher"
	"whale/pkg/modelutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-module/carbon"
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
	// hour 表示把前 X 小时的 offer 通知给用户
	r.Post("/notify-new-motion-offer/{hour}", func(w http.ResponseWriter, r *http.Request) {
		hour := chi.URLParam(r, "hour")
		hourInt, _ := strconv.ParseInt(hour, 10, 64)
		if hourInt == 0 {
			render.JSON(w, r, Resp{Error: "hour must be greater than 0"})
			return
		}
		// 防止机器有一些时差，保证 offer 在是在这个整点（一般是在整点调用通知，比如 10:00，如果机器此时时间差1-2s 没到 10:00，就会有问题，所以拨快一点）
		stopHour := carbon.Now().AddMinutes(5).StartOfHour()
		startHour := stopHour.AddHours(int(-hourInt))

		err := modelutil.NotifyNewMotionOffer(r.Context(), startHour.ToStdTime(), stopHour.ToStdTime())
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
