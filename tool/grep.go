package tool

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type GrepProps struct {
	Flags FlagOptions
	Args  []string
}

type FlagOptions struct {
	OutputFile string
}

type Result struct {
	Lines []string
	Err   error
}

func (grep GrepProps) Search() Result {

	searchText := grep.Args[0]
	var scanner *bufio.Scanner
	res := Result{}

	if len(grep.Args) > 1 {
		fileName := grep.Args[1]

		fileInfo, err := os.Stat(fileName)
		if err != nil {
			res.Err = err
			return res
		}
		if fileInfo.IsDir() {
			res.Err = errors.New("File is a directory")
			return res
		}

		file, err := os.Open(fileName)
		defer file.Close()
		if err != nil {
			res.Err = err
			return res
		}

		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	searchRes := make([]string, 0)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(strings.ToLower(text), strings.ToLower(searchText)) {
			searchRes = append(searchRes, text)
		}
	}

	if err := scanner.Err(); err != nil {
		res.Err = err
		return res
	}

	if grep.Flags.OutputFile == "" {
		res.Lines = searchRes
	} else {

		var outputFile *os.File
		var err error

		info, err := os.Stat(grep.Flags.OutputFile)

		if err == nil && info.IsDir() {
			res.Err = errors.New("Output file cannot be a directory")
			return res
		} else {
			outputFile, err = os.OpenFile(grep.Flags.OutputFile, os.O_WRONLY|os.O_CREATE, 0644)
		}
		for _, res := range searchRes {
			outputFile.WriteString(fmt.Sprintln(res))
		}
	}
	return res
}
