package main

import (
	"context"
	"crypto/tls"
	"net"
	_ "net/http/pprof"
	"time"

	"github.com/quic-go/quic-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:1234")
	if err != nil {
		panic(err)
	}
	udpConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	conn, err := quic.Dial(ctx, udpConn, udpAddr, &tls.Config{InsecureSkipVerify: true}, &quic.Config{})
	if err != nil {
		panic(err)
	}
	println("Connected to server:", conn.RemoteAddr())
	// Use the connection...
	stream, err := conn.AcceptStream(ctx)
	if err != nil {
		panic(err)
	}
	// Read from the stream
	data := make([]byte, 1024)
	size, err := stream.Read(data)
	if err != nil {
		panic(err)
	}
	println("Received data:", string(data[:size]))
	cancel() // Cancel the context to close the connection gracefully
}
