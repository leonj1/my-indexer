# My-Indexer: A High-Performance Full-Text Search Engine in Go

![Robot Indexer](robot-indexer.jpg)

My-Indexer is a lightweight, concurrent full-text search engine written in Go. It provides fast and efficient text indexing and searching capabilities with support for complex queries and Elasticsearch-compatible DSL, making it ideal for applications that need embedded search functionality.

## Features

- **Full-Text Search**: Support for term, match, prefix, and field-specific queries
- **Elasticsearch-Compatible DSL**: Query using familiar Elasticsearch query syntax
- **Complex Query Support**: Boolean queries with must, should, and must_not clauses
- **Range Queries**: Support for numeric and date range searches
- **Concurrent Operations**: Thread-safe design supporting multiple simultaneous readers and writers
- **Transaction Support**: ACID-compliant operations with rollback capability
- **Crash Recovery**: Built-in transaction logging for durability
- **Custom Analyzers**: Flexible text analysis with customizable filters
- **Memory Efficient**: Optimized for low memory footprint
- **Pure Go Implementation**: No external dependencies required

## Comparison with SQLite FTS

While SQLite FTS is an excellent choice for many applications, My-Indexer offers several advantages:

| Feature | My-Indexer | SQLite FTS |
|---------|------------|------------|
| Memory Usage | Optimized in-memory index | Disk-based storage |
| Query Language | Elasticsearch-compatible DSL | SQL-based queries |
| Concurrency | Native Go concurrency with multiple readers/writers | Single writer, multiple readers |
| Query Types | Rich query types (term, match, prefix, range) | Basic full-text search |
| Customization | Extensible analyzers and filters | Limited customization options |
| Integration | Native Go API | Requires SQLite bindings |
| Dependencies | Zero external dependencies | Requires SQLite library |

## Installation

```bash
go get github.com/yourusername/my-indexer
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/yourusername/my-indexer/index"
    "github.com/yourusername/my-indexer/elastic"
    "github.com/yourusername/my-indexer/storage"
)

func main() {
    // Create a new storage with custom index filename
    store, err := storage.NewIndexStorage("/some/path/folder", "index.gob")
    if err != nil {
        panic(err)
    }

    // Create a new index with the storage
    idx := index.NewIndex(store)

    // Add a document
    doc := map[string]interface{}{
        "title": "Go Programming",
        "content": "Go is a statically typed, compiled programming language.",
        "tags": []string{"golang", "programming"},
        "year": 2012,
    }
    idx.AddDocument(doc)

    // Create a bool query
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "match": map[string]interface{}{
                            "content": "programming",
                        },
                    },
                    map[string]interface{}{
                        "range": map[string]interface{}{
                            "year": map[string]interface{}{
                                "gte": 2010,
                            },
                        },
                    },
                },
            },
        },
    }

    // Search using Elasticsearch DSL
    results, err := idx.Search(query)
    if err != nil {
        panic(err)
    }

    // Print results
    for _, result := range results {
        fmt.Printf("Found document: %v\n", result)
    }
}
```

## Storage Configuration

By default, My-Indexer uses `index.gob` as the index filename. You can customize this by providing a filename when creating the storage:

```go
// Use default filename (index.gob)
store, err := storage.NewIndexStorage("/path/to/data", "")

// Use custom filename
store, err := storage.NewIndexStorage("/path/to/data", "custom_index.gob")
```

## Query Examples

My-Indexer supports Elasticsearch-compatible DSL queries:

### Term Query
```json
{
  "query": {
    "term": {
      "title": "golang"
    }
  }
}
```

### Match Query
```json
{
  "query": {
    "match": {
      "content": "go programming"
    }
  }
}
```

### Prefix Query
```json
{
  "query": {
    "prefix": {
      "title": "go"
    }
  }
}
```

### Range Query
```json
{
  "query": {
    "range": {
      "year": {
        "gte": 2010,
        "lt": 2024
      }
    }
  }
}
```

### Boolean Query
```json
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "content": "programming"
          }
        }
      ],
      "should": [
        {
          "term": {
            "tags": "golang"
          }
        }
      ],
      "must_not": [
        {
          "term": {
            "tags": "python"
          }
        }
      ]
    }
  }
}
```

## Building from Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/my-indexer.git
cd my-indexer
```

2. Build the project:
```bash
make build
```

3. Run tests:
```bash
make test
```

## Project Structure

```
my-indexer/
├── analysis/       # Text analysis and tokenization
├── elastic/        # Elasticsearch-compatible DSL implementation
├── index/         # Core indexing and search functionality
├── logger/        # Logging and monitoring
├── query/         # Query parsing and execution
└── router/        # HTTP API endpoints
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.