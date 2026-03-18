package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/HarryYu02/scry/internal/indexer"
)

type Config struct {
}

type Command struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

func commandHelp(config *Config, args []string) error {
	fmt.Println("")
	fmt.Println("Scry is a modular, offline-first, terminal-native search engine.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("\tscry <command> [arguments]")
	fmt.Println("")
	fmt.Println("The commands are:")
	fmt.Println("")
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("\t%s\t%s\n", command.name, command.description)
	}
	fmt.Println("")

	return nil
}

func readDocs(path string) ([]indexer.Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sourceDocs := make([]indexer.Document, 0)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		docBytes := fileScanner.Bytes()
		var docContent indexer.Document
		err := json.Unmarshal(docBytes, &docContent)
		if err != nil {
			return nil, err
		}
		sourceDocs = append(sourceDocs, docContent)
	}

	return sourceDocs, nil
}

func commandSearch(config *Config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("search expects a source and a query")
	}
	sourcePath := args[0]
	query := strings.Join(args[1:], " ")
	numResult := 10

	sourceDocs, err := readDocs(sourcePath)
	if err != nil {
		return err
	}

	docs, err := indexer.Search(sourceDocs, query, numResult)
	if err != nil {
		return err
	}

	for i, doc := range docs {
		fmt.Printf("%d: %s\n", i+1, doc.ID)
	}
	return nil
}

func getCommands() map[string]Command {
	return map[string]Command{
		"help": {
			name:        "help",
			description: "Prints the list of commands",
			callback:    commandHelp,
		},
		"search": {
			name:        "search",
			description: "Search the query from the given source",
			callback:    commandSearch,
		},
	}
}
