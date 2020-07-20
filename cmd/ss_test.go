package cmd

import (
	"testing"
)

func TestSearch(t *testing.T) {
	tmpFile := newTmpFile()
	defer tmpFile.delete()
	list := Search(tmpFile.BaseName)
	if len(list) != 0 {
		t.Fatal("result should be empty")
	}
	_, _ = Remove(tmpFile.Path, false, false, false)
	list = Search(tmpFile.BaseName)
	if len(list) != 1 {
		t.Fatal("result should not be empty")
	}
	list = Search("not exist")
	if len(list) != 0 {
		t.Fatal("result should be empty")
	}
	list = Search(tmpFile.BaseName + "*")
	if len(list) == 0 {
		t.Fatal("result should not be empty")
	}
}
