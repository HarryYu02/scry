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
type TermFreqIndex = map[string]TermFreq

func Index(source []Document) (TermFreqIndex, error) {
	tokenFreqIndex := make(TermFreqIndex)
	for _, doc := range source {
		tokenFreqMap := make(TermFreq)
		lexer := newLexer(doc.Content)
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

func calcInverseDocFreq(term string, index *TermFreqIndex) float64 {
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

	// TODO: move this step to index
	// wordFreq is the map of all unique words (in all the documents) and their total frequency
	wordFreq := make(map[string]int)
	stemmedIndex := make(TermFreqIndex)
	for id, termFreq := range *index {
		for term, freq := range termFreq {
			wordFreq[term] += freq
			stemmed, err := stem(term)
			if err != nil {
				return nil, err
			}
			if _, ok := stemmedIndex[id]; !ok {
				stemmedIndex[id] = make(TermFreq)
			}
			stemmedIndex[id][stemmed] += freq
		}
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
		_, ok := wordFreq[term]
		if ok {
			stemmedTerm, err := stem(term)
			idf := calcInverseDocFreq(stemmedTerm, &stemmedIndex)
			if err != nil {
				return nil, err
			}
			for id, termFreq := range stemmedIndex {
				tf := calcTermFreq(stemmedTerm, &termFreq)
				tfidf := tf * idf
				tfidfMap[id] += tfidf
			}
			i++
		} else {
			closestWord, err := findClosestTerm(&wordFreq, term)
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
