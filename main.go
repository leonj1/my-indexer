package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"my-indexer/router"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.NewRouter()
	
	// Configure server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Server run context
	srvCtx, srvCancel := context.WithCancel(context.Background())
	defer srvCancel()

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signals
	select {
	case err := <-serverErrors:
		log.Printf("Server error: %v", err)
	case s := <-sig:
		log.Printf("Server shutdown initiated by %v signal", s)
	}

	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, shutdownCancel := context.WithTimeout(srvCtx, 30*time.Second)
	defer shutdownCancel()

	// Trigger graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}
}
