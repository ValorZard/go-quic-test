package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
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
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	os.Setenv("QLOGDIR", path+"/qlog")
	conn, err := quic.Dial(ctx, udpConn, udpAddr, &tls.Config{InsecureSkipVerify: true}, &quic.Config{Tracer: qlog.DefaultConnectionTracer})
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
	stream.Write([]byte("exit"))
	cancel() // Cancel the context to close the connection gracefully
}
