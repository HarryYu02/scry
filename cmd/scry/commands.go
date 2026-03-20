package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/HarryYu02/scry/internal/indexer"
)

type Config struct {
	Root       string
	MaxBufSize int
}

type Command struct {
	name        string
	description string
	usage       string
	callback    func(*Config, []string) error
}

func readDocs(config *Config, path string) ([]indexer.Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sourceDocs := make([]indexer.Document, 0)
	fileScanner := bufio.NewScanner(file)

	buf := make([]byte, bufio.MaxScanTokenSize)
	fileScanner.Buffer(buf, config.MaxBufSize)

	for fileScanner.Scan() {
		docBytes := fileScanner.Bytes()
		var docContent indexer.Document
		err := json.Unmarshal(docBytes, &docContent)
		if err != nil {
			return nil, err
		}
		sourceDocs = append(sourceDocs, docContent)
	}

	if err := fileScanner.Err(); err != nil {
		return nil, err
	}

	return sourceDocs, nil
}

func getCommands() map[string]Command {
	return map[string]Command{
		"help": {
			name:        "help",
			description: "Prints the list of commands",
			usage:       "scry help [subcommand]",
			callback:    commandHelp,
		},
		"index": {
			name:        "index",
			description: "Create an index from the given source",
			usage:       "scry index <source>",
			callback:    commandIndex,
		},
		"search": {
			name:        "search",
			description: "Search the query from the given source",
			usage:       "scry search <source> <query>",
			callback:    commandSearch,
		},
	}
}
