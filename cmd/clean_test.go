package cmd

import (
	"github.com/Troublor/go-trash/errs"
	"github.com/Troublor/go-trash/system"
	"testing"
)

func TestClean(t *testing.T) {
	tmpFile := newTmpFile()
	defer tmpFile.delete()
	id, _ := Remove(tmpFile.Path, false, false, false)
	list := List()
	if len(list) <= 0 {
		t.Fatal("trash bin should not be empty")
	}
	err := Clean(true, id)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.GetTrashItemById(id, system.GetUser())
	if err != errs.ItemNotExistError {
		t.Fatal("trash should be cleaned")
	}
}
