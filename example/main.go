package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ridwan414/goexpress"
)

func exampleDefaultEngine() {
	app := goexpress.New()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func exampleCustomConfig() {
	config := &goexpress.Config{
		Port:         ":3001",
		ReadTimeout:  12 * time.Second,
		WriteTimeout: 12 * time.Second,
	}
	app := goexpress.NewWithConfig(config)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func exampleGracefulShutdown() {
	// Create new engine
	app := goexpress.New()

	// Start server in background goroutine
	go func() {
		log.Println("Starting server...")
		if err := app.Run(); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// Create channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)

	// Register to receive SIGINT (Ctrl+C) and SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	<-quit

	log.Println("Shutdown signal received...")

	// Create context with 5 second timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func main() {
	// Uncomment the desired example to run:
	// exampleDefaultEngine()
	// exampleCustomConfig()
	exampleGracefulShutdown()
}
