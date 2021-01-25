package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// define http server and server handler
	vs := ValidatingServerHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", vs.serve)

	// TODO: improve and start as Go routine in background
	log.Printf("Server running listening in port: %s", port)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%s", port), "/tls/tls.crt", "/tls/tls.key", mux)
	if err != nil {
		log.Printf("Failed to listen and serve admission webhook server: %v", err)
	}
}
