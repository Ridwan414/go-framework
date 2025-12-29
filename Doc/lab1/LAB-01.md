# Lab 01: Framework Foundation & Project Setup

## Learning Objectives
- Set up Go project structure for a reusable framework
- Understand how `net/http` works under the hood
- Create the base `Engine` struct that implements `http.Handler`
- Implement basic server lifecycle (start, shutdown)

---

## Part 1: Understanding Go's net/http Package

The `net/http` package is Go's standard library for building HTTP clients and servers. It provides everything you need to create web applications without external dependencies.

### Key Components

| Component | Description |
|-----------|-------------|
| `http.Server` | The HTTP server struct that manages connections, timeouts, and request handling |
| `http.Handler` | An interface with one method: `ServeHTTP(ResponseWriter, *Request)` |
| `http.ResponseWriter` | Interface for writing HTTP response (headers, status code, body) |
| `http.Request` | Struct containing all request data (method, URL, headers, body) |
| `ListenAndServe()` | Starts the server and begins accepting connections |

### The http.Handler Interface

This is the most important interface in Go's HTTP ecosystem:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**Any type that implements this single method can handle HTTP requests.** This is the foundation of our framework.

---

## Part 2: How the Server Listens and Handles Requests

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        HTTP REQUEST LIFECYCLE                                │
└─────────────────────────────────────────────────────────────────────────────┘

    ┌──────────┐          ┌──────────────┐          ┌──────────────┐
    │  Client  │          │   Network    │          │  Go Runtime  │
    │ (Browser)│          │   (TCP/IP)   │          │              │
    └────┬─────┘          └──────┬───────┘          └──────┬───────┘
         │                       │                         │
         │  1. HTTP Request      │                         │
         │ ─────────────────────>│                         │
         │                       │  2. Accept Connection   │
         │                       │ ───────────────────────>│
         │                       │                         │
         │                       │                         │
    ┌────┴─────────────────────────────────────────────────┴───────────────────┐
    │                                                                           │
    │   ┌─────────────────────────────────────────────────────────────────┐    │
    │   │                        http.Server                               │    │
    │   │                                                                  │    │
    │   │   Addr: ":8080"                                                  │    │
    │   │   Handler: engine (implements http.Handler)                      │    │
    │   │   ReadTimeout: 10s                                               │    │
    │   │   WriteTimeout: 10s                                              │    │
    │   │                                                                  │    │
    │   │   3. For each request, calls:                                    │    │
    │   │      ┌─────────────────────────────────────────────────────┐     │    │
    │   │      │  Handler.ServeHTTP(ResponseWriter, *Request)        │     │    │
    │   │      └─────────────────────────────────────────────────────┘     │    │
    │   │                              │                                   │    │
    │   └──────────────────────────────┼───────────────────────────────────┘    │
    │                                  │                                        │
    │                                  ▼                                        │
    │   ┌─────────────────────────────────────────────────────────────────┐    │
    │   │                         Engine                                   │    │
    │   │                                                                  │    │
    │   │   func (e *Engine) ServeHTTP(w ResponseWriter, r *Request) {     │    │
    │   │       // 4. Process the request                                  │    │
    │   │       // 5. Write response headers                               │    │
    │   │       // 6. Write response body                                  │    │
    │   │   }                                                              │    │
    │   │                                                                  │    │
    │   └─────────────────────────────────────────────────────────────────┘    │
    │                                                                           │
    └───────────────────────────────────────────────────────────────────────────┘
         │                                                          │
         │  7. HTTP Response                                        │
         │ <────────────────────────────────────────────────────────│
         │                                                          │
    ┌────┴─────┐
    │  Client  │
    └──────────┘
```

### Request Flow Summary

1. **Client sends request** - Browser/client initiates HTTP request to server address
2. **Server accepts connection** - `http.Server` accepts the TCP connection on the configured port
3. **Server creates goroutine** - Each request is handled in its own goroutine (concurrent by default!)
4. **Handler invoked** - Server calls `ServeHTTP()` on the registered Handler (our Engine)
5. **Process request** - Engine reads request data (`*http.Request`)
6. **Write response** - Engine writes to `http.ResponseWriter`
7. **Response sent** - Response flows back to client

---

## Part 3: Code Walkthrough

### File Structure

```
lab-1/
├── goexpress.go      # Engine struct, New(), Run(), Shutdown()
├── config.go         # Configuration options
├── goexpress_test.go # Unit tests
└── go.mod
```

---

### Step 1: Configuration (config.go)

First, we define configuration options for our server:

```go
package goexpress

import "time"

// Config holds all configuration for the HTTP server
type Config struct {
    // Port is the address and port to listen on
    Port string

    // ReadTimeout is the maximum duration for reading the entire request
    ReadTimeout time.Duration

    // WriteTimeout is the maximum duration before timing out writes of the response
    WriteTimeout time.Duration
}

// DefaultConfig returns a Config with sensible default values
func DefaultConfig() *Config {
    return &Config{
        Port:         ":8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
}
```

**Key Points:**
- `Port` - The address format is `:PORT` (e.g., `:8080` means listen on all interfaces on port 8080)
- `ReadTimeout` - Prevents slow clients from holding connections open indefinitely
- `WriteTimeout` - Prevents slow responses from blocking the server
- `DefaultConfig()` - Factory function providing sensible defaults

---

### Step 2: Engine Struct (goexpress.go)

The Engine is the core of our framework:

```go
package goexpress

import (
    "context"
    "fmt"
    "log"
    "net/http"
)

// Engine is the core of our web framework
type Engine struct {
    config *Config
    server *http.Server
}
```

**Key Points:**
- `config` - Holds our configuration settings
- `server` - The underlying `http.Server` that does the actual network work

---

### Step 3: Constructor Functions

```go
// New creates and returns a new Engine with default configuration
func New() *Engine {
    return NewWithConfig(DefaultConfig())
}

// NewWithConfig creates a new Engine with custom configuration.
func NewWithConfig(config *Config) *Engine {
    engine := &Engine{
        config: config,
    }

    engine.server = &http.Server{
        Addr:         config.Port,
        Handler:      engine,          // Engine itself is the handler!
        ReadTimeout:  config.ReadTimeout,
        WriteTimeout: config.WriteTimeout,
    }

    return engine
}
```

**Key Points:**
- `New()` - Simple constructor using defaults (most common usage)
- `NewWithConfig()` - Flexible constructor for custom settings
- `Handler: engine` - **This is critical!** The Engine itself handles all requests because it implements `http.Handler`

---

### Step 4: Implementing http.Handler Interface

This is the most important method - it makes our Engine an `http.Handler`:

```go
// ServeHTTP is THE most important method!
// By implementing this, our Engine becomes an http.Handler.
// Go's http package will call this method for every incoming request.
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)

    fmt.Fprintf(w, "Hello from GoExpress!\n")
    fmt.Fprintf(w, "You requested: %s %s\n", r.Method, r.URL.Path)
}
```

**Key Points:**
- `w http.ResponseWriter` - Write your response here (headers, status, body)
- `r *http.Request` - Contains all request information
- `w.Header().Set()` - Set response headers BEFORE writing body
- `w.WriteHeader()` - Set HTTP status code (200, 404, 500, etc.)
- `fmt.Fprintf(w, ...)` - Write response body

**Order Matters!**
```
1. Set Headers     →  w.Header().Set(...)
2. Set Status Code →  w.WriteHeader(...)
3. Write Body      →  fmt.Fprintf(w, ...) or w.Write(...)
```

---

### Step 5: Starting the Server

```go
// Run starts the HTTP server and begins listening for requests.
// This is a blocking call - it won't return until the server stops.
func (e *Engine) Run() error {
    log.Printf("GoExpress server starting on http://localhost%s\n", e.config.Port)

    err := e.server.ListenAndServe()

    if err != nil && err != http.ErrServerClosed {
        return fmt.Errorf("server error: %w", err)
    }

    return nil
}
```

**Key Points:**
- `ListenAndServe()` - Opens a TCP socket and starts accepting connections
- **Blocking call** - This function doesn't return until the server stops
- `http.ErrServerClosed` - This is expected when gracefully shutting down (not an error)
- `fmt.Errorf("...: %w", err)` - Error wrapping for better error messages

---

### Step 6: Graceful Shutdown

```go
// Shutdown gracefully shuts down the server.
// It waits for active requests to finish before stopping.
func (e *Engine) Shutdown(ctx context.Context) error {
    log.Println("Shutting down server gracefully...")

    err := e.server.Shutdown(ctx)
    if err != nil {
        return fmt.Errorf("shutdown error: %w", err)
    }

    log.Println("Server stopped successfully")
    return nil
}
```

**Key Points:**
- `context.Context` - Allows setting a timeout for shutdown
- **Graceful shutdown means:**
  1. Stop accepting new connections
  2. Wait for existing requests to complete
  3. Close all connections
- Use `context.WithTimeout()` to prevent waiting forever for slow requests

**Example Usage:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
engine.Shutdown(ctx)
```

---

## Part 4: Understanding the Tests

### Test: Engine Creation

```go
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
    // ... more assertions
}
```

**Tests verify:**
- Engine is created successfully
- Default configuration is applied
- Internal server is initialized

---

### Test: HTTP Request Handling

```go
func TestServeHTTP(t *testing.T) {
    engine := New()
    req := httptest.NewRequest("GET", "/hello", nil)
    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, req)

    if recorder.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", recorder.Code)
    }
    // ... more assertions
}
```

**Key Testing Tools:**
- `httptest.NewRequest()` - Creates a fake HTTP request without network
- `httptest.NewRecorder()` - Captures the response for inspection
- Call `ServeHTTP()` directly - No need to start actual server!

---

## Summary

| Concept | What You Learned |
|---------|------------------|
| `http.Handler` | Interface with `ServeHTTP()` method - the foundation of Go HTTP handling |
| `http.Server` | Manages TCP connections, timeouts, and dispatches requests to handlers |
| Engine Pattern | Create a struct that implements `http.Handler` to build a framework |
| Graceful Shutdown | Use `Shutdown(ctx)` to wait for active requests before stopping |
| Testing | Use `httptest` package to test handlers without network overhead |

---

## Usage Example

```go
package main

import (
    "log"
    "your-module/goexpress"
)

func main() {
    // Create engine with defaults
    app := goexpress.New()

    // Or with custom config
    // app := goexpress.NewWithConfig(&goexpress.Config{
    //     Port:         ":3000",
    //     ReadTimeout:  30 * time.Second,
    //     WriteTimeout: 30 * time.Second,
    // })

    // Start server (blocking)
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## Next Steps (Lab 02)

In the next lab, you will:
- Create a `Context` struct to wrap request/response
- Implement request accessors (Query, Header, Body)
- Implement response methods (Status, SetHeader, Write)
- Add request-scoped storage (Set, Get)
