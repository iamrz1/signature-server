package api

import (
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"net/http"
	cerror "signature-server/error"
	cjson "signature-server/json"
)

type system struct {
	chi.Router
}

// NewSystemHandler ...
func NewSystemHandler() http.Handler {
	h := &system{
		chi.NewRouter(),
	}
	h.registerMiddleware()
	h.registerEndpoints()
	return h
}

func (api *system) registerMiddleware() {
	api.Use(chiware.Logger)
}

func (api *system) registerEndpoints() {
	api.Group(func(r chi.Router) {
		r.Get("/live", api.live)
		r.Get("/ready", api.ready)
		r.Get("/panic", api.systemPanic)
	})
}

func (api *system) live(w http.ResponseWriter, r *http.Request) {
	cjson.ServeData(w, cjson.Object{
		"live": "ok",
	})
}

func (api *system) ready(w http.ResponseWriter, r *http.Request) {
	cjson.ServeData(w, cjson.Object{
		"ready": "ok",
	})
}

func (api *system) systemPanic(w http.ResponseWriter, r *http.Request) {
	cjson.ServeError(w, cerror.NewAPIError(http.StatusInternalServerError, "system panic", nil))
}
