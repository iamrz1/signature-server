package api

import (
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"net/http"
	cerror "signature-server/error"
	cjson "signature-server/json"
	"signature-server/middleware"
	"signature-server/util"
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
	api.Use(util.GenReqID)
	api.Use(chiware.Logger)
	api.Use(middleware.RequestLogger(true))
	api.Use(middleware.ResponseLogger(true))
}

func (api *system) registerEndpoints() {
	api.Group(func(r chi.Router) {
		r.Get("/check", api.systemCheck)
		r.Get("/panic", api.systemPanic)
	})
}

func (api *system) systemCheck(w http.ResponseWriter, r *http.Request) {
	cjson.ServeData(w, cjson.Object{
		"system_status": "ok",
	})
}

func (api *system) systemPanic(w http.ResponseWriter, r *http.Request) {
	cjson.ServeError(w, cerror.NewAPIError(http.StatusInternalServerError, "system panic", nil))
}
