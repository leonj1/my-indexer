package query

import (
	"my-indexer/document"
	"testing"
	"time"
)

func TestTermQuery(t *testing.T) {
	query := NewTermQuery("title", "test")

	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{"Exact match", "test", true},
		{"No match", "other", false},
		{"Non-string value", 123, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := query.Match(tt.value); got != tt.want {
				t.Errorf("TermQuery.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    *RangeQueryImpl
		doc      *document.Document
		expected bool
	}{
		{
			name: "Numeric range",
			query: &RangeQueryImpl{
				field: "age",
				gt:    10.0,
				lt:    20.0,
			},
			doc: func() *document.Document {
				doc := document.NewDocument()
				doc.AddField("age", 15.0)
				return doc
			}(),
			expected: true,
		},
		{
			name: "Inclusive numeric range",
			query: &RangeQueryImpl{
				field: "age",
				gte:   10.0,
				lte:   20.0,
			},
			doc: func() *document.Document {
				doc := document.NewDocument()
				doc.AddField("age", 20.0) // Test inclusive upper bound
				return doc
			}(),
			expected: true,
		},
		{
			name: "Mixed inclusive/exclusive range",
			query: &RangeQueryImpl{
				field: "age",
				gt:    10.0,
				lte:   20.0,
			},
			doc: func() *document.Document {
				doc := document.NewDocument()
				doc.AddField("age", 20.0) // Should match due to lte
				return doc
			}(),
			expected: true,
		},
		{
			name: "Outside range",
			query: &RangeQueryImpl{
				field: "age",
				gte:   10.0,
				lte:   20.0,
			},
			doc: func() *document.Document {
				doc := document.NewDocument()
				doc.AddField("age", 25.0)
				return doc
			}(),
			expected: false,
		},
		{
			name: "Time range",
			query: &RangeQueryImpl{
				field: "timestamp",
				gt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				lt:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			doc: func() *document.Document {
				doc := document.NewDocument()
				doc.AddField("timestamp", time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC))
				return doc
			}(),
			expected: true,
		},
		{
			name: "Inclusive time range",
			query: &RangeQueryImpl{
				field: "timestamp",
				gte:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				lte:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			doc: func() *document.Document {
				doc := document.NewDocument()
				doc.AddField("timestamp", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)) // Test inclusive upper bound
				return doc
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query.Match(tt.doc)
			if result != tt.expected {
				t.Errorf("%s: RangeQuery.Match() = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestBooleanQuery(t *testing.T) {
	t.Run("Must queries", func(t *testing.T) {
		query := NewBooleanQuery()
		query.AddMust(NewTermQuery("title", "test"))
		query.AddMust(NewTermQuery("category", "book"))

		tests := []struct {
			name  string
			value map[string]string
			want  bool
		}{
			{
				"All conditions met",
				map[string]string{"title": "test", "category": "book"},
				true,
			},
			{
				"One condition not met",
				map[string]string{"title": "test", "category": "article"},
				false,
			},
			{
				"No conditions met",
				map[string]string{"title": "other", "category": "article"},
				false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := query.Match(tt.value); got != tt.want {
					t.Errorf("BooleanQuery.Match() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Should queries", func(t *testing.T) {
		query := NewBooleanQuery()
		query.AddShould(NewTermQuery("title", "test"))
		query.AddShould(NewTermQuery("title", "example"))

		tests := []struct {
			name  string
			value string
			want  bool
		}{
			{"First condition met", "test", true},
			{"Second condition met", "example", true},
			{"No conditions met", "other", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := query.Match(tt.value); got != tt.want {
					t.Errorf("BooleanQuery.Match() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Must not queries", func(t *testing.T) {
		query := NewBooleanQuery()
		query.AddMustNot(NewTermQuery("status", "deleted"))

		tests := []struct {
			name  string
			value string
			want  bool
		}{
			{"Condition not met (good)", "active", true},
			{"Condition met (bad)", "deleted", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := query.Match(tt.value); got != tt.want {
					t.Errorf("BooleanQuery.Match() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestQueryMapper(t *testing.T) {
	mapper := NewQueryMapper()

	t.Run("Term query mapping", func(t *testing.T) {
		dslQuery := map[string]interface{}{
			"term": map[string]interface{}{
				"title": "test",
			},
		}

		query, err := mapper.MapQuery(dslQuery)
		if err != nil {
			t.Fatalf("MapQuery() error = %v", err)
		}

		if query.Type() != TermQuery {
			t.Errorf("Expected TermQuery, got %v", query.Type())
		}

		if !query.Match("test") {
			t.Error("Query should match 'test'")
		}
	})

	t.Run("Range query mapping", func(t *testing.T) {
		dslQuery := map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gt": 10.0,
					"lt": 20.0,
				},
			},
		}

		query, err := mapper.MapQuery(dslQuery)
		if err != nil {
			t.Fatalf("MapQuery() error = %v", err)
		}

		if query.Type() != RangeQuery {
			t.Errorf("Expected RangeQuery, got %v", query.Type())
		}

		if !query.Match(15.0) {
			t.Error("Query should match 15.0")
		}
	})

	t.Run("Bool query mapping", func(t *testing.T) {
		dslQuery := map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": "active",
						},
					},
				},
				"must_not": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"deleted": "true",
						},
					},
				},
			},
		}

		query, err := mapper.MapQuery(dslQuery)
		if err != nil {
			t.Fatalf("MapQuery() error = %v", err)
		}

		if query.Type() != BooleanQuery {
			t.Errorf("Expected BooleanQuery, got %v", query.Type())
		}
	})

	t.Run("Invalid query", func(t *testing.T) {
		dslQuery := map[string]interface{}{
			"invalid": map[string]interface{}{},
		}

		if _, err := mapper.MapQuery(dslQuery); err == nil {
			t.Error("Expected error for invalid query type")
		}
	})
}
