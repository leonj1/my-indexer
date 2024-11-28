package analysis

import (
	"strings"
	"unicode"
)

// Token represents a single token in the text
type Token struct {
	Text      string
	Position  int  // Position in the original text
	StartByte int  // Start byte offset
	EndByte   int  // End byte offset
}

// Analyzer defines the interface for text analysis
type Analyzer interface {
	Analyze(text string) []Token
}

// StandardAnalyzer implements a basic analyzer that splits on whitespace,
// converts to lowercase, and removes punctuation
type StandardAnalyzer struct{}

// NewStandardAnalyzer creates a new StandardAnalyzer
func NewStandardAnalyzer() *StandardAnalyzer {
	return &StandardAnalyzer{}
}

// Analyze performs the text analysis process:
// 1. Splits text into tokens based on whitespace
// 2. Converts tokens to lowercase
// 3. Removes punctuation
func (a *StandardAnalyzer) Analyze(text string) []Token {
	if len(strings.TrimSpace(text)) == 0 {
		return []Token{}
	}

	var tokens []Token
	position := 0
	startByte := 0

	// Split on whitespace first
	words := strings.Fields(text)

	for _, word := range words {
		// Skip empty words
		if len(word) == 0 {
			continue
		}

		// Process the word
		cleanWord := strings.Map(func(r rune) rune {
			if unicode.IsPunct(r) || unicode.IsSymbol(r) {
				return -1 // Remove punctuation and symbols
			}
			return unicode.ToLower(r)
		}, word)

		// Skip if the word became empty after cleaning
		if len(cleanWord) == 0 {
			continue
		}

		// Calculate byte offsets
		wordStartByte := startByte
		if startByte > 0 {
			// Find the start of the current word in the original text
			for i := startByte; i < len(text); i++ {
				if strings.HasPrefix(strings.ToLower(text[i:]), strings.ToLower(word)) {
					wordStartByte = i
					break
				}
			}
		}
		wordEndByte := wordStartByte + len(cleanWord)

		tokens = append(tokens, Token{
			Text:      cleanWord,
			Position:  position,
			StartByte: wordStartByte,
			EndByte:   wordEndByte,
		})

		position++
		startByte = wordStartByte + len(word)
	}

	return tokens
}

// CustomAnalyzer allows for configurable analysis with custom filters
type CustomAnalyzer struct {
	filters []TokenFilter
}

// TokenFilter defines an interface for token filtering
type TokenFilter interface {
	Filter(token string) string
}

// NewCustomAnalyzer creates a new CustomAnalyzer with the specified filters
func NewCustomAnalyzer(filters []TokenFilter) *CustomAnalyzer {
	return &CustomAnalyzer{
		filters: filters,
	}
}

// Analyze performs text analysis using the configured filters
func (a *CustomAnalyzer) Analyze(text string) []Token {
	if len(strings.TrimSpace(text)) == 0 {
		return []Token{}
	}

	var tokens []Token
	position := 0
	startByte := 0

	words := strings.Fields(text)

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		processedWord := word
		for _, filter := range a.filters {
			processedWord = filter.Filter(processedWord)
		}

		if len(processedWord) == 0 {
			continue
		}

		// Calculate byte offsets
		wordStartByte := startByte
		if startByte > 0 {
			// Find the start of the current word in the original text
			for i := startByte; i < len(text); i++ {
				if strings.HasPrefix(strings.ToLower(text[i:]), strings.ToLower(word)) {
					wordStartByte = i
					break
				}
			}
		}
		wordEndByte := wordStartByte + len(processedWord)

		tokens = append(tokens, Token{
			Text:      processedWord,
			Position:  position,
			StartByte: wordStartByte,
			EndByte:   wordEndByte,
		})

		position++
		startByte = wordStartByte + len(word)
	}

	return tokens
}
