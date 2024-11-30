from typing import Dict, List, Optional, Union
import requests
from urllib.parse import urljoin

from .models import Document, SearchResponse
from .exceptions import (
    IndexerError,
    DocumentNotFoundError,
    IndexerConnectionError,
    InvalidQueryError,
)

class IndexerClient:
    """Client for interacting with the indexer service."""
    
    def __init__(self, base_url: str = "http://localhost:8080"):
        """Initialize the client.
        
        Args:
            base_url: Base URL of the indexer service
        """
        self.base_url = base_url.rstrip("/")
        
    def _make_request(
        self,
        method: str,
        endpoint: str,
        json: Optional[Dict] = None,
        params: Optional[Dict] = None,
    ) -> Dict:
        """Make HTTP request to indexer service.
        
        Args:
            method: HTTP method
            endpoint: API endpoint
            json: JSON body
            params: Query parameters
            
        Returns:
            Response JSON
            
        Raises:
            IndexerConnectionError: If connection fails
            IndexerError: If request fails
        """
        url = urljoin(self.base_url, endpoint)
        
        try:
            response = requests.request(
                method=method,
                url=url,
                json=json,
                params=params,
            )
            response.raise_for_status()
            return response.json()
        except requests.ConnectionError as e:
            raise IndexerConnectionError(f"Failed to connect to indexer: {e}")
        except requests.HTTPError as e:
            if response.status_code == 404:
                raise DocumentNotFoundError(f"Document not found: {e}")
            raise IndexerError(f"Request failed: {e}")
        except Exception as e:
            raise IndexerError(f"Unexpected error: {e}")
            
    def index_document(self, document: Union[Document, Dict]) -> Dict:
        """Index a document.
        
        Args:
            document: Document to index
            
        Returns:
            Response from indexer
        """
        if isinstance(document, Document):
            document = document.model_dump()
        
        return self._make_request("POST", "/_doc", json=document)
        
    def get_document(self, doc_id: int) -> Document:
        """Get document by ID.
        
        Args:
            doc_id: Document ID
            
        Returns:
            Document
            
        Raises:
            DocumentNotFoundError: If document not found
        """
        response = self._make_request("GET", f"/_doc/{doc_id}")
        return Document(**response)
        
    def search(self, query: Dict) -> SearchResponse:
        """Search documents.
        
        Args:
            query: Search query
            
        Returns:
            Search response
            
        Raises:
            InvalidQueryError: If query is invalid
        """
        try:
            response = self._make_request("POST", "/_search", json=query)
            return SearchResponse(**response)
        except IndexerError as e:
            if "invalid query" in str(e).lower():
                raise InvalidQueryError(f"Invalid query: {e}")
            raise
            
    def bulk_index(self, documents: List[Union[Document, Dict]]) -> List[Dict]:
        """Bulk index documents.
        
        Args:
            documents: List of documents to index
            
        Returns:
            List of responses from indexer
        """
        docs = []
        for doc in documents:
            if isinstance(doc, Document):
                docs.append(doc.model_dump())
            else:
                docs.append(doc)
                
        return self._make_request("POST", "/_bulk", json={"documents": docs})
