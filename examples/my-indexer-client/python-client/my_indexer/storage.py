"""
Storage module for persisting index data.
"""
import os
import pickle
from pathlib import Path
from typing import Dict, Any

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
        self._load_data()
        
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
            
    def _load_data(self):
        """Load data from storage."""
        self.data = {}
        if os.path.exists(self.index_path):
            try:
                with open(self.index_path, 'rb') as f:
                    self.data = pickle.load(f)
            except (pickle.PickleError, EOFError):
                # If file is corrupted, start with empty data
                self.data = {}
                
    def save_data(self, data: Dict[str, Any]):
        """
        Save data to storage.
        
        Args:
            data: Data to save
        """
        with open(self.index_path, 'wb') as f:
            pickle.dump(data, f)
            
    def load_data(self) -> Dict[str, Any]:
        """
        Load data from storage.
        
        Returns:
            Stored data
        """
        return self.data
