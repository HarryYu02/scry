package main

import (
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
		"list": {
			name:        "list",
			description: "List out all created indexes",
			usage:       "scry list",
			callback:    commandList,
		},
		"remove": {
			name:        "remove",
			description: "Remove an existing index",
			usage:       "scry remove <source>",
			callback:    commandRemove,
		},
	}
}
