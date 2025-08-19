package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return walkDir(out, path, printFiles, "")
}

func walkDir(out io.Writer, path string, printFiles bool, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter entries based on printFiles flag
	var filtered []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() || printFiles {
			filtered = append(filtered, entry)
		}
	}

	for i, entry := range filtered {
		isLast := i == len(filtered)-1

		if entry.IsDir() {
			fmt.Fprintf(out, "%s%s%s\n", prefix, getPrefix(isLast), entry.Name())
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			if err := walkDir(out, filepath.Join(path, entry.Name()), printFiles, newPrefix); err != nil {
				return err
			}
		} else if printFiles {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			size := "empty"
			if info.Size() > 0 {
				size = strconv.FormatInt(info.Size(), 10) + "b"
			}
			fmt.Fprintf(out, "%s%s%s (%s)\n", prefix, getPrefix(isLast), entry.Name(), size)
		}
	}
	return nil
}

func getPrefix(isLast bool) string {
	if isLast {
		return "└── "
	}
	return "├── "
}

func main() {
	out := os.Stdout
	if len(os.Args) < 2 || len(os.Args) > 3 {
		panic("usage: go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
