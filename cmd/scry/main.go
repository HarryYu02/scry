package main

import (
	"bufio"
	"fmt"
	"os"
	"encoding/json"

	"github.com/HarryYu02/scry/internal/indexer"
)


func main() {
	file, err := os.Open("scraper/data/stsfandom.jsonl")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer file.Close()

	sourceDocs := make([]indexer.Document, 0)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		docBytes := fileScanner.Bytes()
		var docContent indexer.Document
		err := json.Unmarshal(docBytes, &docContent)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		sourceDocs = append(sourceDocs, docContent)
	}

	docs, err := indexer.Search(sourceDocs, "ectoplasm", 10)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	for i, doc := range docs {
		fmt.Printf("%d: %s\n", i + 1, doc.ID)
	}
}
