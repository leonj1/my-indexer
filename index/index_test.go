package index

import (
	"testing"

	"my-indexer/analysis"
	"my-indexer/document"
)

func TestIndexOperations(t *testing.T) {
	idx := NewIndex(analysis.NewStandardAnalyzer())

	// Test adding documents
	doc1 := document.NewDocument()
	err := doc1.AddField("title", "The quick brown fox")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}
	err = doc1.AddField("content", "jumps over the lazy dog")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}

	doc2 := document.NewDocument()
	err = doc2.AddField("title", "Quick brown foxes")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}
	err = doc2.AddField("content", "are quick and brown")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}

	// Add documents to index
	docID1, err := idx.AddDocument(doc1)
	if err != nil {
		t.Fatalf("Failed to add document 1: %v", err)
	}
	t.Logf("Document 1 (ID: %d) content:", docID1)
	for fieldName, field := range doc1.GetFields() {
		if str, ok := field.Value.(string); ok {
			t.Logf("  %s: %q", fieldName, str)
			tokens := idx.analyzer.Analyze(str)
			t.Logf("  Tokens:")
			for _, token := range tokens {
				t.Logf("    - %q (pos: %d)", token.Text, token.Position)
			}
		}
	}

	docID2, err := idx.AddDocument(doc2)
	if err != nil {
		t.Fatalf("Failed to add document 2: %v", err)
	}
	t.Logf("Document 2 (ID: %d) content:", docID2)
	for fieldName, field := range doc2.GetFields() {
		if str, ok := field.Value.(string); ok {
			t.Logf("  %s: %q", fieldName, str)
			tokens := idx.analyzer.Analyze(str)
			t.Logf("  Tokens:")
			for _, token := range tokens {
				t.Logf("    - %q (pos: %d)", token.Text, token.Position)
			}
		}
	}

	// Test document retrieval
	retrievedDoc1, err := idx.GetDocument(docID1)
	if err != nil {
		t.Errorf("Failed to retrieve document 1: %v", err)
	}
	if retrievedDoc1 != doc1 {
		t.Error("Retrieved document 1 does not match original")
	}

	// Test term frequency
	tests := []struct {
		term     string
		docID    int
		expected int
	}{
		{"quick", docID1, 1},
		{"quick", docID2, 2}, // appears twice in doc2
		{"fox", docID1, 1},
		{"foxes", docID2, 1},
		{"nonexistent", docID1, 0},
	}

	for _, tt := range tests {
		freq, err := idx.GetTermFrequency(tt.term, tt.docID)
		if err != nil {
			t.Errorf("GetTermFrequency(%q, %d) error: %v", tt.term, tt.docID, err)
			continue
		}
		if freq != tt.expected {
			t.Errorf("GetTermFrequency(%q, %d) = %d, want %d", tt.term, tt.docID, freq, tt.expected)
		}
	}

	// Test document frequency
	docFreqTests := []struct {
		term     string
		expected int
	}{
		{"quick", 2},  // appears in both docs
		{"fox", 1},    // appears in doc1
		{"foxes", 1},  // appears in doc2
		{"brown", 2},  // appears in both docs
		{"nonexistent", 0},
	}

	for _, tt := range docFreqTests {
		freq, err := idx.GetDocumentFrequency(tt.term)
		if err != nil {
			t.Errorf("GetDocumentFrequency(%q) error: %v", tt.term, err)
			continue
		}
		if freq != tt.expected {
			t.Logf("Term %q tokens:", tt.term)
			tokens := idx.analyzer.Analyze(tt.term)
			for _, token := range tokens {
				t.Logf("  - %q", token.Text)
			}
			if postingList, err := idx.GetPostingList(tt.term); err == nil && postingList != nil {
				t.Logf("Posting list for %q:", tt.term)
				for docID, posting := range postingList.Postings {
					t.Logf("  - Doc %d: freq=%d, positions=%v", docID, posting.TermFreq, posting.Positions)
				}
			}
			t.Errorf("GetDocumentFrequency(%q) = %d, want %d", tt.term, freq, tt.expected)
		}
	}

	// Test document count
	if count := idx.GetDocumentCount(); count != 2 {
		t.Errorf("GetDocumentCount() = %d, want 2", count)
	}

	// Test posting list retrieval
	postingList, err := idx.GetPostingList("quick")
	if err != nil {
		t.Fatalf("GetPostingList('quick') error: %v", err)
	}
	if postingList.DocFreq != 2 {
		t.Errorf("PostingList.DocFreq = %d, want 2", postingList.DocFreq)
	}
	if len(postingList.Postings) != 2 {
		t.Errorf("len(PostingList.Postings) = %d, want 2", len(postingList.Postings))
	}

	// Test error cases
	if _, err := idx.AddDocument(nil); err == nil {
		t.Error("AddDocument(nil) should return error")
	}
	if _, err := idx.GetDocument(-1); err == nil {
		t.Error("GetDocument(-1) should return error")
	}
	if _, err := idx.GetPostingList(""); err == nil {
		t.Error("GetPostingList('') should return error")
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
