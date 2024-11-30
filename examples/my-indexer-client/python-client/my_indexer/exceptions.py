class IndexerError(Exception):
    """Base exception for indexer client errors."""
    pass

class DocumentNotFoundError(IndexerError):
    """Raised when a document is not found."""
    pass

class IndexerConnectionError(IndexerError):
    """Raised when connection to indexer fails."""
    pass

class InvalidQueryError(IndexerError):
    """Raised when query is invalid."""
    pass
