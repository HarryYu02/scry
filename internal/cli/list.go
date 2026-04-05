package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func commandList(config *Config, args []string) error {
	dataDirPath := filepath.Join(config.Root, "data")
	entries, err := os.ReadDir(dataDirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		before, found := strings.CutSuffix(name, ".jsonl")
		if !found {
			continue
		}
		fmt.Printf("%s\n", before)
	}

	return nil
}
