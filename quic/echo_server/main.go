package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"

	"github.com/quic-go/quic-go"
)

// https://github.com/quic-go/quic-go/blob/master/example/echo/echo.go

func main() {
	echoServer()
}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {
	listener, err := quic.ListenAddr("localhost:4433", generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	for {
		println("start to accept connection")
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}
		println("did accept one connection")
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}
		// Echo through the loggingWriter
		io.Copy(Logger{stream}, stream)
	}
}

type Logger struct {
	io.Writer
}

func (w Logger) Write(b []byte) (int, error) {
	fmt.Printf("Server got: '%s'\n", string(b))
	return w.Writer.Write(b)
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
