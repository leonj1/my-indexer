package storage

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"my-indexer/document"
	"my-indexer/index"
)

// IndexStorage handles persistence of the index
type IndexStorage struct {
	mu           sync.RWMutex
	indexPath    string
	documentsDir string
}

// IndexData represents the serializable form of the index
type IndexData struct {
	Terms    map[string]*index.PostingList
	DocCount int
	NextID   int
}

// DocumentData represents the serializable form of a document
type DocumentData struct {
	Fields map[string]document.Field
}

const (
	// DefaultIndexFilename is the default name for the index file
	DefaultIndexFilename = "index.gob"
)

// NewIndexStorage creates a new index storage
func NewIndexStorage(baseDir string, indexFilename string) (*IndexStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// If no index filename is provided, use the default
	if indexFilename == "" {
		indexFilename = DefaultIndexFilename
	}

	// Create documents directory
	documentsDir := filepath.Join(baseDir, "documents")
	if err := os.MkdirAll(documentsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create documents directory: %w", err)
	}

	return &IndexStorage{
		indexPath:    filepath.Join(baseDir, indexFilename),
		documentsDir: documentsDir,
	}, nil
}

// SaveIndex persists the index to disk
func (s *IndexStorage) SaveIndex(idx *index.Index) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a temporary file for atomic write
	tempPath := s.indexPath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary index file: %w", err)
	}
	defer file.Close()

	// Prepare index data for serialization
	data := &IndexData{
		Terms:    idx.GetTerms(),
		DocCount: idx.GetDocumentCount(),
		NextID:   idx.GetNextDocID(),
	}

	// Serialize index data
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to encode index: %w", err)
	}

	// Ensure all data is written to disk
	if err := file.Sync(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to sync index file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, s.indexPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to save index file: %w", err)
	}

	return nil
}

// LoadIndex loads the index from disk
func (s *IndexStorage) LoadIndex() (*index.Index, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return index.NewIndex(nil), nil
		}
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}
	defer file.Close()

	var data IndexData
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode index: %w", err)
	}

	// Create new index and restore its state
	idx := index.NewIndex(nil)
	if err := idx.RestoreFromData(data.Terms, data.DocCount, data.NextID); err != nil {
		return nil, fmt.Errorf("failed to restore index: %w", err)
	}

	return idx, nil
}

// SaveDocument persists a document to disk
func (s *IndexStorage) SaveDocument(docID int, doc *document.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	docPath := filepath.Join(s.documentsDir, fmt.Sprintf("doc_%d.gob", docID))
	
	// Create a temporary file for atomic write
	tempPath := docPath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary document file: %w", err)
	}
	defer file.Close()

	// Prepare document data for serialization
	data := &DocumentData{
		Fields: doc.GetFields(),
	}

	// Serialize document data
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to encode document: %w", err)
	}

	// Ensure all data is written to disk
	if err := file.Sync(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to sync document file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, docPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to save document file: %w", err)
	}

	return nil
}

// LoadDocument loads a document from disk
func (s *IndexStorage) LoadDocument(docID int) (*document.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docPath := filepath.Join(s.documentsDir, fmt.Sprintf("doc_%d.gob", docID))
	file, err := os.Open(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document file: %w", err)
	}
	defer file.Close()

	var data DocumentData
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	// Create new document and restore its fields
	doc := document.NewDocument()
	for name, field := range data.Fields {
		if err := doc.AddField(name, field.Value); err != nil {
			return nil, fmt.Errorf("failed to restore document field: %w", err)
		}
	}

	return doc, nil
}

// RemoveDocument removes a document from disk
func (s *IndexStorage) RemoveDocument(docID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	docPath := filepath.Join(s.documentsDir, fmt.Sprintf("doc_%d.gob", docID))
	if err := os.Remove(docPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove document file: %w", err)
	}

	return nil
}

// Clear removes all index and document files
func (s *IndexStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove index file
	if err := os.Remove(s.indexPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove index file: %w", err)
	}

	// Remove all document files
	entries, err := os.ReadDir(s.documentsDir)
	if err != nil {
		return fmt.Errorf("failed to read documents directory: %w", err)
	}

	for _, entry := range entries {
		if err := os.Remove(filepath.Join(s.documentsDir, entry.Name())); err != nil {
			return fmt.Errorf("failed to remove document file: %w", err)
		}
	}

	return nil
}
