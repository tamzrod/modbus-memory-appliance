package modbus

import (
	"log"
	"net"

	"modbus-memory-appliance/internal/modbus/ipfilter"
)

// Start starts a Modbus TCP listener with optional IP filtering
func Start(
	addr string,
	resolve MemoryResolver,
	allowIPs []string,
	denyIPs []string,
) error {

	// 1️⃣ Compile IP filter ONCE at startup
	filter, err := ipfilter.Compile(allowIPs, denyIPs)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("Modbus TCP listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		// 2️⃣ Enforce IP filter IMMEDIATELY after Accept
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

		go handleConn(conn, resolve)
	}
}
