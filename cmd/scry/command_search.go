package main

import (
	"encoding/json"
	"flag"
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
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchCmd.Usage = func() {
		commandHelp(config, []string{"search"})
	}

	urlFlag := searchCmd.Bool("url", false, config.Commands["search"].flags["url"])
	stdoutFlag := searchCmd.Bool("stdout", false, config.Commands["search"].flags["stdout"])
	nFlag := searchCmd.Int("n", 10, config.Commands["search"].flags["n"])

	searchCmd.Parse(args)
	args = searchCmd.Args()

	if len(args) < 2 {
		return fmt.Errorf("search expects a source and a query")
	}
	query := strings.Join(args[1:], " ")
	numResult := *nFlag

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

	fmt.Fprintf(os.Stderr, "\nSearch query: %s\n\n", query)
	fmt.Fprintln(os.Stderr, "Results:")
	for i, id := range ids {
		title, err := boltStore.GetTitle(id)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "%0.2d: %s\n", i+1, title)
	}
	fmt.Fprintf(os.Stderr, "\nSelect by typing the number 1-%d (0 to cancel)\n> ", numResult)
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
		fmt.Fprintf(os.Stderr, "Cancel search\n")
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

	if *urlFlag {
		fmt.Print(doc.URL)
		return nil
	} else if *stdoutFlag {
		fmt.Print(doc.Content)
		return nil
	}

	err = render(doc.Content, config.Pager)
	if err != nil {
		return err
	}
	return nil
}
