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
	"math/rand"
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

	rand.Seed(time.Now().UnixNano())

	loggerMiddleware := custommiddleware.NewLoggerMiddleware(logger)
	authMiddleware := custommiddleware.NewAuthMiddleware()
	compressMIddleware := custommiddleware.NewCompressMiddleware()

	r := chi.NewRouter()

	var h handlers.HandlerInterface
	var db *sql.DB

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
		inMemRepo, err := storage.NewInMemStorage(conf.FilePath)
		if err != nil {
			logger.Fatal(err)
		}

		h = handlers.NewHandler(inMemRepo, conf.BaseAddress, logger)
	}

	r.Use(middleware.Compress(5,
		"application/json"+
			"text/html"))
	r.Use(loggerMiddleware.LoggerMiddleware)
	r.Use(authMiddleware.AuthMiddleware)
	r.Use(compressMIddleware.CompressMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenURL)
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", h.ShortenURLJSON)
				r.Post("/batch", h.ShortenURLBatch)
			})
			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", h.GetUsersURLS)
			})
		})
		r.Get("/{id}", h.GetFullURL)
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
	})

	err = http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		logger.Fatal(err)
	}

}
