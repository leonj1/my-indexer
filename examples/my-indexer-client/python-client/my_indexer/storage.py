"""
Storage module for persisting index data.
"""
import os
from pathlib import Path

class IndexStorage:
    """
    Storage backend for the index.
    Wraps the Go IndexStorage functionality.
    """
    
    def __init__(self, index_path: str):
        """
        Initialize storage.
        
        Args:
            index_path: Path to index file
            
        Raises:
            ValueError: If path is invalid
        """
        self.index_path = index_path
        self._validate_path(index_path)
        
    def _validate_path(self, path: str):
        """
        Validate the index path.
        
        Args:
            path: Path to validate
            
        Raises:
            ValueError: If path is invalid
        """
        if not path.endswith('.gob'):
            raise ValueError("Index path must end with .gob")
            
        # Check for directory traversal
        abs_path = os.path.abspath(path)
        if not os.path.normpath(abs_path) == abs_path:
            raise ValueError("Invalid path: possible directory traversal")
            
        # Create directory if it doesn't exist
        directory = os.path.dirname(abs_path)
        if directory:
            os.makedirs(directory, exist_ok=True)
