package elastic

import (
	"encoding/json"
	"fmt"
	"time"
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
	// MatchAll query that matches all documents
	MatchAllQuery QueryType = "match_all"
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
	Field string      // Field to search in
	Value interface{} // Value to search for (must be a string)
}

func (q *MatchQueryClause) MarshalJSON() ([]byte, error) {
	// Validate that Value is a string
	switch v := q.Value.(type) {
	case string:
		// Valid type, proceed with marshaling
	case fmt.Stringer:
		// Also accept types that implement String() string
		q.Value = v.String()
	default:
		return nil, fmt.Errorf("match query value must be a string, got %T", q.Value)
	}

	return json.Marshal(map[string]interface{}{
		"match": map[string]interface{}{
			q.Field: map[string]interface{}{
				"query": q.Value,
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
		"term": map[string]interface{}{
			q.Field: map[string]interface{}{
				"value": q.Value,
			},
		},
	})
}

// RangeQueryClause represents a range query
type RangeQueryClause struct {
	BaseQuery
	Field string
	GT    interface{} `json:"gt,omitempty"`   // Greater than value (must be numeric or time.Time)
	GTE   interface{} `json:"gte,omitempty"`  // Greater than or equal value (must be numeric or time.Time)
	LT    interface{} `json:"lt,omitempty"`   // Less than value (must be numeric or time.Time)
	LTE   interface{} `json:"lte,omitempty"`  // Less than or equal value (must be numeric or time.Time)
}

// validateRangeValue checks if a value is valid for range queries (numeric or time.Time)
func (q *RangeQueryClause) validateRangeValue(val interface{}) error {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return nil
	case json.Number:
		// Try to parse as float64 to validate it's a valid number
		if _, err := v.Float64(); err != nil {
			return fmt.Errorf("invalid numeric value: %v", err)
		}
		return nil
	case string:
		// Check if it's a valid time string
		if _, err := time.Parse(time.RFC3339, v); err != nil {
			return fmt.Errorf("string value must be a valid RFC3339 time format: %v", err)
		}
		return nil
	case time.Time:
		return nil
	default:
		return fmt.Errorf("range value must be numeric or time.Time, got %T", val)
	}
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
		"range": map[string]interface{}{
			q.Field: conditions,
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

	// Build bool query contents
	boolContents := make(map[string]interface{})

	if len(q.Must) > 0 {
		must := make([]interface{}, len(q.Must))
		for i, query := range q.Must {
			data, err := query.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var decoded interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			must[i] = decoded
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
			var decoded interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			should[i] = decoded
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
			var decoded interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			mustNot[i] = decoded
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
			var decoded interface{}
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, err
			}
			filter[i] = decoded
		}
		boolQuery["filter"] = filter
	}

	boolContents["bool"] = boolQuery

	// Wrap in query object
	return json.Marshal(map[string]interface{}{
		"query": boolContents,
	})
}

const maxBoolNestingDepth = 2 // Maximum allowed nesting depth for bool queries

type queryContext struct {
	depth int
	seenFields map[string]map[string]bool // clause type -> field -> seen
}

func newQueryContext() *queryContext {
	return &queryContext{
		depth: 0,
		seenFields: make(map[string]map[string]bool),
	}
}

func (ctx *queryContext) checkAndAddField(clauseType, field string) error {
	if _, exists := ctx.seenFields[clauseType]; !exists {
		ctx.seenFields[clauseType] = make(map[string]bool)
	}
	if ctx.seenFields[clauseType][field] {
		return fmt.Errorf("duplicate field '%s' in clause type '%s'", field, clauseType)
	}
	ctx.seenFields[clauseType][field] = true
	return nil
}

// MatchAllQueryClause represents a query that matches all documents
type MatchAllQueryClause struct {
	BaseQuery
}

func (q *MatchAllQueryClause) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	})
}

func ParseQuery(data []byte) (Query, error) {
	var wrapper struct {
		Query json.RawMessage `json:"query"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse query wrapper: %v", err)
	}

	if len(wrapper.Query) == 0 {
		return nil, fmt.Errorf("query field is required")
	}

	ctx := newQueryContext()
	return parseQueryClause(wrapper.Query, ctx)
}

func parseQueryClause(data []byte, ctx *queryContext) (Query, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse query clause: %v", err)
	}

	// Check bool nesting depth
	if _, ok := raw["bool"]; ok {
		ctx.depth++
		if ctx.depth > maxBoolNestingDepth {
			return nil, fmt.Errorf("bool query nesting depth exceeds maximum of %d", maxBoolNestingDepth)
		}
	}

	// Check for query wrapper
	if queryWrapper, ok := raw["query"]; ok {
		queryBytes, err := json.Marshal(queryWrapper)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query wrapper: %v", err)
		}
		return parseQueryClause(queryBytes, ctx)
	}

	for queryType, value := range raw {
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query value: %v", err)
		}

		switch queryType {
		case "match":
			return parseMatchQuery(valueBytes, ctx)
		case "term":
			return parseTermQuery(valueBytes, ctx)
		case "range":
			return parseRangeQuery(valueBytes, ctx)
		case "bool":
			return parseBoolQuery(raw, ctx)
		case "match_all":
			return parseMatchAllQuery(valueBytes, ctx)
		}
	}

	return nil, fmt.Errorf("invalid or unsupported query type")
}

func parseMatchQuery(data []byte, ctx *queryContext) (Query, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	if len(raw) != 1 {
		return nil, fmt.Errorf("match query must have exactly one field")
	}

	var field string
	var value interface{}

	for f, v := range raw {
		field = f
		switch val := v.(type) {
		case string:
			value = val
		case map[string]interface{}:
			if q, ok := val["query"]; ok {
				value = q
			} else {
				// If no query field is present, use the value directly
				value = val
			}
		default:
			value = val
		}
	}

	if err := ctx.checkAndAddField("match", field); err != nil {
		return nil, err
	}

	return &MatchQueryClause{
		BaseQuery: BaseQuery{queryType: MatchQuery},
		Field:     field,
		Value:     value,
	}, nil
}

func parseTermQuery(data []byte, ctx *queryContext) (Query, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	if len(raw) != 1 {
		return nil, fmt.Errorf("term query must have exactly one field")
	}

	var field string
	var value interface{}

	for f, v := range raw {
		field = f
		switch val := v.(type) {
		case map[string]interface{}:
			if v, ok := val["value"]; ok {
				value = v
			} else {
				// If no value field is present, use the value directly
				value = val
			}
		default:
			value = val
		}
	}

	if err := ctx.checkAndAddField("term", field); err != nil {
		return nil, err
	}

	return &TermQueryClause{
		BaseQuery: BaseQuery{queryType: TermQuery},
		Field:     field,
		Value:     value,
	}, nil
}

func parseRangeQuery(data []byte, ctx *queryContext) (Query, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	if len(raw) != 1 {
		return nil, fmt.Errorf("range query must have exactly one field")
	}

	var field string
	var rangeValues map[string]interface{}

	for f, v := range raw {
		field = f
		if values, ok := v.(map[string]interface{}); ok {
			rangeValues = values
		} else {
			return nil, fmt.Errorf("range query values must be an object")
		}
	}

	if err := ctx.checkAndAddField("range", field); err != nil {
		return nil, err
	}

	clause := &RangeQueryClause{
		BaseQuery: BaseQuery{queryType: RangeQuery},
		Field:     field,
	}

	if gt, ok := rangeValues["gt"]; ok {
		clause.GT = gt
	}
	if gte, ok := rangeValues["gte"]; ok {
		clause.GTE = gte
	}
	if lt, ok := rangeValues["lt"]; ok {
		clause.LT = lt
	}
	if lte, ok := rangeValues["lte"]; ok {
		clause.LTE = lte
	}

	return clause, nil
}

func parseBoolQuery(data map[string]interface{}, ctx *queryContext) (Query, error) {
	boolQuery := &BoolQueryClause{
		BaseQuery: BaseQuery{queryType: BoolQuery},
	}

	boolClauses, ok := data["bool"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid bool query structure")
	}

	// Validate bool query structure
	for key := range boolClauses {
		switch key {
		case "must", "should", "must_not", "filter":
			// Valid keys
		default:
			return nil, fmt.Errorf("invalid bool query clause: %s", key)
		}
	}
	if !ok {
		return nil, fmt.Errorf("invalid bool query structure")
	}

	// Process must clauses
	if mustClauses, ok := boolClauses["must"].([]interface{}); ok {
		for _, clause := range mustClauses {
			clauseBytes, err := json.Marshal(clause)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal must clause: %v", err)
			}
			query, err := parseQueryClause(clauseBytes, ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to parse must clause: %v", err)
			}
			boolQuery.Must = append(boolQuery.Must, query)
		}
	}

	// Process should clauses
	if shouldClauses, ok := boolClauses["should"].([]interface{}); ok {
		for _, clause := range shouldClauses {
			clauseBytes, err := json.Marshal(clause)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal should clause: %v", err)
			}
			query, err := parseQueryClause(clauseBytes, ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to parse should clause: %v", err)
			}
			boolQuery.Should = append(boolQuery.Should, query)
		}
	}

	// Process must_not clauses
	if mustNotClauses, ok := boolClauses["must_not"].([]interface{}); ok {
		for _, clause := range mustNotClauses {
			clauseBytes, err := json.Marshal(clause)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal must_not clause: %v", err)
			}
			query, err := parseQueryClause(clauseBytes, ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to parse must_not clause: %v", err)
			}
			boolQuery.MustNot = append(boolQuery.MustNot, query)
		}
	}

	// Process filter clauses
	if filterClauses, ok := boolClauses["filter"].([]interface{}); ok {
		for _, clause := range filterClauses {
			clauseBytes, err := json.Marshal(clause)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal filter clause: %v", err)
			}
			query, err := parseQueryClause(clauseBytes, ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to parse filter clause: %v", err)
			}
			boolQuery.Filter = append(boolQuery.Filter, query)
		}
	}

	return boolQuery, nil
}

func parseMatchAllQuery(data []byte, ctx *queryContext) (Query, error) {
	return &MatchAllQueryClause{
		BaseQuery: BaseQuery{queryType: MatchAllQuery},
	}, nil
}
