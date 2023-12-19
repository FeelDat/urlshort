package routers

import (
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func ShortenerRouter(h handlers.Handler) http.Handler {

	r := chi.NewRouter()

	// Routing configuration
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenURL)
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", h.ShortenURLJSON)
				r.Post("/batch", h.ShortenURLBatch)
			})
			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", h.GetUsersURLS)
				r.Delete("/urls", h.DeleteURLS)
			})
		})
		r.Get("/{id}", h.GetFullURL)
		//r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		//	if conf.DatabaseAddress == "" {
		//		w.WriteHeader(http.StatusBadRequest)
		//		return
		//	}
		//	if err := db.Ping(); err != nil {
		//		w.WriteHeader(http.StatusInternalServerError)
		//		return
		//	}
		//	w.WriteHeader(http.StatusOK)
		//})
		// Swagger documentation endpoint
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	return r
}
