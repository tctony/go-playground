package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/quic-go/quic-go"
)

func main() {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 4433})
	if err != nil {
		panic(fmt.Sprintf("create udp listener failed: %s", err))
	}

	// WIP config tls certs
	ln, err := quic.Listen(udpConn, &tls.Config{}, nil)
	if err != nil {
		panic(fmt.Sprintf("create quic listener failed: %s", err))
	}

	for {
		conn, err := ln.Accept(context.Background())
		if err != nil {
			panic(fmt.Sprintf("accept quic connection failed: %s", err))
		}

		// WIP impl echo

		conn.CloseWithError(1, "WIP")
	}
}
