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

	// Assign document ID
	docID := idx.nextDocID
	idx.nextDocID++
	idx.docIDMap[docID] = doc

	// Index each field
	for fieldName, field := range doc.GetFields() {
		if field.Type != document.StringType {
			continue // Only index string fields for now
		}

		// Convert field value to string and analyze
		text, ok := field.Value.(string)
		if !ok {
			continue
		}

		tokens := idx.analyzer.Analyze(text)
		
		// Index each token
		for _, token := range tokens {
			postingList, exists := idx.terms[token.Text]
			if !exists {
				postingList = &PostingList{
					Postings: make(map[int]*PostingEntry),
				}
				idx.terms[token.Text] = postingList
			}

			posting, exists := postingList.Postings[docID]
			if !exists {
				posting = &PostingEntry{
					DocID:     docID,
					FieldName: fieldName,
					Positions: make([]int, 0),
				}
				postingList.Postings[docID] = posting
				postingList.DocFreq++
			}

			posting.TermFreq++
			posting.Positions = append(posting.Positions, token.Position)
		}
	}

	idx.docCount++
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

// GetTermFrequency returns the frequency of a term in a specific document
func (idx *Index) GetTermFrequency(term string, docID int) (int, error) {
	postingList, err := idx.GetPostingList(term)
	if err != nil {
		return 0, err
	}
	if postingList == nil {
		return 0, nil
	}

	posting, exists := postingList.Postings[docID]
	if !exists {
		return 0, nil
	}

	return posting.TermFreq, nil
}

// GetDocumentFrequency returns the number of documents containing a term
func (idx *Index) GetDocumentFrequency(term string) (int, error) {
	postingList, err := idx.GetPostingList(term)
	if err != nil {
		return 0, err
	}
	if postingList == nil {
		return 0, nil
	}

	return postingList.DocFreq, nil
}

// GetDocumentCount returns the total number of documents in the index
func (idx *Index) GetDocumentCount() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.docCount
}
