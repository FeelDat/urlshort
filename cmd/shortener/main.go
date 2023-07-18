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
	"time"
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

	mapStorage, err := storage.NewInMemoryStorage(conf.FilePath)
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		if err := mapStorage.Close(); err != nil {
			logger.Error("Failed to close the file", err)
		}
	}()

	db, err := sql.Open("pgx", conf.DatabaseAddress)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

	dbStorage := storage.NewDatabaseStorage(db)

	h := handlers.NewHandler(mapStorage, conf.BaseAddress, dbStorage)

	r.Use(middleware.Compress(5,
		"application/json"+
			"text/html"))
	r.Use(loggerMiddleware.LoggerMiddleware)
	r.Use(compressMiddleware.CompressMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenURL)
		r.Post("/api/shorten", h.ShortenURLJSON)
		r.Get("/{id}", h.GetFullURL)
		r.Get("/ping", h.Ping)
	})

	err = http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		logger.Fatal(err)
	}

}
