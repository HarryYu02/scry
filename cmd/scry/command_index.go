package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/HarryYu02/scry/internal/indexer"
	bolt "go.etcd.io/bbolt"
)

func commandIndex(config *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("index expects a source")
	}
	source := fmt.Sprintf("%s%s", args[0], ".jsonl")
	sourcePath := filepath.Join(config.Root, "data", source)

	sourceDocs, err := readDocs(config, sourcePath)
	if err != nil {
		return err
	}

	index, err := indexer.Index(sourceDocs)
	if err != nil {
		return err
	}

	indexDir := filepath.Join(config.Root, "index")
	if dirStat, err := os.Stat(indexDir); err != nil || !dirStat.IsDir() {
		err := os.MkdirAll(indexDir, 0750)
		if err != nil {
			return err
		}
	}

	indexFileName := fmt.Sprintf("%s%s", args[0], ".db")
	indexPath := filepath.Join(indexDir, indexFileName)
	db, err := bolt.Open(indexPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		titleBucket := tx.Bucket([]byte("Title"))
		if titleBucket == nil {
			titleBucket, err = tx.CreateBucket([]byte("Title"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}

		for id, title := range index.IDTitleMap {
			err := titleBucket.Put([]byte(id), []byte(title))
			if err != nil {
				return err
			}
		}

		wordCountBucket := tx.Bucket([]byte("WordCount"))
		if wordCountBucket == nil {
			wordCountBucket, err = tx.CreateBucket([]byte("WordCount"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		for term, freq := range index.AllTFMap {
			freqStr := strconv.Itoa(freq)
			err := wordCountBucket.Put([]byte(term), []byte(freqStr))
			if err != nil {
				return err
			}
		}

		IDTFBucket := tx.Bucket([]byte("IDTF"))
		if IDTFBucket == nil {
			IDTFBucket, err = tx.CreateBucket([]byte("IDTF"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		for stem, idTFMap := range index.StemIDTFMap {
			mapBytes, err := json.Marshal(idTFMap)
			if err != nil {
				return err
			}
			err = IDTFBucket.Put([]byte(stem), mapBytes)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
