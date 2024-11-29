package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// validateDocumentRequest validates requests for document operations
func validateDocumentRequest(r *http.Request) error {
	// Method validation is now handled in the handler
	if r.Method == http.MethodPut {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return fmt.Errorf("invalid JSON: %v", err)
		}
	}
	return nil
}

// validateBulkRequest validates bulk operation requests
func validateBulkRequest(r *http.Request) error {
	// Method validation is now handled in the handler
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}
	return nil
}

// validateSearchRequest validates search operation requests
func validateSearchRequest(r *http.Request) error {
	// Method validation is now handled in the handler
	if r.Method == http.MethodPost {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return fmt.Errorf("invalid JSON: %v", err)
		}
	}
	return nil
}
