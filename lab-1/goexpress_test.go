package goexpress

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestServeHTTP verifies that the Engine responds correctly to HTTP requests
func TestServeHTTP(t *testing.T) {
	engine := New()
	req := httptest.NewRequest("GET", "/hello", nil)
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Expected Content-Type text/plain, got %s", contentType)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "Hello from GoExpress!") {
		t.Errorf("Expected response to contain 'Hello from GoExpress!', got: %s", body)
	}
	if !strings.Contains(body, "/hello") {
		t.Errorf("Expected response to contain requested path '/hello', got: %s", body)
	}
}

// // TestServeHTTPWithDifferentPaths tests various URL paths
// func TestServeHTTPWithDifferentPaths(t *testing.T) {
// 	engine := New()
// 	paths := []string{
// 		"/",
// 		"/users",
// 		"/users/123",
// 		"/api/v1/posts",
// 		"/path/with/many/segments",
// 	}
// 	for _, path := range paths {
// 		t.Run(path, func(t *testing.T) {
// 			req := httptest.NewRequest("GET", path, nil)
// 			recorder := httptest.NewRecorder()

// 			engine.ServeHTTP(recorder, req)
// 			if recorder.Code != http.StatusOK {
// 				t.Errorf("Expected status 200 for %s, got %d", path, recorder.Code)
// 			}
// 			body := recorder.Body.String()
// 			if !strings.Contains(body, path) {
// 				t.Errorf("Expected response to contain path '%s', got: %s", path, body)
// 			}
// 		})
// 	}
// }

// // TestMultipleRequests simulates multiple concurrent requests
// func TestMultipleRequests(t *testing.T) {
// 	engine := New()
// 	numRequests := 100
// 	results := make(chan error, numRequests)
// 	for i := 0; i < numRequests; i++ {
// 		go func(id int) {
// 			req := httptest.NewRequest("GET", "/test", nil)
// 			recorder := httptest.NewRecorder()
// 			engine.ServeHTTP(recorder, req)
// 			if recorder.Code != http.StatusOK {
// 				results <- http.ErrHandlerTimeout // Use error to signal failure
// 			} else {
// 				results <- nil
// 			}
// 		}(i)
// 	}
// 	failures := 0
// 	for i := 0; i < numRequests; i++ {
// 		if err := <-results; err != nil {
// 			failures++
// 		}
// 	}
// 	if failures > 0 {
// 		t.Errorf("Expected 0 failures, got %d out of %d requests", failures, numRequests)
// 	}
// }
