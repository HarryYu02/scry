package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HarryYu02/scry/internal/indexer"
)

type Config struct {
	Root string
}

type Command struct {
	name        string
	description string
	usage       string
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

func commandIndex(config *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("index expects a source")
	}
	source := fmt.Sprintf("%s%s", args[0], ".jsonl")
	sourcePath := filepath.Join(config.Root, "data", source)

	sourceDocs, err := readDocs(sourcePath)
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

func commandSearch(config *Config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("search expects a source and a query")
	}
	source := fmt.Sprintf("%s%s", args[0], ".jsonl")
	sourcePath := filepath.Join(config.Root, "data", source)
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

	fmt.Println("")
	fmt.Printf("Search query: %s\n", query)
	fmt.Println("")
	fmt.Println("Results:")
	for i, doc := range docs {
		fmt.Printf("%0.2d: %s\n", i+1, doc.Title)
	}
	fmt.Printf("\nSelect by typing the number 1-%d (0 to cancel)\n> ", numResult)
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return err
	}
	choice, err := strconv.Atoi(input)
	if err != nil {
		return err
	}
	if choice < 0 || choice > numResult {
		return fmt.Errorf("invalid choice")
	}
	if choice == 0 {
		fmt.Printf("Cancel search\n")
		return nil
	}

	selectedDoc := docs[choice-1]
	fmt.Printf("\n\n%s\n%s\n\n", selectedDoc.Title, selectedDoc.URL)
	fmt.Printf("%s", selectedDoc.Content)
	return nil
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
