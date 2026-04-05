package cli

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
	Name        string
	Description string
	Usage       string
	Flags       map[string]string
	Callback    func(*Config, []string) error
}

type DocumentWithPos struct {
	document   indexer.Document
	byteOffset int
	length     int
}

func (d DocumentWithPos) GetDocument() (indexer.Document, error) {
	return d.document, nil
}

func GetCommands() map[string]Command {
	return map[string]Command{
		"help": {
			Name:        "help",
			Description: "Prints the list of commands",
			Usage:       "scry help [subcommand]",
			Callback:    commandHelp,
		},
		"index": {
			Name:        "index",
			Description: "Create an index from the given source",
			Usage:       "scry index <source>",
			Callback:    commandIndex,
		},
		"search": {
			Name:        "search",
			Description: "Search the query from the given source",
			Usage:       "scry search [flags] <source> <query>",
			Flags: map[string]string{
				"url":    "Return the url of the search result instead",
				"stdout": "Pipe content of search result to stdout instead",
				"n":      "Limit maximum search result, default to be 10, n=0 means unlimited",
				"docs":   "Only list all search results",
				"meta":   "Only return the metadata of result in json",
			},
			Callback: commandSearch,
		},
		"open": {
			Name:        "open",
			Description: "Open a specific document in a source",
			Usage:       "scry open <source> <doc_id>",
			Callback:    commandOpen,
		},
		"list": {
			Name:        "list",
			Description: "List out all created indexes",
			Usage:       "scry list",
			Callback:    commandList,
		},
		"remove": {
			Name:        "remove",
			Description: "Remove an existing index",
			Usage:       "scry remove <source>",
			Callback:    commandRemove,
		},
		"completion": {
			Name:        "completion",
			Description: "Generate shell completion",
			Usage:       "scry completion <shell>",
			Callback:    commandCompletion,
		},
	}
}
