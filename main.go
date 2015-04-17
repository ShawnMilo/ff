package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var path string // path to search
var args []string
var wg sync.WaitGroup

var filenames chan rec

type rec struct {
	path     string
	filename string
}

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
	filenames = make(chan rec)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func check() {
	for r := range filenames {
		print := true
		for _, arg := range args {
			if !strings.Contains(strings.ToLower(r.filename), strings.ToLower(arg)) {
				print = false
				break
			}
		}
		if print {
			fmt.Println(r.path)
		}

	}
	wg.Done()
}

// walker implements filepath.WalkFunc.
func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println(err)
	}
	filenames <- rec{path: path, filename: info.Name()}
	return nil
}

func main() {
	var max = runtime.NumCPU() * 10
	for i := 0; i < max; i++ {
		wg.Add(1)
		go check()
	}
	filepath.Walk(path, walker)
	close(filenames)
	wg.Wait()
}
