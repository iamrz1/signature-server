package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"signature-server/data"
)

// NewAPIRouter ...
func NewAPIRouter(sStore data.SignatureStore, tStore data.TransactionStore) http.Handler {
	h := chi.NewRouter()
	h.Use(cors.AllowAll().Handler)

	h.Route("/", func(r chi.Router) {
		r.Get("/doc/*", httpSwagger.Handler())
		r.Mount("/", NewSignatureHandler(sStore, tStore))
	})
	return h
}

// NewSystemRouter ...
func NewSystemRouter() http.Handler {
	h := chi.NewRouter()
	h.Mount("/system", NewSystemHandler())
	h.Mount("/debug", middleware.Profiler())
	return h
}
