package indexer

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/HarryYu02/scry/internal/parser"
)

type Document = parser.Document
type TermFreq = map[string]int
type TermFreqIndex = map[string]TermFreq

func Index(source []Document) (TermFreqIndex, error) {
	tokenFreqIndex := make(TermFreqIndex)
	for _, doc := range source {
		tokenFreqMap := make(TermFreq)
		lexer := parser.NewLexer(doc.Content)
		for lexer.Cursor < len(lexer.Content) {
			token := strings.ToUpper(lexer.NextToken())
			if len(token) == 0 {
				break
			}
			tokenFreqMap[token]++
		}
		tokenFreqIndex[doc.ID] = tokenFreqMap
	}
	return tokenFreqIndex, nil
}

func calculateTermFreq(term string, termFreq *TermFreq) float64 {
	totalTerms := 0.0
	for _, freq := range *termFreq {
		totalTerms += float64(freq)
	}
	termCount := float64((*termFreq)[term])

	if totalTerms == 0 {
		return 0.0
	}
	return termCount / totalTerms
}

func calculateInverseDocFreq(term string, index *TermFreqIndex) float64 {
	numDocs := float64(len(*index))
	numDocsHaveTerm := 0.0
	for _, termFreq := range *index {
		if _, ok := termFreq[term]; ok {
			numDocsHaveTerm++
		}
	}
	if numDocsHaveTerm == 0.0 {
		return 0.0
	}
	return math.Log10(numDocs / numDocsHaveTerm)
}

func Search(source []Document, query string, count int) ([]Document, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("Search expects len(source) > 0")
	}

	index, err := Index(source)
	if err != nil {
		return nil, err
	}

	terms := strings.Fields(strings.ToUpper(query))
	if len(terms) == 0 {
		return nil, fmt.Errorf("Search expects query contains non-whitespace character")
	}

	tfidfMap := make(map[string]float64)
	for _, term := range terms {
		idf := calculateInverseDocFreq(term, &index)
		for id, termFreq := range index {
			tf := calculateTermFreq(term, &termFreq)
			tfidfMap[id] += (tf * idf)
		}
	}

	sort.SliceStable(source, func(i, j int) bool {
		return tfidfMap[source[j].ID] < tfidfMap[source[i].ID]
	})

	return source[:count], nil
}
