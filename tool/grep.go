package tool

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type GrepProps struct {
	Flags FlagOptions
	Args  []string
}

type FlagOptions struct {
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

	res.Lines = searchRes
	return res
}
