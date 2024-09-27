// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"
)

func serveTLS(certPubPath string, certKeyPath string) {
	// Create TLS configuration (for development or production)
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	// Create HTTPS server using TLS
	server := &http.Server{
		Addr:      ":8443", // TLS typically uses port 443 (or 8443 for development)
		TLSConfig: tlsConfig,
	}
	// Start the server with TLS
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

	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, xApiKey)
	})

	serveTLS(certPathPub, certPathKey)
}
