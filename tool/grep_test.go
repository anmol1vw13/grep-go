package tool

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSearchWithOneFileAndOneSearchParam(t *testing.T) {
	grep := GrepProps{
		Args: []string{"search string", "../test_assets/testFile.txt"},
	}

	result := grep.Search()
	assert.Equal(t, result.Lines, []string {"I am a File with a search string and I don't know what to do.", 
	"Plus I don't know what a search string looks like"})
}

func TestSearchWithReadFromStandardInput(t *testing.T) {
	oldStdIn := os.Stdin
	defer func() { os.Stdout = oldStdIn }()
	r, w, _ := os.Pipe()
	os.Stdin = r

	data := []string {"Writing search string on Stdin\n", "I am dumb\n", "I don't know what a search string looks like\n"}
	grep := GrepProps{
		Args: []string{"search string"},
	}

	go func() {for _, d := range data {
		w.WriteString(d)
	}
	w.Close()
	}()
	result := grep.Search()
	
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)	
	assert.Equal(t, result.Lines, []string {"Writing search string on Stdin", "I don't know what a search string looks like"})
}