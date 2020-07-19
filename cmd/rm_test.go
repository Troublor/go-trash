package cmd

import (
	"github.com/Troublor/go-trash/errs"
	"github.com/Troublor/go-trash/system"
	"io/ioutil"
	"os"
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
	tmpFile, err := ioutil.TempFile(os.TempDir(), "0644")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpFile.Name())
	}()
	// test rm the file
	_, err = Remove(tmpFile.Name(), true, true, false)
	if err != errs.IsFileError {
		t.Fatal(err)
	}
	id, err := Remove(tmpFile.Name(), false, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 6 {
		t.Fatal("id format not right ", id)
	}
	// file should already be deleted
	if _, err := os.Stat(tmpFile.Name()); !os.IsNotExist(err) {
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
	if item.OriginalPath != tmpFile.Name() {
		t.Fatal("original path not right ", item.OriginalPath)
	}
}

func TestRemoveDir(t *testing.T) {
	// create a tmp file
	tmpDir, err := ioutil.TempDir(os.TempDir(), "0644")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	// test rm the dir
	_, err = Remove(tmpDir, false, false, false)
	if err != errs.IsDirectoryError {
		t.Fatal(err)
	}
	id, err := Remove(tmpDir, true, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 6 {
		t.Fatal("id format not right ", id)
	}
	// dir should already be deleted
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
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
	if item.OriginalPath != tmpDir {
		t.Fatal("original path not right ", item.OriginalPath)
	}
}

func TestRemoveDirRecursive(t *testing.T) {
	// create a tmp file
	tmpDir, err := ioutil.TempDir(os.TempDir(), "0644")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	_, err = ioutil.TempDir(tmpDir, "0644")
	if err != nil {
		panic(err)
	}
	// test rm dir
	_, err = Remove(tmpDir, true, false, false)
	if err != errs.DirectoryNotEmptyError {
		t.Fatal(err)
	}
	id, err := Remove(tmpDir, true, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 6 {
		t.Fatal("id format not right ", id)
	}
	// dir should already be deleted
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
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
	if item.OriginalPath != tmpDir {
		t.Fatal("original path not right ", item.OriginalPath)
	}
}

func TestRemovePermanently(t *testing.T) {
	// create a tmp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "0644")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpFile.Name())
	}()
	// test rm the file
	id, err := Remove(tmpFile.Name(), false, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) > 0 {
		t.Fatal("permanent delete should not return id")
	}
	// file should already be deleted
	if _, err := os.Stat(tmpFile.Name()); !os.IsNotExist(err) {
		t.Fatal("file not deleted")
	}
	list := db.ListTrashItems(system.GetUser())
	for _, item := range list {
		if item.OriginalPath == tmpFile.Name() {
			t.Fatal("permanent delete not work")
		}

	}
}
