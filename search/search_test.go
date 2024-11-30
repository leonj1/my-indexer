package search

import (
	"sort"
	"testing"

	"my-indexer/analysis"
	"my-indexer/document"
	"my-indexer/index"
	"fmt"
)

// mockDocumentStore implements DocumentStore for testing
type mockDocumentStore struct {
	docs map[int]*document.Document
}

func (m *mockDocumentStore) LoadDocument(docID int) (*document.Document, error) {
	if doc, ok := m.docs[docID]; ok {
		return doc, nil
	}
	return nil, fmt.Errorf("document not found: %d", docID)
}

func (m *mockDocumentStore) LoadAllDocuments() ([]*document.Document, error) {
	docs := make([]*document.Document, 0, len(m.docs))
	for _, doc := range m.docs {
		docs = append(docs, doc)
	}
	return docs, nil
}

func newMockStore() *mockDocumentStore {
	return &mockDocumentStore{
		docs: make(map[int]*document.Document),
	}
}

func TestSearch(t *testing.T) {
	// Create analyzer
	analyzer := analysis.NewStandardAnalyzer()

	// Create index
	idx := index.NewIndex(analyzer)

	// Create document store
	store := newMockStore()

	// Create search
	search := NewSearch(idx, store)

	// Create test documents
	docs := []*document.Document{
		func() *document.Document {
			doc := document.NewDocument()
			doc.AddField("title", "The quick brown fox")
			doc.AddField("content", "The quick brown fox jumps over the lazy dog")
			return doc
		}(),
		func() *document.Document {
			doc := document.NewDocument()
			doc.AddField("title", "Lazy dog sleeping")
			doc.AddField("content", "The lazy dog is sleeping in the sun")
			return doc
		}(),
		func() *document.Document {
			doc := document.NewDocument()
			doc.AddField("title", "Quick rabbit running")
			doc.AddField("content", "A quick rabbit is running in the field")
			return doc
		}(),
	}

	// Add documents to index and store
	for i, doc := range docs {
		docID, err := idx.AddDocument(doc)
		if err != nil {
			t.Fatalf("Failed to add document %d: %v", i, err)
		}
		store.docs[docID] = doc
	}

	tests := []struct {
		name          string
		terms         []string
		op            Operator
		expectedDocs  int
		expectedTerms []string
	}{
		{
			name:          "Single term search",
			terms:         []string{"quick"},
			op:            OR,
			expectedDocs:  2,
			expectedTerms: []string{"quick"},
		},
		{
			name:          "AND search",
			terms:         []string{"quick", "fox"},
			op:            AND,
			expectedDocs:  1,
			expectedTerms: []string{"quick", "fox"},
		},
		{
			name:          "OR search",
			terms:         []string{"quick", "lazy"},
			op:            OR,
			expectedDocs:  3,
			expectedTerms: []string{"quick", "lazy"},
		},
		{
			name:          "Empty search",
			terms:         []string{},
			op:            OR,
			expectedDocs:  0,
			expectedTerms: []string{},
		},
		{
			name:          "No results",
			terms:         []string{"nonexistent"},
			op:            OR,
			expectedDocs:  0,
			expectedTerms: []string{"nonexistent"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := search.Search(tt.terms, tt.op)
			if err != nil {
				t.Fatalf("Search failed: %v", err)
			}

			if len(results.GetHits()) != tt.expectedDocs {
				t.Errorf("Expected %d docs, got %d", tt.expectedDocs, len(results.GetHits()))
			}

			// Verify results are sorted by score
			hits := results.GetHits()
			if len(hits) > 1 {
				if !sort.SliceIsSorted(hits, func(i, j int) bool {
					return hits[i].Score > hits[j].Score
				}) {
					t.Error("Results are not sorted by score")
				}
			}

			// Verify operator behavior
			if tt.op == AND && len(hits) > 0 {
				// For AND, all terms should be present in each result
				for _, hit := range hits {
					doc := hit.Doc
					for _, term := range tt.terms {
						found := false
						for _, field := range doc.GetFields() {
							tokens := analyzer.Analyze(field.Value.(string))
							for _, token := range tokens {
								if token.Text == term {
									found = true
									break
								}
							}
							if found {
								break
							}
						}
						if !found {
							t.Errorf("Term %q not found in document for AND search", term)
						}
					}
				}
			}
		})
	}
}

func TestSearchScoring(t *testing.T) {
	// Create analyzer
	analyzer := analysis.NewStandardAnalyzer()

	// Create index
	idx := index.NewIndex(analyzer)

	// Create document store
	store := newMockStore()

	// Create search
	search := NewSearch(idx, store)

	// Create test documents with varying term frequencies
	docs := []*document.Document{
		func() *document.Document {
			doc := document.NewDocument()
			doc.AddField("title", "test document")
			doc.AddField("content", "test test test") // 4 occurrences of "test"
			return doc
		}(),
		func() *document.Document {
			doc := document.NewDocument()
			doc.AddField("title", "test")
			doc.AddField("content", "test document") // 2 occurrences of "test"
			return doc
		}(),
		func() *document.Document {
			doc := document.NewDocument()
			doc.AddField("content", "test") // 1 occurrence of "test"
			return doc
		}(),
	}

	// Add documents to index and store
	for i, doc := range docs {
		docID, err := idx.AddDocument(doc)
		if err != nil {
			t.Fatalf("Failed to add document %d: %v", i, err)
		}
		store.docs[docID] = doc
	}

	// Search for "test"
	results, err := search.Search([]string{"test"}, OR)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	hits := results.GetHits()
	if len(hits) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(hits))
	}

	// Verify scoring order (document with more "test" occurrences should score higher)
	frequencies := make([]float64, len(hits))
	for i, hit := range hits {
		frequencies[i] = hit.Score
	}

	// Check that scores are in descending order
	for i := 1; i < len(frequencies); i++ {
		if frequencies[i-1] <= frequencies[i] {
			t.Errorf("Results not properly scored by term frequency: %v", frequencies)
			break
		}
	}
}

func TestConcurrentSearch(t *testing.T) {
	// Create analyzer
	analyzer := analysis.NewStandardAnalyzer()

	// Create index
	idx := index.NewIndex(analyzer)

	// Create document store
	store := newMockStore()

	// Create search
	search := NewSearch(idx, store)

	// Add test documents
	doc := document.NewDocument()
	doc.AddField("content", "test document")
	docID, err := idx.AddDocument(doc)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}
	store.docs[docID] = doc

	// Run concurrent searches
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := search.Search([]string{"test"}, OR)
			if err != nil {
				t.Errorf("Concurrent search failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all searches to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
