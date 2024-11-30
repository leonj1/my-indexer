"""
Index module for document storage and search.
"""
from typing import Dict, List, Tuple, Optional
from .analysis import StandardAnalyzer
from .document import Document
from .storage import IndexStorage

class Index:
    """
    Index for storing and searching documents.
    Wraps the Go Index functionality.
    """
    
    def __init__(self, analyzer: StandardAnalyzer, storage: Optional[IndexStorage] = None):
        """
        Initialize the index.
        
        Args:
            analyzer: Text analyzer
            storage: Optional storage backend
        """
        self.analyzer = analyzer
        self.storage = storage or IndexStorage("index.gob")
        self._documents = self.storage.load_data()  # Load documents from storage
        
    def add_document(self, doc: Document) -> int:
        """
        Add a document to the index.
        
        Args:
            doc: Document to add
            
        Returns:
            Document ID
        """
        doc_id = len(self._documents)
        doc.id = doc_id  # Set document ID
        self._documents[doc_id] = doc
        return doc_id
        
    def get_document(self, doc_id: int) -> Document:
        """
        Get a document by ID.
        
        Args:
            doc_id: Document ID
            
        Returns:
            Document
            
        Raises:
            KeyError: If document not found
        """
        if doc_id not in self._documents:
            raise KeyError(f"Document {doc_id} not found")
        return self._documents[doc_id]
        
    def get_all_documents(self) -> List[Document]:
        """
        Get all documents in the index.
        
        Returns:
            List of documents
        """
        return list(self._documents.values())
        
    def search(self, query: Dict) -> List[Tuple[int, float]]:
        """
        Search for documents.
        
        Args:
            query: Search query
            
        Returns:
            List of (document_id, score) tuples
        """
        # This is just a placeholder - the actual implementation is in Go
        return [(doc_id, 1.0) for doc_id in self._documents.keys()]
        
    def save(self):
        """Save the index to storage."""
        self.storage.save_data(self._documents)
