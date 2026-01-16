// internal/rawingest/errors.go
// PURPOSE: Sentinel errors for Raw Ingest.
// ALLOWED: error definitions only
// FORBIDDEN: branching, logging, wrapping, semantics

package rawingest

import "errors"

// All failures collapse to rejection.
// No descriptive errors cross the transport boundary.
var (
	ErrRejected = errors.New("raw ingest rejected")
)
