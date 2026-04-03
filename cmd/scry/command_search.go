package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/HarryYu02/scry/internal/indexer"
)

type Meta struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func parseChoice(input string, numChoice int) (int, error) {
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	if choice < 0 || choice > numChoice {
		return 0, fmt.Errorf("choice out of bound")
	}
	return choice, nil
}

func promptSelectID(query string, ids []string, boltStore *BoltStore, numResult int) (string, error) {
	fmt.Fprintf(os.Stderr, "\nSearch query: %s\n\n", query)
	fmt.Fprintln(os.Stderr, "Results:")
	for i, id := range ids {
		title, err := boltStore.GetTitle(id)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(os.Stderr, "%0.2d: %s\n", i+1, title)
	}
	fmt.Fprintf(os.Stderr, "\nSelect by typing the number 1-%d (0 to cancel)\n> ", numResult)

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}

	choice, err := parseChoice(input, len(ids))
	if err != nil {
		return "", err
	}

	if choice == 0 {
		fmt.Fprintf(os.Stderr, "Cancel search\n")
		return "", nil
	}
	return ids[choice-1], nil
}

func commandSearch(config *Config, args []string) error {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchCmd.Usage = func() {
		commandHelp(config, []string{"search"})
	}

	urlFlag := searchCmd.Bool("url", false, config.Commands["search"].flags["url"])
	stdoutFlag := searchCmd.Bool("stdout", false, config.Commands["search"].flags["stdout"])
	docsFlag := searchCmd.Bool("docs", false, config.Commands["search"].flags["docs"])
	metaFlag := searchCmd.Bool("meta", false, config.Commands["search"].flags["meta"])
	nFlag := searchCmd.Int("n", 10, config.Commands["search"].flags["n"])

	searchCmd.Parse(args)
	args = searchCmd.Args()

	if len(args) < 2 {
		return fmt.Errorf("search expects a source and a query")
	}
	source := args[0]
	query := strings.Join(args[1:], " ")
	numResult := *nFlag

	db, err := openIndex(config, source)
	if err != nil {
		return err
	}
	defer db.Close()

	boltStore := BoltStore{
		db: db,
	}
	ids, err := indexer.Search(&boltStore, query, numResult)
	if err != nil {
		return err
	}

	if *docsFlag {
		for _, id := range ids {
			title, err := boltStore.GetTitle(id)
			if err != nil {
				return err
			}
			if *metaFlag {
				res := Meta{
					ID:    id,
					Title: title,
				}
				resBytes, err := json.Marshal(&res)
				if err != nil {
					return err
				}
				fmt.Printf("%v\n", string(resBytes))
			} else {
				fmt.Println(id)
			}
		}
		return nil
	}

	selectedID, err := promptSelectID(query, ids, &boltStore, numResult)
	if err != nil {
		return err
	}

	doc, err := open(config, source, selectedID, &boltStore)
	if err != nil {
		return err
	}

	if *metaFlag {
		if *urlFlag {
			res := struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				URL   string `json:"url"`
			}{
				ID:    doc.ID,
				Title: doc.Title,
				URL:   doc.URL,
			}
			resBytes, err := json.Marshal(&res)
			if err != nil {
				return err
			}
			fmt.Println(string(resBytes))
		} else {
			res := Meta{
				ID:    doc.ID,
				Title: doc.Title,
			}
			resBytes, err := json.Marshal(&res)
			if err != nil {
				return err
			}
			fmt.Println(string(resBytes))
		}
		return nil
	} else if *urlFlag {
		fmt.Println(doc.URL)
		return nil
	} else if *stdoutFlag {
		fmt.Println(doc.Content)
		return nil
	}

	err = render(doc.Content, config.Pager)
	if err != nil {
		return err
	}

	return nil
}
