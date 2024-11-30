package query

import (
	"testing"
	"reflect"
)

func TestQueryParser(t *testing.T) {
	parser := NewParser("content")

	tests := []struct {
		name    string
		input   string
		want    *ParsedQuery
		wantErr bool
	}{
		{
			name:  "Simple term query",
			input: "test",
			want: &ParsedQuery{
				Type:  TermQuery,
				Field: "content",
				Terms: []string{"test"},
			},
			wantErr: false,
		},
		{
			name:  "Field query",
			input: "title:test",
			want: &ParsedQuery{
				Type:  FieldQuery,
				Field: "title",
				Terms: []string{"test"},
			},
			wantErr: false,
		},
		{
			name:  "Phrase query",
			input: "\"quick brown fox\"",
			want: &ParsedQuery{
				Type:     PhraseQuery,
				Field:    "content",
				Terms:    []string{"quick", "brown", "fox"},
				IsPhrase: true,
			},
			wantErr: false,
		},
		{
			name:  "Field phrase query",
			input: "title:\"quick brown fox\"",
			want: &ParsedQuery{
				Type:     PhraseQuery,
				Field:    "title",
				Terms:    []string{"quick", "brown", "fox"},
				IsPhrase: true,
			},
			wantErr: false,
		},
		{
			name:  "AND query",
			input: "quick AND fox",
			want: &ParsedQuery{
				Type: TermQuery,
				SubQueries: []ParsedQuery{
					{
						Type:  TermQuery,
						Field: "content",
						Terms: []string{"quick"},
					},
					{
						Type:  TermQuery,
						Field: "content",
						Terms: []string{"fox"},
					},
				},
				Operator: "AND",
			},
			wantErr: false,
		},
		{
			name:  "OR query",
			input: "quick OR fox",
			want: &ParsedQuery{
				Type: TermQuery,
				SubQueries: []ParsedQuery{
					{
						Type:  TermQuery,
						Field: "content",
						Terms: []string{"quick"},
					},
					{
						Type:  TermQuery,
						Field: "content",
						Terms: []string{"fox"},
					},
				},
				Operator: "OR",
			},
			wantErr: false,
		},
		{
			name:    "Empty query",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid field query",
			input:   "title:",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid phrase query",
			input:   "\"\"",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryValidation(t *testing.T) {
	parser := NewParser("content")

	tests := []struct {
		name    string
		query   *ParsedQuery
		wantErr bool
	}{
		{
			name: "Valid term query",
			query: &ParsedQuery{
				Type:  TermQuery,
				Field: "content",
				Terms: []string{"test"},
			},
			wantErr: false,
		},
		{
			name:    "Nil query",
			query:   nil,
			wantErr: true,
		},
		{
			name: "Empty query",
			query: &ParsedQuery{
				Type:  TermQuery,
				Field: "content",
			},
			wantErr: true,
		},
		{
			name: "Invalid phrase query",
			query: &ParsedQuery{
				Type:     PhraseQuery,
				Field:    "content",
				Terms:    []string{"test"},
				IsPhrase: true,
			},
			wantErr: true,
		},
		{
			name: "Valid AND query",
			query: &ParsedQuery{
				Type: TermQuery,
				SubQueries: []ParsedQuery{
					{
						Type:  TermQuery,
						Field: "content",
						Terms: []string{"quick"},
					},
					{
						Type:  TermQuery,
						Field: "content",
						Terms: []string{"fox"},
					},
				},
				Operator: "AND",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.Validate(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
