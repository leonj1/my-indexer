package search

import (
	"strings"
	"testing"

	"my-indexer/analysis"
	"my-indexer/document"
	"my-indexer/index"
	"my-indexer/query"
)

// mockAnalyzer implements a simple analyzer for testing
type mockAnalyzer struct{}

func (a *mockAnalyzer) Analyze(text string) []analysis.Token {
	// Split text into tokens by whitespace
	words := strings.Fields(strings.ToLower(text))
	tokens := make([]analysis.Token, len(words))
	for i, word := range words {
		tokens[i] = analysis.Token{
			Text:      word,
			Position:  i,
			StartByte: 0, // Simplified for testing
			EndByte:   0, // Simplified for testing
		}
	}
	return tokens
}

// MockDocumentStore implements DocumentStore for testing
type MockDocumentStore struct {
	docs map[int]*document.Document
}

func newMockDocumentStore() *MockDocumentStore {
	return &MockDocumentStore{
		docs: make(map[int]*document.Document),
	}
}

func (s *MockDocumentStore) LoadDocument(docID int) (*document.Document, error) {
	return s.docs[docID], nil
}

func TestQueryExecutor(t *testing.T) {
	// Setup test environment
	analyzer := &mockAnalyzer{}
	idx := index.NewIndex(analyzer)
	store := newMockDocumentStore()
	search := NewSearch(idx, store)
	executor := NewQueryExecutor(search)

	// Add test documents
	doc1 := document.NewDocument()
	doc1.AddField("title", "The quick brown fox")
	doc1.AddField("content", "jumps over the lazy dog")
	doc1.AddField("age", 5)
	store.docs[0] = doc1

	doc2 := document.NewDocument()
	doc2.AddField("title", "The lazy brown dog")
	doc2.AddField("content", "sleeps in the sun")
	doc2.AddField("age", 3)
	store.docs[1] = doc2

	doc3 := document.NewDocument()
	doc3.AddField("title", "A brown fox")
	doc3.AddField("content", "jump over fences")
	doc3.AddField("age", 7)
	store.docs[2] = doc3

	// Add documents to index
	idx.AddDocument(doc1)
	idx.AddDocument(doc2)
	idx.AddDocument(doc3)

	// Test cases
	t.Run("Term Query", func(t *testing.T) {
		q := query.NewTermQuery("title", "quick")
		results, err := executor.Execute(q)
		if err != nil {
			t.Errorf("Failed to execute term query: %v", err)
		}
		if len(results.hits) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results.hits))
		}
		if len(results.hits) > 0 && results.hits[0].DocID != 0 {
			t.Errorf("Expected document 0 to match, got document %d", results.hits[0].DocID)
		}
	})

	t.Run("Range Query", func(t *testing.T) {
		q := query.NewRangeQuery("age")
		q.GreaterThan(2.0)
		q.LessThan(4.0)
		results, err := executor.Execute(q)
		if err != nil {
			t.Errorf("Failed to execute range query: %v", err)
		}
		if len(results.hits) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results.hits))
		}
	})

	t.Run("Boolean Query", func(t *testing.T) {
		bq := query.NewBooleanQuery()
		bq.AddMust(query.NewTermQuery("title", "quick"))
		bq.AddMust(query.NewTermQuery("content", "jumps"))
		
		results, err := executor.Execute(bq)
		if err != nil {
			t.Errorf("Failed to execute boolean query: %v", err)
		}
		if len(results.hits) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results.hits))
		}
	})
}
