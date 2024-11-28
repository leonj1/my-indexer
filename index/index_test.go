package index

import (
	"testing"

	"my-indexer/analysis"
	"my-indexer/document"
)

func TestIndexOperations(t *testing.T) {
	// Create analyzer and index
	analyzer := analysis.NewStandardAnalyzer()
	idx := NewIndex(analyzer)

	// Test document 1
	doc1 := document.NewDocument()
	doc1.AddField("title", "The quick brown fox")
	doc1.AddField("content", "jumps over the lazy dog")

	docID1, err := idx.AddDocument(doc1)
	if err != nil {
		t.Fatalf("Failed to add document 1: %v", err)
	}

	// Test document 2
	doc2 := document.NewDocument()
	doc2.AddField("title", "Quick brown foxes")
	doc2.AddField("content", "are quick and brown")

	docID2, err := idx.AddDocument(doc2)
	if err != nil {
		t.Fatalf("Failed to add document 2: %v", err)
	}

	// Test term frequency
	tests := []struct {
		term   string
		docID  int
		expect int
	}{
		{"quick", docID1, 1},  // appears once in doc1
		{"quick", docID2, 2},  // appears twice in doc2
		{"brown", docID1, 1},  // appears once in doc1
		{"brown", docID2, 2},  // appears twice in doc2
		{"fox", docID1, 1},    // appears once in doc1
		{"foxes", docID2, 1},  // appears once in doc2
		{"nonexistent", docID1, 0},
	}

	for _, tt := range tests {
		tf, err := idx.GetTermFrequency(tt.term, tt.docID)
		if err != nil {
			t.Errorf("GetTermFrequency(%q, %d) returned error: %v", tt.term, tt.docID, err)
			continue
		}
		if tf != tt.expect {
			t.Errorf("GetTermFrequency(%q, %d) = %d, want %d", tt.term, tt.docID, tf, tt.expect)
		}
	}

	// Test document frequency
	dfTests := []struct {
		term   string
		expect int
	}{
		{"quick", 2},  // appears in both docs
		{"brown", 2},  // appears in both docs
		{"fox", 1},    // appears only in doc1
		{"foxes", 1},  // appears only in doc2
		{"nonexistent", 0},
	}

	for _, tt := range dfTests {
		df, err := idx.GetDocumentFrequency(tt.term)
		if err != nil {
			t.Errorf("GetDocumentFrequency(%q) returned error: %v", tt.term, err)
			continue
		}
		if df != tt.expect {
			t.Errorf("GetDocumentFrequency(%q) = %d, want %d", tt.term, df, tt.expect)
		}
	}

	// Test document count
	if count := idx.GetDocumentCount(); count != 2 {
		t.Errorf("GetDocumentCount() = %d, want 2", count)
	}

	// Test nil document
	if _, err := idx.AddDocument(nil); err == nil {
		t.Error("AddDocument(nil) should return error")
	}
}

func TestConcurrentAccess(t *testing.T) {
	idx := NewIndex(analysis.NewStandardAnalyzer())
	done := make(chan bool)

	// Add initial document
	doc := document.NewDocument()
	err := doc.AddField("content", "test document")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}
	docID, err := idx.AddDocument(doc)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			_, err := idx.GetDocument(docID)
			if err != nil {
				t.Errorf("Concurrent GetDocument failed: %v", err)
			}
			_, err = idx.GetTermFrequency("test", docID)
			if err != nil {
				t.Errorf("Concurrent GetTermFrequency failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
