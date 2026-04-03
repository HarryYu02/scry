package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/HarryYu02/scry/internal/indexer"
	bolt "go.etcd.io/bbolt"
)

func open(config *Config, source, id string, boltStore *BoltStore) (indexer.Document, error) {
	var pos Position
	boltStore.db.View(func(tx *bolt.Tx) error {
		posBucket := tx.Bucket([]byte("Position"))
		if posBucket == nil {
			return fmt.Errorf("posBucket not found in index")
		}
		posBytes := posBucket.Get([]byte(id))
		if posBytes == nil {
			return fmt.Errorf("cannot find %s in posBucket", id)
		}

		err := json.Unmarshal(posBytes, &pos)
		if err != nil {
			return err
		}

		return nil
	})

	dataFileName := fmt.Sprintf("%s%s", source, ".jsonl")
	dataPath := filepath.Join(config.Root, "data", dataFileName)
	dataFile, err := os.Open(dataPath)
	if err != nil {
		return indexer.Document{}, err
	}
	defer dataFile.Close()
	_, err = dataFile.Seek(int64(pos.Offset), io.SeekStart)
	if err != nil {
		return indexer.Document{}, err
	}
	docBytes := make([]byte, pos.Len)
	n, err := dataFile.Read(docBytes)
	if err != nil && err != io.EOF {
		return indexer.Document{}, err
	}
	docBytes = docBytes[:n]

	var doc indexer.Document
	err = json.Unmarshal(docBytes, &doc)
	if err != nil {
		return indexer.Document{}, err
	}

	return doc, nil
}

func commandOpen(config *Config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("open expects a source and a doc_id")
	}

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

	doc, err := open(config, args[0], args[1], &boltStore)
	if err != nil {
		return err
	}

	err = render(doc.Content, config.Pager)
	if err != nil {
		return err
	}

	return nil
}
