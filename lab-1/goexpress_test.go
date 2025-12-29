package goexpress

import (
	"context"
	"net/http"
	"testing"
	"time"
)

// TestNew verifies that New() creates an Engine with default configuration
func TestNew(t *testing.T) {
	engine := New()

	if engine == nil {
		t.Fatal("New() returned nil, expected an Engine")
	}
	if engine.config == nil {
		t.Fatal("Engine config is nil, expected default config")
	}
	if engine.config.Port != ":8080" {
		t.Errorf("Expected default port :8080, got %s", engine.config.Port)
	}

	expectedTimeout := 10 * time.Second
	if engine.config.ReadTimeout != expectedTimeout {
		t.Errorf("Expected ReadTimeout %v, got %v", expectedTimeout, engine.config.ReadTimeout)
	}
	if engine.config.WriteTimeout != expectedTimeout {
		t.Errorf("Expected WriteTimeout %v, got %v", expectedTimeout, engine.config.WriteTimeout)
	}
	if engine.server == nil {
		t.Fatal("Engine server is nil, expected http.Server")
	}
}

// TestWithConfig verifies that NewWithConfig() creates an Engine with custom configuration,
// starts the server, and checks it listens on the given port. It also tests ServeHTTP.
func TestWithConfig(t *testing.T) {
	config := &Config{
		Port:         ":8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 6 * time.Second,
	}
	engine := NewWithConfig(config)

	if engine == nil {
		t.Fatal("NewWithConfig() returned nil, expected an Engine")
	}
	if engine.config == nil {
		t.Fatal("Engine config is nil, expected custom config")
	}
	if engine.config.Port != ":8082" {
		t.Errorf("Expected custom port :8082, got %s", engine.config.Port)
	}
	if engine.config.ReadTimeout != 5*time.Second {
		t.Errorf("Expected custom ReadTimeout 5s, got %v", engine.config.ReadTimeout)
	}
	if engine.config.WriteTimeout != 6*time.Second {
		t.Errorf("Expected custom WriteTimeout 6s, got %v", engine.config.WriteTimeout)
	}
	if engine.server == nil {
		t.Fatal("Engine server is nil, expected http.Server")
	}

	// Start server in a goroutine
	done := make(chan struct{})
	go func() {
		if err := engine.Run(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server failed: %v", err)
		}
		close(done)
	}()

	// Give the server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Make a real HTTP request (to the custom port)
	resp, err := http.Get("http://localhost:8082")
	if err != nil {
		t.Fatalf("Failed to GET from server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := engine.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
	<-done // wait for server goroutine to finish
}

// TestGracefulShutdown simulates a long-running request and verifies
// that the server waits for it to complete during shutdown
func TestGracefulShutdown(t *testing.T) {
	engine := New()

	// Track if the long-running request completed
	requestCompleted := false

	// Override ServeHTTP to simulate a long-running task
	engine.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a long-running task (2 seconds)
		t.Log("long-running task started")
		time.Sleep(2 * time.Second)
		requestCompleted = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Long task completed"))
		t.Log("Long-running task completed")
	})

	// Start server in a goroutine
	go func() {
		if err := engine.Run(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server failed: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a request in the background (it will take 2 seconds)
	requestDone := make(chan struct{})
	go func() {
		resp, err := http.Get("http://localhost:8080/long-task")
		if err != nil {
			t.Errorf("Request failed: %v", err)
		} else {
			resp.Body.Close()
		}
		close(requestDone)
	}()

	// Give the request time to start processing
	time.Sleep(100 * time.Millisecond)

	// Now initiate shutdown with a 5-second timeout
	t.Log("Initiating graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownErr := engine.Shutdown(ctx)
	if shutdownErr != nil {
		t.Errorf("Shutdown failed: %v", shutdownErr)
	}

	// Wait for the request to finish
	<-requestDone

	// Verify the long-running request was completed before shutdown
	if !requestCompleted {
		t.Error("Expected long-running request to complete during graceful shutdown")
	}
	t.Log("Graceful shutdown test passed")
}
