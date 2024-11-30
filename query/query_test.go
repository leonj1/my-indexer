package query

import (
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
	t.Run("Numeric range", func(t *testing.T) {
		q := NewRangeQuery("price")
		q.GreaterThan(10.0)
		q.LessThan(20.0)

		if !q.Match(15.0) {
			t.Error("Expected 15.0 to match range query")
		}
		if q.Match(5.0) {
			t.Error("Expected 5.0 to not match range query")
		}
		if q.Match(25.0) {
			t.Error("Expected 25.0 to not match range query")
		}
	})

	t.Run("Time range", func(t *testing.T) {
		now := time.Now()
		before := now.Add(-1 * time.Hour)
		after := now.Add(1 * time.Hour)

		q := NewRangeQuery("timestamp")
		q.GreaterThan(before)
		q.LessThan(after)

		if !q.Match(now) {
			t.Error("Expected now to match range query")
		}
		if q.Match(before.Add(-1 * time.Hour)) {
			t.Error("Expected time before range to not match")
		}
		if q.Match(after.Add(1 * time.Hour)) {
			t.Error("Expected time after range to not match")
		}
	})
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
