package query

import (
	"fmt"
	"strings"
)

// QueryType represents the type of query
type QueryType int

const (
	TermQuery QueryType = iota
	PhraseQuery
	FieldQuery
)

// Query represents a parsed search query
type Query struct {
	Type      QueryType
	Field     string
	Terms     []string
	IsPhrase  bool
	SubQueries []Query
	Operator   string // "AND" or "OR"
}

// Parser handles query parsing
type Parser struct {
	defaultField string
}

// NewParser creates a new query parser
func NewParser(defaultField string) *Parser {
	return &Parser{
		defaultField: defaultField,
	}
}

// Parse parses a query string into a Query object
func (p *Parser) Parse(queryStr string) (*Query, error) {
	queryStr = strings.TrimSpace(queryStr)
	if queryStr == "" {
		return nil, fmt.Errorf("empty query")
	}

	// Handle field-specific queries (field:value)
	if strings.Contains(queryStr, ":") {
		parts := strings.SplitN(queryStr, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field query syntax")
		}
		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		if value == "" {
			return nil, fmt.Errorf("empty field value")
		}
		
		// Handle phrase queries
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
			terms := strings.Fields(value)
			if len(terms) < 2 {
				return nil, fmt.Errorf("phrase query must contain at least two terms")
			}
			return &Query{
				Type:     PhraseQuery,
				Field:    field,
				Terms:    terms,
				IsPhrase: true,
			}, nil
		}
		
		terms := strings.Fields(value)
		if len(terms) == 0 {
			return nil, fmt.Errorf("empty field value")
		}
		
		return &Query{
			Type:  FieldQuery,
			Field: field,
			Terms: terms,
		}, nil
	}

	// Handle phrase queries
	if strings.HasPrefix(queryStr, "\"") && strings.HasSuffix(queryStr, "\"") {
		queryStr = strings.Trim(queryStr, "\"")
		terms := strings.Fields(queryStr)
		if len(terms) < 2 {
			return nil, fmt.Errorf("phrase query must contain at least two terms")
		}
		return &Query{
			Type:     PhraseQuery,
			Field:    p.defaultField,
			Terms:    terms,
			IsPhrase: true,
		}, nil
	}

	// Handle AND/OR queries
	if strings.Contains(queryStr, " AND ") {
		parts := strings.Split(queryStr, " AND ")
		subQueries := make([]Query, 0, len(parts))
		for _, part := range parts {
			subQuery, err := p.Parse(part)
			if err != nil {
				return nil, err
			}
			subQueries = append(subQueries, *subQuery)
		}
		return &Query{
			Type:       TermQuery,
			SubQueries: subQueries,
			Operator:   "AND",
		}, nil
	}

	if strings.Contains(queryStr, " OR ") {
		parts := strings.Split(queryStr, " OR ")
		subQueries := make([]Query, 0, len(parts))
		for _, part := range parts {
			subQuery, err := p.Parse(part)
			if err != nil {
				return nil, err
			}
			subQueries = append(subQueries, *subQuery)
		}
		return &Query{
			Type:       TermQuery,
			SubQueries: subQueries,
			Operator:   "OR",
		}, nil
	}

	// Simple term query
	terms := strings.Fields(queryStr)
	if len(terms) == 0 {
		return nil, fmt.Errorf("empty query")
	}

	return &Query{
		Type:  TermQuery,
		Field: p.defaultField,
		Terms: terms,
	}, nil
}

// Validate checks if a query is valid
func (p *Parser) Validate(query *Query) error {
	if query == nil {
		return fmt.Errorf("nil query")
	}

	if len(query.Terms) == 0 && len(query.SubQueries) == 0 {
		return fmt.Errorf("query must contain at least one term or subquery")
	}

	if query.IsPhrase && len(query.Terms) < 2 {
		return fmt.Errorf("phrase query must contain at least two terms")
	}

	// Validate subqueries recursively
	for _, subQuery := range query.SubQueries {
		if err := p.Validate(&subQuery); err != nil {
			return err
		}
	}

	return nil
}
