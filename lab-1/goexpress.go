package goexpress

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

/*
Engine is the core of our web framework
*/
type Engine struct {
	config *Config
	server *http.Server
}

/*
New creates and returns a new Engine with default configuration
*/
func New() *Engine {
	return NewWithConfig(DefaultConfig())
}

/*
NewWithConfig creates a new Engine with custom configuration.
The Engine itself handles requests! because it implements the http.Handler interface
It will call our ServeHTTP method for each request
*/
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

/*
 ServeHTTP is THE most important method! By implementing this, our Engine becomes an http.Handler.
 Go's http package will call this method for every incoming request
*/

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Hello from GoExpress!\n")
	fmt.Fprintf(w, "You requested: %s %s\n", r.Method, r.URL.Path)
}

/*
Run starts the HTTP server and begins listening for requests. This is a blocking call - it won't return until the server stops
It will call our ServeHTTP method for each request
*/
func (e *Engine) Run() error {
	log.Printf("GoExpress server starting on http://localhost%s\n", e.config.Port)

	err := e.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

/*
Shutdown gracefully shuts down the server. It waits for active requests to finish before stopping.
The context allows you to set a timeout for shutdown
It first closes all open listeners, then waits for connections to finish
*/
func (e *Engine) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server gracefully...")

	err := e.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Println("Server stopped successfully")
	return nil
}

/*
main is the entry point when running this file directly. This demonstrates how to use the Engine
*/
// func main() {
// 	app := New()

// 	if err := app.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }
