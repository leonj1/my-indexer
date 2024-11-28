package index

import (
	"fmt"
	"sync"

	"my-indexer/analysis"
	"my-indexer/document"
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

// AddDocument adds a document to the index
func (idx *Index) AddDocument(doc *document.Document) (int, error) {
	if doc == nil {
		return 0, fmt.Errorf("cannot index nil document")
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	docID := idx.nextDocID
	idx.nextDocID++
	idx.docCount++

	// Store document in map
	idx.docIDMap[docID] = doc

	// Track total term frequencies across all fields
	docTermFreqs := make(map[string]int)

	// First pass: collect term frequencies across all fields
	for _, field := range doc.GetFields() {
		// Skip non-string fields
		fieldValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		// Analyze the text
		tokens := idx.analyzer.Analyze(fieldValue)

		// Count term frequencies
		for _, token := range tokens {
			docTermFreqs[token.Text]++
		}
	}

	// Second pass: update posting lists with total frequencies
	for term, totalFreq := range docTermFreqs {
		postingList, exists := idx.terms[term]
		if !exists {
			postingList = &PostingList{
				Postings: make(map[int]*PostingEntry),
			}
			idx.terms[term] = postingList
		}

		// Add or update posting entry
		entry, exists := postingList.Postings[docID]
		if !exists {
			entry = &PostingEntry{
				DocID:     docID,
				Positions: make([]int, 0),
			}
			postingList.Postings[docID] = entry
			postingList.DocFreq++
		}
		entry.TermFreq = totalFreq
	}

	return docID, nil
}

// GetDocument retrieves a document by its ID
func (idx *Index) GetDocument(docID int) (*document.Document, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

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
	idx.mu.RLock()
	defer idx.mu.RUnlock()

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

// UpdateDocument updates an existing document in the index
func (idx *Index) UpdateDocument(docID int, doc *document.Document) error {
	if doc == nil {
		return fmt.Errorf("cannot update with nil document")
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Check if document exists
	oldDoc, exists := idx.docIDMap[docID]
	if !exists {
		return fmt.Errorf("document with ID %d does not exist", docID)
	}

	// Remove old document's terms from posting lists
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

	// Index new document's terms
	docTermFreqs := make(map[string]int)

	// First pass: collect term frequencies
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
		if _, exists := postingList.Postings[docID]; !exists {
			postingList.DocFreq++
		}
	}

	// Update document in map
	idx.docIDMap[docID] = doc

	return nil
}

// DeleteDocument removes a document from the index
func (idx *Index) DeleteDocument(docID int) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Check if document exists
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

	// Remove document from map and decrement count
	delete(idx.docIDMap, docID)
	idx.docCount--

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
