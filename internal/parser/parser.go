package parser

import (
	// "fmt"
	"strings"
	"unicode"
)


type Document struct {
	ID      string   `json:"id"`
	Source  string   `json:"source"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	URL     string   `json:"url,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

type Lexer struct {
	Content string
	Cursor  int
}

func NewLexer(content string) *Lexer {
	preparedContent := strings.TrimSpace(content)
	return &Lexer{
		Content: preparedContent,
		Cursor: 0,
	}
}

func isAlpha(c byte) bool {
	return unicode.IsLetter(rune(c))
}

func isNumeric(c byte) bool {
	return unicode.IsNumber(rune(c))
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isNumeric(c)
}

func isWhiteSpace(c byte) bool {
	return unicode.IsSpace(rune(c))
}

func (l *Lexer)NextToken() string {
	for isWhiteSpace(l.Content[l.Cursor]) {
		l.Cursor++
	}
	start := l.Cursor
	end := start

	if end >= len(l.Content) {
		return ""
	}

	if isAlpha(l.Content[start]) {
		for end < len(l.Content) && isAlphaNumeric(l.Content[end]) {
			end++
		}
		token := l.Content[start:end]
		l.Cursor = end
		return token
	}

	if isNumeric(l.Content[start]) {
		for end < len(l.Content) && isNumeric(l.Content[end]) {
			end++
		}
		token := l.Content[start:end]
		l.Cursor = end
		return token
	}

	token := l.Content[end]
	l.Cursor++
	return string(token)
}
