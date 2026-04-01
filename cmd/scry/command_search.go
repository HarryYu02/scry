package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HarryYu02/scry/internal/indexer"
	bolt "go.etcd.io/bbolt"
)

func commandSearch(config *Config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("search expects a source and a query")
	}
	query := strings.Join(args[1:], " ")
	// TODO: add -n flag
	numResult := 10

	indexFileName := fmt.Sprintf("%s%s", args[0], ".db")
	indexPath := filepath.Join(config.Root, "index", indexFileName)
	db, err := bolt.Open(indexPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	boltStore := BoltStore{
		db: db,
	}
	ids, err := indexer.Search(&boltStore, query, numResult)
	if err != nil {
		return err
	}

	fmt.Printf("\nSearch query: %s\n\n", query)
	fmt.Println("Results:")
	for i, id := range ids {
		title, err := boltStore.GetTitle(id)
		if err != nil {
			return err
		}
		fmt.Printf("%0.2d: %s\n", i+1, title)
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
	if choice < 0 || choice > len(ids) {
		return fmt.Errorf("choice out of bound")
	}
	if choice == 0 {
		fmt.Printf("Cancel search\n")
		return nil
	}

	selectedID := ids[choice-1]
	var pos Position
	boltStore.db.View(func(tx *bolt.Tx) error {
		posBucket := tx.Bucket([]byte("Position"))
		if posBucket == nil {
			return fmt.Errorf("posBucket not found in index")
		}
		posBytes := posBucket.Get([]byte(selectedID))
		if posBytes == nil {
			return fmt.Errorf("cannot find %s in posBucket", selectedID)
		}

		err := json.Unmarshal(posBytes, &pos)
		if err != nil {
			return err
		}

		return nil
	})

	dataFileName := fmt.Sprintf("%s%s", args[0], ".jsonl")
	dataPath := filepath.Join(config.Root, "data", dataFileName)
	dataFile, err := os.Open(dataPath)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	_, err = dataFile.Seek(int64(pos.Offset), io.SeekStart)
	if err != nil {
		return err
	}
	docBytes := make([]byte, pos.Len)
	n, err := dataFile.Read(docBytes)
	if err != nil && err != io.EOF {
		return err
	}
	docBytes = docBytes[:n]

	var doc indexer.Document
	err = json.Unmarshal(docBytes, &doc)
	if err != nil {
		return err
	}

	err = render(doc.Content)
	if err != nil {
		return err
	}
	return nil
}
