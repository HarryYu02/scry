package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("ERROR: command not found\n")
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	config := &Config{}
	commands := getCommands()
	command, ok := commands[cmd]
	if !ok {
		fmt.Printf("ERROR: unknown command\n")
		os.Exit(1)
	}

	err := command.callback(config, args)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
