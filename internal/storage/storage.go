package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Storage interface {
	SaveUrl(urlAbsolute string, alias string) error
	GetUrl(alias string) (string, error)
}
