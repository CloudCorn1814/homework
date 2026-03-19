package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, "")
}

func printDir(out io.Writer, path string, printFiles bool, prefix string) error {
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
		var nextPrefix string

		if i == len(entries)-1 {
			fmt.Fprint(out, prefix, "└───", entry.Name())
			nextPrefix = prefix + "\t"
		} else {
			fmt.Fprint(out, prefix, "├───", entry.Name())
			nextPrefix = prefix + "│\t"
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

		if entry.IsDir() {
			fmt.Fprint(out, "\n")
			printDir(out, filepath.Join(path, entry.Name()), printFiles, nextPrefix)
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
