package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	requestLogger *log.Logger
)

// Initialize sets up the loggers with the specified output files
func Initialize() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Open log files
	infoFile, err := os.OpenFile("logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open info log file: %v", err)
	}

	errorFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}

	requestFile, err := os.OpenFile("logs/request.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open request log file: %v", err)
	}

	// Initialize loggers with prefix and flags
	infoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	requestLogger = log.New(requestFile, "REQUEST: ", log.Ldate|log.Ltime)

	return nil
}

// Info logs an informational message
func Info(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
}

// LogRequest logs HTTP request details
func LogRequest(r *http.Request, statusCode int, duration time.Duration) {
	if requestLogger != nil {
		requestLogger.Printf(
			"Method=%s Path=%s StatusCode=%d Duration=%v RemoteAddr=%s UserAgent=%s",
			r.Method,
			r.URL.Path,
			statusCode,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
		)
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware is a middleware that logs request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		
		// Call the next handler
		next.ServeHTTP(rw, r)
		
		// Log the request after it's handled
		duration := time.Since(start)
		LogRequest(r, rw.statusCode, duration)
	})
}
