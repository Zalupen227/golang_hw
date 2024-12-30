package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	items, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, item := range items {
		if !item.IsDir() && !printFiles {
			continue

		}

		isLast := i == len(items)-1

		printItem(out, item, path, printFiles, "", isLast)
	}
	return nil
}

func printItem(out io.Writer, item os.DirEntry, path string, printFiles bool, prefix string, isLast bool) {
	var currentPrefix string
	if isLast {
		currentPrefix = "└───"
	} else {
		currentPrefix = "├───"
	}

	fullPath := filepath.Join(path, item.Name())

	if item.IsDir() {
		fmt.Fprintf(out, "%s%s%s\n", prefix, currentPrefix, item.Name())
		subItems, err := os.ReadDir(fullPath)
		
		if err != nil {
			return
		}

		if !printFiles{
			var clone []fs.DirEntry
			for _, item := range subItems {
				if item.IsDir(){
					clone = append(clone, item)
				}
			}
			subItems = clone
		}
		newPrefix := prefix
		if isLast {
			newPrefix += "\t"
		 } else {
		 	newPrefix += "│\t"
		}
		
		for i, subItem := range subItems {
			subIsLast := i == len(subItems)-1
			printItem(out, subItem, fullPath, printFiles, newPrefix, subIsLast)
		}
		return
	}

	if printFiles {

		fileInfo, err := item.Info()
		if err != nil {
			return
		}

        fileSize := strconv.FormatInt(fileInfo.Size(), 10) + "b"
		fmt.Fprintf(out, "%s%s%s (%s)\n", prefix, currentPrefix, item.Name(), fileSize)
	}
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printfiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printfiles)
	if err != nil {
		panic(err.Error())
	}
}
