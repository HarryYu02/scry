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
	Commands   map[string]Command
	Pager      string
}

type Command struct {
	name        string
	description string
	usage       string
	flags       map[string]string
	callback    func(*Config, []string) error
}

type Position struct {
	Offset int
	Len    int
}

type DocumentWithPos struct {
	document   indexer.Document
	byteOffset int
	length     int
}

func (d DocumentWithPos) GetDocument() (indexer.Document, error) {
	return d.document, nil
}

func readDocs(config *Config, path string) ([]DocumentWithPos, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sourceDocs := make([]DocumentWithPos, 0)
	fileScanner := bufio.NewScanner(file)

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
			usage:       "scry search [flags] <source> <query>",
			flags: map[string]string{
				"url":    "Return the url of the search result instead",
				"stdout": "Pipe content of search result to stdout instead",
				"n":      "Limit maximum search result, default to be 10, n=0 means unlimited",
				"docs":   "Only list all search results",
				"meta":   "Only return the metadata of result in json",
			},
			callback: commandSearch,
		},
		"open": {
			name:        "open",
			description: "Open a specific document in a source",
			usage:       "scry open <source> <doc_id>",
			callback:    commandOpen,
		},
	}
}
