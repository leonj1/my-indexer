# My Indexer Python Client

A Python client library for using My Indexer as a library for full-text search and document indexing. This client provides direct access to the indexing and search capabilities without requiring a running server.

## Purpose

This client allows you to:
- Create and manage document indices with persistent storage
- Index documents with custom fields
- Perform full-text search with Elasticsearch-like query syntax
- Store indices in custom locations with `.gob` files
- Work with the indexer directly in your Python code without HTTP overhead

## Features

- Full-text search with Elasticsearch-like query syntax
- Document indexing with custom fields
- Persistent storage with custom filenames
- Built-in text analysis with StandardAnalyzer
- Type-safe document handling
- Zero external dependencies (no server required)
- Custom index file locations

## Installation

Since this package is not yet published to PyPI, install it directly from the source:

```bash
# Clone the repository
git clone https://github.com/yourusername/my-indexer.git
cd my-indexer

# Install the package in development mode
pip install -e examples/my-indexer-client/python-client
```

## Quick Start

```python
from my_indexer.analysis import StandardAnalyzer
from my_indexer.index import Index
from my_indexer.document import Document
from my_indexer.storage import IndexStorage

# Initialize with custom storage location
index_filename = "my_index.gob"
analyzer = StandardAnalyzer()
storage = IndexStorage(index_filename)
index = Index(analyzer, storage=storage)

# Index a document
doc = Document(fields={
    "title": "Test Document",
    "content": "This is a test document about Python"
})
doc_id = index.add_document(doc)

# Search for documents
query = {
    "match": {
        "content": "python"
    }
}
results = index.search(query)

# Print results
for doc_id, score in results:
    doc = index.get_document(doc_id)
    print(f"{doc.fields['title']} (score: {score})")

# Save index to disk
index.save()
```

## Usage Guide

### Document Management

#### Creating Documents
```python
doc = Document(fields={
    "title": "My Document",
    "content": "Document content",
    "tags": "python, search, example"
})
```

#### Indexing Documents
```python
# Single document
doc_id = index.add_document(doc)

# Multiple documents
docs = [doc1, doc2, doc3]
for doc in docs:
    index.add_document(doc)
```

#### Retrieving Documents
```python
# By ID
doc = index.get_document(1)

# All documents
all_docs = index.get_all_documents()
```

### Storage Management

#### Custom Storage Location
```python
# Create index with custom filename
storage = IndexStorage("custom_path/index.gob")
index = Index(analyzer, storage=storage)

# Save changes
index.save()

# Load existing index
loaded_storage = IndexStorage("custom_path/index.gob")
loaded_index = Index(analyzer, storage=loaded_storage)
```

### Running the Example

The repository includes a complete example that demonstrates all features:

```bash
# Make sure you're in the repository root
cd my-indexer

# Run the example
python examples/my-indexer-client/python-client/examples/basic_usage.py
```

The example will:
1. Create an index with a custom filename
2. Add several test documents
3. Perform search operations
4. Save the index to disk
5. Load the index from disk to verify persistence

## Best Practices

1. **Storage Management**
   - Always call `index.save()` after making changes
   - Use descriptive filenames for different indices
   - Store indices in appropriate locations (e.g., data directory)

2. **Document Structure**
   - Keep field names consistent across documents
   - Use meaningful field names
   - Consider field content when designing searches

3. **Error Handling**
   ```python
   try:
       doc = index.get_document(999)
   except KeyError:
       print("Document not found")

   try:
       storage = IndexStorage("invalid/path/index.gob")
   except ValueError:
       print("Invalid storage path")
   ```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Write tests for new features
4. Submit a pull request

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the [Elasticsearch Python Client](https://github.com/elastic/elasticsearch-py)
- Built with [FastAPI](https://fastapi.tiangolo.com/)
