package storage

import (
	"errors"
	"sync"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgresql"
)

const (
	POSTGRES = "postgres"
	MEMORY   = "memory"
)

type Provider interface {
	SetNext(provider Provider)
	Provide(storageType string, config config.Config) (Storage, error)
}

type BaseProvider struct {
	next Provider
}

var (
	instance Provider
	once     sync.Once
)

func GetDefaultProvider() Provider {
	once.Do(func() {
		postgresqlProvider := &postgresql.Provider{}
		memoryProvider := &memory.Provider{}

		postgresqlProvider.SetNext(memoryProvider)
		instance = postgresqlProvider
	})
	return instance
}

func (p *BaseProvider) Provide(storageType string, config config.Config) (Storage, error) {
	if p.next != nil {
		return p.next.Provide(storageType, config)
	}
	return nil, errors.New("unknown storage type")
}

func (b *BaseProvider) SetNext(provider Provider) {
	b.next = provider
}
