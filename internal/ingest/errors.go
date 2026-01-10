// internal/ingest/errors.go
package ingest

import "errors"

var (
	ErrUnknownMemory     = errors.New("unknown memory")
	ErrInvalidArea       = errors.New("invalid area")
	ErrInvalidPayload    = errors.New("invalid payload")
	ErrInvalidBoolean    = errors.New("invalid numeric boolean")
	ErrPayloadMismatch   = errors.New("payload does not match area")
)
