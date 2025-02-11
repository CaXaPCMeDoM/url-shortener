package interfaces

import (
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
)

type Provider interface {
	SetNext(provider Provider)
	Provide(storageType string, config config.Config) (storage.Storage, error)
}
