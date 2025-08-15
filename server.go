package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"net"
	_ "net/http/pprof"
	"time"

	"github.com/quic-go/quic-go"
)

func main() {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 1234})
	if err != nil {
		panic(err)
	}
	tr := quic.Transport{
		Conn: udpConn,
	}
	tlsConf, err := generateTLSConfig()
	if err != nil {
		panic(err)
	}
	quicConf := quic.Config{}
	ln, err := tr.Listen(tlsConf, &quicConf)
	if err != nil {
		panic(err)
	}
	println("Listening on:", ln.Addr())
	for {
		conn, err := ln.Accept(context.Background())
		if err != nil {
			println("Accept error:", err)
			continue
		}
		go func(c *quic.Conn) {
			println("New connection accepted:", c.RemoteAddr())
			stream, err := c.OpenStream()
			if err != nil {
				println("Error opening stream:", err)
				return
			}
			stream.Write([]byte("Hello from server!"))
			// Uncomment the following line to make the code work
			time.Sleep(5 * time.Second)
			stream.Close()
			c.CloseWithError(0, "Connection closed by server")
			println("Connection closed:", c.RemoteAddr())
		}(conn)
	}
}

// generateTLSConfig creates a self-signed TLS config for testing
func generateTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  key,
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}, nil
}
