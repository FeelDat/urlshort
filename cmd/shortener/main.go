package main

import (
	"context"
	"database/sql"
	"github.com/FeelDat/urlshort/internal/app/config"
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/FeelDat/urlshort/internal/custommiddleware"
	log "github.com/FeelDat/urlshort/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func main() {

	logger, err := log.InitLogger("Info")
	if err != nil {
		panic(err)
	}

	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	loggerMiddleware := custommiddleware.NewLoggerMiddleware(logger)
	compressMiddleware := custommiddleware.NewCompressMiddleware()

	r := chi.NewRouter()

	var h handlers.HandlerInterface
	var db *sql.DB

	if conf.DatabaseAddress != "" {
		db, err = sql.Open("pgx", conf.DatabaseAddress)
		if err != nil {
			logger.Fatal(err)
		}
		defer db.Close()

		dbRepo, err := storage.NewDBStorage(context.Background(), db)
		if err != nil {
			logger.Fatal(err)
		}

		h = handlers.NewHandler(dbRepo, conf.BaseAddress)

	} else {
		inMemRepo, err := storage.NewInMemStorage(conf.FilePath)
		if err != nil {
			logger.Fatal(err)
		}

		h = handlers.NewHandler(inMemRepo, conf.BaseAddress)
	}

	r.Use(middleware.Compress(5,
		"application/json"+
			"text/html"))
	r.Use(loggerMiddleware.LoggerMiddleware)
	r.Use(compressMiddleware.CompressMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenURL)
		r.Route("/api/shorten", func(r chi.Router) {
			r.Post("/", h.ShortenURLJSON)
			r.Post("/batch", h.ShortenURLBatch)
		})
		r.Get("/{id}", h.GetFullURL)
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			if err = db.Ping(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
	})

	err = http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		logger.Fatal(err)
	}

}
