package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
)

func main() {
	storageType := *flag.String("storage", "memory", "Тип хранилища (memory или postgres)")
	flag.Parse()

	cfg := config.MustLoad()

	store, err := storage.GetDefaultProvider().Provide(storageType, *cfg)

	if err != nil {
		slog.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}
