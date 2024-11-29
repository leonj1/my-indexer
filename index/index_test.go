package index

import (
	"fmt"
	"sync"
	"testing"
	"time"

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
	startTime := time.Now()
	t.Log("Initializing test...")
	idx := NewIndex(analysis.NewStandardAnalyzer())
	var wg sync.WaitGroup
	var errors []string
	var errorMu sync.Mutex

	// Add initial document
	t.Log("Creating initial document...")
	doc := document.NewDocument()
	err := doc.AddField("content", "test document")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}
	
	t.Log("Attempting to add document to index...")
	var docID int
	done := make(chan struct{})
	go func() {
		defer close(done)
		var err error
		t.Log("Starting AddDocument operation...")
		docID, err = idx.AddDocument(doc)
		if err != nil {
			t.Errorf("Failed to add document: %v", err)
			return
		}
		t.Log("AddDocument operation completed successfully")
	}()

	select {
	case <-done:
		t.Logf("Successfully added document with ID: %d", docID)
	case <-time.After(5 * time.Second):
		t.Fatal("AddDocument operation timed out after 5 seconds")
	}

	t.Log("Starting concurrent operations...")

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(routineNum int) {
			defer wg.Done()
			
			startGet := time.Now()
			t.Logf("Routine %d: Starting GetDocument at %v", routineNum, startGet)
			_, err := idx.GetDocument(docID)
			if err != nil {
				errorMu.Lock()
				errors = append(errors, fmt.Sprintf("GetDocument failed in routine %d: %v", routineNum, err))
				errorMu.Unlock()
				return
			}
			t.Logf("Routine %d: Completed GetDocument in %v", routineNum, time.Since(startGet))
			
			startFreq := time.Now()
			t.Logf("Routine %d: Starting GetTermFrequency at %v", routineNum, startFreq)
			_, err = idx.GetTermFrequency("test", docID)
			if err != nil {
				errorMu.Lock()
				errors = append(errors, fmt.Sprintf("GetTermFrequency failed in routine %d: %v", routineNum, err))
				errorMu.Unlock()
				return
			}
			t.Logf("Routine %d: Completed GetTermFrequency in %v", routineNum, time.Since(startFreq))
		}(i)
	}

	t.Log("All routines launched, waiting for completion...")

	// Wait with timeout
	waitDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		t.Logf("All routines completed successfully in %v", time.Since(startTime))
		if len(errors) > 0 {
			for _, err := range errors {
				t.Error(err)
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out - possible deadlock")
	}
}
