package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var dir string // path to search
var args []string
var wg sync.WaitGroup

var skip []string

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
	flag.StringVar(&dir, "p", ".", "path")
	flag.Parse()
	// Validate flags.
	if !isDir(dir) {
		log.Fatal(dir, "is not a valid path.")
	}
	args = flag.Args()
	if len(args) == 0 {
		log.Fatal("no arguments passed")
	}
	filenames = make(chan rec)

	// Read rc file if available.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rc, err := ioutil.ReadFile(path.Join(cwd, ".ffrc"))
	skip = []string{}
	if err != nil {
		return
	}
	for _, bad := range strings.Split(string(rc), "\n") {
		if bad != "" {
			skip = append(skip, bad)
		}
	}
}

func check() {
	for r := range filenames {
		print := true
		for _, arg := range args {
			if !strings.Contains(strings.ToLower(r.filename), strings.ToLower(arg)) {
				print = false
				break
			}

			for _, bad := range skip {
				if strings.Contains(r.path, bad) {
					print = false
					break
				}
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
	filepath.Walk(dir, walker)
	close(filenames)
	wg.Wait()
}
