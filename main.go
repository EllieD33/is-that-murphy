package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ellied33/is-that-murphy/handlers"
	"github.com/ellied33/is-that-murphy/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.VerifyHandler(w, r)
		case http.MethodPost:
			handlers.AddHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ----- Rate limiter -----
	limiter := middleware.NewIPRateLimiter(5, 10) 
	stopCleanup := make(chan struct{})
	go limiter.CleanupExpired(5*time.Minute, stopCleanup)

	handler := limiter.Middleware(mux)

	// ----- HTTP server with timeouts -----
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// ----- Graceful shutdown -----
	go func() {
		// Capture OS signals
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c

		log.Printf("Shutdown signal received: %v", sig)

		// Stop cleanup goroutine
		close(stopCleanup)

		// Shutdown HTTP server gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		} else {
			log.Println("Server stopped gracefully")
		}
	}()

	// ----- Start server -----
	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
