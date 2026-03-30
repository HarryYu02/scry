package indexer

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Document struct {
	ID      string   `json:"id"`
	Source  string   `json:"source"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	URL     string   `json:"url,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

type TermFreq = map[string]int
type TermFreqIndex struct {
	StemmedIDTFMap map[string]TermFreq
	AllTFMap       TermFreq
}

func Index(source []Document) (TermFreqIndex, error) {
	index := TermFreqIndex{
		StemmedIDTFMap: make(map[string]TermFreq),
		AllTFMap: make(TermFreq),
	}
	allTFMap := make(TermFreq)
	for _, doc := range source {
		stemmedTFMap := make(TermFreq)
		lexer := newLexer(doc.Content)
		for lexer.Cursor < len(lexer.Content) {
			token := strings.ToUpper(lexer.NextToken())
			if isStopWord(token) {
				continue
			}
			stemmedToken, err := stem(token)
			if err != nil {
				return TermFreqIndex{}, err
			}
			if len(token) == 0 {
				break
			}
			allTFMap[token]++
			stemmedTFMap[stemmedToken]++
		}
		index.StemmedIDTFMap[doc.ID] = stemmedTFMap
	}
	index.AllTFMap = allTFMap
	return index, nil
}

func calcTermFreq(term string, termFreq *TermFreq) float64 {
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

func calcInverseDocFreq(term string, index *map[string]TermFreq) float64 {
	numDocs := float64(len(*index))
	numDocsHaveTerm := 0.0
	for _, termFreq := range *index {
		if _, ok := termFreq[term]; ok {
			numDocsHaveTerm++
		}
	}
	return math.Log10(numDocs / (1 + numDocsHaveTerm)) + 1
}

func Search(source []Document, index *TermFreqIndex, query string, count int) ([]Document, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("Search expects len(source) > 0")
	}
	if count == 0 {
		return nil, fmt.Errorf("Search expects count > 0")
	}

	terms := strings.Fields(strings.ToUpper(query))
	if len(terms) == 0 {
		return nil, fmt.Errorf("Search expects query contains non-whitespace character")
	}
	correctedTerms := make([]string, len(terms))
	_ = copy(correctedTerms, terms)

	tfidfMap := make(map[string]float64)
	i := 0
	for i < len(correctedTerms) {
		term := correctedTerms[i]
		if isStopWord(term) {
			i++
			continue
		}
		if _, ok := index.AllTFMap[term]; ok {
			stemmedTerm, err := stem(term)
			idf := calcInverseDocFreq(stemmedTerm, &index.StemmedIDTFMap)
			if err != nil {
				return nil, err
			}
			for id, termFreq := range index.StemmedIDTFMap {
				tf := calcTermFreq(stemmedTerm, &termFreq)
				tfidf := tf * idf
				tfidfMap[id] += tfidf
			}
			i++
		} else {
			closestWord, err := findClosestTerm(&index.AllTFMap, term)
			if err != nil {
				return nil, err
			}
			correctedTerms[i] = closestWord
		}
	}

	sort.SliceStable(source, func(i, j int) bool {
		return tfidfMap[source[j].ID] < tfidfMap[source[i].ID]
	})

	numResult := min(len(source), count)
	// for i := range numResult {
	// 	fmt.Printf("%s: %f\n", source[i].ID, tfidfMap[source[i].ID])
	// }

	return source[:numResult], nil
}
