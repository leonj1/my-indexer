package index

import (
	"my-indexer/document"
	"sync"
	"testing"
)

// func TestDocumentUpdate(t *testing.T) {
// 	idx := NewIndex(nil)

// 	// Add initial document
// 	doc1 := document.NewDocument()
// 	doc1.AddField("title", "initial title")
// 	doc1.AddField("content", "initial content")
// 	docID, err := idx.AddDocument(doc1)
// 	if err != nil {
// 		t.Fatalf("Failed to add document: %v", err)
// 	}

// 	// Update document
// 	doc2 := document.NewDocument()
// 	doc2.AddField("title", "updated title")
// 	doc2.AddField("content", "updated content")
// 	err = idx.UpdateDocument(docID, doc2)
// 	if err != nil {
// 		t.Fatalf("Failed to update document: %v", err)
// 	}

// 	// Verify update
// 	updatedDoc, err := idx.GetDocument(docID)
// 	if err != nil {
// 		t.Fatalf("Failed to get document: %v", err)
// 	}
// 	if title, _ := updatedDoc.GetField("title"); title.Value != "updated title" {
// 		t.Errorf("Expected updated title, got %v", title.Value)
// 	}

// 	// Verify term frequencies are updated
// 	tf, err := idx.GetTermFrequency("initial", docID)
// 	if err == nil && tf > 0 {
// 		t.Error("Old terms should not exist in updated document")
// 	}
// 	tf, err = idx.GetTermFrequency("updated", docID)
// 	if err != nil || tf == 0 {
// 		t.Error("New terms should exist in updated document")
// 	}
// }

// func TestDocumentDeletion(t *testing.T) {
// 	idx := NewIndex(nil)

// 	// Add document
// 	doc := document.NewDocument()
// 	doc.AddField("title", "test title")
// 	doc.AddField("content", "test content")
// 	docID, err := idx.AddDocument(doc)
// 	if err != nil {
// 		t.Fatalf("Failed to add document: %v", err)
// 	}

// 	// Delete document
// 	err = idx.DeleteDocument(docID)
// 	if err != nil {
// 		t.Fatalf("Failed to delete document: %v", err)
// 	}

// 	// Verify deletion
// 	_, err = idx.GetDocument(docID)
// 	if err == nil {
// 		t.Error("Document should not exist after deletion")
// 	}

// 	// Verify posting lists are updated
// 	tf, err := idx.GetTermFrequency("test", docID)
// 	if err == nil && tf > 0 {
// 		t.Error("Terms from deleted document should not exist in index")
// 	}
// }

// func TestIndexOptimization(t *testing.T) {
// 	idx := NewIndex(nil)

// 	// Add and delete documents to create gaps
// 	for i := 0; i < 10; i++ {
// 		doc := document.NewDocument()
// 		doc.AddField("content", "test content")
// 		docID, _ := idx.AddDocument(doc)
// 		if i%2 == 0 {
// 			idx.DeleteDocument(docID)
// 		}
// 	}

// 	// Optimize index
// 	err := idx.Optimize()
// 	if err != nil {
// 		t.Fatalf("Failed to optimize index: %v", err)
// 	}

// 	// Verify optimization
// 	if idx.GetDocumentCount() != 5 {
// 		t.Errorf("Expected 5 documents after optimization, got %d", idx.GetDocumentCount())
// 	}
// }

func TestConcurrentModifications(t *testing.T) {
	idx := NewIndex(nil)
	var wg sync.WaitGroup
	numOps := 100

	// Concurrent additions
	wg.Add(numOps)
	for i := 0; i < numOps; i++ {
		go func(i int) {
			defer wg.Done()
			doc := document.NewDocument()
			doc.AddField("content", "test content")
			idx.AddDocument(doc)
		}(i)
	}
	wg.Wait()

	// Verify document count
	if count := idx.GetDocumentCount(); count != numOps {
		t.Errorf("Expected %d documents, got %d", numOps, count)
	}
}
