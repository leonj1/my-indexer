package search

import (
	"fmt"
	"math"
	"my-indexer/query"
	"sort"
)

// QueryExecutor executes internal queries and returns search results
type QueryExecutor struct {
	search *Search
}

// NewQueryExecutor creates a new query executor
func NewQueryExecutor(search *Search) *QueryExecutor {
	return &QueryExecutor{
		search: search,
	}
}

// Execute executes an internal query and returns search results
func (e *QueryExecutor) Execute(q query.Query) (*Results, error) {
	e.search.mu.RLock()
	defer e.search.mu.RUnlock()

	// Handle different query types
	switch q.Type() {
	case query.TermQuery:
		return e.executeTermQuery(q)
	case query.PhraseQuery:
		return e.executePhraseQuery(q)
	case query.RangeQuery:
		return e.executeRangeQuery(q)
	case query.BooleanQuery:
		return e.executeBooleanQuery(q)
	case query.MatchQuery:
		return e.executeMatchQuery(q)
	default:
		return nil, fmt.Errorf("unsupported query type: %v", q.Type())
	}
}

// executeTermQuery executes a term query
func (e *QueryExecutor) executeTermQuery(q query.Query) (*Results, error) {
	tq, ok := q.(*query.TermQueryImpl)
	if !ok {
		return nil, fmt.Errorf("invalid term query type")
	}

	// Get the analyzer from the search instance
	tokens := e.search.idx.Analyzer().Analyze(tq.Term())
	if len(tokens) == 0 {
		return &Results{hits: make([]*Result, 0)}, nil
	}
	
	// Use the first token's text as our search term
	term := tokens[0].Text
	
	// Get posting list for the term
	postings := e.search.idx.GetPostings(term)
	
	// Create results
	results := &Results{
		hits: make([]*Result, 0, len(postings)),
	}

	// Process each document
	for docID, posting := range postings {
		// Check if the term appears in the specified field
		fieldFound := false
		for _, field := range posting.Fields {
			if field == tq.Field() {
				fieldFound = true
				break
			}
		}
		if !fieldFound {
			continue
		}

		// Load document
		doc, err := e.search.store.LoadDocument(docID)
		if err != nil {
			return nil, fmt.Errorf("failed to load document %d: %w", docID, err)
		}

		// Calculate score using TF-IDF
		score := e.calculateScore(docID, []string{term})

		results.hits = append(results.hits, &Result{
			ID:     fmt.Sprintf("%d", docID),
			Score:  score,
			Source: doc,
		})
	}

	// Sort results by score
	sort.Sort(results)

	return results, nil
}

// executePhraseQuery executes a phrase query
func (e *QueryExecutor) executePhraseQuery(q query.Query) (*Results, error) {
	// For now, treat phrase queries as term queries
	// TODO: Implement proper phrase query execution using term positions
	return e.executeTermQuery(q)
}

// executeRangeQuery executes a range query
func (e *QueryExecutor) executeRangeQuery(q query.Query) (*Results, error) {
	// Get all documents and filter by range
	results := &Results{
		hits: make([]*Result, 0),
	}

	// Scan all documents (inefficient, but works for now)
	// TODO: Implement field indexing for efficient range queries
	for docID := 0; docID < e.search.idx.GetDocumentCount(); docID++ {
		doc, err := e.search.store.LoadDocument(docID)
		if err != nil {
			continue
		}

		// Check if document matches range criteria
		field, err := doc.GetField(q.Field())
		if err != nil {
			continue
		}

		// Convert field value to float64 for comparison
		var fieldValue float64
		switch v := field.Value.(type) {
		case int:
			fieldValue = float64(v)
		case int32:
			fieldValue = float64(v)
		case int64:
			if v > 9223372036854775807 || v < -9223372036854775808 {
				return nil, fmt.Errorf("int64 value %d exceeds float64 range", v)
			}
			fieldValue = float64(v)
		case float32:
			fieldValue = float64(v)
		case float64:
			fieldValue = v
		case uint:
			fieldValue = float64(v)
		case uint32:
			fieldValue = float64(v)
		case uint64:
			if v > 18446744073709551615 {
				return nil, fmt.Errorf("uint64 value %d exceeds float64 range", v)
			}
			fieldValue = float64(v)
		default:
			// Skip non-numeric fields
			continue
		}

		rq := q.(*query.RangeQueryImpl)
		if rq.Gt() != nil {
			if gt, ok := rq.Gt().(float64); ok {
				if fieldValue <= gt {
					continue
				}
			} else {
				return nil, fmt.Errorf("gt value is not a float64")
			}
		}
		if rq.Lt() != nil {
			if lt, ok := rq.Lt().(float64); ok {
				if fieldValue >= lt {
					continue
				}
			} else {
				return nil, fmt.Errorf("lt value is not a float64")
			}
		}

		results.hits = append(results.hits, &Result{
			ID:     fmt.Sprintf("%d", docID),
			Score:  1.0, // Default score for range queries
			Source: doc,
		})
	}

	return results, nil
}

// executeBooleanQuery executes a boolean query
func (e *QueryExecutor) executeBooleanQuery(q query.Query) (*Results, error) {
	bq, ok := q.(*query.BooleanQueryImpl)
	if !ok {
		return nil, fmt.Errorf("invalid boolean query type")
	}

	// Execute must queries
	var mustResults *Results
	if len(bq.Must()) > 0 {
		var err error
		mustResults, err = e.executeMustClauses(bq.Must())
		if err != nil {
			return nil, err
		}
	}

	// Execute should queries
	var shouldResults *Results
	if len(bq.Should()) > 0 {
		var err error
		shouldResults, err = e.executeShouldClauses(bq.Should())
		if err != nil {
			return nil, err
		}
	}

	// If both must and should clauses are empty, return empty results
	if mustResults == nil && shouldResults == nil {
		return &Results{hits: make([]*Result, 0)}, nil
	}

	// Combine results
	return e.combineResults(mustResults, shouldResults), nil
}

// executeMatchQuery executes a match query
func (e *QueryExecutor) executeMatchQuery(q query.Query) (*Results, error) {
	mq, ok := q.(*query.MatchQueryImpl)
	if !ok {
		return nil, fmt.Errorf("invalid match query type")
	}

	// Get the analyzer from the search instance
	tokens := e.search.idx.Analyzer().Analyze(mq.Text())
	if len(tokens) == 0 {
		return &Results{hits: make([]*Result, 0)}, nil
	}

	// Process each token
	results := &Results{
		hits: make([]*Result, 0),
	}

	// Track seen documents to avoid duplicates
	seenDocs := make(map[int]bool)

	for _, token := range tokens {
		// Get posting list for the term
		postings := e.search.idx.GetPostings(token.Text)

		// Process each document
		for docID, posting := range postings {
			// Skip if we've already processed this document
			if seenDocs[docID] {
				continue
			}

			// Check if the term appears in the specified field
			fieldFound := false
			for _, field := range posting.Fields {
				if field == mq.Field() {
					fieldFound = true
					break
				}
			}
			if !fieldFound {
				continue
			}

			// Load document
			doc, err := e.search.store.LoadDocument(docID)
			if err != nil {
				return nil, fmt.Errorf("failed to load document %d: %w", docID, err)
			}

			// Calculate score using TF-IDF
			score := e.calculateScore(docID, []string{token.Text})

			results.hits = append(results.hits, &Result{
				ID:     fmt.Sprintf("%d", docID),
				Score:  score,
				Source: doc,
			})

			// Mark document as seen
			seenDocs[docID] = true
		}
	}

	// Sort results by score
	sort.Sort(results)

	return results, nil
}

// executeMustClauses executes must clauses of a boolean query
func (e *QueryExecutor) executeMustClauses(queries []query.Query) (*Results, error) {
	if len(queries) == 0 {
		return nil, nil
	}

	// Execute first query
	results, err := e.Execute(queries[0])
	if err != nil {
		return nil, err
	}

	// Filter results through remaining queries
	for _, q := range queries[1:] {
		nextResults, err := e.Execute(q)
		if err != nil {
			return nil, err
		}

		// Keep only documents that appear in both result sets
		filteredHits := make([]*Result, 0)
		docMap := make(map[string]*Result)
		for _, hit := range nextResults.hits {
			docMap[hit.ID] = hit
		}

		for _, hit := range results.hits {
			if _, exists := docMap[hit.ID]; exists {
				filteredHits = append(filteredHits, hit)
			}
		}

		results.hits = filteredHits
	}

	return results, nil
}

// executeShouldClauses executes should clauses of a boolean query
func (e *QueryExecutor) executeShouldClauses(queries []query.Query) (*Results, error) {
	if len(queries) == 0 {
		return nil, nil
	}

	// Create a map to track unique documents and their highest scores
	docMap := make(map[string]*Result)

	// Execute each query and merge results
	for _, q := range queries {
		results, err := e.Execute(q)
		if err != nil {
			return nil, err
		}

		for _, hit := range results.hits {
			if existing, exists := docMap[hit.ID]; exists {
				// Keep the higher score
				if hit.Score > existing.Score {
					docMap[hit.ID] = hit
				}
			} else {
				docMap[hit.ID] = hit
			}
		}
	}

	// Convert map to results
	results := &Results{
		hits: make([]*Result, 0, len(docMap)),
	}
	for _, hit := range docMap {
		results.hits = append(results.hits, hit)
	}

	// Sort by score
	sort.Sort(results)

	return results, nil
}

// combineResults combines must and should results
func (e *QueryExecutor) combineResults(must, should *Results) *Results {
	if must == nil && should == nil {
		return &Results{hits: make([]*Result, 0)}
	}

	if must == nil {
		return should
	}

	if should == nil {
		return must
	}

	// Combine scores from must and should clauses
	docMap := make(map[string]*Result)
	
	// Add all must results
	for _, hit := range must.hits {
		docMap[hit.ID] = hit
	}

	// Add scores from should results
	for _, hit := range should.hits {
		if existing, exists := docMap[hit.ID]; exists {
			// Combine scores
			existing.Score += hit.Score
		}
	}

	// Convert map back to results
	results := &Results{
		hits: make([]*Result, 0, len(docMap)),
	}
	for _, hit := range docMap {
		results.hits = append(results.hits, hit)
	}

	// Sort by combined score
	sort.Sort(results)

	return results
}

// calculateScore calculates TF-IDF score for a document
func (e *QueryExecutor) calculateScore(docID int, terms []string) float64 {
	var score float64

	// Calculate TF-IDF score for each term
	for _, term := range terms {
		postings := e.search.idx.GetPostings(term)
		if entry, exists := postings[docID]; exists {
			tf := float64(entry.TermFreq)  // Using TermFreq field from PostingEntry
			df := float64(len(postings))
			if df > 0 {
				// TF-IDF scoring: tf * idf
				// idf = log(1 + N/df) where N is total number of documents
				// Adding 1 ensures IDF is always positive
				N := float64(e.search.idx.GetDocumentCount())
				idf := math.Log1p(N / df)
				score += tf * idf
			}
		}
	}

	return score
}
