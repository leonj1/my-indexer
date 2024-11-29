package txlog

import (
	"os"
	"testing"

	"my-indexer/document"
)

func TestTransactionLogging(t *testing.T) {
	// Create temporary directory for test logs
	tmpDir, err := os.MkdirTemp("", "txlog_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create transaction log
	txLog, err := NewTransactionLog(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create transaction log: %v", err)
	}
	defer txLog.Close()

	// Test logging an add operation
	doc := document.NewDocument()
	doc.AddField("title", "test document")
	err = txLog.LogOperation(OpAdd, 1, doc)
	if err != nil {
		t.Errorf("Failed to log add operation: %v", err)
	}

	// Verify uncommitted operations
	uncommitted := txLog.GetUncommittedOperations()
	if len(uncommitted) != 1 {
		t.Errorf("Expected 1 uncommitted operation, got %d", len(uncommitted))
	}

	// Test commit
	err = txLog.Commit(1)
	if err != nil {
		t.Errorf("Failed to commit operation: %v", err)
	}

	uncommitted = txLog.GetUncommittedOperations()
	if len(uncommitted) != 0 {
		t.Errorf("Expected 0 uncommitted operations after commit, got %d", len(uncommitted))
	}
}

func TestRecovery(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "txlog_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create and populate transaction log
	txLog, err := NewTransactionLog(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create transaction log: %v", err)
	}

	// Add some operations
	doc1 := document.NewDocument()
	doc1.AddField("title", "doc1")
	doc2 := document.NewDocument()
	doc2.AddField("title", "doc2")

	txLog.LogOperation(OpAdd, 1, doc1)
	txLog.LogOperation(OpAdd, 2, doc2)
	txLog.Commit(1)

	// Close the log
	txLog.Close()

	// Create new transaction log instance and recover
	recoveredLog, err := NewTransactionLog(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create new transaction log: %v", err)
	}
	defer recoveredLog.Close()

	entries, err := recoveredLog.Recover()
	if err != nil {
		t.Fatalf("Failed to recover log: %v", err)
	}

	// Verify recovered entries
	if len(entries) != 3 { // 2 adds + 1 commit
		t.Errorf("Expected 3 recovered entries, got %d", len(entries))
	}

	var committed, uncommitted int
	for _, entry := range entries {
		if entry.Committed {
			committed++
		} else {
			uncommitted++
		}
	}

	if committed != 1 {
		t.Errorf("Expected 1 committed entry, got %d", committed)
	}
	if uncommitted != 2 {
		t.Errorf("Expected 2 uncommitted entries, got %d", uncommitted)
	}
}

func TestRollback(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "txlog_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	txLog, err := NewTransactionLog(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create transaction log: %v", err)
	}
	defer txLog.Close()

	// Log an operation
	doc := document.NewDocument()
	doc.AddField("title", "test rollback")
	err = txLog.LogOperation(OpAdd, 1, doc)
	if err != nil {
		t.Errorf("Failed to log operation: %v", err)
	}

	// Verify operation is uncommitted
	uncommitted := txLog.GetUncommittedOperations()
	if len(uncommitted) != 1 {
		t.Errorf("Expected 1 uncommitted operation, got %d", len(uncommitted))
	}

	// Test rollback
	err = txLog.Rollback(1)
	if err != nil {
		t.Errorf("Failed to rollback operation: %v", err)
	}

	// Verify operation was rolled back
	uncommitted = txLog.GetUncommittedOperations()
	if len(uncommitted) != 0 {
		t.Errorf("Expected 0 uncommitted operations after rollback, got %d", len(uncommitted))
	}
}

func TestCrashRecovery(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "txlog_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create and populate transaction log
	txLog, err := NewTransactionLog(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create transaction log: %v", err)
	}

	// Add some operations
	doc1 := document.NewDocument()
	doc1.AddField("title", "doc1")
	doc2 := document.NewDocument()
	doc2.AddField("title", "doc2")

	txLog.LogOperation(OpAdd, 1, doc1)
	txLog.Commit(1)
	txLog.LogOperation(OpAdd, 2, doc2)

	// Simulate crash by not closing properly
	txLog.file.Close()

	// Recover after crash
	recoveredLog, err := NewTransactionLog(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create new transaction log: %v", err)
	}
	defer recoveredLog.Close()

	entries, err := recoveredLog.Recover()
	if err != nil {
		t.Fatalf("Failed to recover log: %v", err)
	}

	// Verify recovered entries
	if len(entries) != 3 { // 2 adds + 1 commit
		t.Errorf("Expected 3 recovered entries, got %d", len(entries))
	}

	// Check that doc1 is committed and doc2 is not
	var foundCommitted, foundUncommitted bool
	for _, entry := range entries {
		if entry.DocumentID == 1 && entry.Committed {
			foundCommitted = true
		}
		if entry.DocumentID == 2 && !entry.Committed {
			foundUncommitted = true
		}
	}

	if !foundCommitted {
		t.Error("Failed to find committed document 1")
	}
	if !foundUncommitted {
		t.Error("Failed to find uncommitted document 2")
	}
}
