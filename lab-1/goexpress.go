package goexpress

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Engine is the core type of the web framework,
// holding configuration and the underlying HTTP server.
type Engine struct {
	config *Config
	server *http.Server
}

// New returns a new Engine instance using the default configuration.
func New() *Engine {
	return NewWithConfig(DefaultConfig())
}

// NewWithConfig creates a new Engine with the provided custom configuration.
// The Engine implements http.Handler: the ServeHTTP method is invoked for each request.
func NewWithConfig(config *Config) *Engine {
	engine := &Engine{
		config: config,
	}

	engine.server = &http.Server{
		Addr:         config.Port,
		Handler:      engine,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return engine
}

// ServeHTTP implements the http.Handler interface for Engine.
// It is invoked by the net/http package for every HTTP request.
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from GoExpress!\n")
	fmt.Fprintf(w, "You requested: %s %s\n", r.Method, r.URL.Path)
}

// Run starts the HTTP server and begins serving requests.
// This is a blocking call; it only returns when the server shuts down
// or encounters an error.
func (e *Engine) Run() error {
	log.Printf("GoExpress server starting on http://localhost%s\n", e.config.Port)
	err := e.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the HTTP server with the given context.
// It waits for active requests to finish before shutting down.
func (e *Engine) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server gracefully...")
	err := e.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	log.Println("Server stopped successfully")
	return nil
}
