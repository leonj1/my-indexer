package elastic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple prefix query",
			input:    `{"query": {"prefix": {"title": "test"}}}`,
			expected: `{"prefix":{"title":{"value":"test"}}}`,
			wantErr:  false,
		},
		{
			name:     "structured prefix query",
			input:    `{"query": {"prefix": {"title": {"value": "test"}}}}`,
			expected: `{"prefix":{"title":{"value":"test"}}}`,
			wantErr:  false,
		},
		{
			name:    "invalid prefix query - multiple fields",
			input:   `{"query": {"prefix": {"title": "test", "body": "test"}}}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := ParseQuery([]byte(tt.input))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			result, err := json.Marshal(query)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}
