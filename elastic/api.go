package elastic

import (
	"context"
	"encoding/json"
)

// Document represents an Elasticsearch document
type Document struct {
	Index      string                 `json:"_index"`
	ID         string                 `json:"_id,omitempty"`
	Source     map[string]interface{} `json:"_source"`
	Version    int64                  `json:"_version,omitempty"`
	SeqNo      int64                  `json:"_seq_no,omitempty"`
	PrimaryTerm int64                 `json:"_primary_term,omitempty"`
}

// SearchResponse represents an Elasticsearch search response
type SearchResponse struct {
	Took     int64      `json:"took"`
	TimedOut bool       `json:"timed_out"`
	Hits     SearchHits `json:"hits"`
}

// SearchHits contains search results and handles null max_score values
type SearchHits struct {
	Total    Total      `json:"total"`
	MaxScore *float64   `json:"max_score,omitempty"`
	Hits     []Document `json:"hits"`
}

// Total represents the total number of hits
type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

// API defines the Elasticsearch-compatible API interface
type API interface {
	// Document APIs
	Index(ctx context.Context, doc *Document) (*Document, error)
	Get(ctx context.Context, index, id string) (*Document, error)
	Update(ctx context.Context, doc *Document) (*Document, error)
	Delete(ctx context.Context, index, id string) error

	// Search APIs
	Search(ctx context.Context, query map[string]interface{}) (*SearchResponse, error)
	MultiSearch(ctx context.Context, queries []map[string]interface{}) ([]*SearchResponse, error)

	// Index APIs
	ListIndices(ctx context.Context) ([]string, error)

	// Bulk Operations
	Bulk(ctx context.Context, operations []json.RawMessage) error
}
