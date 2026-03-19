package main

import (
	"fmt"
	"os"
	"strings"
)


const ROOT_PATH = "~/.local/share/scry"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("ERROR: command not found\n")
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("ERROR: user home dir failed to resolve\n")
		os.Exit(1)
	}
	root := ROOT_PATH
	if root[0] == '~' {
		root = strings.Replace(root, "~", homeDir, 1)
	}

	config := &Config{
		Root: root,
	}

	commands := getCommands()
	command, ok := commands[cmd]
	if !ok {
		fmt.Printf("ERROR: unknown command\n")
		os.Exit(1)
	}

	err = command.callback(config, args)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
