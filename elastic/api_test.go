package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

// MockAPI implements the API interface for testing
type MockAPI struct {
	documents map[string]*Document
	mu        sync.RWMutex
}

func NewMockAPI() *MockAPI {
	return &MockAPI{
		documents: make(map[string]*Document),
	}
}

func (m *MockAPI) Index(ctx context.Context, doc *Document) (*Document, error) {
	if doc.Source == nil || len(doc.Source) == 0 {
		return nil, fmt.Errorf("document source is empty")
	}
	
	// Check for malformed document
	if _, hasSource := doc.Source["_source"]; hasSource {
		return nil, fmt.Errorf("malformed document: contains reserved field '_source'")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	key := doc.Index + ":" + doc.ID
	m.documents[key] = doc
	return doc, nil
}

func (m *MockAPI) Get(ctx context.Context, index, id string) (*Document, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := index + ":" + id
	if doc, ok := m.documents[key]; ok {
		return doc, nil
	}
	return nil, nil
}

func (m *MockAPI) Update(ctx context.Context, doc *Document) (*Document, error) {
	return m.Index(ctx, doc)
}

func (m *MockAPI) Delete(ctx context.Context, index, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := index + ":" + id
	delete(m.documents, key)
	return nil
}

func (m *MockAPI) Search(ctx context.Context, query map[string]interface{}) (*SearchResponse, error) {
	if query == nil {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	}
	
	// If no query specified, treat as match_all
	if _, hasQuery := query["query"]; !hasQuery {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	// Validate query structure
	queryObj, ok := query["query"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid query format: 'query' must be an object")
	}

	// Check for empty query object
	if len(queryObj) == 0 {
		return nil, fmt.Errorf("invalid query format: empty query object")
	}

	// Must have exactly one query type
	validTypes := []string{"match", "match_all", "term", "range", "bool"}
	foundType := ""
	for _, qType := range validTypes {
		if queryValue, exists := queryObj[qType]; exists {
			if foundType != "" {
				return nil, fmt.Errorf("invalid query format: multiple query types specified")
			}
			// Validate the query value is a map
			if _, ok := queryValue.(map[string]interface{}); !ok {
				return nil, fmt.Errorf("invalid query format: %s value must be an object", qType)
			}
			foundType = qType
		}
	}
	
	// Check if any invalid query types are present
	for qType := range queryObj {
		isValid := false
		for _, validType := range validTypes {
			if qType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, fmt.Errorf("invalid query format: unsupported query type '%s'", qType)
		}
	}

	if foundType == "" {
		return nil, fmt.Errorf("invalid query format: no valid query type found")
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Simple mock implementation that returns all documents
	hits := make([]Document, 0)
	for _, doc := range m.documents {
		hits = append(hits, *doc)
	}
	return &SearchResponse{
		Took: 1,
		Hits: SearchHits{
			Total: Total{Value: int64(len(hits)), Relation: TotalRelationEq},
			Hits:  hits,
		},
	}, nil
}

func (m *MockAPI) MultiSearch(ctx context.Context, queries []map[string]interface{}) ([]*SearchResponse, error) {
	responses := make([]*SearchResponse, len(queries))
	for i, query := range queries {
		resp, err := m.Search(ctx, query)
		if err != nil {
			return nil, err
		}
		responses[i] = resp
	}
	return responses, nil
}

func (m *MockAPI) ListIndices(ctx context.Context) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	indices := make(map[string]bool)
	for _, doc := range m.documents {
		indices[doc.Index] = true
	}
	result := make([]string, 0, len(indices))
	for index := range indices {
		result = append(result, index)
	}
	return result, nil
}

func (m *MockAPI) Bulk(ctx context.Context, operations []json.RawMessage) error {
	return nil
}

func TestElasticAPI(t *testing.T) {
	ctx := context.Background()
	api := NewMockAPI()

	// Test document operations
	doc := &Document{
		Index:  "test",
		ID:     "1",
		Source: map[string]interface{}{"field": "value"},
	}

	// Test Index
	indexed, err := api.Index(ctx, doc)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}
	if indexed.ID != doc.ID {
		t.Errorf("Expected ID %s, got %s", doc.ID, indexed.ID)
	}

	// Test Get
	retrieved, err := api.Get(ctx, "test", "1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved.ID != doc.ID {
		t.Errorf("Expected ID %s, got %s", doc.ID, retrieved.ID)
	}

	// Test Search
	searchResp, err := api.Search(ctx, map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if searchResp.Hits.Total.Value != 1 {
		t.Errorf("Expected 1 hit, got %d", searchResp.Hits.Total.Value)
	}

	// Test Delete
	err = api.Delete(ctx, "test", "1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	deleted, err := api.Get(ctx, "test", "1")
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if deleted != nil {
		t.Error("Document still exists after deletion")
	}
}

func TestTotalRelation(t *testing.T) {
	tests := []struct {
		name     string
		relation TotalRelation
		want     bool
	}{
		{
			name:     "Valid eq relation",
			relation: TotalRelationEq,
			want:     true,
		},
		{
			name:     "Valid gte relation",
			relation: TotalRelationGte,
			want:     true,
		},
		{
			name:     "Invalid relation",
			relation: TotalRelation("invalid"),
			want:     false,
		},
		{
			name:     "Empty relation",
			relation: TotalRelation(""),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.relation.IsValid(); got != tt.want {
				t.Errorf("TotalRelation.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElasticSearchCompatibility(t *testing.T) {
	ctx := context.Background()
	api := NewMockAPI()

	t.Run("ElasticSearch Request Format", func(t *testing.T) {
		// Test standard ElasticSearch query format
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"field": "value",
				},
			},
			"size": 10,
			"from": 0,
		}
		
		resp, err := api.Search(ctx, query)
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		
		// Validate ElasticSearch response format
		if resp.Took == 0 {
			t.Error("Expected non-zero took value")
		}
		if resp.Hits.Total.Relation != TotalRelationEq {
			t.Error("Expected eq relation for total")
		}
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		const numOps = 100
		errCh := make(chan error, numOps)
		doneCh := make(chan bool, numOps)

		for i := 0; i < numOps; i++ {
			go func(idx int) {
				doc := &Document{
					Index:  "concurrent-test",
					ID:     fmt.Sprintf("doc-%d", idx),
					Source: map[string]interface{}{"value": idx},
				}
				_, err := api.Index(ctx, doc)
				if err != nil {
					errCh <- err
					return
				}
				doneCh <- true
			}(i)
		}

		// Wait for all operations
		for i := 0; i < numOps; i++ {
			select {
			case err := <-errCh:
				t.Errorf("Concurrent operation failed: %v", err)
			case <-doneCh:
				// Operation successful
			}
		}
	})

	t.Run("Large Payload", func(t *testing.T) {
		// Create a large document (>1MB)
		largeData := make([]string, 100000)
		for i := range largeData {
			largeData[i] = "test data string that takes up space"
		}

		doc := &Document{
			Index:  "large-test",
			ID:     "large-1",
			Source: map[string]interface{}{"data": largeData},
		}

		indexed, err := api.Index(ctx, doc)
		if err != nil {
			t.Fatalf("Large document indexing failed: %v", err)
		}

		// Verify the indexed document has the correct ID
		if indexed.ID != "large-1" {
			t.Errorf("Expected document ID 'large-1', got '%s'", indexed.ID)
		}

		// Verify we can retrieve the large document
		retrieved, err := api.Get(ctx, "large-test", "large-1")
		if err != nil {
			t.Fatalf("Failed to retrieve large document: %v", err)
		}
		if retrieved == nil {
			t.Error("Large document not found after indexing")
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		// Test invalid query format
		invalidQuery := map[string]interface{}{
			"query": map[string]interface{}{
				"invalid_type": map[string]interface{}{
					"field": "value",
				},
			},
		}
		_, err := api.Search(ctx, invalidQuery)
		if err == nil {
			t.Error("Expected error for invalid query format")
		}

		// Test multiple query types
		multipleTypesQuery := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"field": "value",
				},
				"term": map[string]interface{}{
					"field": "value",
				},
			},
		}
		_, err = api.Search(ctx, multipleTypesQuery)
		if err == nil {
			t.Error("Expected error for multiple query types")
		}

		// Test invalid query value type
		invalidValueQuery := map[string]interface{}{
			"query": map[string]interface{}{
				"match": "invalid_value",
			},
		}
		_, err = api.Search(ctx, invalidValueQuery)
		if err == nil {
			t.Error("Expected error for invalid query value type")
		}

		// Test empty document source
		emptyDoc := &Document{
			Index:  "test",
			ID:     "empty",
			Source: map[string]interface{}{},
		}
		_, err = api.Index(ctx, emptyDoc)
		if err == nil {
			t.Error("Expected error for empty document source")
		}

		// Test malformed document
		malformedDoc := &Document{
			Index: "test",
			ID:    "malformed",
			Source: map[string]interface{}{
				"_source": nil,
			},
		}
		_, err = api.Index(ctx, malformedDoc)
		if err == nil {
			t.Error("Expected error for malformed document")
		}
	})
}
