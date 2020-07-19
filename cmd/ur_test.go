package cmd

import (
	"github.com/Troublor/go-trash/errs"
	"os"
	"path"
	"testing"
)

func TestUnRemove(t *testing.T) {
	// recreate a tmp file
	tmpFile := newTmpFile()
	defer tmpFile.delete()
	// delete the file
	id, _ := Remove(tmpFile.Path, false, false, false)
	// un-remove the file by id
	_, err := UnRemove(id, true, false, "/original", false)
	if err != nil {
		t.Fatal(err)
	}
	if !tmpFile.exists() {
		t.Fatal("file not un-removed")
	}
	// delete the file
	_, _ = Remove(tmpFile.Path, false, false, false)
	// un-remove file by name
	_, err = UnRemove(tmpFile.BaseName, false, false, "/original", false)
	if err != nil {
		t.Fatal(err)
	}
	// delete the file
	id, _ = Remove(tmpFile.Path, false, false, false)
	// un-remove file to a different place
	// recreate another file with the same name
	tmpFile.recreate()
	_, err = UnRemove(id, true, false, "/original", false)
	if err != errs.ItemExistError {
		t.Fatal("should not un-remove")
	}
	tmpFile.delete()
	// un-remove to a different place
	tmpDir := newTmpDir()
	defer tmpDir.delete()
	tmpDir.delete()
	_, err = UnRemove(id, true, false, tmpDir.Path, false)
	if _, ok := err.(errs.FileOrDirNotExistError); !ok {
		t.Fatal(err)
	}
	_, err = UnRemove(id, true, false, tmpDir.Path, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(path.Join(tmpDir.Path, tmpFile.BaseName)); os.IsNotExist(err) {
		t.Fatal("fit not un-removed")
	}
}
