package memory

import (
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
	constants "url-shortener/internal/storage/constans"
	"url-shortener/internal/storage/errdef"
	"url-shortener/internal/storage/interfaces"
)

type Provider struct {
	next interfaces.Provider
}

func (p *Provider) SetNext(provider interfaces.Provider) {
	p.next = provider
}

func (p *Provider) Provide(storageType string, config config.Config) (storage.Storage, error) {
	if storageType == constants.MEMORY {
		return New(), nil
	}
	if p.next != nil {
		return p.next.Provide(storageType, config)
	}
	return nil, errdef.ErrUnknownStorage
}
