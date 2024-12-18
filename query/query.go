package query

import (
	"fmt"
	"my-indexer/document"
	"strings"
	"time"
)

// QueryType represents the type of internal query
type QueryType int

const (
	// TermQuery for exact term matches
	TermQuery QueryType = iota
	// PhraseQuery for phrase matches
	PhraseQuery
	// FieldQuery for field-specific queries
	FieldQuery
	// PrefixQuery for prefix matches
	PrefixQuery
	// RangeQuery for range comparisons
	RangeQuery
	// BooleanQuery for combining multiple queries
	BooleanQuery
	// MatchQuery for analyzed text queries
	MatchQuery
	// MatchPhraseQuery for exact phrase matches
	MatchPhraseQuery
	// MatchAllQuery for matching all documents
	MatchAllQuery
)

// Query represents the internal query interface
type Query interface {
	Type() QueryType
	Field() string
	Match(value interface{}) bool
}

// TermQueryImpl represents an exact term match query
type TermQueryImpl struct {
	field string
	term  string
}

func NewTermQuery(field, term string) *TermQueryImpl {
	return &TermQueryImpl{field: field, term: term}
}

func (q *TermQueryImpl) Type() QueryType { return TermQuery }
func (q *TermQueryImpl) Field() string   { return q.field }
func (q *TermQueryImpl) Term() string    { return q.term }
func (q *TermQueryImpl) Match(value interface{}) bool {
	if str, ok := value.(string); ok {
		return str == q.term
	}
	return false
}

// RangeQueryImpl implements a range query
type RangeQueryImpl struct {
	field string
	gt    interface{} // Exclusive greater than
	gte   interface{} // Inclusive greater than or equal
	lt    interface{} // Exclusive less than
	lte   interface{} // Inclusive less than or equal
}

func NewRangeQuery(field string) *RangeQueryImpl {
	return &RangeQueryImpl{field: field}
}

func (q *RangeQueryImpl) Type() QueryType { return RangeQuery }
func (q *RangeQueryImpl) Field() string   { return q.field }
func (q *RangeQueryImpl) Gt() interface{} { return q.gt }
func (q *RangeQueryImpl) Lt() interface{} { return q.lt }
func (q *RangeQueryImpl) Gte() interface{} { return q.gte }
func (q *RangeQueryImpl) Lte() interface{} { return q.lte }

// GreaterThan sets the exclusive greater than value for the range query
func (q *RangeQueryImpl) GreaterThan(val interface{}) {
	q.gt = val
	q.gte = nil // Clear inclusive operator
}

// GreaterThanOrEqual sets the inclusive greater than or equal value for the range query
func (q *RangeQueryImpl) GreaterThanOrEqual(val interface{}) {
	q.gte = val
	q.gt = nil // Clear exclusive operator
}

// LessThan sets the exclusive less than value for the range query
func (q *RangeQueryImpl) LessThan(val interface{}) {
	q.lt = val
	q.lte = nil // Clear inclusive operator
}

// LessThanOrEqual sets the inclusive less than or equal value for the range query
func (q *RangeQueryImpl) LessThanOrEqual(val interface{}) {
	q.lte = val
	q.lt = nil // Clear exclusive operator
}

func (q *RangeQueryImpl) matchNumeric(val float64) bool {
	if q.gt != nil {
		if gt, ok := q.gt.(float64); ok && val <= gt {
			return false
		}
	}
	if q.gte != nil {
		if gte, ok := q.gte.(float64); ok && val < gte {
			return false
		}
	}
	if q.lt != nil {
		if lt, ok := q.lt.(float64); ok && val >= lt {
			return false
		}
	}
	if q.lte != nil {
		if lte, ok := q.lte.(float64); ok && val > lte {
			return false
		}
	}
	return true
}

func (q *RangeQueryImpl) matchTime(val time.Time) bool {
	if q.gt != nil {
		if gt, ok := q.gt.(time.Time); ok && val.Before(gt) || val.Equal(gt) {
			return false
		}
	}
	if q.gte != nil {
		if gte, ok := q.gte.(time.Time); ok && val.Before(gte) {
			return false
		}
	}
	if q.lt != nil {
		if lt, ok := q.lt.(time.Time); ok && val.After(lt) || val.Equal(lt) {
			return false
		}
	}
	if q.lte != nil {
		if lte, ok := q.lte.(time.Time); ok && val.After(lte) {
			return false
		}
	}
	return true
}

func (q *RangeQueryImpl) Match(value interface{}) bool {
	// Handle direct value comparison first
	switch v := value.(type) {
	case float64:
		return q.matchNumeric(v)
	case time.Time:
		return q.matchTime(v)
	}

	// Handle document case
	doc, ok := value.(*document.Document)
	if !ok {
		return false
	}

	field, err := doc.GetField(q.field)
	if err != nil {
		return false
	}

	// Handle field value comparison
	switch v := field.Value.(type) {
	case float64:
		return q.matchNumeric(v)
	case time.Time:
		return q.matchTime(v)
	default:
		return false
	}
}

// BooleanQueryImpl represents a boolean combination of queries
type BooleanQueryImpl struct {
	field    string
	must     []Query
	should   []Query
	mustNot  []Query
	minMatch int
}

func NewBooleanQuery() *BooleanQueryImpl {
	return &BooleanQueryImpl{minMatch: 1}
}

func (q *BooleanQueryImpl) Type() QueryType { return BooleanQuery }
func (q *BooleanQueryImpl) Field() string   { return q.field }

func (q *BooleanQueryImpl) Must() []Query   { return q.must }
func (q *BooleanQueryImpl) Should() []Query { return q.should }

func (q *BooleanQueryImpl) AddMust(query Query)    { q.must = append(q.must, query) }
func (q *BooleanQueryImpl) AddShould(query Query)  { q.should = append(q.should, query) }
func (q *BooleanQueryImpl) AddMustNot(query Query) { q.mustNot = append(q.mustNot, query) }

func (q *BooleanQueryImpl) Match(value interface{}) bool {
	// Handle map values for field-specific queries
	valueMap, ok := value.(map[string]string)
	if !ok {
		// If not a map, treat as a single value
		// Must match all MUST queries
		for _, must := range q.must {
			if !must.Match(value) {
				return false
			}
		}

		// Must not match any MUST NOT queries
		for _, mustNot := range q.mustNot {
			if mustNot.Match(value) {
				return false
			}
		}

		// Must match at least minMatch of SHOULD queries if any exist
		if len(q.should) > 0 {
			matches := 0
			for _, should := range q.should {
				if should.Match(value) {
					matches++
					if matches >= q.minMatch {
						return true
					}
				}
			}
			return false
		}

		return true
	}

	// Handle map values for field-specific queries
	// Must match all MUST queries
	for _, must := range q.must {
		fieldValue, exists := valueMap[must.Field()]
		if !exists || !must.Match(fieldValue) {
			return false
		}
	}

	// Must not match any MUST NOT queries
	for _, mustNot := range q.mustNot {
		fieldValue, exists := valueMap[mustNot.Field()]
		if exists && mustNot.Match(fieldValue) {
			return false
		}
	}

	// Must match at least minMatch of SHOULD queries if any exist
	if len(q.should) > 0 {
		matches := 0
		for _, should := range q.should {
			fieldValue, exists := valueMap[should.Field()]
			if exists && should.Match(fieldValue) {
				matches++
				if matches >= q.minMatch {
					return true
				}
			}
		}
		return false
	}

	return true
}

// MatchQueryImpl represents a match query that matches analyzed text
type MatchQueryImpl struct {
	field string
	text  string
}

func NewMatchQuery(field, text string) *MatchQueryImpl {
	return &MatchQueryImpl{field: field, text: text}
}

func (q *MatchQueryImpl) Type() QueryType { return MatchQuery }
func (q *MatchQueryImpl) Field() string   { return q.field }
func (q *MatchQueryImpl) Text() string    { return q.text }
func (q *MatchQueryImpl) Match(value interface{}) bool {
	if str, ok := value.(string); ok {
		// For now, we'll do a simple case-insensitive contains check
		// In a real implementation, this would use the analyzer
		return strings.Contains(strings.ToLower(str), strings.ToLower(q.text))
	}
	return false
}

// MatchPhraseQueryImpl represents a match_phrase query that matches exact phrases
type MatchPhraseQueryImpl struct {
	field  string
	phrase string
}

func NewMatchPhraseQuery(field, phrase string) *MatchPhraseQueryImpl {
	return &MatchPhraseQueryImpl{field: field, phrase: phrase}
}

func (q *MatchPhraseQueryImpl) Type() QueryType { return MatchPhraseQuery }
func (q *MatchPhraseQueryImpl) Field() string   { return q.field }
func (q *MatchPhraseQueryImpl) Match(value interface{}) bool {
	if str, ok := value.(string); ok {
		// For now, we'll do a simple case-insensitive exact match
		// In a real implementation, this would use the analyzer
		return strings.EqualFold(str, q.phrase)
	}
	return false
}

// MatchAllQueryImpl represents a match_all query that matches all documents
type MatchAllQueryImpl struct{}

func NewMatchAllQuery() *MatchAllQueryImpl {
	return &MatchAllQueryImpl{}
}

func (q *MatchAllQueryImpl) Type() QueryType { return MatchAllQuery }
func (q *MatchAllQueryImpl) Field() string   { return "" }
func (q *MatchAllQueryImpl) Match(value interface{}) bool {
	return true
}

// QueryMapper maps ElasticSearch DSL queries to internal query representations
type QueryMapper struct{}

func NewQueryMapper() *QueryMapper {
	return &QueryMapper{}
}

// MapQuery maps an ElasticSearch DSL query to our internal query representation
func (m *QueryMapper) MapQuery(dslQuery map[string]interface{}) (Query, error) {
	if len(dslQuery) != 1 {
		return nil, fmt.Errorf("invalid query structure: expected exactly one root query type")
	}

	for queryType, queryBody := range dslQuery {
		switch queryType {
		case "term":
			return m.mapTermQuery(queryBody)
		case "match":
			return m.mapMatchQuery(queryBody)
		case "match_phrase":
			return m.mapMatchPhraseQuery(queryBody)
		case "match_all":
			return NewMatchAllQuery(), nil
		case "range":
			return m.mapRangeQuery(queryBody)
		case "bool":
			return m.mapBoolQuery(queryBody)
		default:
			return nil, fmt.Errorf("unsupported query type: %s", queryType)
		}
	}

	return nil, fmt.Errorf("invalid query structure")
}

func (m *QueryMapper) mapTermQuery(body interface{}) (Query, error) {
	termBody, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid term query structure")
	}

	if len(termBody) != 1 {
		return nil, fmt.Errorf("term query must specify exactly one field")
	}

	for field, value := range termBody {
		switch v := value.(type) {
		case string:
			return NewTermQuery(field, v), nil
		case map[string]interface{}:
			if termValue, ok := v["value"].(string); ok {
				return NewTermQuery(field, termValue), nil
			}
			if termValue, ok := v["term"].(string); ok {
				return NewTermQuery(field, termValue), nil
			}
		}
		return nil, fmt.Errorf("term query value must be a string or {value: string}")
	}

	return nil, fmt.Errorf("invalid term query structure")
}

func (m *QueryMapper) mapRangeQuery(body interface{}) (Query, error) {
	rangeBody, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid range query structure")
	}

	if len(rangeBody) != 1 {
		return nil, fmt.Errorf("range query must specify exactly one field")
	}

	for field, conditions := range rangeBody {
		condMap, ok := conditions.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid range conditions structure")
		}

		query := NewRangeQuery(field)
		for op, val := range condMap {
			switch op {
			case "gt":
				query.GreaterThan(val)
			case "gte":
				query.GreaterThanOrEqual(val)
			case "lt":
				query.LessThan(val)
			case "lte":
				query.LessThanOrEqual(val)
			default:
				return nil, fmt.Errorf("unsupported range operator: %s", op)
			}
		}
		return query, nil
	}

	return nil, fmt.Errorf("invalid range query structure")
}

func (m *QueryMapper) mapBoolQuery(body interface{}) (Query, error) {
	boolBody, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid bool query structure")
	}

	query := NewBooleanQuery()

	for clause, queries := range boolBody {
		queryList, ok := queries.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid bool clause structure for %s", clause)
		}

		for _, q := range queryList {
			queryMap, ok := q.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid query in bool clause")
			}

			subQuery, err := m.MapQuery(queryMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse query in bool clause: %v", err)
			}

			switch clause {
			case "must":
				query.AddMust(subQuery)
			case "should":
				query.AddShould(subQuery)
			case "must_not":
				query.AddMustNot(subQuery)
			default:
				return nil, fmt.Errorf("unsupported bool clause: %s", clause)
			}
		}
	}

	return query, nil
}

func (m *QueryMapper) mapMatchQuery(body interface{}) (Query, error) {
	matchBody, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid match query structure")
	}

	if len(matchBody) != 1 {
		return nil, fmt.Errorf("match query must specify exactly one field")
	}

	for field, value := range matchBody {
		switch v := value.(type) {
		case string:
			return NewMatchQuery(field, v), nil
		case map[string]interface{}:
			if query, ok := v["query"].(string); ok {
				return NewMatchQuery(field, query), nil
			}
			if query, ok := v["value"].(string); ok {
				return NewMatchQuery(field, query), nil
			}
		}
		return nil, fmt.Errorf("match query value must be a string or {query: string}")
	}

	return nil, fmt.Errorf("invalid match query structure")
}

func (m *QueryMapper) mapMatchPhraseQuery(body interface{}) (Query, error) {
	phraseBody, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid match_phrase query structure")
	}

	if len(phraseBody) != 1 {
		return nil, fmt.Errorf("match_phrase query must specify exactly one field")
	}

	for field, value := range phraseBody {
		switch v := value.(type) {
		case string:
			return NewMatchPhraseQuery(field, v), nil
		case map[string]interface{}:
			if query, ok := v["query"].(string); ok {
				return NewMatchPhraseQuery(field, query), nil
			}
		}
		return nil, fmt.Errorf("invalid match_phrase query value")
	}

	return nil, fmt.Errorf("invalid match_phrase query structure")
}
