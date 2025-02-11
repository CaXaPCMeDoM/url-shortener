package errdef

import "errors"

var (
	ErrURLNotFound    = errors.New("url not found")
	ErrURLExists      = errors.New("url exists")
	ErrUnknownStorage = errors.New("unknown storage")
)
