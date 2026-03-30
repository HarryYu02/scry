package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

const ROOT_PATH = "~/.local/share/scry"
const MAX_BUF_SIZE = 10 * 1024 * 1024 // 10MB

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal("could not create CPU profile: ", err)
        }
        defer f.Close() // error handling omitted for example
        if err := pprof.StartCPUProfile(f); err != nil {
            log.Fatal("could not start CPU profile: ", err)
        }
        defer pprof.StopCPUProfile()
    }

	if len(flag.Args()) < 2 {
		fmt.Printf("ERROR: command not found\n")
		os.Exit(1)
	}
	cmd := flag.Args()[0]
	args := flag.Args()[1:]

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
		Root:       root,
		MaxBufSize: MAX_BUF_SIZE,
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

	if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        defer f.Close() // error handling omitted for example
        runtime.GC() // get up-to-date statistics
        // Lookup("allocs") creates a profile similar to go test -memprofile.
        // Alternatively, use Lookup("heap") for a profile
        // that has inuse_space as the default index.
        if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
    }
}
