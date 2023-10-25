package tool

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type GrepProps struct {
	Flags FlagOptions
	Args  []string
}

type FlagOptions struct {
	OutputFile      string
	CaseInsensitive bool
	Recursive       bool
}

type Result struct {
	File  string
	Lines []string
	Err   error
}

const chunkSize = 1024 * 1024

func (grep GrepProps) Search() {

	searchText := grep.Args[0]
	if grep.Flags.CaseInsensitive {
		searchText = strings.ToLower(searchText)
	}
	maxFileBuffer := make(chan int, 10)
	wg := &sync.WaitGroup{}

	if len(grep.Args) > 1 && grep.Flags.Recursive {
		filepath.Walk(grep.Args[1], func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !info.IsDir() {
				wg.Add(1)
				go readFromFile(path, grep, searchText, maxFileBuffer, wg)
			}
			return nil
		})
	} else {

		var file string
		if len(grep.Args) > 1 {
			fileInfo, err := os.Stat(grep.Args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			if fileInfo.IsDir() {
				fmt.Println(err)
				return
			}
			file = grep.Args[1]
		} else {
			file = "-"
		}
		wg.Add(1)
		go readFromFile(file, grep, searchText, maxFileBuffer, wg)
	}

	wg.Wait()
}

func readFromFile(fileName string, grep GrepProps, searchText string, maxFileBuffer chan int, wg *sync.WaitGroup) {

	maxFileBuffer <- 1
	defer func() {
		<-maxFileBuffer
		wg.Done()
	}()
	lineChan := make(chan string)
	errChan := make(chan error)
	var scanner *bufio.Scanner

	if fileName != "-" {
		file, err := os.Open(fileName)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		scanner = bufio.NewScanner(file)
		scanner.Buffer(make([]byte, chunkSize), chunkSize)

	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	go func() {
		defer close(lineChan)
		defer close(errChan)
		for scanner.Scan() {
			lineChan <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
		}
	}()

	result := grep.search(lineChan, errChan, searchText)
	result.File = fileName
	if result.Err != nil {
		fmt.Println(result.Err)
	} else {
		if grep.Flags.OutputFile == "" {
			for _, line := range result.Lines {
				fmt.Println(line)
			}

		} else {
			writeOutputToFile(grep, result.Lines)
		}
	}
}

func (grep GrepProps) search(lineChan chan string, errChan chan error, searchText string) Result {
	searchRes := Result{Lines: make([]string, 0)}
	for {
		select {
		case err := <-errChan:
			if err != nil {
				searchRes.Err = err
				return searchRes
			}
		case line, ok := <-lineChan:
			if !ok {
				return searchRes
			}

			textToSearchOn := line
			if grep.Flags.CaseInsensitive {
				textToSearchOn = strings.ToLower(line)
			}

			if strings.Contains(textToSearchOn, searchText) {
				searchRes.Lines = append(searchRes.Lines, line)
			}

		}
	}
}

func writeOutputToFile(grep GrepProps, resLines []string) {
	var outputFile *os.File
	var err error

	info, err := os.Stat(grep.Flags.OutputFile)

	if err == nil && info.IsDir() {
		fmt.Println(err)
	} else {
		outputFile, err = os.OpenFile(grep.Flags.OutputFile, os.O_WRONLY|os.O_CREATE, 0644)
		for _, res := range resLines {
			outputFile.WriteString(fmt.Sprintln(res))
		}
	}
}
