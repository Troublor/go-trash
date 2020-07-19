package cmd

import "testing"

func TestList(t *testing.T) {
	list := List()
	if len(list) != 0 {
		t.Fatal("trash bin should be empty")
	}
	tmpFile := newTmpFile()
	defer tmpFile.delete()
	_, _ = Remove(tmpFile.Path, false, false, false)
	list = List()
	if len(list) != 1 {
		t.Fatal("trash bin should only have one")
	}
	tmpDir := newTmpDir()
	defer tmpDir.delete()
	_, _ = Remove(tmpDir.Path, true, false, false)
	list = List()
	if len(list) != 2 {
		t.Fatal("trash bin should only have two")
	}
}
