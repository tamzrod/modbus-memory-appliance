package modbus

import (
	"log"
	"net"
)

func Start(addr string, resolve MemoryResolver) error {
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
		go handleConn(conn, resolve)
	}
}
