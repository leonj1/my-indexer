package index

import (
	"fmt"
	"sync"

	"my-indexer/analysis"
	"my-indexer/document"
	"my-indexer/txlog"
)

// PostingList represents a list of documents containing a term
type PostingList struct {
	DocFreq int                    // Number of documents containing the term
	Postings map[int]*PostingEntry // Map of document ID to posting entry
}

// PostingEntry represents a single document entry in a posting list
type PostingEntry struct {
	DocID         int       // Document ID
	TermFreq      int       // Frequency of term in document
	Positions     []int     // Positions of term in document
	FieldName     string    // Name of the field containing the term
	Fields        []string  // Names of the fields containing the term
}

// Index represents an inverted index
type Index struct {
	mu            sync.RWMutex
	terms         map[string]*PostingList
	docCount      int
	analyzer      analysis.Analyzer
	nextDocID     int
	docIDMap      map[int]*document.Document // Maps document IDs to documents
	txLog         *txlog.TransactionLog      // Transaction log for crash recovery
}

// NewIndex creates a new inverted index
func NewIndex(analyzer analysis.Analyzer) *Index {
	if analyzer == nil {
		analyzer = analysis.NewStandardAnalyzer()
	}
	return &Index{
		terms:     make(map[string]*PostingList),
		analyzer:  analyzer,
		docIDMap:  make(map[int]*document.Document),
	}
}

// InitTransactionLog initializes the transaction log
func (idx *Index) InitTransactionLog(logDir string) error {
	txLog, err := txlog.NewTransactionLog(logDir)
	if err != nil {
		return fmt.Errorf("failed to create transaction log: %v", err)
	}
	idx.txLog = txLog

	// Recover any pending operations
	return idx.recover()
}

// recover processes any pending operations from the transaction log
func (idx *Index) recover() error {
	fmt.Printf("recover: Starting recovery process\n")
	if idx.txLog == nil {
		fmt.Printf("recover: No transaction log present, skipping recovery\n")
		return nil
	}

	fmt.Printf("recover: Attempting to recover entries from transaction log\n")
	entries, err := idx.txLog.Recover()
	if err != nil {
		return fmt.Errorf("failed to recover from transaction log: %v", err)
	}

	// Reset index state
	fmt.Printf("recover: Resetting index state\n")
	idx.terms = make(map[string]*PostingList)
	idx.docIDMap = make(map[int]*document.Document)
	idx.docCount = 0
	idx.nextDocID = 0

	fmt.Printf("recover: Processing %d entries in chronological order\n", len(entries))
	// Process entries in chronological order
	for _, entry := range entries {
		if !entry.Committed {
			continue
		}

		fmt.Printf("recover: Processing entry [Operation: %s, DocID: %d, Committed: %v]\n", 
			entry.Operation, entry.DocumentID, entry.Committed)
		switch entry.Operation {
		case txlog.OpAdd:
			if entry.Document != nil {
				// Create a new document and copy all fields
				newDoc := document.NewDocument()
				for _, field := range entry.Document.GetFields() {
					newDoc.AddField(field.Name, field.Value)
				}
				
				// Use the original document ID from the log entry
				idx.docIDMap[entry.DocumentID] = newDoc
				idx.docCount++
				
				// Index the document terms
				docTermFreqs := make(map[string]int)
				for _, field := range newDoc.GetFields() {
					fieldValue, ok := field.Value.(string)
					if !ok {
						continue
					}
					
					tokens := idx.analyzer.Analyze(fieldValue)
					for _, token := range tokens {
						docTermFreqs[token.Text]++
					}
				}
				
				// Update posting lists
				for term, freq := range docTermFreqs {
					postingList, exists := idx.terms[term]
					if !exists {
						postingList = &PostingList{
							Postings: make(map[int]*PostingEntry),
						}
						idx.terms[term] = postingList
					}
					
					postingEntry := &PostingEntry{
						DocID:    entry.DocumentID,
						TermFreq: freq,
					}
					postingList.Postings[entry.DocumentID] = postingEntry
					postingList.DocFreq++
				}
			}
		case txlog.OpUpdate:
			if entry.Document != nil {
				// Create a new document and copy all fields
				newDoc := document.NewDocument()
				for _, field := range entry.Document.GetFields() {
					if err := newDoc.AddField(field.Name, field.Value); err != nil {
						return fmt.Errorf("failed to restore field %s: %v", field.Name, err)
					}
				}
			
				// Store document directly in map since we're recovering
				idx.docIDMap[entry.DocumentID] = newDoc
				
				// Index the document terms
				docTermFreqs := make(map[string]int)
				for _, field := range newDoc.GetFields() {
					fieldValue, ok := field.Value.(string)
					if !ok {
						continue
					}
					
					tokens := idx.analyzer.Analyze(fieldValue)
					for _, token := range tokens {
						docTermFreqs[token.Text]++
					}
				}
				
				// Update posting lists
				for term, freq := range docTermFreqs {
					postingList, exists := idx.terms[term]
					if !exists {
						postingList = &PostingList{
							Postings: make(map[int]*PostingEntry),
						}
						idx.terms[term] = postingList
					}
					
					postingEntry := &PostingEntry{
						DocID:    entry.DocumentID,
						TermFreq: freq,
					}
					postingList.Postings[entry.DocumentID] = postingEntry
					postingList.DocFreq++
				}
			}
		case txlog.OpDelete:
			// Only attempt delete if document exists
			if _, exists := idx.docIDMap[entry.DocumentID]; exists {
				if err := idx.deleteDocumentInternal(entry.DocumentID); err != nil {
					return fmt.Errorf("failed to replay delete operation: %v", err)
				}
			}
		}
	}

	// Update nextDocID to be after the highest used ID
	maxID := -1
	for docID := range idx.docIDMap {
		if docID > maxID {
			maxID = docID
		}
	}
	idx.nextDocID = maxID + 1
	fmt.Printf("recover: Set nextDocID to %d after scanning existing documents\n", idx.nextDocID)

	fmt.Printf("recover: Recovery completed successfully, truncating log\n")
	// Clear the log after successful recovery
	return idx.txLog.Truncate()
}

// addDocumentInternal adds a document without transaction logging
func (idx *Index) addDocumentInternal(doc *document.Document) (int, error) {
	if doc == nil {
		return 0, fmt.Errorf("cannot index nil document")
	}

	// Note: Caller must hold write lock
	docID := idx.nextDocID
	idx.nextDocID++
	idx.docCount++

	// Store document in map
	idx.docIDMap[docID] = doc

	// Track total term frequencies across all fields
	docTermFreqs := make(map[string]int)

	// First pass: collect term frequencies across all fields
	for _, field := range doc.GetFields() {
		fieldValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		tokens := idx.analyzer.Analyze(fieldValue)
		for _, token := range tokens {
			docTermFreqs[token.Text]++
		}
	}

	// Second pass: update posting lists
	for term, freq := range docTermFreqs {
		postingList, exists := idx.terms[term]
		if !exists {
			postingList = &PostingList{
				Postings: make(map[int]*PostingEntry),
			}
			idx.terms[term] = postingList
		}

		entry := &PostingEntry{
			DocID:    docID,
			TermFreq: freq,
		}
		postingList.Postings[docID] = entry
		postingList.DocFreq++
	}

	return docID, nil
}

// AddDocument adds a document to the index with transaction logging
func (idx *Index) AddDocument(doc *document.Document) (int, error) {
	fmt.Printf("AddDocument: Starting...\n")
	if doc == nil {
		return 0, fmt.Errorf("cannot index nil document")
	}

	fmt.Printf("AddDocument: Attempting to acquire write lock\n")
	idx.mu.Lock()
	fmt.Printf("AddDocument: Write lock acquired\n")
	defer func() {
		idx.mu.Unlock()
		fmt.Printf("AddDocument: Released write lock\n")
	}()

	// Get the next document ID under the lock
	docID := idx.nextDocID

	// Handle transaction logging if enabled
	if idx.txLog != nil {
		fmt.Printf("AddDocument: Using transaction log\n")
		if err := idx.txLog.LogOperation(txlog.OpAdd, docID, doc); err != nil {
			return 0, fmt.Errorf("failed to log add operation: %v", err)
		}

		// Add the document with transaction logging
		id, err := idx.addDocumentInternal(doc)
		if err != nil {
			idx.txLog.Rollback(docID)
			return 0, err
		}

		// Commit the operation
		if err := idx.txLog.Commit(docID); err != nil {
			return 0, fmt.Errorf("failed to commit add operation: %v", err)
		}

		return id, nil
	}

	// If no transaction log, add document directly
	return idx.addDocumentInternal(doc)
}

// updateDocumentInternal updates a document without transaction logging
func (idx *Index) updateDocumentInternal(docID int, doc *document.Document) error {
	if doc == nil {
		return fmt.Errorf("cannot update with nil document")
	}

	// Note: Caller must hold write lock
	oldDoc, exists := idx.docIDMap[docID]
	if !exists {
		return fmt.Errorf("document with ID %d does not exist", docID)
	}

	// Remove old document's terms
	for _, field := range oldDoc.GetFields() {
		fieldValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		tokens := idx.analyzer.Analyze(fieldValue)
		for _, token := range tokens {
			if postingList, exists := idx.terms[token.Text]; exists {
				if _, exists := postingList.Postings[docID]; exists {
					delete(postingList.Postings, docID)
					postingList.DocFreq--
					if postingList.DocFreq == 0 {
						delete(idx.terms, token.Text)
					}
				}
			}
		}
	}

	// Add new document's terms
	docTermFreqs := make(map[string]int)
	for _, field := range doc.GetFields() {
		fieldValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		tokens := idx.analyzer.Analyze(fieldValue)
		for _, token := range tokens {
			docTermFreqs[token.Text]++
		}
	}

	for term, freq := range docTermFreqs {
		postingList, exists := idx.terms[term]
		if !exists {
			postingList = &PostingList{
				Postings: make(map[int]*PostingEntry),
			}
			idx.terms[term] = postingList
		}

		entry := &PostingEntry{
			DocID:    docID,
			TermFreq: freq,
		}
		postingList.Postings[docID] = entry
		if _, exists := postingList.Postings[docID]; !exists {
			postingList.DocFreq++
		}
	}

	idx.docIDMap[docID] = doc
	return nil
}

// UpdateDocument updates a document with transaction logging
func (idx *Index) UpdateDocument(docID int, doc *document.Document) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Log the operation first
	if idx.txLog != nil {
		if err := idx.txLog.LogOperation(txlog.OpUpdate, docID, doc); err != nil {
			return fmt.Errorf("failed to log update operation: %v", err)
		}

		// Update the document
		if err := idx.updateDocumentInternal(docID, doc); err != nil {
			idx.txLog.Rollback(docID)
			return err
		}

		// Commit the operation
		if err := idx.txLog.Commit(docID); err != nil {
			return fmt.Errorf("failed to commit update operation: %v", err)
		}

		return nil
	}

	// If no transaction log, just update the document
	return idx.updateDocumentInternal(docID, doc)
}

// deleteDocumentInternal deletes a document without transaction logging
func (idx *Index) deleteDocumentInternal(docID int) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	doc, exists := idx.docIDMap[docID]
	if !exists {
		return fmt.Errorf("document with ID %d does not exist", docID)
	}

	// Remove document's terms from posting lists
	for _, field := range doc.GetFields() {
		fieldValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		tokens := idx.analyzer.Analyze(fieldValue)
		for _, token := range tokens {
			if postingList, exists := idx.terms[token.Text]; exists {
				if _, exists := postingList.Postings[docID]; exists {
					delete(postingList.Postings, docID)
					postingList.DocFreq--
					if postingList.DocFreq == 0 {
						delete(idx.terms, token.Text)
					}
				}
			}
		}
	}

	delete(idx.docIDMap, docID)
	idx.docCount--
	return nil
}

// DeleteDocument deletes a document with transaction logging
func (idx *Index) DeleteDocument(docID int) error {
	// Log the operation first if transaction logging is enabled
	if idx.txLog != nil {
		if err := idx.txLog.LogOperation(txlog.OpDelete, docID, nil); err != nil {
			return fmt.Errorf("failed to log delete operation: %v", err)
		}

		// Delete the document
		if err := idx.deleteDocumentInternal(docID); err != nil {
			idx.txLog.Rollback(docID)
			return err
		}

		// Commit the operation
		if err := idx.txLog.Commit(docID); err != nil {
			return fmt.Errorf("failed to commit delete operation: %v", err)
		}

		return nil
	}

	// If no transaction log, just delete the document
	return idx.deleteDocumentInternal(docID)
}

// Close closes the index and its transaction log
func (idx *Index) Close() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.txLog != nil {
		if err := idx.txLog.Close(); err != nil {
			return fmt.Errorf("failed to close transaction log: %v", err)
		}
	}
	return nil
}

// GetDocument retrieves a document by its ID
func (idx *Index) GetDocument(docID int) (*document.Document, error) {
	fmt.Printf("GetDocument: Attempting to acquire read lock for docID %d\n", docID)
	idx.mu.RLock()
	defer func() {
		idx.mu.RUnlock()
		fmt.Printf("GetDocument: Released read lock for docID %d\n", docID)
	}()

	doc, exists := idx.docIDMap[docID]
	if !exists {
		return nil, fmt.Errorf("document with ID %d not found", docID)
	}
	return doc, nil
}

// GetPostingList retrieves the posting list for a term
func (idx *Index) GetPostingList(term string) (*PostingList, error) {
	if term == "" {
		return nil, fmt.Errorf("empty term")
	}

	idx.mu.RLock()
	defer idx.mu.RUnlock()

	// Analyze the term using the same analyzer
	tokens := idx.analyzer.Analyze(term)
	if len(tokens) == 0 {
		return nil, nil
	}

	// Use the first token as the term
	analyzedTerm := tokens[0].Text
	postingList, exists := idx.terms[analyzedTerm]
	if !exists {
		return nil, nil
	}

	return postingList, nil
}

// GetTermFrequency returns the frequency of a term in a document
func (idx *Index) GetTermFrequency(term string, docID int) (int, error) {
	fmt.Printf("GetTermFrequency: Attempting to acquire read lock for term '%s' docID %d\n", term, docID)
	idx.mu.RLock()
	defer func() {
		idx.mu.RUnlock()
		fmt.Printf("GetTermFrequency: Released read lock for term '%s' docID %d\n", term, docID)
	}()

	if postingList, exists := idx.terms[term]; exists {
		if entry, exists := postingList.Postings[docID]; exists {
			return entry.TermFreq, nil
		}
	}
	return 0, nil
}

// GetDocumentFrequency returns the number of documents containing a term
func (idx *Index) GetDocumentFrequency(term string) (int, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if postingList, exists := idx.terms[term]; exists {
		return postingList.DocFreq, nil
	}
	return 0, nil
}

// GetPostings returns the posting list entries for a term
func (idx *Index) GetPostings(term string) map[int]*PostingEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if postingList, exists := idx.terms[term]; exists {
		// Create a copy to avoid concurrent access issues
		entries := make(map[int]*PostingEntry, len(postingList.Postings))
		for docID, entry := range postingList.Postings {
			entries[docID] = entry
		}
		return entries
	}
	return make(map[int]*PostingEntry)
}

// GetDocumentCount returns the total number of documents in the index
func (idx *Index) GetDocumentCount() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.docCount
}

// GetTerms returns a copy of the terms map for serialization
func (idx *Index) GetTerms() map[string]*PostingList {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	terms := make(map[string]*PostingList, len(idx.terms))
	for term, postingList := range idx.terms {
		terms[term] = postingList
	}
	return terms
}

// GetNextDocID returns the next document ID
func (idx *Index) GetNextDocID() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.nextDocID
}

// RestoreFromData restores the index state from serialized data
func (idx *Index) RestoreFromData(terms map[string]*PostingList, docCount, nextDocID int) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.terms = terms
	idx.docCount = docCount
	idx.nextDocID = nextDocID
	return nil
}

// Optimize performs index optimization by removing gaps in document IDs
// and cleaning up unused terms
func (idx *Index) Optimize() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Create new document ID mapping
	newDocIDMap := make(map[int]*document.Document)
	oldToNewID := make(map[int]int)
	newID := 0

	// Reassign document IDs sequentially
	for oldID, doc := range idx.docIDMap {
		newDocIDMap[newID] = doc
		oldToNewID[oldID] = newID
		newID++
	}

	// Update posting lists with new document IDs
	newTerms := make(map[string]*PostingList)
	for term, postingList := range idx.terms {
		newPostings := make(map[int]*PostingEntry)
		for oldID, entry := range postingList.Postings {
			if newID, exists := oldToNewID[oldID]; exists {
				entry.DocID = newID
				newPostings[newID] = entry
			}
		}
		if len(newPostings) > 0 {
			postingList.Postings = newPostings
			newTerms[term] = postingList
		}
	}

	// Update index state
	idx.docIDMap = newDocIDMap
	idx.terms = newTerms
	idx.nextDocID = len(newDocIDMap)

	return nil
}

// IndexDocument indexes an ElasticSearch-compatible document
func (idx *Index) IndexDocument(indexName string, docID string, doc map[string]interface{}) error {
    // Create new document
    internalDoc := document.NewDocument()

    // Copy all fields from the ElasticSearch document
    for field, value := range doc {
        // Skip metadata fields that start with _
        if len(field) > 0 && field[0] == '_' {
            continue
        }
        if err := internalDoc.AddField(field, value); err != nil {
            return fmt.Errorf("failed to add field %s: %v", field, err)
        }
    }

    // If docID is provided, try to update existing document
    if docID != "" {
        // Convert string docID to int
        var intDocID int
        _, err := fmt.Sscanf(docID, "%d", &intDocID)
        if err != nil {
            return fmt.Errorf("invalid document ID format: %v", err)
        }

        // Check if document exists
        existingDoc, err := idx.GetDocument(intDocID)
        if err == nil && existingDoc != nil {
            // Update existing document
            return idx.UpdateDocument(intDocID, internalDoc)
        }
    }

    // Add as new document
    _, err := idx.AddDocument(internalDoc)
    return err
}
