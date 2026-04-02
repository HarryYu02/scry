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

type TermFreqIndex struct {
	StemIDTFMap map[string]map[string]float64
	AllTFMap    map[string]int
	IDTitleMap  map[string]string
}

type DocumentStore interface {
	GetDocument() (Document, error)
}

type IndexStore interface {
	GetIDTFMap(string) (map[string]float64, error)
	GetTitle(string) (string, error)
	WordFreq
}

func Index[T DocumentStore](docsStore []T) (TermFreqIndex, error) {
	index := TermFreqIndex{
		StemIDTFMap: make(map[string]map[string]float64),
		AllTFMap:    make(map[string]int),
		IDTitleMap:  make(map[string]string),
	}

	// total num of docs
	numDocs := float64(len(docsStore))
	// stem -> num of documents it appears
	stemDocFreqMap := make(map[string]int)
	// doc id -> stem -> count
	rawTFMap := make(map[string]map[string]int)
	// doc id -> total word count in doc
	wordCountMap := make(map[string]int)

	for _, docStore := range docsStore {
		doc, err := docStore.GetDocument()
		if err != nil {
			return TermFreqIndex{}, err
		}
		lexer := newLexer(doc.Content)
		index.IDTitleMap[doc.ID] = doc.Title

		for lexer.Cursor < len(lexer.Content) {
			token := strings.ToUpper(lexer.NextToken())
			if len(token) == 0 {
				break
			}
			if isStopWord(token) {
				continue
			}
			index.AllTFMap[token]++

			stemmedToken, err := stem(token)
			if err != nil {
				return TermFreqIndex{}, err
			}
			if _, ok := rawTFMap[doc.ID]; !ok {
				rawTFMap[doc.ID] = make(map[string]int)
			}
			rawTFMap[doc.ID][stemmedToken]++
			wordCountMap[doc.ID]++
		}

		for term := range rawTFMap {
			stemDocFreqMap[term]++
		}
	}

	for id, tfMap := range rawTFMap {
		for term := range tfMap {
			token := strings.ToUpper(term)
			if len(token) == 0 {
				break
			}
			if isStopWord(token) {
				continue
			}
			stemmedToken, err := stem(token)
			if err != nil {
				return TermFreqIndex{}, err
			}

			// calc TF
			termCount := rawTFMap[id][stemmedToken]
			totalTerms := wordCountMap[id]

			tf := 0.0
			if totalTerms != 0 {
				tf = float64(termCount) / float64(totalTerms)
			}

			// calc IDF
			numDocsHaveTerm := stemDocFreqMap[stemmedToken]
			idf := math.Log10(float64(numDocs)/float64(1+numDocsHaveTerm)) + 1

			// calc TF-IDF and save to index
			tfidf := tf * idf
			if _, ok := index.StemIDTFMap[stemmedToken]; !ok {
				index.StemIDTFMap[stemmedToken] = make(map[string]float64)
			}
			index.StemIDTFMap[stemmedToken][id] = tfidf
		}
	}

	return index, nil
}

func Search(index IndexStore, query string, count int) ([]string, error) {
	if count == 0 {
		return nil, fmt.Errorf("Search expects count > 0")
	}

	terms := strings.Fields(strings.ToUpper(query))
	if len(terms) == 0 {
		return nil, fmt.Errorf("Search expects query contains non-whitespace character")
	}
	correctedTerms := make([]string, len(terms))
	_ = copy(correctedTerms, terms)

	// doc id -> tfidf of entire query
	tfidfMap := make(map[string]float64)
	i := 0
	for i < len(correctedTerms) {
		term := correctedTerms[i]
		if isStopWord(term) {
			i++
			continue
		}
		wordCount, err := index.GetWordCount(term)
		if err != nil {
			return nil, err
		}
		if wordCount > 0 {
			stemmedTerm, err := stem(term)
			if err != nil {
				return nil, err
			}
			idTFMap, err := index.GetIDTFMap(stemmedTerm)
			for id, tfidf := range idTFMap {
				tfidfMap[id] += tfidf
			}
			i++
		} else {
			closestWord, err := findClosestTerm(index, term)
			if err != nil {
				return nil, err
			}
			correctedTerms[i] = closestWord
		}
	}

	ids := make([]string, len(tfidfMap))
	idx := 0
	for id := range tfidfMap {
		ids[idx] = id
		idx++

		title, err := index.GetTitle(id)
		if err == nil {
			if strings.EqualFold(title, query) {
				tfidfMap[id] *= 10
			} else if strings.HasPrefix(strings.ToUpper(title), strings.ToUpper(query)) {
				tfidfMap[id] *= (float64(len(query)) / float64(len(title))) * 10.0
			} else if strings.Contains(strings.ToUpper(title), strings.ToUpper(query)) {
				tfidfMap[id] *= (float64(len(query)) / float64(len(title))) * 8.0
			}
		}
	}

	sort.SliceStable(ids, func(i, j int) bool {
		iScore := tfidfMap[ids[i]]
		jScore := tfidfMap[ids[j]]
		return iScore > jScore
	})

	numResult := min(len(ids), count)
	return ids[:numResult], nil
}
