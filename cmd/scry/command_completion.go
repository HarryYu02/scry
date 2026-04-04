package main

import "fmt"

func genZshComp(config *Config) string {
	commands := "\n"
	for _, key := range sortStringMap(config.Commands) {
		command := config.Commands[key]
		commands += fmt.Sprintf("                '%s:%s'\n", command.name, command.description)
	}

	searchFlags := "\n"
	for _, flagName := range sortStringMap(config.Commands["search"].flags) {
		// '--docs[Only list results]' \
		flagDesc := config.Commands["search"].flags[flagName]
		searchFlags += fmt.Sprintf("                        '--%s[%s]' \\\n", flagName, flagDesc)
	}

	completion := `
_scry() {
    local line state
    _arguments -C "1: :->command" "*::arg:->args"

    case $state in
        command)
            local -a cmds
            cmds=(` +
	commands +
`            )
            _describe -t commands 'scry command' cmds
            ;;
        args)
            case $line[1] in
                search)
                    _arguments \` +
	searchFlags +
`                        '1:sources:->sources' \
                        '2:query: '

                    if [[ "$state" == "sources" ]]; then
                        local -a sources
                        sources=(${(f)"$(scry list 2>/dev/null)"})
                        _describe -t sources 'sources' sources
                    fi
                    ;;
                open)
                    _arguments \
                    '1:sources:->sources' \
                    '2:doc_id: '

                    if [[ "$state" == "sources" ]]; then
                        local -a sources
                        sources=(${(f)"$(scry list 2>/dev/null)"})
                        _describe -t sources 'indexed sources' sources
                    fi
                    ;;
                remove)
                    _arguments \
                    '1:sources:->sources' \
                    '*: :()'

                    if [[ "$state" == "sources" ]]; then
                        local -a sources
                        sources=(${(f)"$(scry list 2>/dev/null)"})
                        _describe -t sources 'indexed sources' sources
                    fi
                    ;;
                completion)
                    _arguments \
                    '1:shell:(zsh bash)' \
                    '*: :()'

                    if [[ "$state" == "sources" ]]; then
                        local -a sources
                        sources=(${(f)"$(scry list 2>/dev/null)"})
                        _describe -t sources 'indexed sources' sources
                    fi
                    ;;
            esac
            ;;
    esac
}
compdef _scry scry
`
	return completion
}

func commandCompletion(config *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("completion expects a shell")
	}
	shell := args[0]

	switch shell {
	case "zsh":
		fmt.Print(genZshComp(config))
	default:
		return fmt.Errorf("unrecognized shell")
	}
	return nil
}
