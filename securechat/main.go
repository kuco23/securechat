package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"
)

func serveWithTLS(certPubPath string, certKeyPath string) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}
	log.Printf("Starting TLS server on %v", server.Addr)
	err := server.ListenAndServeTLS(certPubPath, certKeyPath)
	if err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	xApiKey := os.Getenv("X_API_KEY")
	certPathPub := os.Getenv("CERT_PUB_PATH")
	certPathKey := os.Getenv("CERT_KEY_PATH")

	// define wss handle
	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r, xApiKey)
	})

	// serve http server with TLS
	serveWithTLS(certPathPub, certPathKey)
}
