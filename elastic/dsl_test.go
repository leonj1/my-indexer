package elastic

import (
	"encoding/json"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name: "match query",
			query: `{
				"query": {
					"match": {
						"title": "golang programming"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "term query",
			query: `{
				"query": {
					"term": {
						"status": "active"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "range query",
			query: `{
				"query": {
					"range": {
						"age": {
							"gte": 18,
							"lte": 65
						}
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "bool query",
			query: `{
				"query": {
					"bool": {
						"must": [
							{
								"match": {
									"title": "golang"
								}
							}
						],
						"filter": [
							{
								"term": {
									"status": "published"
								}
							}
						]
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "invalid query structure",
			query: `{
				"invalid": {}
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := ParseQuery([]byte(tt.query))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && query == nil {
				t.Error("ParseQuery() returned nil query for valid input")
			}

			// Additional type-specific validations
			if !tt.wantErr {
				switch q := query.(type) {
				case *MatchQueryClause:
					if q.Type() != MatchQuery {
						t.Errorf("Expected MatchQuery type, got %v", q.Type())
					}
				case *TermQueryClause:
					if q.Type() != TermQuery {
						t.Errorf("Expected TermQuery type, got %v", q.Type())
					}
				case *RangeQueryClause:
					if q.Type() != RangeQuery {
						t.Errorf("Expected RangeQuery type, got %v", q.Type())
					}
				case *BoolQueryClause:
					if q.Type() != BoolQuery {
						t.Errorf("Expected BoolQuery type, got %v", q.Type())
					}
				}
			}
		})
	}
}

func TestParseComplexQueries(t *testing.T) {
	query := `{
		"query": {
			"bool": {
				"must": [
					{
						"match": {
							"title": "golang"
						}
					},
					{
						"range": {
							"year": {
								"gte": 2020
							}
						}
					}
				],
				"should": [
					{
						"term": {
							"tags": "programming"
						}
					}
				],
				"must_not": [
					{
						"term": {
							"status": "draft"
						}
					}
				],
				"filter": [
					{
						"term": {
							"published": true
						}
					}
				]
			}
		}
	}`

	q, err := ParseQuery([]byte(query))
	if err != nil {
		t.Fatalf("ParseQuery() error = %v", err)
	}

	boolQuery, ok := q.(*BoolQueryClause)
	if !ok {
		t.Fatal("Expected BoolQueryClause")
	}

	// Validate must clauses
	if len(boolQuery.Must) != 2 {
		t.Errorf("Expected 2 must clauses, got %d", len(boolQuery.Must))
	}

	// Validate should clauses
	if len(boolQuery.Should) != 1 {
		t.Errorf("Expected 1 should clause, got %d", len(boolQuery.Should))
	}

	// Validate must_not clauses
	if len(boolQuery.MustNot) != 1 {
		t.Errorf("Expected 1 must_not clause, got %d", len(boolQuery.MustNot))
	}

	// Validate filter clauses
	if len(boolQuery.Filter) != 1 {
		t.Errorf("Expected 1 filter clause, got %d", len(boolQuery.Filter))
	}
}

func TestQueryToJSON(t *testing.T) {
	// Create a complex query
	query := &BoolQueryClause{
		BaseQuery: BaseQuery{queryType: BoolQuery},
		Must: []Query{
			&MatchQueryClause{
				BaseQuery: BaseQuery{queryType: MatchQuery},
				Field:     "title",
				Value:     "golang",
			},
		},
		Filter: []Query{
			&TermQueryClause{
				BaseQuery: BaseQuery{queryType: TermQuery},
				Field:     "status",
				Value:     "active",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(query)
	if err != nil {
		t.Fatalf("Failed to marshal query: %v", err)
	}

	// Parse back
	parsedQuery, err := ParseQuery(data)
	if err != nil {
		t.Fatalf("Failed to parse marshaled query: %v", err)
	}

	// Validate parsed query
	boolQuery, ok := parsedQuery.(*BoolQueryClause)
	if !ok {
		t.Fatal("Expected BoolQueryClause")
	}

	if len(boolQuery.Must) != 1 {
		t.Errorf("Expected 1 must clause, got %d", len(boolQuery.Must))
	}

	if len(boolQuery.Filter) != 1 {
		t.Errorf("Expected 1 filter clause, got %d", len(boolQuery.Filter))
	}
}
