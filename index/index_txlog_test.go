package index

// import (
// 	"os"
// 	"testing"

// 	"my-indexer/document"
// )

// func TestTransactionLogIntegration(t *testing.T) {
// 	// Create temporary directory for test logs
// 	tmpDir, err := os.MkdirTemp("", "index_txlog_test")
// 	if err != nil {
// 		t.Fatalf("Failed to create temp dir: %v", err)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	// Create index with transaction log
// 	idx := NewIndex(nil)
// 	err = idx.InitTransactionLog(tmpDir)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize transaction log: %v", err)
// 	}

// 	// Test document addition with transaction logging
// 	doc1 := document.NewDocument()
// 	doc1.AddField("title", "test document 1")
// 	docID, err := idx.AddDocument(doc1)
// 	if err != nil {
// 		t.Fatalf("Failed to add document: %v", err)
// 	}

// 	// Verify document was added
// 	retrievedDoc, err := idx.GetDocument(docID)
// 	if err != nil {
// 		t.Fatalf("Failed to retrieve document: %v", err)
// 	}
// 	if title, _ := retrievedDoc.GetField("title"); title.Value != "test document 1" {
// 		t.Errorf("Expected title 'test document 1', got '%v'", title.Value)
// 	}

// 	// Test document update with transaction logging
// 	doc2 := document.NewDocument()
// 	doc2.AddField("title", "updated document 1")
// 	err = idx.UpdateDocument(docID, doc2)
// 	if err != nil {
// 		t.Fatalf("Failed to update document: %v", err)
// 	}

// 	// Verify update
// 	retrievedDoc, err = idx.GetDocument(docID)
// 	if err != nil {
// 		t.Fatalf("Failed to retrieve updated document: %v", err)
// 	}
// 	if title, _ := retrievedDoc.GetField("title"); title.Value != "updated document 1" {
// 		t.Errorf("Expected title 'updated document 1', got '%v'", title.Value)
// 	}

// 	// Test crash recovery
// 	idx.Close()

// 	// Create new index instance
// 	newIdx := NewIndex(nil)
// 	err = newIdx.InitTransactionLog(tmpDir)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize transaction log for recovery: %v", err)
// 	}

// 	// Verify document state after recovery
// 	retrievedDoc, err = newIdx.GetDocument(docID)
// 	if err != nil {
// 		t.Fatalf("Failed to retrieve document after recovery: %v", err)
// 	}
// 	if title, _ := retrievedDoc.GetField("title"); title.Value != "updated document 1" {
// 		t.Errorf("Expected title 'updated document 1' after recovery, got '%v'", title.Value)
// 	}

// 	// Test document deletion with transaction logging
// 	err = newIdx.DeleteDocument(docID)
// 	if err != nil {
// 		t.Fatalf("Failed to delete document: %v", err)
// 	}

// 	// Verify deletion
// 	_, err = newIdx.GetDocument(docID)
// 	if err == nil {
// 		t.Error("Document should not exist after deletion")
// 	}

// 	// Test recovery after deletion
// 	newIdx.Close()

// 	finalIdx := NewIndex(nil)
// 	err = finalIdx.InitTransactionLog(tmpDir)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize transaction log for final recovery: %v", err)
// 	}
// 	defer finalIdx.Close()

// 	// Verify document remains deleted after recovery
// 	_, err = finalIdx.GetDocument(docID)
// 	if err == nil {
// 		t.Error("Document should not exist after recovery of deletion")
// 	}
// }

// func TestConcurrentTransactions(t *testing.T) {
// 	tmpDir, err := os.MkdirTemp("", "index_txlog_test")
// 	if err != nil {
// 		t.Fatalf("Failed to create temp dir: %v", err)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	idx := NewIndex(nil)
// 	err = idx.InitTransactionLog(tmpDir)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize transaction log: %v", err)
// 	}

// 	// Test concurrent additions
// 	done := make(chan bool)
// 	for i := 0; i < 10; i++ {
// 		go func(i int) {
// 			doc := document.NewDocument()
// 			doc.AddField("title", "concurrent document")
// 			_, err := idx.AddDocument(doc)
// 			if err != nil {
// 				t.Errorf("Failed to add document concurrently: %v", err)
// 			}
// 			done <- true
// 		}(i)
// 	}

// 	// Wait for all operations to complete
// 	for i := 0; i < 10; i++ {
// 		<-done
// 	}

// 	// Verify document count
// 	if count := idx.GetDocumentCount(); count != 10 {
// 		t.Errorf("Expected 10 documents after concurrent additions, got %d", count)
// 	}

// 	// Close and recover
// 	idx.Close()

// 	// Create new index instance
// 	newIdx := NewIndex(nil)
// 	err = newIdx.InitTransactionLog(tmpDir)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize transaction log for recovery: %v", err)
// 	}
// 	defer newIdx.Close()

// 	// Verify document count after recovery
// 	if count := newIdx.GetDocumentCount(); count != 10 {
// 		t.Errorf("Expected 10 documents after recovery, got %d", count)
// 	}
// }
