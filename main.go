package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	//"path/filepath"
)

func stringBuilder(fileName string, fileSize int64, isLast bool, isDirectory bool, level int, isLastDirectory bool, lastParentIndex int, printFiles bool) string {
	var printString strings.Builder

	if level != 0 && !(lastParentIndex == 0 && !printFiles) {
		printString.WriteString("│")
	}

	for i := 1; i <= level; i++ {
		printString.WriteString("\t")

		// ставим только на текущем уровне
		if i <= level-1 && !(isLastDirectory && i >= lastParentIndex) {
			printString.WriteString("│")
		}
	}

	separator := "├"

	if isLast {
		separator = "└"
	}

	printString.WriteString(separator)
	printString.WriteString("───")
	printString.WriteString(fileName)

	if !isDirectory {
		printString.WriteString(" (")

		if fileSize == 0 {
			printString.WriteString("empty")
		} else {
			printString.WriteString(strconv.FormatInt(fileSize, 10))
			printString.WriteString("b")
		}

		printString.WriteString(")")
	}

	return printString.String()
}

func readDir(dirName string, level int, printFiles bool, out io.Writer, isLastDirectory bool, lastParentIndex int) error {
	if dirName == "-f" || dirName[0] == '-' {
		dirName = "."
		printFiles = true
	}

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var directoriesAmount = 0
	for _, file := range files {
		if file.IsDir() {
			directoriesAmount += 1
		}
	}

	directoryIndex := 0

	for index, file := range files {
		isDirectory := file.IsDir()
		fileName := file.Name()

		if isDirectory {
			directoryIndex += 1
		}

		if !isDirectory && !printFiles {
			continue
		}

		miu := isLastDirectory
		isLast := index == len(files)-1

		if !printFiles && isDirectory {
			isLast = directoryIndex == directoriesAmount || directoriesAmount <= 1
		}

		printStr := stringBuilder(fileName, file.Size(), isLast, isDirectory, level, miu, lastParentIndex, printFiles)
		fmt.Fprintln(out, printStr)

		newIndex := lastParentIndex

		if directoriesAmount != directoryIndex {
			newIndex += 1
		}

		if file.IsDir() {
			readDir(dirName+string(os.PathSeparator)+fileName, level+1, printFiles, out, directoriesAmount == directoryIndex, newIndex)
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	readDir(path, 0, printFiles, out, false, 0)

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
