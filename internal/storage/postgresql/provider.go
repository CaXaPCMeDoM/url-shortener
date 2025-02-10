package postgresql

import (
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
)

type Provider struct {
	storage.BaseProvider
}

func (p Provider) Provide(storageType string, config config.Config) (storage.Storage, error) {
	if storageType == storage.POSTGRES {
		return New(config.Database.DSN())
	}
	return p.BaseProvider.Provide(storageType, config)
}
