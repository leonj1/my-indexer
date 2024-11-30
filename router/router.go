package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"my-indexer/logger"
	"my-indexer/index"
	"my-indexer/analysis"
	"my-indexer/document"
	"my-indexer/search"
)

// IndexDocumentStore adapts Index to implement search.DocumentStore
type IndexDocumentStore struct {
	idx *index.Index
}

// LoadDocument implements search.DocumentStore
func (s *IndexDocumentStore) LoadDocument(docID int) (*document.Document, error) {
	return s.idx.GetDocument(docID)
}

// Router handles HTTP requests for the indexer
type Router struct {
	mux    *http.ServeMux
	index  *index.Index
	search *search.Search
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	analyzer := analysis.NewStandardAnalyzer()
	idx := index.NewIndex(analyzer)
	store := &IndexDocumentStore{idx: idx}
	
	router := &Router{
		mux:    http.NewServeMux(),
		index:  idx,
		search: search.NewSearch(idx, store),
	}

	// Initialize the logger
	logger.Initialize()

	// Register handlers
	router.RegisterElasticSearchHandlers()

	return router
}

// Close performs cleanup of router resources
func (r *Router) Close() {
	logger.Close()
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Log the request
	logger.Info("Received request: %s %s", req.Method, req.URL.Path)

	// Handle the request based on the path
	if strings.Contains(req.URL.Path, "/_doc/") {
		r.handleDocument(w, req)
		return
	}

	if strings.HasSuffix(req.URL.Path, "/_bulk") {
		r.handleBulk(w, req)
		return
	}

	if strings.Contains(req.URL.Path, "/_search") {
		r.handleSearch(w, req)
		return
	}

	if strings.HasSuffix(req.URL.Path, "/_msearch") {
		r.handleMultiSearch(w, req)
		return
	}

	if strings.HasSuffix(req.URL.Path, "/_cat/indices") {
		r.handleListIndices(w, req)
		return
	}

	if strings.HasSuffix(req.URL.Path, "/_scroll") {
		r.handleScroll(w, req)
		return
	}

	if strings.HasSuffix(req.URL.Path, "/_index") {
		r.handleIndex(w, req)
		return
	}

	// Not found
	http.NotFound(w, req)
}

// RegisterElasticSearchHandlers registers all ElasticSearch-compatible endpoints
func (r *Router) RegisterElasticSearchHandlers() {
	// Document API endpoints
	r.mux.HandleFunc("/", r.handleDocument)                // Single document operations (matches /index/_doc/id)
	r.mux.HandleFunc("/_index", r.handleIndex)            // Index API endpoint
	r.mux.HandleFunc("/_bulk", r.handleBulk)              // Bulk operations
	r.mux.HandleFunc("/_search", r.handleSearch)          // Search
	r.mux.HandleFunc("/_msearch", r.handleMultiSearch)    // Multi-search
	r.mux.HandleFunc("/_cat/indices", r.handleListIndices) // List indices
	r.mux.HandleFunc("/_scroll", r.handleScroll)          // Scroll API
}

// ElasticSearchResponse represents a standard ES response format
type ElasticSearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Result string `json:"result,omitempty"`
	Status int    `json:"status,omitempty"`
}

// errorResponse sends an error response in JSON format
func (r *Router) errorResponse(w http.ResponseWriter, code int, message string) {
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
		r.errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Validate the request
	if err := validateDocumentRequest(req); err != nil {
		r.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Extract index name and document ID from path
	parts := strings.Split(req.URL.Path, "/")
	indexName := parts[1]
	docID := parts[3]

	switch req.Method {
	case http.MethodPut:
		// TODO: Implement document creation/update
		logger.Info("Creating/updating document: index=%s, id=%s", indexName, docID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_index": indexName,
			"_id":    docID,
			"result": "created",
			"status": http.StatusOK,
		})

	case http.MethodGet:
		// TODO: Implement document retrieval
		logger.Info("Retrieving document: index=%s, id=%s", indexName, docID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_index": indexName,
			"_id":    docID,
			"found":  true,
			"status": http.StatusOK,
		})

	case http.MethodDelete:
		// TODO: Implement document deletion
		logger.Info("Deleting document: index=%s, id=%s", indexName, docID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_index": indexName,
			"_id":    docID,
			"result": "deleted",
			"status": http.StatusOK,
		})
	}
}

func (r *Router) handleSearch(w http.ResponseWriter, req *http.Request) {
	logger.Info("Handling search request: %s %s", req.Method, req.URL.Path)

	// Validate request method
	if req.Method != http.MethodGet && req.Method != http.MethodPost {
		r.errorResponse(w, http.StatusMethodNotAllowed, "only GET and POST methods are allowed")
		return
	}

	// Validate the request
	if err := validateSearchRequest(req); err != nil {
		r.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Extract index name from path
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) < 3 {
		r.errorResponse(w, http.StatusBadRequest, "invalid index name")
		return
	}
	indexName := parts[1]

	// Start timing the search operation
	startTime := time.Now()

	// Default query for GET requests
	var terms []string
	if req.Method == http.MethodGet {
		terms = []string{"*:*"} // Match all query
	} else {
		// Validate Content-Type for POST requests
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			r.errorResponse(w, http.StatusBadRequest, "Content-Type must be application/json")
			return
		}

		// Parse request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			// If JSON parsing fails, default to match_all query
			terms = append(terms, "*:*")
		} else {
			// For POST requests, parse the query
			if queryRaw, ok := requestBody["query"]; ok {
				// Handle match_all query
				if queryMap, ok := queryRaw.(map[string]interface{}); ok {
					if _, ok := queryMap["match_all"]; ok {
						terms = append(terms, "*:*")
					} else {
						// Convert other query types to terms
						for field, value := range queryMap {
							if valueMap, ok := value.(map[string]interface{}); ok {
								for k, v := range valueMap {
									terms = append(terms, fmt.Sprintf("%s:%v", k, v))
								}
							} else {
								terms = append(terms, fmt.Sprintf("%s:%v", field, value))
							}
						}
					}
				}
			}

			// If no terms were added, default to match_all
			if len(terms) == 0 {
				terms = append(terms, "*:*")
			}
		}
	}

	// Execute search
	results, err := r.search.Search(terms, search.OR) // Use OR for better recall
	if err != nil {
		r.errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to execute search: %v", err))
		return
	}

	// Format results
	took := time.Since(startTime)
	esResponse := search.FormatESResponse(results, took, indexName)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(esResponse); err != nil {
		logger.Error("Failed to encode response: %v", err)
		r.errorResponse(w, http.StatusInternalServerError, "failed to encode response")
		return
	}
}

type validationError struct {
	status  int
	message string
}

func (e *validationError) Error() string {
	return e.message
}

func (r *Router) handleMultiSearch(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement multi-search
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (r *Router) handleListIndices(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.errorResponse(w, http.StatusMethodNotAllowed, "only GET method is allowed")
		return
	}
	// TODO: Implement list indices
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (r *Router) handleScroll(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement scroll API
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (r *Router) handleIndex(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodPut {
		r.errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse the request body
	var doc map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&doc); err != nil {
		r.errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Extract index name and document ID from URL path
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 1 {
		r.errorResponse(w, http.StatusBadRequest, "invalid index path")
		return
	}

	indexName := parts[0]
	docID := ""
	if len(parts) > 1 {
		docID = parts[1]
	}

	// Index the document
	startTime := time.Now()
	err := r.index.IndexDocument(indexName, docID, doc)
	if err != nil {
		r.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare ElasticSearch-compatible response
	resp := ElasticSearchResponse{
		Took:     int(time.Since(startTime).Milliseconds()),
		TimedOut: false,
		Shards: struct {
			Total      int `json:"total"`
			Successful int `json:"successful"`
			Failed     int `json:"failed"`
		}{
			Total:      1,
			Successful: 1,
			Failed:     0,
		},
		Result: "created",
		Status: http.StatusCreated,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
