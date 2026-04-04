package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func commandRemove(config *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("remove expects a source")
	}
	source := args[0]

	indexFileName := fmt.Sprintf("%s%s", source, ".db")
	indexPath := filepath.Join(config.Root, "index", indexFileName)
	err := os.Remove(indexPath)
	if err != nil {
		return err
	}

	return nil
}
