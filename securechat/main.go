package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

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

func ServeHome(w http.ResponseWriter, r *http.Request, xak string) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		log.Printf("Wrong url path from %v", r.RemoteAddr)
		return
	}
	auth := chatAuth(r)
	if auth != xak {
		log.Printf("Failed auth for home request from %v", r.RemoteAddr)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, xak string) {
	auth := chatAuth(r)
	if auth != xak {
		log.Printf("Forbidden auth for wss request from %v", r.RemoteAddr)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		id:     clientId(r),
		roomID: auth,
		hub:    hub,
		conn:   conn,
		send:   make(chan Message, 256),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func chatAuth(r *http.Request) string {
	queryParams := r.URL.Query()
	return queryParams.Get("x-api-key")
}

func clientId(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ServeHome(w, r, xApiKey)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r, xApiKey)
	})

	// serve http server with TLS
	serveWithTLS(certPathPub, certPathKey)
}
