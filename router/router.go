package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"my-indexer/logger"
)

// Router handles both ElasticSearch-compatible and custom API endpoints
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	router := &Router{
		mux: http.NewServeMux(),
	}

	// Initialize the logger
	if err := logger.Initialize(); err != nil {
		logger.Error("Failed to initialize logger: %v", err)
	}

	return router
}

// Close performs cleanup of router resources
func (r *Router) Close() {
	logger.Close()
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Wrap the response writer to capture status code
	wrapped := logger.LoggingMiddleware(http.HandlerFunc(r.route))
	wrapped.ServeHTTP(w, req)
}

// route handles the actual routing logic
func (r *Router) route(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Handle document operations (e.g., /test-index/_doc/1)
	if strings.Contains(path, "/_doc/") {
		r.handleDocument(w, req)
		return
	}

	// Handle bulk operations (e.g., /test-index/_bulk)
	if strings.Contains(path, "/_bulk") {
		r.handleBulk(w, req)
		return
	}

	// Handle search operations (e.g., /test-index/_search)
	if strings.Contains(path, "/_search") {
		r.handleSearch(w, req)
		return
	}

	logger.Error("Invalid path: %s", path)
	http.Error(w, "Not found", http.StatusNotFound)
}

// RegisterElasticSearchHandlers registers all ElasticSearch-compatible endpoints
func (r *Router) RegisterElasticSearchHandlers() {
	// Document API endpoints
	r.mux.HandleFunc("/_doc/", r.handleDocument)           // Single document operations
	r.mux.HandleFunc("/_bulk", r.handleBulk)              // Bulk operations
	r.mux.HandleFunc("/_search", r.handleSearch)          // Search
	r.mux.HandleFunc("/_msearch", r.handleMultiSearch)    // Multi-search
	r.mux.HandleFunc("/_cat/indices", r.handleListIndices) // List indices
	r.mux.HandleFunc("/_scroll", r.handleScroll)          // Scroll API
}

// errorResponse sends an error response in JSON format
func errorResponse(w http.ResponseWriter, code int, message string) {
	logger.Error("Error response: %s (code: %d)", message, code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Handler functions for ElasticSearch-compatible endpoints
func (r *Router) handleDocument(w http.ResponseWriter, req *http.Request) {
	logger.Info("Handling document request: %s %s", req.Method, req.URL.Path)

	// Check method first
	if req.Method != http.MethodPut && req.Method != http.MethodGet && req.Method != http.MethodDelete {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := validateDocumentRequest(req); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	switch req.Method {
	case http.MethodPut:
		// TODO: Implement document creation/update
		logger.Info("Creating/updating document: %s", req.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": "created",
			"status": http.StatusOK,
		})
	case http.MethodGet:
		// TODO: Implement document retrieval
		logger.Info("Retrieving document: %s", req.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"found": true,
			"status": http.StatusOK,
		})
	case http.MethodDelete:
		// TODO: Implement document deletion
		logger.Info("Deleting document: %s", req.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": "deleted",
			"status": http.StatusOK,
		})
	}
}

func (r *Router) handleBulk(w http.ResponseWriter, req *http.Request) {
	logger.Info("Handling bulk request: %s", req.URL.Path)

	// Check method first
	if req.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := validateBulkRequest(req); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement bulk operations
	logger.Info("Processing bulk request")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"took": 0,
		"errors": false,
		"items": []interface{}{},
	})
}

func (r *Router) handleSearch(w http.ResponseWriter, req *http.Request) {
	logger.Info("Handling search request: %s %s", req.Method, req.URL.Path)

	// Check method first
	if req.Method != http.MethodGet && req.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := validateSearchRequest(req); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement search
	logger.Info("Processing search request")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"took": 0,
		"hits": map[string]interface{}{
			"total": map[string]interface{}{
				"value": 0,
				"relation": "eq",
			},
			"hits": []interface{}{},
		},
	})
}

func (r *Router) handleMultiSearch(w http.ResponseWriter, req *http.Request) {
	if err := validateSearchRequest(req); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Implement multi-search
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (r *Router) handleListIndices(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "only GET method is allowed")
		return
	}
	// TODO: Implement list indices
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (r *Router) handleScroll(w http.ResponseWriter, req *http.Request) {
	if err := validateSearchRequest(req); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// TODO: Implement scroll API
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
