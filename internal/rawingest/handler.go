// internal/rawingest/handler.go
// PURPOSE: One connection loop: read frame -> resolve memory -> apply -> respond.
// ALLOWED: framing + resolver + apply + response
// FORBIDDEN: logging outside debug markers, retries, semantics

package rawingest

import (
	"bufio"
	"context"
	"net"
	"time"
)

// handleConn serves one TCP connection as a stream of frames.
// Any rejection closes the connection immediately after sending ResponseRejected.
func handleConn(
	ctx context.Context,
	conn net.Conn,
	resolver MemoryResolver,
	maxPacketBytes int,
	readTimeout time.Duration,
) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	for {
		if ctx.Err() != nil {
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(readTimeout))

		f, err := readFrame(r, maxPacketBytes)
		if err != nil {
			_ = writeRejected(conn)
			return
		}

		mem, ok := resolver.ResolveMemoryByID(f.Header.MemID)
		if !ok {
			_ = writeRejected(conn)
			return
		}

		if err := applyPayload(mem, f.Header, f.Payload); err != nil {
			_ = writeRejected(conn)
			return
		}

		if err := writeOK(conn); err != nil {
			return
		}

		// debug
		// optional low-level instrumentation here
		// debug ends
	}
}
