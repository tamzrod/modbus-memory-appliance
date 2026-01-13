package modbus

import (
	"log"
	"net"

	"modbus-memory-appliance/internal/modbus/ipfilter"
)

const defaultMaxConnections = 32

func Start(
	addr string,
	resolve MemoryResolver,
	allowIPs []string,
	denyIPs []string,
	maxConns int,
) error {

	filter, err := ipfilter.Compile(allowIPs, denyIPs)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if maxConns <= 0 {
		maxConns = defaultMaxConnections
	}

	log.Println("Modbus TCP listening on", addr, "max_connections =", maxConns)

	sem := make(chan struct{}, maxConns)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		if filter.Enabled() {
			host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
			if err != nil {
				conn.Close()
				continue
			}

			ip := net.ParseIP(host)
			if !filter.Allowed(ip) {
				conn.Close()
				continue
			}
		}

		select {
		case sem <- struct{}{}:
		default:
			conn.Close()
			continue
		}

		go func() {
			defer func() { <-sem }()
			handleConn(conn, resolve)
		}()
	}
}
