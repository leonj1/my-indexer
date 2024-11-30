package storage

import (
	"os"
	"path/filepath"
	"testing"

	"my-indexer/document"
	"my-indexer/index"
)

func TestIndexStorage(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "indexer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create storage
	storage, err := NewIndexStorage(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create and populate test index
	idx := index.NewIndex(nil)

	// Add test documents
	doc1 := document.NewDocument()
	err = doc1.AddField("title", "Test Document 1")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}
	err = doc1.AddField("content", "This is a test document")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}

	doc2 := document.NewDocument()
	err = doc2.AddField("title", "Test Document 2")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}
	err = doc2.AddField("content", "Another test document")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}

	// Add documents to index
	docID1, err := idx.AddDocument(doc1)
	if err != nil {
		t.Fatalf("Failed to add document 1: %v", err)
	}

	docID2, err := idx.AddDocument(doc2)
	if err != nil {
		t.Fatalf("Failed to add document 2: %v", err)
	}

	// Test saving documents
	err = storage.SaveDocument(docID1, doc1)
	if err != nil {
		t.Errorf("Failed to save document 1: %v", err)
	}

	err = storage.SaveDocument(docID2, doc2)
	if err != nil {
		t.Errorf("Failed to save document 2: %v", err)
	}

	// Test saving index
	err = storage.SaveIndex(idx)
	if err != nil {
		t.Errorf("Failed to save index: %v", err)
	}

	// Test loading index
	loadedIdx, err := storage.LoadIndex()
	if err != nil {
		t.Errorf("Failed to load index: %v", err)
	}

	// Verify loaded index
	if loadedIdx.GetDocumentCount() != idx.GetDocumentCount() {
		t.Errorf("Loaded index document count = %d, want %d",
			loadedIdx.GetDocumentCount(), idx.GetDocumentCount())
	}

	// Test loading documents
	loadedDoc1, err := storage.LoadDocument(docID1)
	if err != nil {
		t.Errorf("Failed to load document 1: %v", err)
	}

	loadedDoc2, err := storage.LoadDocument(docID2)
	if err != nil {
		t.Errorf("Failed to load document 2: %v", err)
	}

	// Verify loaded documents
	doc1Fields := doc1.GetFields()
	loadedDoc1Fields := loadedDoc1.GetFields()
	if len(doc1Fields) != len(loadedDoc1Fields) {
		t.Errorf("Loaded document 1 field count = %d, want %d",
			len(loadedDoc1Fields), len(doc1Fields))
	}

	doc2Fields := doc2.GetFields()
	loadedDoc2Fields := loadedDoc2.GetFields()
	if len(doc2Fields) != len(loadedDoc2Fields) {
		t.Errorf("Loaded document 2 field count = %d, want %d",
			len(loadedDoc2Fields), len(doc2Fields))
	}

	// Test document removal
	err = storage.RemoveDocument(docID1)
	if err != nil {
		t.Errorf("Failed to remove document 1: %v", err)
	}

	_, err = storage.LoadDocument(docID1)
	if err == nil {
		t.Error("Expected error loading removed document")
	}

	// Test clearing storage
	err = storage.Clear()
	if err != nil {
		t.Errorf("Failed to clear storage: %v", err)
	}

	// Verify storage is cleared
	if _, err := os.Stat(filepath.Join(tempDir, "index.gob")); !os.IsNotExist(err) {
		t.Error("Index file still exists after clear")
	}

	entries, err := os.ReadDir(filepath.Join(tempDir, "documents"))
	if err != nil {
		t.Errorf("Failed to read documents directory: %v", err)
	}
	if len(entries) > 0 {
		t.Error("Documents directory not empty after clear")
	}
}

func TestIndexStorageCustomFilename(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "indexer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test invalid filenames
	invalidNames := []string{
		"../index.gob",          // Path traversal
		"./index.gob",           // Path traversal
		"subdir/index.gob",      // Path separator
		"index<.gob",            // Invalid character <
		"index>.gob",            // Invalid character >
		"index:.gob",            // Invalid character :
		"index\".gob",           // Invalid character "
		"index?.gob",            // Invalid character ?
		"index*.gob",            // Invalid character *
		"index|.gob",            // Invalid character |
		"index",                 // Missing extension
		"index.txt",             // Wrong extension
	}

	for _, name := range invalidNames {
		if _, err := NewIndexStorage(tempDir, name); err == nil {
			t.Errorf("Expected error for invalid filename %q but got none", name)
		}
	}

	// Test valid filenames
	validNames := []string{
		"",                  // Empty string (should use default)
		"my-index.gob",     // Valid characters with hyphen
		"index_test_123.gob", // Valid characters with underscore and numbers
	}

	for _, name := range validNames {
		if _, err := NewIndexStorage(tempDir, name); err != nil {
			t.Errorf("Unexpected error for valid filename %q: %v", name, err)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "indexer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create storage
	storage, err := NewIndexStorage(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create test index and document
	idx := index.NewIndex(nil)
	doc := document.NewDocument()
	err = doc.AddField("content", "test document")
	if err != nil {
		t.Fatalf("Failed to add field to document: %v", err)
	}

	docID, err := idx.AddDocument(doc)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}

	// Save initial state
	err = storage.SaveIndex(idx)
	if err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}
	err = storage.SaveDocument(docID, doc)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			// Load index
			_, err := storage.LoadIndex()
			if err != nil {
				t.Errorf("Concurrent LoadIndex failed: %v", err)
			}

			// Load document
			_, err = storage.LoadDocument(docID)
			if err != nil {
				t.Errorf("Concurrent LoadDocument failed: %v", err)
			}

			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
