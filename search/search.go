package search

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"my-indexer/document"
	"my-indexer/index"
	"my-indexer/query"
)

// Operator represents a boolean operator
type Operator int

const (
	// AND requires all terms to be present
	AND Operator = iota
	// OR requires at least one term to be present
	OR
)

// Result represents a search result with its score
type Result struct {
	Index  string             `json:"_index"`
	Type   string             `json:"_type"`
	ID     string             `json:"_id"`
	DocID  int               `json:"doc_id"`
	Score  float64            `json:"_score"`
	Source *document.Document `json:"_source"`
	Doc    *document.Document `json:"doc"` // Alias for Source for backward compatibility
}

// Results represents a sorted list of search results
type Results struct {
	hits   []*Result
	maxDoc int
}

// Len returns the number of results
func (r *Results) Len() int { return len(r.hits) }

// Less compares results by score
func (r *Results) Less(i, j int) bool {
	// Sort by score in descending order
	return r.hits[i].Score > r.hits[j].Score
}

// Swap swaps two results
func (r *Results) Swap(i, j int) { r.hits[i], r.hits[j] = r.hits[j], r.hits[i] }

// GetHits returns the sorted results
func (r *Results) GetHits() []*Result {
	return r.hits
}

// Search performs a search operation on the index
type Search struct {
	idx    *index.Index
	mu     sync.RWMutex
	store  DocumentStore
	maxDoc int
}

// DocumentStore is an interface for loading documents
type DocumentStore interface {
	LoadDocument(docID int) (*document.Document, error)
	LoadAllDocuments() ([]*document.Document, error)
}

// NewSearch creates a new search instance
func NewSearch(idx *index.Index, store DocumentStore) *Search {
	return &Search{
		idx:   idx,
		store: store,
	}
}

// calculateScore calculates the score for a document based on term frequencies
func (s *Search) calculateScore(docID int, terms []string) float64 {
	var score float64

	// Calculate TF-IDF score for each term
	for _, term := range terms {
		tf, err := s.idx.GetTermFrequency(term, docID)
		if err != nil {
			continue
		}
		df, err := s.idx.GetDocumentFrequency(term)
		if err != nil {
			continue
		}
		if df > 0 {
			// TF-IDF scoring: tf * idf
			// idf = log(1 + N/df) where N is total number of documents
			// Adding 1 ensures IDF is always positive
			N := float64(s.idx.GetDocumentCount())
			idf := math.Log1p(N / float64(df))
			score += float64(tf) * idf
		}
	}

	return score
}

// Search performs a search with the given terms and operator
func (s *Search) Search(terms []string, op Operator) (*Results, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(terms) == 0 {
		return &Results{}, nil
	}

	// Get document IDs based on operator
	docIDs := make(map[int]bool)
	
	// Handle first term
	postings := s.idx.GetPostings(terms[0])
	for docID := range postings {
		docIDs[docID] = true
	}

	// Process remaining terms based on operator
	for _, term := range terms[1:] {
		postings := s.idx.GetPostings(term)
		
		switch op {
		case AND:
			// Remove documents that don't contain the term
			for docID := range docIDs {
				if _, exists := postings[docID]; !exists {
					delete(docIDs, docID)
				}
			}
		case OR:
			// Add documents that contain the term
			for docID := range postings {
				docIDs[docID] = true
			}
		}
	}

	// Calculate scores and create results
	results := &Results{
		hits: make([]*Result, 0, len(docIDs)),
	}

	for docID := range docIDs {
		score := s.calculateScore(docID, terms)
		doc, err := s.store.LoadDocument(docID)
		if err != nil {
			return nil, fmt.Errorf("failed to load document %d: %w", docID, err)
		}

		results.hits = append(results.hits, &Result{
			Index:  "",
			Type:   "",
			ID:     fmt.Sprintf("%d", docID),
			DocID:  docID,
			Score:  score,
			Source: doc,
			Doc:    doc,
		})
	}

	// Sort results by score
	sort.Sort(results)

	return results, nil
}

// SearchWithQuery performs a search using a Query object
func (s *Search) SearchWithQuery(query query.Query) (*Results, error) {
	results := &Results{
		hits: make([]*Result, 0),
	}

	// Get all documents from the store
	docs, err := s.store.LoadAllDocuments()
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %v", err)
	}

	// Filter documents based on the query
	for _, doc := range docs {
		// For each field in the document
		for field, value := range doc.GetFields() {
			// Check if the query matches this field
			if query.Field() == "" || query.Field() == field {
				if query.Match(value) {
					result := &Result{
						Index:  "",
						Type:   "",
						ID:     fmt.Sprintf("%d", doc.ID),
						DocID:  doc.ID,
						Score:  1.0,
						Source: doc,
						Doc:    doc,
					}
					results.hits = append(results.hits, result)
					break // Once we find a match in any field, we can stop checking other fields
				}
			}
		}
	}

	// Sort results by score
	sort.Sort(results)

	return results, nil
}
