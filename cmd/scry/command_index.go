package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HarryYu02/scry/internal/indexer"
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

	var indexGob bytes.Buffer
	enc := gob.NewEncoder(&indexGob)
	err = enc.Encode(index)
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

	indexFileName := fmt.Sprintf("%s%s", args[0], ".gob")
	indexPath := filepath.Join(indexDir, indexFileName)
	file, err := os.OpenFile(indexPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(indexGob.Bytes())
	if err != nil {
		return err
	}

	return nil
}
