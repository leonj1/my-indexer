import pytest
from unittest.mock import Mock, patch
import requests

from my_indexer.client import IndexerClient
from my_indexer.models import Document, SearchResponse
from my_indexer.exceptions import (
    IndexerError,
    DocumentNotFoundError,
    IndexerConnectionError,
    InvalidQueryError,
)

@pytest.fixture
def client():
    return IndexerClient("http://localhost:8080")

@pytest.fixture
def mock_response():
    mock = Mock()
    mock.json.return_value = {"fields": {"title": "Test", "content": "Content"}}
    mock.status_code = 200
    return mock

def test_client_initialization():
    client = IndexerClient()
    assert client.base_url == "http://localhost:8080"
    
    client = IndexerClient("http://example.com/")
    assert client.base_url == "http://example.com"

@patch("requests.request")
def test_index_document(mock_request, client, mock_response):
    mock_request.return_value = mock_response
    
    doc = Document(fields={"title": "Test", "content": "Content"})
    response = client.index_document(doc)
    
    assert response == mock_response.json()
    mock_request.assert_called_once_with(
        method="POST",
        url="http://localhost:8080/_doc",
        json={"fields": {"title": "Test", "content": "Content"}},
        params=None,
    )

@patch("requests.request")
def test_get_document(mock_request, client, mock_response):
    mock_request.return_value = mock_response
    
    doc = client.get_document(1)
    
    assert isinstance(doc, Document)
    assert doc.fields["title"] == "Test"
    mock_request.assert_called_once_with(
        method="GET",
        url="http://localhost:8080/_doc/1",
        json=None,
        params=None,
    )

@patch("requests.request")
def test_search(mock_request, client):
    mock_response = Mock()
    mock_response.json.return_value = {
        "took": 1,
        "timed_out": False,
        "_shards": {"total": 1, "successful": 1, "failed": 0},
        "hits": [
            {
                "_id": "1",
                "_score": 1.0,
                "_source": {"fields": {"title": "Test", "content": "Content"}},
            }
        ],
        "total": 1,
    }
    mock_response.status_code = 200
    mock_request.return_value = mock_response
    
    query = {"query": {"match": {"content": "test"}}}
    response = client.search(query)
    
    assert isinstance(response, SearchResponse)
    assert len(response.hits) == 1
    assert response.hits[0]._source.fields["title"] == "Test"
    mock_request.assert_called_once_with(
        method="POST",
        url="http://localhost:8080/_search",
        json=query,
        params=None,
    )

@patch("requests.request")
def test_connection_error(mock_request, client):
    mock_request.side_effect = requests.ConnectionError()
    
    with pytest.raises(IndexerConnectionError):
        client.get_document(1)

@patch("requests.request")
def test_document_not_found(mock_request, client):
    mock_response = Mock()
    mock_response.status_code = 404
    mock_response.raise_for_status.side_effect = requests.HTTPError()
    mock_request.return_value = mock_response
    
    with pytest.raises(DocumentNotFoundError):
        client.get_document(1)

@patch("requests.request")
def test_invalid_query(mock_request, client):
    mock_response = Mock()
    mock_response.status_code = 400
    mock_response.raise_for_status.side_effect = requests.HTTPError("invalid query")
    mock_request.return_value = mock_response
    
    with pytest.raises(InvalidQueryError):
        client.search({"invalid": "query"})
