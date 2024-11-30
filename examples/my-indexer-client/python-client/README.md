# My Indexer Python Client

A Python client library for interacting with the My Indexer service.

## Installation

```bash
pip install my-indexer-client
```

## Quick Start

```python
from my_indexer.client import IndexerClient
from my_indexer.models import Document

# Initialize client
client = IndexerClient("http://localhost:8080")

# Index a document
doc = Document(fields={
    "title": "Test Document",
    "content": "This is a test document"
})
response = client.index_document(doc)

# Search for documents
query = {
    "query": {
        "match": {
            "content": "test"
        }
    }
}
results = client.search(query)

# Print results
for hit in results.hits:
    print(f"{hit._source.fields['title']} (score: {hit._score})")
```

## Features

- Document indexing and retrieval
- Full-text search with query DSL support
- Bulk document operations
- Type-safe with Pydantic models
- Comprehensive error handling

## API Reference

### IndexerClient

The main client class for interacting with the indexer service.

#### Methods

- `index_document(document: Union[Document, Dict]) -> Dict`
  Index a single document

- `get_document(doc_id: int) -> Document`
  Retrieve a document by ID

- `search(query: Dict) -> SearchResponse`
  Search for documents using query DSL

- `bulk_index(documents: List[Union[Document, Dict]]) -> List[Dict]`
  Index multiple documents in bulk

### Models

- `Document`: Represents a document in the index
- `SearchResponse`: Represents a search response
- `SearchHit`: Represents a single search result

### Exceptions

- `IndexerError`: Base exception class
- `DocumentNotFoundError`: Raised when a document is not found
- `IndexerConnectionError`: Raised when connection fails
- `InvalidQueryError`: Raised when search query is invalid

## Development

1. Clone the repository
2. Install dependencies: `pip install -r requirements.txt`
3. Run tests: `pytest tests/`

## License

MIT
