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

// This file demonstrates how to use the GoExpress framework
// Run with: go run goexpress.go config.go example_usage.go

func exampleBasicUsage() {
	// Create a new engine with default settings
	app := goexpress.New()

	// Run the server (this blocks until server stops)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func exampleCustomConfig() {
	// Create custom configuration
	config := &goexpress.Config{
		Port:         ":3000",          // Use port 3000 instead of 8080
		ReadTimeout:  15 * time.Second, // Wait up to 15 seconds for requests
		WriteTimeout: 15 * time.Second, // Wait up to 15 seconds to send response
	}

	// Create engine with custom config
	app := goexpress.NewWithConfig(config)

	// Run the server
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

// Uncomment the function you want to run and execute:
// go run main.go

func main() {
	// exampleBasicUsage()
	// exampleCustomConfig()
	exampleGracefulShutdown()
}
