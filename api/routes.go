package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"signature-server/config"
	memDB "signature-server/data/memory"
	"signature-server/logger"
)

// NewAPIRouter ...
func NewAPIRouter(appConfig *config.App) http.Handler {
	h := chi.NewRouter()
	h.Use(cors.AllowAll().Handler)
	sStore, err := memDB.NewSignatureStore(appConfig.DaemonKey)
	if err != nil {
		logger.Fatal(err.Error())
	}
	tStore := memDB.NewTransactionStore()

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
