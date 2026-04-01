package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

const ROOT_PATH = "~/.local/share/scry"
const MAX_BUF_SIZE = 10 * 1024 * 1024 // 10MB

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("ERROR: could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("ERROR: could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
    }

	if len(flag.Args()) < 1 {
		log.Fatal("ERROR: command not provided")
	}
	cmd := flag.Args()[0]
	args := flag.Args()[1:]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("ERROR: user home dir failed to resolve: ", err)
	}
	root := ROOT_PATH
	if root[0] == '~' {
		root = strings.Replace(root, "~", homeDir, 1)
	}

	config := &Config{
		Root:       root,
		MaxBufSize: MAX_BUF_SIZE,
		Commands: getCommands(),
	}
	command, ok := config.Commands[cmd]
	if !ok {
		log.Fatal("unknown command")
	}

	err = command.callback(config, args)
	if err != nil {
		log.Fatal("ERROR: error in executing command: ", err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("ERROR: could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Fatal("ERROR: could not write memory profile: ", err)
        }
    }
}
