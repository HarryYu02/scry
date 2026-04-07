package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/HarryYu02/scry/internal/indexer"
	"github.com/HarryYu02/scry/internal/render"
	"github.com/HarryYu02/scry/internal/store"
	bolt "go.etcd.io/bbolt"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func openIndex(config *Config, source string, createIfNotExist bool) (*bolt.DB, error) {
	indexFileName := fmt.Sprintf("%s%s", source, ".db")
	indexPath := filepath.Join(config.Root, "index", indexFileName)

	if !createIfNotExist {
		indexExists, err := fileExists(indexPath)
		if err != nil {
			return nil, err
		}
		if !indexExists {
			return nil, fmt.Errorf("index not found, please run 'scry index <source>' first")
		}
	}

	db, err := bolt.Open(indexPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openDataFile(config *Config, source string) (*os.File, error) {
	dataFileName := fmt.Sprintf("%s%s", source, ".jsonl")
	dataPath := filepath.Join(config.Root, "data", dataFileName)
	dataFile, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}
	return dataFile, nil
}

func readFrom(file *os.File, pos store.Position) ([]byte, error) {
	_, err := file.Seek(int64(pos.Offset), io.SeekStart)
	if err != nil {
		return nil, err
	}
	docBytes := make([]byte, pos.Len)
	n, err := file.Read(docBytes)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return docBytes[:n], nil
}

func unmarshalDoc(docBytes []byte) (indexer.Document, error) {
	var doc indexer.Document
	err := json.Unmarshal(docBytes, &doc)
	if err != nil {
		return indexer.Document{}, err
	}
	return doc, nil
}

func open(config *Config, source, id string, boltStore *store.BoltStore) (indexer.Document, error) {
	pos, err := boltStore.GetPosition(id)
	if err != nil {
		return indexer.Document{}, err
	}

	dataFile, err := openDataFile(config, source)
	if err != nil {
		return indexer.Document{}, err
	}
	defer dataFile.Close()

	dataStat, err := dataFile.Stat()
	if err != nil {
		return indexer.Document{}, err
	}
	dataLastMod := dataStat.ModTime()

	indexFileName := fmt.Sprintf("%s%s", source, ".db")
	indexPath := filepath.Join(config.Root, "index", indexFileName)
	indexStat, err := os.Stat(indexPath)
	if err != nil {
		return indexer.Document{}, err
	}
	indexLastMod := indexStat.ModTime()

	if dataLastMod.After(indexLastMod) {
		return indexer.Document{},
			fmt.Errorf("data might be modified since you last indexed it, please index it again.\n")
	}

	docBytes, err := readFrom(dataFile, pos)
	if err != nil {
		return indexer.Document{}, err
	}

	doc, err := unmarshalDoc(docBytes)
	if err != nil {
		return indexer.Document{}, err
	}

	return doc, nil
}

func commandOpen(config *Config, args []string) error {
	openCmd := flag.NewFlagSet("open", flag.ExitOnError)
	openCmd.Usage = func() {
		commandHelp(config, []string{"open"})
	}

	stdoutFlag := openCmd.Bool("stdout", false, config.Commands["open"].Flags["stdout"])

	openCmd.Parse(args)
	args = openCmd.Args()

	if len(args) < 2 {
		return fmt.Errorf("open expects a source and a doc_id")
	}

	db, err := openIndex(config, args[0], false)
	if err != nil {
		return err
	}
	defer db.Close()

	boltStore := store.BoltStore{
		DB: db,
	}

	doc, err := open(config, args[0], args[1], &boltStore)
	if err != nil {
		return err
	}

	if *stdoutFlag {
		fmt.Printf("%s", doc.Content)
		return nil
	}
	err = render.Render(doc.Content, config.Pager)
	if err != nil {
		return err
	}

	return nil
}
