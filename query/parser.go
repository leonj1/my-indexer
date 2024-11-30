package query

import (
	"fmt"
	"strings"
)

// Use the QueryType from query.go

// ParsedQuery represents a parsed search query
type ParsedQuery struct {
	Type       QueryType
	Field      string
	Terms      []string
	IsPhrase   bool
	SubQueries []ParsedQuery
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

// Parse parses a query string into a ParsedQuery object
func (p *Parser) Parse(queryStr string) (*ParsedQuery, error) {
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
			return &ParsedQuery{
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
		
		return &ParsedQuery{
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
		return &ParsedQuery{
			Type:     PhraseQuery,
			Field:    p.defaultField,
			Terms:    terms,
			IsPhrase: true,
		}, nil
	}

	// Handle AND/OR queries
	if strings.Contains(queryStr, " AND ") {
		parts := strings.Split(queryStr, " AND ")
		subQueries := make([]ParsedQuery, 0, len(parts))
		for _, part := range parts {
			subQuery, err := p.Parse(part)
			if err != nil {
				return nil, err
			}
			subQueries = append(subQueries, *subQuery)
		}
		return &ParsedQuery{
			Type:       TermQuery,
			SubQueries: subQueries,
			Operator:   "AND",
		}, nil
	}

	if strings.Contains(queryStr, " OR ") {
		parts := strings.Split(queryStr, " OR ")
		subQueries := make([]ParsedQuery, 0, len(parts))
		for _, part := range parts {
			subQuery, err := p.Parse(part)
			if err != nil {
				return nil, err
			}
			subQueries = append(subQueries, *subQuery)
		}
		return &ParsedQuery{
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

	return &ParsedQuery{
		Type:  TermQuery,
		Field: p.defaultField,
		Terms: terms,
	}, nil
}

// Validate checks if a query is valid
func (p *Parser) Validate(query *ParsedQuery) error {
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
