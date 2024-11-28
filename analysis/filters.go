package analysis

import (
	"strings"
	"unicode"
)

// LowercaseFilter converts tokens to lowercase
type LowercaseFilter struct{}

func NewLowercaseFilter() *LowercaseFilter {
	return &LowercaseFilter{}
}

func (f *LowercaseFilter) Filter(token string) string {
	return strings.ToLower(token)
}

// PunctuationFilter removes punctuation from tokens
type PunctuationFilter struct{}

func NewPunctuationFilter() *PunctuationFilter {
	return &PunctuationFilter{}
}

func (f *PunctuationFilter) Filter(token string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return -1
		}
		return r
	}, token)
}

// TrimSpaceFilter removes leading and trailing whitespace
type TrimSpaceFilter struct{}

func NewTrimSpaceFilter() *TrimSpaceFilter {
	return &TrimSpaceFilter{}
}

func (f *TrimSpaceFilter) Filter(token string) string {
	return strings.TrimSpace(token)
}
