package elastic

import (
	"encoding/json"
	"fmt"
)

// QueryType represents the type of query
type QueryType string

const (
	// Match query for full text search
	MatchQuery QueryType = "match"
	// Term query for exact term matches
	TermQuery QueryType = "term"
	// Range query for numeric/date ranges
	RangeQuery QueryType = "range"
	// Bool query for combining multiple queries
	BoolQuery QueryType = "bool"
)

// Query represents the base query interface
type Query interface {
	Type() QueryType
	MarshalJSON() ([]byte, error)
}

// BaseQuery provides common query fields
type BaseQuery struct {
	queryType QueryType
}

func (q BaseQuery) Type() QueryType {
	return q.queryType
}

// MatchQueryClause represents a full text query
type MatchQueryClause struct {
	BaseQuery
	Field string
	Value interface{}
}

func (q *MatchQueryClause) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				q.Field: q.Value,
			},
		},
	})
}

// TermQueryClause represents an exact term query
type TermQueryClause struct {
	BaseQuery
	Field string
	Value interface{}
}

func (q *TermQueryClause) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				q.Field: q.Value,
			},
		},
	})
}

// RangeQueryClause represents a range query
type RangeQueryClause struct {
	BaseQuery
	Field string
	GT    interface{} `json:"gt,omitempty"`
	GTE   interface{} `json:"gte,omitempty"`
	LT    interface{} `json:"lt,omitempty"`
	LTE   interface{} `json:"lte,omitempty"`
}

func (q *RangeQueryClause) MarshalJSON() ([]byte, error) {
	conditions := make(map[string]interface{})
	if q.GT != nil {
		conditions["gt"] = q.GT
	}
	if q.GTE != nil {
		conditions["gte"] = q.GTE
	}
	if q.LT != nil {
		conditions["lt"] = q.LT
	}
	if q.LTE != nil {
		conditions["lte"] = q.LTE
	}

	return json.Marshal(map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				q.Field: conditions,
			},
		},
	})
}

// BoolQueryClause represents a boolean combination of queries
type BoolQueryClause struct {
	BaseQuery
	Must    []Query `json:"must,omitempty"`
	Should  []Query `json:"should,omitempty"`
	MustNot []Query `json:"must_not,omitempty"`
	Filter  []Query `json:"filter,omitempty"`
}

func (q *BoolQueryClause) MarshalJSON() ([]byte, error) {
	boolQuery := make(map[string]interface{})

	if len(q.Must) > 0 {
		must := make([]interface{}, len(q.Must))
		for i, query := range q.Must {
			data, err := query.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var decoded map[string]interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			must[i] = decoded["query"]
		}
		boolQuery["must"] = must
	}

	if len(q.Should) > 0 {
		should := make([]interface{}, len(q.Should))
		for i, query := range q.Should {
			data, err := query.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var decoded map[string]interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			should[i] = decoded["query"]
		}
		boolQuery["should"] = should
	}

	if len(q.MustNot) > 0 {
		mustNot := make([]interface{}, len(q.MustNot))
		for i, query := range q.MustNot {
			data, err := query.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var decoded map[string]interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			mustNot[i] = decoded["query"]
		}
		boolQuery["must_not"] = mustNot
	}

	if len(q.Filter) > 0 {
		filter := make([]interface{}, len(q.Filter))
		for i, query := range q.Filter {
			data, err := query.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var decoded map[string]interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			filter[i] = decoded["query"]
		}
		boolQuery["filter"] = filter
	}

	return json.Marshal(map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
	})
}

// ParseQuery parses an Elasticsearch query DSL into our internal query representation
func ParseQuery(data []byte) (Query, error) {
	var rawQuery map[string]interface{}
	if err := json.Unmarshal(data, &rawQuery); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query: %w", err)
	}

	queryBody, ok := rawQuery["query"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid query structure: missing 'query' object")
	}

	return parseQueryBody(queryBody)
}

// parseQueryBody parses the actual query body
func parseQueryBody(queryBody map[string]interface{}) (Query, error) {
	// Check for each query type
	if matchQuery, ok := queryBody["match"].(map[string]interface{}); ok {
		return parseMatchQuery(matchQuery)
	}
	if termQuery, ok := queryBody["term"].(map[string]interface{}); ok {
		return parseTermQuery(termQuery)
	}
	if rangeQuery, ok := queryBody["range"].(map[string]interface{}); ok {
		return parseRangeQuery(rangeQuery)
	}
	if boolQuery, ok := queryBody["bool"].(map[string]interface{}); ok {
		return parseBoolQuery(boolQuery)
	}

	return nil, fmt.Errorf("unsupported query type")
}

func parseMatchQuery(matchQuery map[string]interface{}) (*MatchQueryClause, error) {
	if len(matchQuery) != 1 {
		return nil, fmt.Errorf("invalid match query structure")
	}

	for field, value := range matchQuery {
		return &MatchQueryClause{
			BaseQuery: BaseQuery{queryType: MatchQuery},
			Field:     field,
			Value:     value,
		}, nil
	}

	return nil, fmt.Errorf("invalid match query structure")
}

func parseTermQuery(termQuery map[string]interface{}) (*TermQueryClause, error) {
	if len(termQuery) != 1 {
		return nil, fmt.Errorf("invalid term query structure")
	}

	for field, value := range termQuery {
		return &TermQueryClause{
			BaseQuery: BaseQuery{queryType: TermQuery},
			Field:     field,
			Value:     value,
		}, nil
	}

	return nil, fmt.Errorf("invalid term query structure")
}

func parseRangeQuery(rangeQuery map[string]interface{}) (*RangeQueryClause, error) {
	if len(rangeQuery) != 1 {
		return nil, fmt.Errorf("invalid range query structure")
	}

	for field, conditions := range rangeQuery {
		condMap, ok := conditions.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid range query conditions")
		}

		clause := &RangeQueryClause{
			BaseQuery: BaseQuery{queryType: RangeQuery},
			Field:     field,
		}

		if gt, ok := condMap["gt"]; ok {
			clause.GT = gt
		}
		if gte, ok := condMap["gte"]; ok {
			clause.GTE = gte
		}
		if lt, ok := condMap["lt"]; ok {
			clause.LT = lt
		}
		if lte, ok := condMap["lte"]; ok {
			clause.LTE = lte
		}

		return clause, nil
	}

	return nil, fmt.Errorf("invalid range query structure")
}

func parseBoolQuery(boolQuery map[string]interface{}) (*BoolQueryClause, error) {
	clause := &BoolQueryClause{
		BaseQuery: BaseQuery{queryType: BoolQuery},
	}

	if must, ok := boolQuery["must"].([]interface{}); ok {
		queries, err := parseQueryList(must)
		if err != nil {
			return nil, fmt.Errorf("failed to parse must clauses: %w", err)
		}
		clause.Must = queries
	}

	if should, ok := boolQuery["should"].([]interface{}); ok {
		queries, err := parseQueryList(should)
		if err != nil {
			return nil, fmt.Errorf("failed to parse should clauses: %w", err)
		}
		clause.Should = queries
	}

	if mustNot, ok := boolQuery["must_not"].([]interface{}); ok {
		queries, err := parseQueryList(mustNot)
		if err != nil {
			return nil, fmt.Errorf("failed to parse must_not clauses: %w", err)
		}
		clause.MustNot = queries
	}

	if filter, ok := boolQuery["filter"].([]interface{}); ok {
		queries, err := parseQueryList(filter)
		if err != nil {
			return nil, fmt.Errorf("failed to parse filter clauses: %w", err)
		}
		clause.Filter = queries
	}

	return clause, nil
}

func parseQueryList(queries []interface{}) ([]Query, error) {
	result := make([]Query, 0, len(queries))
	for _, q := range queries {
		queryMap, ok := q.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid query in list")
		}

		query, err := parseQueryBody(queryMap)
		if err != nil {
			return nil, err
		}
		result = append(result, query)
	}
	return result, nil
}
