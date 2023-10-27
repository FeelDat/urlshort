// Package main provides the entry point for the urlshort application.
// This application provides utilities for URL shortening.
package main

import (
	"context"
	"database/sql"
	_ "github.com/FeelDat/urlshort/docs"
	"github.com/FeelDat/urlshort/internal/app/config"
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/FeelDat/urlshort/internal/custommiddleware"
	log "github.com/FeelDat/urlshort/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/swaggo/http-swagger"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// main is the entry point of the urlshort application. It initializes and configures
// logger, database, middlewares, routes, and starts the HTTP server.
func main() {

	// Initialize logger
	logger, err := log.InitLogger("Info")
	if err != nil {
		panic(err)
	}

	// Load configuration
	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// Deprecated: Using rand.Seed with time.Now().UnixNano() is deprecated.
	// Consider using a more robust seed mechanism or another source of randomness if needed.
	rand.Seed(time.Now().UnixNano())

	// Middleware initialization
	loggerMiddleware := custommiddleware.NewLoggerMiddleware(logger)
	authMiddleware := custommiddleware.NewAuthMiddleware()

	// Router initialization
	r := chi.NewRouter()

	var h handlers.Handler
	var db *sql.DB

	// Database initialization
	if conf.DatabaseAddress != "" {
		db, err = sql.Open("pgx", conf.DatabaseAddress)
		if err != nil {
			logger.Fatal(err)
		}
		defer db.Close()

		err = storage.InitDB(context.Background(), db)
		if err != nil {
			logger.Fatal(err)
		}

		dbRepo := storage.NewDBStorage(db)
		h = handlers.NewHandler(dbRepo, conf.BaseAddress, logger)

	} else {
		// In-memory storage initialization
		inMemRepo, err := storage.NewInMemStorage(conf.FilePath)
		if err != nil {
			logger.Fatal(err)
		}

		h = handlers.NewHandler(inMemRepo, conf.BaseAddress, logger)
	}

	// Middleware registration
	r.Use(middleware.Compress(5,
		"application/json"+
			"text/html"))
	r.Use(loggerMiddleware.LoggerMiddleware)
	r.Use(authMiddleware.AuthMiddleware)

	// Debug profiler mount
	r.Mount("/debug", middleware.Profiler())

	// Routing configuration

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
	// Ping route for health check
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if conf.DatabaseAddress == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = db.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	// Swagger documentation endpoint
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Starting HTTP server
	err = http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		logger.Fatal(err)
	}
}
