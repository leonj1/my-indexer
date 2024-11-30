"""
Document module for representing indexed documents.
"""
from typing import Dict, Optional

class Document:
    """
    Document class for storing document fields.
    """
    
    def __init__(self, fields: Dict[str, str], doc_id: Optional[int] = None):
        """
        Initialize a document.
        
        Args:
            fields: Document fields
            doc_id: Optional document ID
        """
        self.fields = fields
        self.id = doc_id
        
    def __repr__(self) -> str:
        """String representation of the document."""
        fields_str = ", ".join(f"{k}: {v}" for k, v in self.fields.items())
        return f"Document(id={self.id}, fields={{{fields_str}}})"
        
    def __getstate__(self):
        """Get state for pickling."""
        return {'fields': self.fields, 'id': self.id}
        
    def __setstate__(self, state):
        """Set state for unpickling."""
        self.fields = state['fields']
        self.id = state['id']
