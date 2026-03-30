package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HarryYu02/scry/internal/indexer"
)

func commandSearch(config *Config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("search expects a source and a query")
	}
	source := fmt.Sprintf("%s%s", args[0], ".jsonl")
	sourcePath := filepath.Join(config.Root, "data", source)
	query := strings.Join(args[1:], " ")
	// TODO: add -n flag
	numResult := 10

	sourceDocs, err := readDocs(config, sourcePath)
	if err != nil {
		return err
	}
	idDocMap := make(map[string]indexer.Document)
	for _, doc := range sourceDocs {
		idDocMap[doc.ID] = doc
	}

	indexFileName := fmt.Sprintf("%s%s", args[0], ".gob")
	indexPath := filepath.Join(config.Root, "index", indexFileName)
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	_, err = b.Write(indexContent)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(&b)
	var index indexer.TermFreqIndex
	err = dec.Decode(&index)
	if err != nil {
		return err
	}

	ids, err := indexer.Search(&index, query, numResult)
	if err != nil {
		return err
	}

	fmt.Printf("\nSearch query: %s\n\n", query)
	fmt.Println("Results:")
	for i, id := range ids {
		fmt.Printf("%0.2d: %s\n", i+1, idDocMap[id].Title)
	}
	fmt.Printf("\nSelect by typing the number 1-%d (0 to cancel)\n> ", numResult)
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return err
	}
	choice, err := strconv.Atoi(input)
	if err != nil {
		return err
	}
	if choice < 0 || choice > numResult {
		return fmt.Errorf("invalid choice")
	}
	if choice == 0 {
		fmt.Printf("Cancel search\n")
		return nil
	}

	selectedDoc := idDocMap[ids[choice-1]]
	// TODO: add --stdout flag
	err = render(selectedDoc.Content)
	if err != nil {
		return err
	}

	return nil
}
