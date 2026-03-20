package indexer

import (
	"strings"
	"unicode"
)

type lexer struct {
	Content string
	Cursor  int
}

func newLexer(content string) *lexer {
	preparedContent := strings.TrimSpace(content)
	return &lexer{
		Content: preparedContent,
		Cursor:  0,
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

func (l *lexer) NextToken() string {
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
