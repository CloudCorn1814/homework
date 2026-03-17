package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	content, err := os.ReadDir(path)
	if err != nil {
		return errors.New("invalid path")
	}

	entries := make([]os.DirEntry, 0)

	for _, dir := range content {
		if dir.IsDir() || printFiles {
			entries = append(entries, dir)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i, entry := range entries {

		if i == len(entries)-1 {
			fmt.Fprint(out, "└───", entry.Name())
		} else {
			fmt.Fprint(out, "├───", entry.Name())
		}

		if !entry.IsDir() {
			fileInfo, err := entry.Info()
			if err != nil {
				return errors.New("invalid file")
			}
			size := fileInfo.Size()
			if size == 0 {
				fmt.Fprint(out, " (empty)\n")
			} else {
				fmt.Fprintf(out, " (%db)\n", size)
			}
		}
	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
