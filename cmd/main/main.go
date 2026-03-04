package main

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/atlant1da-404/internal/cfg"
	"github.com/atlant1da-404/internal/infra/handler/rest"
	"github.com/atlant1da-404/internal/infra/repo/cache"
	"github.com/atlant1da-404/internal/usecase"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	config, err := cfg.NewConfig()
	if err != nil {
		log.Fatalf("cfg.NewConfig(): %s", err.Error())
		return
	}

	rCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10_000_000,
		MaxCost:     256 << 20,
		BufferItems: 64,
	})

	repo := cache.New(rCache)
	uc := usecase.NewUsecase(repo)
	apiHandler := rest.NewAPIHandler(ctx, uc, 64, 1*time.Millisecond, wg)

	http := &fasthttp.Server{
		Handler: apiHandler.Router,
		// Set to allow more simultaneous connections per IP
		MaxConnsPerIP: 1000, // Adjust as needed
		// Handle high volume of requests per second
		MaxRequestBodySize: 10 * 1024 * 1024, // 10 MB for large payloads
		IdleTimeout:        10 * time.Minute, // 10 minutes idle timeout
		// You can adjust the read/write buffer sizes based on your needs
		ReadBufferSize:  1024 * 1024, // 1MB buffer
		WriteBufferSize: 1024 * 1024, // 1MB buffer
		// Configure TCP keepalive for long-lived connections
		TCPKeepalive: true,
	}

	// Graceful Shutdown Logic
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		if err := http.ListenAndServe(config.Server.Port); err != nil {
			log.Fatalf("HTTP server ListenAndServe failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	sigReceived := <-signalChan
	log.Printf("Received signal: %s. Initiating graceful shutdown...", sigReceived)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := http.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("HTTP server Shutdown failed: %v", err)
	} else {
		log.Println("HTTP server successfully shut down")
	}

	cancel()

	wg.Wait()
}
