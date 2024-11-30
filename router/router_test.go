package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateDocumentRequest(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		body        string
		wantErr     error
	}{
		{
			name:    "Valid PUT request",
			method:  http.MethodPut,
			path:    "/test-index/_doc/123",
			body:    `{"field": "value"}`,
			wantErr: nil,
		},
		{
			name:    "Valid GET request",
			method:  http.MethodGet,
			path:    "/test-index/_doc/123",
			body:    "",
			wantErr: nil,
		},
		{
			name:    "Valid DELETE request",
			method:  http.MethodDelete,
			path:    "/test-index/_doc/123",
			body:    "",
			wantErr: nil,
		},
		{
			name:    "Missing index name",
			method:  http.MethodPut,
			path:    "/_doc/123",
			body:    `{"field": "value"}`,
			wantErr: ErrInvalidIndex,
		},
		{
			name:    "Missing document ID",
			method:  http.MethodPut,
			path:    "/test-index/_doc/",
			body:    `{"field": "value"}`,
			wantErr: ErrInvalidDocID,
		},
		{
			name:    "Invalid path format",
			method:  http.MethodPut,
			path:    "/test-index/123",
			body:    `{"field": "value"}`,
			wantErr: ErrInvalidIndex,
		},
		{
			name:    "PUT request without body",
			method:  http.MethodPut,
			path:    "/test-index/_doc/123",
			body:    "",
			wantErr: ErrMissingBody,
		},
		{
			name:    "PUT request with invalid JSON",
			method:  http.MethodPut,
			path:    "/test-index/_doc/123",
			body:    `{"field": invalid}`,
			wantErr: ErrInvalidJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.method == http.MethodPut && tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			
			err := validateDocumentRequest(req)
			
			if err != tt.wantErr {
				t.Errorf("validateDocumentRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocumentEndpoint(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Valid PUT request",
			method:         http.MethodPut,
			path:          "/test-index/_doc/1",
			body:          `{"field": "value"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid GET request",
			method:         http.MethodGet,
			path:          "/test-index/_doc/1",
			body:          "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid DELETE request",
			method:         http.MethodDelete,
			path:          "/test-index/_doc/1",
			body:          "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			path:          "/test-index/_doc/1",
			body:          "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON in PUT request",
			method:         http.MethodPut,
			path:          "/test-index/_doc/1",
			body:          `{"field": invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestBulkEndpoint(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "Valid bulk request",
			method:         http.MethodPost,
			body:          `{"index": {"_index": "test", "_id": "1"}}
{"field1": "value1", "field2": "value2"}
{"index": {"_index": "test", "_id": "2"}}
{"field1": "value3", "field2": "value4"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPut,
			body:          "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			body:          `{"invalid`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test/_bulk", bytes.NewBufferString(tt.body))
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/x-ndjson")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestSearchEndpoint(t *testing.T) {
	router := NewRouter()

	// Add test data
	req := httptest.NewRequest(http.MethodPut, "/test-index/_doc/1", strings.NewReader(`{"field": "value"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("failed to set up test data: %d", w.Code)
	}

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "Valid GET request",
			method:         http.MethodGet,
			body:          "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid POST request",
			method:         http.MethodPost,
			body:          `{"query": {"match_all": {}}}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPut,
			body:          "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON in POST request",
			method:         http.MethodPost,
			body:          `{"query": invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unsupported query type",
			method:         http.MethodPost,
			body:          `{"query": {"unknown_type": {}}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty query object",
			method:         http.MethodPost,
			body:          `{"query": {}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid query structure",
			method:         http.MethodPost,
			body:          `{"query": "not an object"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Match query",
			method:         http.MethodPost,
			body:          `{"query": {"match": {"field": "value"}}}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Term query",
			method:         http.MethodPost,
			body:          `{"query": {"term": {"field": "value"}}}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Range query",
			method:         http.MethodPost,
			body:          `{"query": {"range": {"field": {"gt": 5}}}}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Bool query",
			method:         http.MethodPost,
			body:          `{"query": {"bool": {"must": [{"match": {"field": "value"}}]}}}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test-index/_search", strings.NewReader(tt.body))
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d for test %s", tt.expectedStatus, w.Code, tt.name)
			}
		})
	}
}
