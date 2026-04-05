package cli

import (
	"fmt"
	"strings"
)

func genZshComp(config *Config) string {
	commands := "\n"
	for _, key := range sortStringMap(config.Commands) {
		command := config.Commands[key]
		commands += fmt.Sprintf("                '%s:%s'\n", command.Name, command.Description)
	}

	searchFlags := "\n"
	for _, flagName := range sortStringMap(config.Commands["search"].Flags) {
		flagDesc := config.Commands["search"].Flags[flagName]
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

func genBashComp(config *Config) string {
	commands := strings.Join(sortStringMap(config.Commands), " ")
	searchFlags := sortStringMap(config.Commands["search"].Flags)
	for i, flag := range searchFlags {
		searchFlags[i] = "--" + flag
	}
	searchFlagsStr := strings.Join(searchFlags, " ")

	completion := `
_scry_completions() {
	local cur prev words cword

    words=("${COMP_WORDS[@]}")
    cword=$COMP_CWORD
    cur="${words[cword]}"
    prev="${words[cword-1]}"

	if [[ $cword -eq 1 ]]; then
        local opts="` + commands + `"
        COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
        return 0
    fi

	local cmd="${words[1]}"

    case "${cmd}" in
        search)
            local search_opts="` + searchFlagsStr + ` $(scry list 2>/dev/null)"
            COMPREPLY=( $(compgen -W "${search_opts}" -- "${cur}") )
            ;;
        open)
            local sources=$(scry list 2>/dev/null)
            COMPREPLY=( $(compgen -W "${sources}" -- "${cur}") )
            ;;
        remove)
            local sources=$(scry list 2>/dev/null)
            COMPREPLY=( $(compgen -W "${sources}" -- "${cur}") )
            ;;
    esac
}

complete -F _scry_completions scry
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
	case "bash":
		fmt.Print(genBashComp(config))
	default:
		return fmt.Errorf("unrecognized shell")
	}
	return nil
}
