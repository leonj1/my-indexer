package router

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"my-indexer/document"
)

// handleBulk handles bulk operations
func (r *Router) handleBulk(w http.ResponseWriter, req *http.Request) {
	// Only allow POST method
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate content type
	if req.Header.Get("Content-Type") != "application/x-ndjson" {
		http.Error(w, "Content-Type must be application/x-ndjson", http.StatusBadRequest)
		return
	}

	// Parse path to get index name
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) < 2 || parts[1] == "" {
		http.Error(w, "Invalid index name", http.StatusBadRequest)
		return
	}
	indexName := parts[1]

	// Process bulk request
	scanner := bufio.NewScanner(req.Body)
	defer req.Body.Close()

	var responses []map[string]interface{}
	var currentAction map[string]interface{}
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Skip empty lines
		}

		lineNum++
		if lineNum%2 == 1 {
			// Action line
			if err := json.Unmarshal([]byte(line), &currentAction); err != nil {
				http.Error(w, fmt.Sprintf("Invalid JSON at line %d: %v", lineNum, err), http.StatusBadRequest)
				return
			}

			// Validate action
			if len(currentAction) != 1 {
				http.Error(w, fmt.Sprintf("Invalid action at line %d: exactly one action type expected", lineNum), http.StatusBadRequest)
				return
			}

			// Check for valid action types
			validAction := false
			for _, actionType := range []string{"index", "create", "update", "delete"} {
				if _, ok := currentAction[actionType]; ok {
					validAction = true
					break
				}
			}
			if !validAction {
				http.Error(w, fmt.Sprintf("Invalid action type at line %d: must be one of index, create, update, or delete", lineNum), http.StatusBadRequest)
				return
			}
		} else {
			// Document line (for index/create/update operations)
			var doc map[string]interface{}
			if err := json.Unmarshal([]byte(line), &doc); err != nil {
				http.Error(w, fmt.Sprintf("Invalid JSON at line %d: %v", lineNum, err), http.StatusBadRequest)
				return
			}

			// Process the action
			response := make(map[string]interface{})
			switch {
			case currentAction["index"] != nil:
				// Create a new document
				newDoc := document.NewDocument()
				for field, value := range doc {
					newDoc.AddField(field, value)
				}

				// Add the document to the index
				docID, err := r.index.AddDocument(newDoc)
				if err != nil {
					response["index"] = map[string]interface{}{
						"_index":  indexName,
						"_id":     fmt.Sprintf("%d", docID),
						"status":  "error",
						"message": err.Error(),
					}
				} else {
					response["index"] = map[string]interface{}{
						"_index": indexName,
						"_id":    fmt.Sprintf("%d", docID),
						"status": "success",
					}
				}
			// Add other action types (create, update, delete) here
			default:
				http.Error(w, "Unsupported action type", http.StatusBadRequest)
				return
			}
			responses = append(responses, response)
		}
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"took":      0, // TODO: Add timing
		"errors":    false,
		"responses": responses,
	})
}
