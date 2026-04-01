package main

import (
	"fmt"
)

func padRight(s string, c byte, l int) string {
	if len(s) >= l {
		return s
	}
	output := make([]byte, l)
	for i := range len(s) {
		output[i] = s[i]
	}
	for i := len(s); i < l; i++ {
		output[i] = c
	}
	return string(output)
}

func commandHelp(config *Config, args []string) error {
	subCmd := ""
	if len(args) > 0 {
		subCmd = args[0]
		subCommand, ok := config.Commands[subCmd]
		if !ok {
			return fmt.Errorf("sub-command not recognized")
		}
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Println("\t", subCommand.usage)
		fmt.Println("")
		fmt.Println("The flags are:")
		fmt.Println("")
		lenLongestFlagName := 0
		for flagName := range subCommand.flags {
			lenLongestFlagName = max(lenLongestFlagName, len(flagName))
		}
		for flagName, flagDesc := range subCommand.flags {
			paddedName := padRight(flagName, ' ', lenLongestFlagName+1)
			fmt.Printf("\t--%s-- %s\n", paddedName, flagDesc)
		}
		fmt.Println("")
	} else {
		fmt.Println("")
		fmt.Println("Scry is a modular, offline-first, terminal-native search engine.")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Println("\tscry <command> [arguments]")
		fmt.Println("")
		fmt.Println("The commands are:")
		fmt.Println("")
		lenLongestCommand := 0
		for command := range config.Commands {
			lenLongestCommand = max(lenLongestCommand, len(command))
		}
		for _, command := range config.Commands {
			paddedCommand := padRight(command.name, ' ', lenLongestCommand+1)
			fmt.Printf("\t%s\t%s\n", paddedCommand, command.description)
		}
		fmt.Println("")
		fmt.Println("Use \"scry help <command>\" for more information about a command.")
		fmt.Println("")
	}
	return nil
}
