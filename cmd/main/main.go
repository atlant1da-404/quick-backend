package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/atlant1da-404/internal/cfg"
	"github.com/atlant1da-404/internal/infra/client/dragonfly"
	"github.com/atlant1da-404/internal/infra/handler/rest"
	"github.com/atlant1da-404/internal/infra/repo/cache"
	"github.com/atlant1da-404/internal/usecase"
)

func main() {
	config, err := cfg.NewConfig()
	if err != nil {
		log.Fatalf("cfg.NewConfig(): %s", err.Error())
		return
	}

	dfClient, err := dragonfly.New(
		dragonfly.WithAddress(config.DragonFly.Address),
		dragonfly.WithPassword(config.DragonFly.Password),
		dragonfly.WithDB(config.DragonFly.DB),
		dragonfly.WithReadTimeout(config.DragonFly.ReadTimeout),
		dragonfly.WithWriteTimeout(config.DragonFly.WriteTimeout),
		dragonfly.WithMaxRetries(config.DragonFly.MaxRetries),
		dragonfly.WithMinRetryBackoff(config.DragonFly.MinRetryBackoff),
		dragonfly.WithMaxRetryBackoff(config.DragonFly.MaxRetryBackoff),
	)
	if err != nil {
		log.Fatalf("dragonfly.New: %s", err.Error())
		return
	}
	defer func() {
		if err := dfClient.Close(); err != nil {
			log.Fatalf("app.Run - dfClient.Close: %s", err.Error())
			return
		}
	}()

	repo := cache.New(dfClient)
	uc := usecase.NewUsecase(repo)
	apiHandler := rest.NewAPIHandler(uc)

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

	// Attempt graceful shutdown
	if err := http.Shutdown(); err != nil {
		log.Fatalf("HTTP server Shutdown failed: %v", err)
	} else {
		log.Println("HTTP server successfully shut down")
	}
}
