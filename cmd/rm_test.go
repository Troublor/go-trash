package cmd

import (
	"github.com/Troublor/go-trash/errs"
	"github.com/Troublor/go-trash/system"
	"io/ioutil"
	"testing"
)

func TestGenId(t *testing.T) {
	id1 := genId()
	id2 := genId()
	if id1 == id2 {
		t.Fatal("id is the same")
	}
}

func TestRemoveFile(t *testing.T) {
	// create a tmp file
	tmpFile := newTmpFile()
	defer tmpFile.delete()
	// test rm the file
	_, err := Remove(tmpFile.Path, true, true, false)
	if err != errs.IsFileError {
		t.Fatal(err)
	}
	id, err := Remove(tmpFile.Path, false, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 6 {
		t.Fatal("id format not right ", id)
	}
	// file should already be deleted
	if tmpFile.exists() {
		t.Fatal("file not deleted")
	}
	// check database
	item, err := db.GetTrashItemById(id, system.GetUser())
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != id {
		t.Fatal("id in db not right", item.ID, id)
	}
	if item.OriginalPath != tmpFile.Path {
		t.Fatal("original path not right ", item.OriginalPath)
	}
}

func TestRemoveDir(t *testing.T) {
	// create a tmp dir
	tmpDir := newTmpDir()
	defer tmpDir.delete()
	// test rm the dir
	_, err := Remove(tmpDir.Path, false, false, false)
	if err != errs.IsDirectoryError {
		t.Fatal(err)
	}
	id, err := Remove(tmpDir.Path, true, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 6 {
		t.Fatal("id format not right ", id)
	}
	// dir should already be deleted
	if tmpDir.exists() {
		t.Fatal("directory not deleted")
	}
	// check database
	item, err := db.GetTrashItemById(id, system.GetUser())
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != id {
		t.Fatal("id in db not right")
	}
	if item.OriginalPath != tmpDir.Path {
		t.Fatal("original path not right ", item.OriginalPath)
	}
}

func TestRemoveDirRecursive(t *testing.T) {
	// create a tmp dir
	tmpDir := newTmpDir()
	defer tmpDir.delete()
	_, err := ioutil.TempDir(tmpDir.Path, "0644")
	if err != nil {
		panic(err)
	}
	// test rm dir
	_, err = Remove(tmpDir.Path, true, false, false)
	if err != errs.DirectoryNotEmptyError {
		t.Fatal(err)
	}
	id, err := Remove(tmpDir.Path, true, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 6 {
		t.Fatal("id format not right ", id)
	}
	// dir should already be deleted
	if tmpDir.exists() {
		t.Fatal("directory not deleted")
	}
	// check database
	item, err := db.GetTrashItemById(id, system.GetUser())
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != id {
		t.Fatal("id in db not right")
	}
	if item.OriginalPath != tmpDir.Path {
		t.Fatal("original path not right ", item.OriginalPath)
	}
}

func TestRemovePermanently(t *testing.T) {
	// create a tmp file
	tmpFile := newTmpFile()
	defer tmpFile.delete()
	// test rm the file
	id, err := Remove(tmpFile.Path, false, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) > 0 {
		t.Fatal("permanent delete should not return id")
	}
	// file should already be deleted
	if tmpFile.exists() {
		t.Fatal("file not deleted")
	}
	list := db.ListTrashItems(system.GetUser())
	for _, item := range list {
		if item.OriginalPath == tmpFile.Path {
			t.Fatal("permanent delete not work")
		}

	}
}
