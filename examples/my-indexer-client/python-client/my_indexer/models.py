from typing import Dict, List, Optional
from pydantic import BaseModel, Field

class Document(BaseModel):
    """Represents a document in the index."""
    fields: Dict[str, str] = Field(..., description="Document fields")
    id: Optional[int] = Field(None, description="Document ID")

class SearchShards(BaseModel):
    """Represents shard information in search response."""
    total: int = Field(0, description="Total number of shards")
    successful: int = Field(0, description="Number of successful shards")
    failed: int = Field(0, description="Number of failed shards")

class SearchHit(BaseModel):
    """Represents a single search hit."""
    _id: str = Field(..., description="Document ID")
    _source: Document = Field(..., description="Document source")
    _score: float = Field(..., description="Search score")

class SearchResponse(BaseModel):
    """Represents a search response."""
    took: int = Field(..., description="Time taken in milliseconds")
    timed_out: bool = Field(False, description="Whether search timed out")
    _shards: SearchShards = Field(..., description="Shard information")
    hits: List[SearchHit] = Field(default_factory=list, description="Search hits")
    total: int = Field(0, description="Total number of hits")
