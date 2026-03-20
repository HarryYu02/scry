package main

import "fmt"

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
