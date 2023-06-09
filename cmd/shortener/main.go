package main

import (
	"github.com/FeelDat/urlshort/internal/app/config"
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/FeelDat/urlshort/internal/custommiddleware"
	logger2 "github.com/FeelDat/urlshort/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {

	logger, err := logger2.InitLogger("Info")
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
	defer func() {
		if err := mapStorage.Close(); err != nil {
			logger.Error("Failed to close the file", err)
		}
	}()

	if err != nil {
		logger.Fatal(err)
	}

	h := handlers.NewHandler(mapStorage, conf)

	r.Use(middleware.Compress(5,
		"application/json"+
			"text/html"))
	r.Use(loggerMiddleware.LoggerMiddleware)
	r.Use(compressMiddleware.CompressMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenURL)
		r.Post("/api/shorten", h.ShortenURLJSON)
		r.Get("/{id}", h.GetFullURL)
	})

	err = http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		logger.Fatal(err)
	}

}
