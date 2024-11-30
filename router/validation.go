package router

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// MaxRequestBodySize is 10MB
	MaxRequestBodySize = 10 * 1024 * 1024
)

var (
	ErrInvalidJSON     = errors.New("invalid JSON in request body")
	ErrBodyTooLarge    = errors.New("request body too large")
	ErrEmptyBody       = errors.New("request body is empty")
	ErrMissingBody     = errors.New("request body is required")
	ErrInvalidIndex    = errors.New("invalid index name")
	ErrInvalidDocID    = errors.New("invalid document ID")
	ErrInvalidBulkData = errors.New("invalid bulk request data")
)

// validateRequestBody checks if the request body is present and not too large
func validateRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, ErrMissingBody
	}
	defer r.Body.Close()

	// Set size limit on request body
	r.Body = http.MaxBytesReader(nil, r.Body, MaxRequestBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			return nil, ErrBodyTooLarge
		}
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	if len(body) == 0 {
		return nil, ErrEmptyBody
	}

	return body, nil
}

// validateJSONBody validates that the request body contains valid JSON
func validateJSONBody(body []byte) error {
	var js json.RawMessage
	if err := json.Unmarshal(body, &js); err != nil {
		return ErrInvalidJSON
	}
	return nil
}

// validateDocumentRequest validates a document API request
func validateDocumentRequest(r *http.Request) error {
	// Extract and validate index name and document ID from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		return ErrInvalidIndex
	}

	indexName := parts[1]
	if indexName == "" {
		return ErrInvalidIndex
	}

	if len(parts) >= 4 && parts[3] == "" {
		return ErrInvalidDocID
	}

	// For PUT requests, validate the request body
	if r.Method == http.MethodPut {
		body, err := validateRequestBody(r)
		if err != nil {
			return err
		}

		if err := validateJSONBody(body); err != nil {
			return err
		}
	}

	return nil
}

// validateBulkRequest validates a bulk API request
func validateBulkRequest(r *http.Request) error {
	// Validate Content-Type for NDJSON format
	if r.Header.Get("Content-Type") != "application/x-ndjson" {
		return fmt.Errorf("invalid Content-Type, expected application/x-ndjson")
	}

	// Limit request body size to 10MB
	r.Body = http.MaxBytesReader(nil, r.Body, 10<<20)
	defer r.Body.Close()

	// Read and validate each line as a separate JSON object
	scanner := bufio.NewScanner(r.Body)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Skip empty lines
		}

		var action map[string]interface{}
		if err := json.Unmarshal([]byte(line), &action); err != nil {
			return fmt.Errorf("invalid JSON at line %d: %v", lineCount+1, err)
		}

		// Validate action type
		if len(action) != 1 {
			return fmt.Errorf("invalid action at line %d: exactly one action type expected", lineCount+1)
		}

		// Check for valid action types
		validAction := false
		for _, actionType := range []string{"index", "create", "update", "delete"} {
			if _, ok := action[actionType]; ok {
				validAction = true
				break
			}
		}
		if !validAction {
			return fmt.Errorf("invalid action type at line %d: must be one of index, create, update, or delete", lineCount+1)
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading request body: %v", err)
	}

	if lineCount == 0 {
		return fmt.Errorf("empty bulk request")
	}

	return nil
}

// validateSearchRequest validates a search API request
func validateSearchRequest(r *http.Request) error {
	// Extract and validate index name from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return ErrInvalidIndex
	}

	indexName := parts[1]
	if indexName == "" {
		return ErrInvalidIndex
	}

	// For POST requests, validate the request body
	if r.Method == http.MethodPost {
		// Limit request body size to 10MB to prevent memory exhaustion
		r.Body = http.MaxBytesReader(nil, r.Body, 10<<20)
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			if err.Error() == "http: request body too large" {
				return fmt.Errorf("request body exceeds 10MB limit")
			}
			return fmt.Errorf("invalid JSON: %v", err)
		}
	}

	return nil
}
