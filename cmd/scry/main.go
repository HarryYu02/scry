package main

import (
	"bufio"
	"fmt"
	"os"
	"encoding/json"

	"github.com/HarryYu02/scry/internal/indexer"
	"github.com/HarryYu02/scry/internal/parser"
)


func main() {
	file, err := os.Open("scraper/data/sts.jsonl")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer file.Close()

	sourceDocs := make([]parser.Document, 0)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		docBytes := fileScanner.Bytes()
		var docContent parser.Document
		err := json.Unmarshal(docBytes, &docContent)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		sourceDocs = append(sourceDocs, docContent)
	}

	docs, err := indexer.Search([]parser.Document{}, "", 3)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}
