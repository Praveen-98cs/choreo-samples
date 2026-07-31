package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type EchoResponse struct {
	Headers map[string][]string `json:"headers"`
	Payload string              `json:"payload"`
}

func main() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", echoHandler)

	serverPort := 9090
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: serverMux,
	}
	go func() {
		log.Printf("Starting HTTP Echo Service on port %d\n", serverPort)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP ListenAndServe error: %v", err)
		}
		log.Println("HTTP server stopped serving new requests.")
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh // Wait for shutdown signal

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Shutdown complete.")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string][]string)
	for name, values := range r.Header {
		headers[name] = values
		for _, value := range values {
			log.Printf("%s: %s", name, value)
		}
	}

	var payload string
	if r.Method == http.MethodGet {
		payload = "echo-get"
	} else if r.Method == http.MethodPost {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		payload = string(bodyBytes)
	} else {
		payload = "echo-other"
	}

	response := EchoResponse{
		Headers: headers,
		Payload: payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
	}
}
