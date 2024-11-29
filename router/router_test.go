package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test-index/_search", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
