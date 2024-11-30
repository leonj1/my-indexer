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

func ParseQuery(data []byte) (Query, error) {
	var wrapper struct {
		Query json.RawMessage `json:"query"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse query wrapper: %v", err)
	}

	ctx := newQueryContext()
	return parseQueryClause(wrapper.Query, ctx)
}

func parseQueryClause(data []byte, ctx *queryContext) (Query, error) {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, fmt.Errorf("failed to parse query clause: %v", err)
	}

	for queryType, value := range temp {
		switch queryType {
		case "match":
			var matchQuery struct {
				Field string `json:"-"`
				Value string
			}
			var raw map[string]string
			if err := json.Unmarshal(value, &raw); err != nil {
				return nil, fmt.Errorf("failed to parse match query: %v", err)
			}
			for field, val := range raw {
				if err := ctx.checkAndAddField("match", field); err != nil {
					return nil, err
				}
				matchQuery.Field = field
				matchQuery.Value = val
				return &MatchQueryClause{
					BaseQuery: BaseQuery{queryType: MatchQuery},
					Field: matchQuery.Field,
					Value: matchQuery.Value,
				}, nil
			}
			
		case "term":
			var termQuery struct {
				Field string `json:"-"`
				Value interface{}
			}
			var raw map[string]interface{}
			if err := json.Unmarshal(value, &raw); err != nil {
				return nil, fmt.Errorf("failed to parse term query: %v", err)
			}
			for field, val := range raw {
				if err := ctx.checkAndAddField("term", field); err != nil {
					return nil, err
				}
				termQuery.Field = field
				termQuery.Value = val
				return &TermQueryClause{
					BaseQuery: BaseQuery{queryType: TermQuery},
					Field: termQuery.Field,
					Value: termQuery.Value,
				}, nil
			}

		case "range":
			var raw map[string]map[string]interface{}
			if err := json.Unmarshal(value, &raw); err != nil {
				return nil, fmt.Errorf("failed to parse range query: %v", err)
			}

			for field, conditions := range raw {
				if err := ctx.checkAndAddField("range", field); err != nil {
					return nil, err
				}

				rangeQuery := &RangeQueryClause{
					BaseQuery: BaseQuery{queryType: RangeQuery},
					Field:    field,
				}

				for op, val := range conditions {
					switch op {
					case "gt":
						rangeQuery.GT = val
					case "gte":
						rangeQuery.GTE = val
					case "lt":
						rangeQuery.LT = val
					case "lte":
						rangeQuery.LTE = val
					default:
						return nil, fmt.Errorf("invalid range operator: %s", op)
					}
				}

				return rangeQuery, nil
			}

			return nil, fmt.Errorf("invalid range query structure")

		case "bool":
			ctx.depth++
			if ctx.depth > maxBoolNestingDepth {
				return nil, fmt.Errorf("bool query nesting depth exceeds maximum allowed (%d)", maxBoolNestingDepth)
			}

			var boolQuery struct {
				Must    []json.RawMessage `json:"must,omitempty"`
				Should  []json.RawMessage `json:"should,omitempty"`
				MustNot []json.RawMessage `json:"must_not,omitempty"`
				Filter  []json.RawMessage `json:"filter,omitempty"`
			}
			if err := json.Unmarshal(value, &boolQuery); err != nil {
				return nil, fmt.Errorf("failed to parse bool query: %v", err)
			}

			result := &BoolQueryClause{
				BaseQuery: BaseQuery{queryType: BoolQuery},
			}

			// Parse Must clauses
			for _, clause := range boolQuery.Must {
				q, err := parseQueryClause(clause, ctx)
				if err != nil {
					return nil, err
				}
				result.Must = append(result.Must, q)
			}

			// Parse Should clauses
			for _, clause := range boolQuery.Should {
				q, err := parseQueryClause(clause, ctx)
				if err != nil {
					return nil, err
				}
				result.Should = append(result.Should, q)
			}

			// Parse MustNot clauses
			for _, clause := range boolQuery.MustNot {
				q, err := parseQueryClause(clause, ctx)
				if err != nil {
					return nil, err
				}
				result.MustNot = append(result.MustNot, q)
			}

			// Parse Filter clauses
			for _, clause := range boolQuery.Filter {
				q, err := parseQueryClause(clause, ctx)
				if err != nil {
					return nil, err
				}
				result.Filter = append(result.Filter, q)
			}

			ctx.depth--
			return result, nil
		}
	}

	return nil, fmt.Errorf("invalid or unsupported query type")
}
