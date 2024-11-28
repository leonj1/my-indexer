package analysis

import (
	"reflect"
	"testing"
)

func TestStandardAnalyzer(t *testing.T) {
	analyzer := NewStandardAnalyzer()

	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []Token{},
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: []Token{},
		},
		{
			name:  "Simple text",
			input: "Hello World",
			expected: []Token{
				{Text: "hello", Position: 0, StartByte: 0, EndByte: 5},
				{Text: "world", Position: 1, StartByte: 6, EndByte: 11},
			},
		},
		{
			name:  "Text with punctuation",
			input: "Hello, World!",
			expected: []Token{
				{Text: "hello", Position: 0, StartByte: 0, EndByte: 5},
				{Text: "world", Position: 1, StartByte: 7, EndByte: 12},
			},
		},
		{
			name:  "Multiple spaces",
			input: "Hello    World",
			expected: []Token{
				{Text: "hello", Position: 0, StartByte: 0, EndByte: 5},
				{Text: "world", Position: 1, StartByte: 9, EndByte: 14},
			},
		},
		{
			name:  "Mixed case",
			input: "HeLLo WoRLD",
			expected: []Token{
				{Text: "hello", Position: 0, StartByte: 0, EndByte: 5},
				{Text: "world", Position: 1, StartByte: 6, EndByte: 11},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyzer.Analyze(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Analyze() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCustomAnalyzer(t *testing.T) {
	filters := []TokenFilter{
		NewLowercaseFilter(),
		NewPunctuationFilter(),
		NewTrimSpaceFilter(),
	}
	analyzer := NewCustomAnalyzer(filters)

	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []Token{},
		},
		{
			name:  "Complex text",
			input: "Hello, World! This is a TEST.",
			expected: []Token{
				{Text: "hello", Position: 0, StartByte: 0, EndByte: 5},
				{Text: "world", Position: 1, StartByte: 7, EndByte: 12},
				{Text: "this", Position: 2, StartByte: 14, EndByte: 18},
				{Text: "is", Position: 3, StartByte: 19, EndByte: 21},
				{Text: "a", Position: 4, StartByte: 22, EndByte: 23},
				{Text: "test", Position: 5, StartByte: 24, EndByte: 28},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyzer.Analyze(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Analyze() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilters(t *testing.T) {
	tests := []struct {
		name     string
		filter   TokenFilter
		input    string
		expected string
	}{
		{
			name:     "Lowercase filter",
			filter:   NewLowercaseFilter(),
			input:    "HeLLo",
			expected: "hello",
		},
		{
			name:     "Punctuation filter",
			filter:   NewPunctuationFilter(),
			input:    "Hello, World!",
			expected: "Hello World",
		},
		{
			name:     "Trim space filter",
			filter:   NewTrimSpaceFilter(),
			input:    "  hello  ",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Filter(tt.input)
			if got != tt.expected {
				t.Errorf("Filter() = %v, want %v", got, tt.expected)
			}
		})
	}
}
