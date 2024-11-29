# My-Indexer: A High-Performance Full-Text Search Engine in Go

![Robot Indexer](robot-indexer.jpg)

My-Indexer is a lightweight, concurrent full-text search engine written in Go. It provides fast and efficient text indexing and searching capabilities with support for complex queries, making it ideal for applications that need embedded search functionality.

## Features

- **Full-Text Search**: Support for term, phrase, and field-specific queries
- **Complex Query Support**: Boolean operators (AND, OR) for advanced search combinations
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
| Concurrency | Native Go concurrency with multiple readers/writers | Single writer, multiple readers |
| Query Types | Custom query language with field-specific searches | SQL-based queries |
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
)

func main() {
    // Create a new index
    idx := index.NewIndex()

    // Add a document
    doc := index.NewDocument()
    doc.AddField("title", "Go Programming")
    doc.AddField("content", "Go is a statically typed, compiled programming language.")
    idx.AddDocument(doc)

    // Search the index
    query := "programming language"
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

## Query Syntax

My-Indexer supports a rich query syntax:

- **Term Query**: `programming`
- **Phrase Query**: `"go programming"`
- **Field Query**: `title:golang`
- **Boolean Queries**:
  - AND: `go AND programming`
  - OR: `golang OR rust`

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
├── document/      # Document representation
├── index/         # Core indexing functionality
├── query/         # Query parsing and execution
├── search/        # Search implementation
├── storage/       # Storage management
└── txlog/         # Transaction logging
```

## Performance

My-Indexer is designed for high performance:

- Concurrent read/write operations
- Optimized in-memory index structure
- Efficient query execution
- Low memory footprint

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Lucene's architecture and SQLite FTS
- Built with Go's excellent concurrency primitives
- Special thanks to all contributors

## Contact

If you have any questions or suggestions, please open an issue on GitHub.
