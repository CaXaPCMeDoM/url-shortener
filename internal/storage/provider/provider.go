package provider

import (
	"errors"
	"sync"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/interfaces"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgresql"
)

type BaseProvider struct {
	next interfaces.Provider
}

var (
	instance interfaces.Provider
	once     sync.Once
)

func GetDefaultProvider() interfaces.Provider {
	once.Do(func() {
		postgresqlProvider := &postgresql.Provider{}
		memoryProvider := &memory.Provider{}

		postgresqlProvider.SetNext(memoryProvider)
		instance = postgresqlProvider
	})
	return instance
}

func (p *BaseProvider) Provide(storageType string, config config.Config) (storage.Storage, error) {
	if p.next != nil {
		return p.next.Provide(storageType, config)
	}
	return nil, errors.New("unknown storage type")
}

func (b *BaseProvider) SetNext(provider interfaces.Provider) {
	b.next = provider
}
