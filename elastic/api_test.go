package elastic

import (
	"context"
	"encoding/json"
	"testing"
)

// MockAPI implements the API interface for testing
type MockAPI struct {
	documents map[string]*Document
}

func NewMockAPI() *MockAPI {
	return &MockAPI{
		documents: make(map[string]*Document),
	}
}

func (m *MockAPI) Index(ctx context.Context, doc *Document) (*Document, error) {
	key := doc.Index + ":" + doc.ID
	m.documents[key] = doc
	return doc, nil
}

func (m *MockAPI) Get(ctx context.Context, index, id string) (*Document, error) {
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
	key := index + ":" + id
	delete(m.documents, key)
	return nil
}

func (m *MockAPI) Search(ctx context.Context, query map[string]interface{}) (*SearchResponse, error) {
	// Simple mock implementation that returns all documents
	hits := make([]Document, 0)
	for _, doc := range m.documents {
		hits = append(hits, *doc)
	}
	return &SearchResponse{
		Took: 1,
		Hits: SearchHits{
			Total: Total{Value: int64(len(hits)), Relation: "eq"},
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
	searchResp, err := api.Search(ctx, map[string]interface{}{})
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
