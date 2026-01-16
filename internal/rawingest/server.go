// internal/rawingest/server.go
// PURPOSE: Raw Ingest TCP listener and accept loop.
// ALLOWED: net listen/accept, goroutine spawning, context cancellation
// FORBIDDEN: protocol parsing, memory access, business logic

package rawingest

import (
	"context"
	"net"
	"time"
)

// Server owns the Raw Ingest TCP socket.
// It does NOT understand protocol or memory semantics.
type Server struct {
	listenAddr     string
	maxPacketBytes int
	readTimeout    time.Duration
	resolver       MemoryResolver
}

// NewServer constructs a Raw Ingest server.
// All validation here is socket-safety only.
func NewServer(
	listenAddr string,
	resolver MemoryResolver,
	maxPacketBytes int,
	readTimeout time.Duration,
) (*Server, error) {

	if listenAddr == "" {
		return nil, ErrRejected
	}
	if resolver == nil {
		return nil, ErrRejected
	}
	if maxPacketBytes <= v1HeaderSize {
		return nil, ErrRejected
	}
	if readTimeout <= 0 {
		return nil, ErrRejected
	}

	return &Server{
		listenAddr:     listenAddr,
		maxPacketBytes: maxPacketBytes,
		readTimeout:    readTimeout,
		resolver:       resolver,
	}, nil
}

// Start begins listening and serving until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	// Unblock Accept on context cancellation
	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			continue
		}

		go handleConn(
			ctx,
			conn,
			s.resolver,
			s.maxPacketBytes,
			s.readTimeout,
		)
	}
}
