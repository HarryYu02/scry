package indexer

import (
	"fmt"
	"math"
)

type WordFreq interface {
	GetWords() ([]string, error)
	GetWordCount(word string) (int, error)
}

func calcLevenshteinDist(source, target string) int {
	dpRowNum := len(source) + 1
	dpColNum := len(target) + 1
	dp := make([][]int, dpRowNum)
	for i := range dpRowNum {
		row := make([]int, dpColNum)
		dp[i] = row
	}

	for i := 1; i < dpRowNum; i++ {
		dp[i][0] = i
	}
	for j := 1; j < dpColNum; j++ {
		dp[0][j] = j
	}

	for j := 1; j < dpColNum; j++ {
		for i := 1; i < dpRowNum; i++ {
			cost := 0
			if source[i-1] != target[j-1] {
				cost = 1
			}

			dp[i][j] = min(
				dp[i-1][j] + 1,
				dp[i][j-1] + 1,
				dp[i-1][j-1] + cost,
			)
		}
	}

	return dp[len(source)][len(target)]
}

func findClosestTerm(wordFreq WordFreq, term string) (string, error) {
	if len(term) == 0 {
		return "", fmt.Errorf("findClosestTerm() expects term to be non-empty")
	}
	closestWord := ""
	minLevDist := math.MaxInt
	closestWordFreq := 0
	words, err := wordFreq.GetWords()
	if err != nil {
		return "", err
	}
	for _, word := range words {
		freq, err := wordFreq.GetWordCount(word)
		if err != nil {
			return "", err
		}
		levDist := calcLevenshteinDist(word, term)
		if levDist < minLevDist {
			closestWord = word
			minLevDist = levDist
			closestWordFreq = freq
		} else if levDist == minLevDist && freq > closestWordFreq {
			// tie breaker
			closestWord = word
			minLevDist = levDist
			closestWordFreq = freq
		}
	}
	if closestWord == "" {
		return "", fmt.Errorf("findClosestTerm() cannot find a closest term")
	}
	return closestWord, nil
}

