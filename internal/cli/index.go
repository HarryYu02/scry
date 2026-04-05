package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/HarryYu02/scry/internal/indexer"
	"github.com/HarryYu02/scry/internal/store"
	bolt "go.etcd.io/bbolt"
)

func readDocs(config *Config, dataFile *os.File) ([]DocumentWithPos, error) {
	sourceDocs := make([]DocumentWithPos, 0)
	fileScanner := bufio.NewScanner(dataFile)

	buf := make([]byte, bufio.MaxScanTokenSize)
	fileScanner.Buffer(buf, config.MaxBufSize)

	cursor := 0
	for fileScanner.Scan() {
		docBytes := fileScanner.Bytes()
		var docContent indexer.Document
		err := json.Unmarshal(docBytes, &docContent)
		if err != nil {
			return nil, err
		}

		sourceDocs = append(sourceDocs, DocumentWithPos{
			document:   docContent,
			byteOffset: cursor,
			length:     len(docBytes) + 1,
		})
		cursor += len(docBytes) + 1
	}

	if err := fileScanner.Err(); err != nil {
		return nil, err
	}

	return sourceDocs, nil
}


func commandIndex(config *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("index expects a source")
	}
	source := args[0]

	dataFile, err := openDataFile(config, source)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	sourceDocs, err := readDocs(config, dataFile)
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

	db, err := openIndex(config, source, true)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		posBucket, err := store.OpenOrCreateBucket(tx, "Position")
		if err != nil {
			return err
		}
		for _, docWithPos := range sourceDocs {
			posBytes, err := json.Marshal(store.Position{
				Offset: docWithPos.byteOffset,
				Len:    docWithPos.length,
			})
			if err != nil {
				return err
			}
			err = posBucket.Put([]byte(docWithPos.document.ID), posBytes)
			if err != nil {
				return err
			}
		}

		titleBucket, err := store.OpenOrCreateBucket(tx, "Title")
		if err != nil {
			return err
		}
		for id, title := range index.IDTitleMap {
			err := titleBucket.Put([]byte(id), []byte(title))
			if err != nil {
				return err
			}
		}

		wordCountBucket, err := store.OpenOrCreateBucket(tx, "WordCount")
		if err != nil {
			return err
		}
		for term, freq := range index.AllTFMap {
			freqStr := strconv.Itoa(freq)
			err := wordCountBucket.Put([]byte(term), []byte(freqStr))
			if err != nil {
				return err
			}
		}

		IDTFBucket, err := store.OpenOrCreateBucket(tx, "IDTF")
		if err != nil {
			return err
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
