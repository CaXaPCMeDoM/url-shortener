package main

import (
	"flag"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/gateway"
	"url-shortener/internal/service/shortener"
	"url-shortener/internal/storage/provider"
	"url-shortener/internal/transport"
)

func main() {
	storageType := *flag.String("storage", "memory", "Тип хранилища (memory или postgres)")
	flag.Parse()

	cfg := *config.MustLoad()

	store, err := provider.GetDefaultProvider().Provide(storageType, cfg)

	if err != nil {
		slog.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	shortenerService := shortener.New(store, cfg)

	go func() {
		if err := transport.RunGRPCServer(cfg.Grpc.Address, shortenerService, cfg); err != nil {
			slog.Error("gRPC server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	go func() {
		if err := gateway.RunHTTPServer(cfg.Http.Address, cfg.Grpc.Address); err != nil {
			slog.Error("HTTP server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	select {}
}
