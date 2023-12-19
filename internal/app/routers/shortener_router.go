package routers

import (
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// ShortenerRouter creates and returns a new HTTP handler that implements
// the routing logic for the URL shortener application. It defines routes
// for URL shortening operations and associated API endpoints.
// The router includes the following routes:
// - POST /: Shortens a given URL.
// - POST /api/shorten: Shortens a URL via JSON request.
// - POST /api/shorten/batch: Shortens multiple URLs in a batch.
// - GET /user/urls: Retrieves URLs associated with a user.
// - DELETE /user/urls: Deletes URLs associated with a user.
// - GET /{id}: Redirects to the full URL corresponding to the given ID.
// - GET /swagger/*: Serves Swagger documentation.
// The function requires a handlers.Handler struct that contains the necessary
// handlers for the above routes.
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

		// Swagger documentation endpoint
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	return r
}
