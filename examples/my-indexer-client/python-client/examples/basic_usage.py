from my_indexer.analysis import StandardAnalyzer
from my_indexer.index import Index
from my_indexer.document import Document
from my_indexer.storage import IndexStorage

def main():
    # Initialize analyzer and storage with custom filename
    index_filename = "custom_index.gob"
    analyzer = StandardAnalyzer()
    storage = IndexStorage(index_filename)  # Custom filename for the index
    index = Index(analyzer, storage=storage)
    
    # Create some test documents
    documents = [
        Document(
            fields={
                "title": "First Document",
                "content": "This is the first test document about Python programming"
            }
        ),
        Document(
            fields={
                "title": "Second Document",
                "content": "This is another document about Go programming"
            }
        ),
        Document(
            fields={
                "title": "Third Document",
                "content": "This document discusses both Python and Go"
            }
        )
    ]
    
    print("1. Indexing documents...")
    for doc in documents:
        doc_id = index.add_document(doc)
        print(f"Indexed document with ID: {doc_id}")
    
    print("\n2. Retrieving a document...")
    try:
        doc = index.get_document(1)
        print(f"Retrieved document: {doc}")
    except KeyError:
        print("Document not found")
    
    print("\n3. Searching for documents...")
    # Search for documents about Python adhering to ElasticSearch query syntax
    query = {
        "match": {
            "content": "python"
        }
    }
    
    results = index.search(query)
    print(f"\nFound {len(results)} documents:")
    for doc_id, score in results:
        doc = index.get_document(doc_id)
        print(f"- {doc.fields['title']} (score: {score})")
    
    print("\n4. Getting all documents...")
    all_docs = index.get_all_documents()
    print(f"Total documents in index: {len(all_docs)}")
    for doc in all_docs:
        print(f"- {doc.fields['title']}")
    
    # Save the index to disk
    print("\n5. Saving index to disk...")
    index.save()
    print(f"Index saved to: {index.storage.index_path}")
    
    # Load the index from disk
    print("\n6. Loading index from disk...")
    loaded_storage = IndexStorage(index_filename)
    loaded_index = Index(analyzer, storage=loaded_storage)
    loaded_docs = loaded_index.get_all_documents()
    print(f"Loaded {len(loaded_docs)} documents from disk")

if __name__ == "__main__":
    main()
