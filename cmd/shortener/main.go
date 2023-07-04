package main

import (
	"github.com/FeelDat/urlshort/internal/app/config"
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/FeelDat/urlshort/internal/app/storage"
	logger2 "github.com/FeelDat/urlshort/internal/logger"
	"github.com/FeelDat/urlshort/internal/middleware"
	"github.com/go-chi/chi/v5"
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

	loggerMiddleware := middleware.NewLoggerMiddleware(logger)
	compressMiddleware := middleware.NewCompressMiddleware()

	r := chi.NewRouter()

	mapStorage := storage.NewInMemoryStorage()

	h := handlers.NewHandler(mapStorage, conf.BaseAddress)

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
