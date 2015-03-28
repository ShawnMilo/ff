package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var path string // path to search
var args []string

// isDir accepts a string (file path) and returns
// a boolean which indicates if the path is
// a valid directory.
func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return stat.IsDir()
}

func init() {
	// Set up command-line flags.
	flag.StringVar(&path, "p", ".", "path")
	flag.Parse()
	// Validate flags.
	if !isDir(path) {
		log.Fatal(path, "is not a valid path.")
	}
	args = flag.Args()
	if len(args) == 0 {
		log.Fatal("no arguments passed")
	}
}

// walker implements filepath.WalkFunc.
func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println(err)
	}
	for _, arg := range args {
		if !strings.Contains(strings.ToLower(info.Name()), strings.ToLower(arg)) {
			return nil
		}
	}
	fmt.Println(path)
	return nil
}

func main() {
	filepath.Walk(path, walker)
}
