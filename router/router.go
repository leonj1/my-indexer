package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"my-indexer/logger"
	"my-indexer/index"
	"my-indexer/analysis"
	"my-indexer/document"
	"my-indexer/search"
	"my-indexer/query"
)

// IndexDocumentStore adapts Index to implement search.DocumentStore
type IndexDocumentStore struct {
	idx *index.Index
}

// LoadDocument implements search.DocumentStore
func (s *IndexDocumentStore) LoadDocument(docID int) (*document.Document, error) {
	return s.idx.GetDocument(docID)
}

// LoadAllDocuments implements search.DocumentStore
func (s *IndexDocumentStore) LoadAllDocuments() ([]*document.Document, error) {
	return s.idx.GetAllDocuments()
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
	// Only allow GET and POST methods
	if req.Method != http.MethodGet && req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var queryMapObj map[string]interface{}
	var err error

	if req.Method == http.MethodGet {
		// For GET requests without a query parameter, use match_all query
		queryStr := req.URL.Query().Get("q")
		if queryStr == "" {
			queryMapObj = map[string]interface{}{
				"match_all": map[string]interface{}{},
			}
		} else {
			// Create a match query for the q parameter
			queryMapObj = map[string]interface{}{
				"match": map[string]interface{}{
					"_all": queryStr,
				},
			}
		}
	} else {
		// Parse query from request body for POST
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		var searchRequest struct {
			Query map[string]interface{} `json:"query"`
		}

		if err := json.Unmarshal(body, &searchRequest); err != nil {
			http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
			return
		}

		if searchRequest.Query == nil {
			http.Error(w, "Query object is required", http.StatusBadRequest)
			return
		}

		queryMapObj = searchRequest.Query
	}

	// Initialize query mapper
	queryMapper := query.NewQueryMapper()

	// Prepare query wrapper
	var queryWrapper map[string]interface{}

	// If the query is a direct match/term/range/bool query
	if queryType, ok := getQueryType(queryMapObj); ok {
		switch queryType {
		case "match", "term", "match_phrase", "match_all", "range", "bool":
			// For match queries, ensure proper structure
			if queryType == "match" {
				if fieldMap, ok := queryMapObj[queryType].(map[string]interface{}); ok {
					for field, fieldValue := range fieldMap {
						// If the field value is a string, wrap it in a query object
						if _, ok := fieldValue.(string); ok {
							fieldMap[field] = map[string]interface{}{
								"query": fieldValue,
							}
						}
					}
				}
			}
			queryWrapper = queryMapObj
			goto processQuery
		}
	}

	// If no valid query type found, treat the entire query object as a match query
	queryWrapper = map[string]interface{}{"match": queryMapObj}

processQuery:
	// Pass the query object to the mapper
	queryObj, err := queryMapper.MapQuery(queryWrapper)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to map query: %v", err), http.StatusBadRequest)
		return
	}

	// Execute the query
	results, err := r.search.SearchWithQuery(queryObj)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute search: %v", err), http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func getQueryType(query map[string]interface{}) (string, bool) {
	for queryType := range query {
		return queryType, true
	}
	return "", false
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
