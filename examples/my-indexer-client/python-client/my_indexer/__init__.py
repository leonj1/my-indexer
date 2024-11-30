"""
My Indexer Python Client Library

A Python library for full-text search and document indexing.
"""

from .analysis import StandardAnalyzer
from .index import Index
from .document import Document
from .storage import IndexStorage

__version__ = "0.1.0"

__all__ = [
    "StandardAnalyzer",
    "Index",
    "Document",
    "IndexStorage",
]
