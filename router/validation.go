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
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	
	// Check basic path structure
	if len(parts) < 3 {
		if len(parts) >= 2 && parts[1] == "_doc" {
			return ErrInvalidDocID
		}
		return ErrInvalidIndex
	}

	// Validate index name first
	indexName := parts[0]
	if indexName == "" || indexName == "_doc" {
		return ErrInvalidIndex
	}

	// Validate _doc part
	if parts[1] != "_doc" {
		return ErrInvalidIndex
	}

	// Validate document ID
	if parts[2] == "" {
		return ErrInvalidDocID
	}

	// For PUT requests, validate the request body
	if r.Method == http.MethodPut {
		if r.Body == nil {
			return ErrMissingBody
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return ErrMissingBody
		}

		if !json.Valid(body) {
			return ErrInvalidJSON
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

	// Validate method
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		return &validationError{
			status:  http.StatusMethodNotAllowed,
			message: "only GET and POST methods are allowed",
		}
	}

	// For POST requests, validate the Content-Type and request body
	if r.Method == http.MethodPost {
		// Validate Content-Type header
		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			mediaType := strings.ToLower(strings.Split(contentType, ";")[0])
			if mediaType != "application/json" {
				return &validationError{
					status:  http.StatusBadRequest,
					message: fmt.Sprintf("invalid Content-Type: expected 'application/json', got '%s'", contentType),
				}
			}
		}

		// Limit request body size to 10MB to prevent memory exhaustion
		r.Body = http.MaxBytesReader(nil, r.Body, 10<<20)
		
		// Read and validate JSON structure
		var body map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		
		if err := decoder.Decode(&body); err != nil {
			switch {
			case err.Error() == "http: request body too large":
				return fmt.Errorf("request body exceeds 10MB limit")
			case strings.Contains(err.Error(), "cannot unmarshal"):
				return fmt.Errorf("malformed JSON: %v", err)
			default:
				return fmt.Errorf("invalid JSON: %v", err)
			}
		}

		// Validate query structure if present
		if query, exists := body["query"]; exists {
			queryMap, ok := query.(map[string]interface{})
			if !ok {
				return &validationError{
					status:  http.StatusBadRequest,
					message: "'query' must be an object",
				}
			}

			// Validate supported query types
			if len(queryMap) == 0 {
				return &validationError{
					status:  http.StatusBadRequest,
					message: "empty query object",
				}
			}

			// Check for valid query types
			validQueryTypes := map[string]bool{
				"match_all": true,
				"match":     true,
				"term":      true,
				"range":     true,
				"bool":      true,
			}

			// Helper function to validate field values
			validateFieldValue := func(value interface{}) error {
				switch v := value.(type) {
				case map[string]interface{}, string, float64, int, bool, []interface{}:
					return nil // Allow arrays, objects and primitive types
				default:
					return fmt.Errorf("invalid field value type: %T", v)
				}
			}

			// Helper function to validate query clauses recursively
			var validateQueryClause func(clause map[string]interface{}) error
			validateQueryClause = func(clause map[string]interface{}) error {
				for queryType, value := range clause {
					if !validQueryTypes[queryType] {
						return &validationError{
							status:  http.StatusBadRequest,
							message: fmt.Sprintf("unsupported query type: %s", queryType),
						}
					}

					switch queryType {
					case "match_all":
						// match_all can be empty object or null
						if value != nil {
							_, ok := value.(map[string]interface{})
							if !ok {
								return &validationError{
									status:  http.StatusBadRequest,
									message: "match_all value must be an object or null",
								}
							}
						}
					case "match", "term":
						valueMap, ok := value.(map[string]interface{})
						if !ok {
							return &validationError{
								status:  http.StatusBadRequest,
								message: fmt.Sprintf("%s query must be an object", queryType),
							}
						}
						for field, fieldValue := range valueMap {
							if err := validateFieldValue(fieldValue); err != nil {
								return &validationError{
									status:  http.StatusBadRequest,
									message: fmt.Sprintf("invalid value for field %s: %v", field, err),
								}
							}
						}
					case "range":
						rangeMap, ok := value.(map[string]interface{})
						if !ok {
							return &validationError{
								status:  http.StatusBadRequest,
								message: "range query must be an object",
							}
						}
						for field, conditions := range rangeMap {
							condMap, ok := conditions.(map[string]interface{})
							if !ok {
								return &validationError{
									status:  http.StatusBadRequest,
									message: fmt.Sprintf("range conditions for field %s must be an object", field),
								}
							}
							for op, val := range condMap {
								switch op {
								case "gt", "gte", "lt", "lte", "eq":
									if err := validateFieldValue(val); err != nil {
										return &validationError{
											status:  http.StatusBadRequest,
											message: fmt.Sprintf("invalid range value for %s: %v", op, err),
										}
									}
								default:
									return &validationError{
										status:  http.StatusBadRequest,
										message: fmt.Sprintf("unsupported range operator: %s", op),
									}
								}
							}
						}
					case "bool":
						boolMap, ok := value.(map[string]interface{})
						if !ok {
							return &validationError{
								status:  http.StatusBadRequest,
								message: "bool query must be an object",
							}
						}
						for clause, clauseValue := range boolMap {
							switch clause {
							case "must", "should", "must_not", "filter":
								switch clauses := clauseValue.(type) {
								case []interface{}:
									for _, subQuery := range clauses {
										subQueryMap, ok := subQuery.(map[string]interface{})
										if !ok {
											return &validationError{
												status:  http.StatusBadRequest,
												message: fmt.Sprintf("bool %s array elements must be objects", clause),
											}
										}
										if err := validateQueryClause(subQueryMap); err != nil {
											return err
										}
									}
								case map[string]interface{}:
									if err := validateQueryClause(clauses); err != nil {
										return err
									}
								default:
									return &validationError{
										status:  http.StatusBadRequest,
										message: fmt.Sprintf("bool %s clauses must be an object or array of objects", clause),
									}
								}
							case "minimum_should_match":
								switch v := clauseValue.(type) {
								case float64, int:
									// These are valid types
								default:
									return &validationError{
										status:  http.StatusBadRequest,
										message: fmt.Sprintf("minimum_should_match must be a number, got %T", v),
									}
								}
							default:
								return &validationError{
									status:  http.StatusBadRequest,
									message: fmt.Sprintf("unsupported bool operation: %s", clause),
								}
							}
						}
					}
				}
				return nil
			}

			// Validate the query structure
			if err := validateQueryClause(queryMap); err != nil {
				return err
			}
		}
	}

	return nil
}

type validationError struct {
	status  int
	message string
}

func (e *validationError) Error() string {
	return e.message
}
