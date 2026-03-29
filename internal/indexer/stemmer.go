package indexer

import (
	"strings"

	"github.com/kljensen/snowball"
)

func stem(word string) (string, error) {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil{
		return "", err
	}
	upperStemmed := strings.ToUpper(stemmed)
	return upperStemmed, nil
}
