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
	tests := []struct {
		name    string
		query   string
		wantErr bool
		validate func(*testing.T, Query)
	}{
		{
			name: "Valid complex bool query",
			query: `{
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
			}`,
			wantErr: false,
			validate: func(t *testing.T, q Query) {
				boolQuery, ok := q.(*BoolQueryClause)
				if !ok {
					t.Fatal("Expected BoolQueryClause")
				}

				// Validate must clauses
				if len(boolQuery.Must) != 2 {
					t.Errorf("Expected 2 must clauses, got %d", len(boolQuery.Must))
				}

				// Validate match query in must clause
				if matchQuery, ok := boolQuery.Must[0].(*MatchQueryClause); ok {
					if matchQuery.Field != "title" || matchQuery.Value != "golang" {
						t.Errorf("Expected match query with field='title', value='golang', got field='%s', value='%v'",
							matchQuery.Field, matchQuery.Value)
					}
				} else {
					t.Error("First must clause should be a MatchQueryClause")
				}

				// Validate range query in must clause
				if rangeQuery, ok := boolQuery.Must[1].(*RangeQueryClause); ok {
					if rangeQuery.Field != "year" {
						t.Errorf("Expected range query with field='year', got field='%s'", rangeQuery.Field)
					}
					if rangeQuery.GTE != float64(2020) {
						t.Errorf("Expected range query with gte=2020, got %v", rangeQuery.GTE)
					}
				} else {
					t.Error("Second must clause should be a RangeQueryClause")
				}

				// Validate should clause
				if len(boolQuery.Should) != 1 {
					t.Errorf("Expected 1 should clause, got %d", len(boolQuery.Should))
				}
				if termQuery, ok := boolQuery.Should[0].(*TermQueryClause); ok {
					if termQuery.Field != "tags" || termQuery.Value != "programming" {
						t.Errorf("Expected term query with field='tags', value='programming', got field='%s', value='%v'",
							termQuery.Field, termQuery.Value)
					}
				} else {
					t.Error("Should clause should be a TermQueryClause")
				}

				// Validate must_not clause
				if len(boolQuery.MustNot) != 1 {
					t.Errorf("Expected 1 must_not clause, got %d", len(boolQuery.MustNot))
				}
				if termQuery, ok := boolQuery.MustNot[0].(*TermQueryClause); ok {
					if termQuery.Field != "status" || termQuery.Value != "draft" {
						t.Errorf("Expected term query with field='status', value='draft', got field='%s', value='%v'",
							termQuery.Field, termQuery.Value)
					}
				} else {
					t.Error("Must_not clause should be a TermQueryClause")
				}

				// Validate filter clause
				if len(boolQuery.Filter) != 1 {
					t.Errorf("Expected 1 filter clause, got %d", len(boolQuery.Filter))
				}
				if termQuery, ok := boolQuery.Filter[0].(*TermQueryClause); ok {
					if termQuery.Field != "published" || termQuery.Value != true {
						t.Errorf("Expected term query with field='published', value=true, got field='%s', value=%v",
							termQuery.Field, termQuery.Value)
					}
				} else {
					t.Error("Filter clause should be a TermQueryClause")
				}
			},
		},
		{
			name: "Invalid - duplicate must clauses",
			query: `{
				"query": {
					"bool": {
						"must": [
							{
								"match": {
									"title": "golang"
								}
							},
							{
								"match": {
									"title": "golang"
								}
							}
						]
					}
				}
			}`,
			wantErr: true,
		},
		{
			name: "Invalid - nested bool exceeds depth",
			query: `{
				"query": {
					"bool": {
						"must": [
							{
								"bool": {
									"must": [
										{
											"bool": {
												"must": [
													{
														"match": {
															"title": "golang"
														}
													}
												]
											}
										}
									]
								}
							}
						]
					}
				}
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := ParseQuery([]byte(tt.query))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, q)
			}
		})
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