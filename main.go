package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var args []string // arguments passed by user to partially match filename
var wg sync.WaitGroup
var maxWorkers = 512

var filenames = make(chan rec)

// rec is a convenience for passing both the full path and the filename around together
type rec struct {
	path     string
	filename string
}

func check() {
	defer wg.Done()
	for r := range filenames {
		match := true
		for _, arg := range args {
			if !strings.Contains(strings.ToLower(r.filename), strings.ToLower(arg)) {
				match = false
				break
			}
		}
		if match == false {
			continue
		}
		fmt.Println(r.path)
	}
}

// walker implements filepath.WalkFunc.
func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println(err)
		return nil
	}
	filenames <- rec{path: path, filename: info.Name()}
	return nil
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("no arguments passed")
	}
	args = os.Args[1:]
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go check()
	}
	filepath.Walk(".", walker)
	close(filenames)
	wg.Wait()
	os.Stdout.Sync()
}
