package tool

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchWithOneFileAndOneSearchParam(t *testing.T) {
	grep := GrepProps{
		Args: []string{"search string", "../test_assets/testFile.txt"},
	}
	r, w, oldStdout := replaceStdout()
	grep.Search()
	w.Close()
	readFromStdoutAndCompare(r, oldStdout,
		[]string{"I am a File with a search string and I don't know what to do.",
			"Plus I don't know what a search string looks like"}, t)

}

func TestSearchWithReadFromStandardInput(t *testing.T) {
	oldStdIn := os.Stdin
	defer func() { os.Stdin = oldStdIn }()
	r, w, _ := os.Pipe()
	os.Stdin = r

	stdOutReader, stdWriter, oldStdout := replaceStdout()
	data := []string{"Writing search string on Stdin\n", "I am dumb\n", "I don't know what a search string looks like\n"}
	grep := GrepProps{
		Args: []string{"search string"},
	}

	go func() {
		for _, d := range data {
			w.WriteString(d)
		}
		w.Close()
	}()
	grep.Search()
	stdWriter.Close()
	readFromStdoutAndCompare(stdOutReader, oldStdout, []string{"Writing search string on Stdin", "I don't know what a search string looks like"}, t)
}

func TestSearchWithOutputAsFile(t *testing.T) {
	grep := GrepProps{
		Args:  []string{"search string", "../test_assets/testFile.txt"},
		Flags: FlagOptions{OutputFile: "../test_assets/outputFile.txt"},
	}

	grep.Search()
	f, err := os.Open("../test_assets/outputFile.txt")
	defer f.Close()
	defer os.Remove("../test_assets/outputFile.txt")
	assert.Equal(t, err, nil)
	scanner := bufio.NewScanner(f)
	output := []string{}

	for scanner.Scan() {
		output = append(output, scanner.Text())
	}
	assert.Equal(t, output, []string{"I am a File with a search string and I don't know what to do.",
		"Plus I don't know what a search string looks like"})
}

func TestCaseInsensitiveSearch(t *testing.T) {
	grep := GrepProps{
		Args:  []string{"Search String", "../test_assets/testFile.txt"},
		Flags: FlagOptions{OutputFile: "../test_assets/outputFile.txt", CaseInsensitive: true},
	}

	grep.Search()
	f, err := os.Open("../test_assets/outputFile.txt")
	defer f.Close()
	defer os.Remove("../test_assets/outputFile.txt")
	assert.Equal(t, err, nil)
	scanner := bufio.NewScanner(f)
	output := []string{}

	for scanner.Scan() {
		output = append(output, scanner.Text())
	}
	assert.Equal(t, output, []string{"I am a File with a search string and I don't know what to do.",
		"Plus I don't know what a search string looks like"})

}

func replaceStdout() (*os.File, *os.File, *os.File) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	return r, w, oldStdout
}

func readFromStdoutAndCompare(r, oldStdout *os.File, expected []string, t *testing.T) {
	defer r.Close()
	defer func() { os.Stdout = oldStdout }()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	var actual = strings.Split(buf.String(), "\n")

	assert.Equal(t, expected, actual[0:len(actual)-1])
}
