package indexer


import (
	"fmt"

	"github.com/HarryYu02/scry/internal/parser"
)


type Document = parser.Document

func Search(source []Document, query string, count int) ([]Document, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("Search expects len(source) > 0")
	}

	lexer := parser.NewLexer(source[0].Content)
	for lexer.Cursor < len(lexer.Content) {
		token := lexer.NextToken()
		if len(token) == 0 {
			break
		}
		fmt.Printf("token: %v (%v)\n", token, []byte(token))
	}

	return nil, nil
}
