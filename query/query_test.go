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
		query := NewRangeQuery("price")
		query.SetGT(10.0)
		query.SetLTE(20.0)

		tests := []struct {
			value float64
			want  bool
		}{
			{5.0, false},
			{10.0, false},
			{15.0, true},
			{20.0, true},
			{25.0, false},
		}

		for _, tt := range tests {
			if got := query.Match(tt.value); got != tt.want {
				t.Errorf("RangeQuery.Match(%v) = %v, want %v", tt.value, got, tt.want)
			}
		}
	})

	t.Run("Time range", func(t *testing.T) {
		query := NewRangeQuery("timestamp")
		now := time.Now()
		query.SetGTE(now)
		query.SetLT(now.Add(24 * time.Hour))

		tests := []struct {
			value time.Time
			want  bool
		}{
			{now.Add(-1 * time.Hour), false},
			{now, true},
			{now.Add(12 * time.Hour), true},
			{now.Add(24 * time.Hour), false},
			{now.Add(25 * time.Hour), false},
		}

		for _, tt := range tests {
			if got := query.Match(tt.value); got != tt.want {
				t.Errorf("RangeQuery.Match(%v) = %v, want %v", tt.value, got, tt.want)
			}
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
					"gt":  10.0,
					"lte": 20.0,
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
