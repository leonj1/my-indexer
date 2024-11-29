package logger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInitialize(t *testing.T) {
	// Clean up any existing log files
	os.RemoveAll("logs")
	defer os.RemoveAll("logs")

	// Test initialization
	err := Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer Close()

	// Check if log files were created
	files := []string{"info.log", "error.log", "request.log"}
	for _, file := range files {
		if _, err := os.Stat("logs/" + file); os.IsNotExist(err) {
			t.Errorf("Log file %s was not created", file)
		}
	}
}

func TestLogging(t *testing.T) {
	// Clean up any existing log files
	os.RemoveAll("logs")
	defer os.RemoveAll("logs")

	// Initialize loggers
	err := Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer Close()

	// Test info logging
	Info("Test info message")
	content, err := os.ReadFile("logs/info.log")
	if err != nil {
		t.Errorf("Failed to read info log: %v", err)
	}
	if !strings.Contains(string(content), "Test info message") {
		t.Error("Info message was not logged correctly")
	}

	// Test error logging
	Error("Test error message")
	content, err = os.ReadFile("logs/error.log")
	if err != nil {
		t.Errorf("Failed to read error log: %v", err)
	}
	if !strings.Contains(string(content), "Test error message") {
		t.Error("Error message was not logged correctly")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Clean up any existing log files
	os.RemoveAll("logs")
	defer os.RemoveAll("logs")

	// Initialize loggers
	err := Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer Close()

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a test server with the logging middleware
	server := httptest.NewServer(LoggingMiddleware(handler))
	defer server.Close()

	// Make a test request
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Wait a bit for the log to be written
	time.Sleep(100 * time.Millisecond)

	// Verify the request was logged
	content, err := os.ReadFile("logs/request.log")
	if err != nil {
		t.Errorf("Failed to read request log: %v", err)
	}
	logContent := string(content)
	if !strings.Contains(logContent, "Method=GET") {
		t.Error("Request method was not logged")
	}
	if !strings.Contains(logContent, "StatusCode=200") {
		t.Error("Status code was not logged")
	}
}

func TestResponseWriter(t *testing.T) {
	// Create a test response writer
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Test default status code
	if rw.statusCode != http.StatusOK {
		t.Errorf("Default status code should be 200, got %d", rw.statusCode)
	}

	// Test writing a different status code
	rw.WriteHeader(http.StatusNotFound)
	if rw.statusCode != http.StatusNotFound {
		t.Errorf("Status code should be 404, got %d", rw.statusCode)
	}
}
