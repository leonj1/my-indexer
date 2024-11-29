package txlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"my-indexer/document"
)

// Operation types for transaction log entries
const (
	OpAdd    = "add"
	OpUpdate = "update"
	OpDelete = "delete"
)

// LogEntry represents a single operation in the transaction log
type LogEntry struct {
	Operation   string              `json:"operation"`
	Timestamp   time.Time           `json:"timestamp"`
	DocumentID  int                 `json:"document_id"`
	Document    *document.Document  `json:"document,omitempty"`
	Committed   bool               `json:"committed"`
}

// TransactionLog manages write-ahead logging and recovery
type TransactionLog struct {
	mu           sync.RWMutex
	file         *os.File
	logPath      string
	encoder      *json.Encoder
	decoder      *json.Decoder
	uncommitted  map[int]*LogEntry
}

// NewTransactionLog creates a new transaction log
func NewTransactionLog(logDir string) (*TransactionLog, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	logPath := filepath.Join(logDir, "transaction.log")
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	txLog := &TransactionLog{
		file:        file,
		logPath:     logPath,
		encoder:     json.NewEncoder(file),
		decoder:     json.NewDecoder(file),
		uncommitted: make(map[int]*LogEntry),
	}

	return txLog, nil
}

// LogOperation logs an operation to the transaction log
func (t *TransactionLog) LogOperation(op string, docID int, doc *document.Document) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	entry := &LogEntry{
		Operation:  op,
		Timestamp:  time.Now(),
		DocumentID: docID,
		Document:   doc,
		Committed:  false,
	}

	if err := t.encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode log entry: %v", err)
	}

	t.uncommitted[docID] = entry
	return nil
}

// Commit marks an operation as committed
func (t *TransactionLog) Commit(docID int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	entry, exists := t.uncommitted[docID]
	if !exists {
		return fmt.Errorf("no uncommitted operation found for document ID %d", docID)
	}

	entry.Committed = true
	if err := t.encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode commit entry: %v", err)
	}

	delete(t.uncommitted, docID)
	return nil
}

// Rollback removes an uncommitted operation
func (t *TransactionLog) Rollback(docID int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.uncommitted[docID]; !exists {
		return fmt.Errorf("no uncommitted operation found for document ID %d", docID)
	}

	delete(t.uncommitted, docID)
	return nil
}

// GetUncommittedOperations returns all uncommitted operations
func (t *TransactionLog) GetUncommittedOperations() []*LogEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	entries := make([]*LogEntry, 0, len(t.uncommitted))
	for _, entry := range t.uncommitted {
		entries = append(entries, entry)
	}
	return entries
}

// Recover processes the transaction log and returns operations that need to be replayed
func (t *TransactionLog) Recover() ([]*LogEntry, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, err := t.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek to start of log: %v", err)
	}

	var entries []*LogEntry
	t.uncommitted = make(map[int]*LogEntry) // Reset uncommitted map

	for {
		var entry LogEntry
		if err := t.decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode log entry: %v", err)
		}

		entries = append(entries, &entry)
		
		// Track uncommitted entries in memory
		if !entry.Committed {
			t.uncommitted[entry.DocumentID] = &entry
		} else {
			// If we see a commit entry, remove the corresponding uncommitted entry
			delete(t.uncommitted, entry.DocumentID)
		}
	}

	return entries, nil
}

// Close closes the transaction log file
func (t *TransactionLog) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.file.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %v", err)
	}
	return nil
}

// Truncate removes all entries from the log file
func (t *TransactionLog) Truncate() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate log file: %v", err)
	}
	if _, err := t.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek to start of log: %v", err)
	}
	t.uncommitted = make(map[int]*LogEntry)
	return nil
}
