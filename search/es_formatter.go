package search

import (
	"fmt"
	"time"
)

// ESResponse represents an ElasticSearch-compatible response
type ESResponse struct {
	Took     int        `json:"took"`
	TimedOut bool       `json:"timed_out"`
	Shards   ESShards   `json:"_shards"`
	Hits     ESHits     `json:"hits"`
}

// ESShards represents shard information in an ES response
type ESShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// ESHits represents the hits section of an ES response
type ESHits struct {
	Total    ESTotal   `json:"total"`
	MaxScore float64   `json:"max_score"`
	Hits     []ESHit   `json:"hits"`
}

// ESTotal represents the total hits information
type ESTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// ESHit represents a single hit in an ES response
type ESHit struct {
	Index  string                 `json:"_index"`
	ID     string                 `json:"_id"`
	Score  float64               `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// FormatESResponse formats search results into an ElasticSearch-compatible response
func FormatESResponse(results *Results, took time.Duration, index string) *ESResponse {
	hits := make([]ESHit, 0, len(results.hits))
	var maxScore float64

	for _, hit := range results.hits {
		if hit.Score > maxScore {
			maxScore = hit.Score
		}

		// Convert document fields to map
		source := make(map[string]interface{})
		for name, field := range hit.Doc.GetFields() {
			source[name] = field.Value
		}

		hits = append(hits, ESHit{
			Index:  index,
			ID:     fmt.Sprintf("%d", hit.DocID),
			Score:  hit.Score,
			Source: source,
		})
	}

	return &ESResponse{
		Took:     int(took.Milliseconds()),
		TimedOut: false,
		Shards: ESShards{
			Total:      1,
			Successful: 1,
			Skipped:    0,
			Failed:     0,
		},
		Hits: ESHits{
			Total: ESTotal{
				Value:    len(hits),
				Relation: "eq",
			},
			MaxScore: maxScore,
			Hits:     hits,
		},
	}
}
