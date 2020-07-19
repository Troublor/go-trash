package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	GOTRASH_PATH = "./tmp/"
	// recreate GOTRASH_PATH if not exist
	if _, err := os.Stat(GOTRASH_PATH); os.IsNotExist(err) {
		err := os.MkdirAll(GOTRASH_PATH, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	defer func() {
		// remove the GOTRASH_PATH
		err := os.RemoveAll(GetDbPath())
		if err != nil {
			panic(err)
		}
	}()

	// do initialization work
	initialize()
	defer func() {
		// remove the db
		err := os.RemoveAll(GetDbPath())
		if err != nil {
			panic(err)
		}
	}()
	defer func() {
		// close db
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_ = os.MkdirAll(GOTRASH_PATH, os.ModePerm)
	defer func() {
		_ = os.RemoveAll(GOTRASH_PATH)
	}()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_ = m.Run()
}

type tmpItem interface {
	recreate()    // recreate the item if not exist
	delete()      // delete the item
	exists() bool // whether the item exists
}

type tmpFile struct {
	Path     string
	BaseName string
}

func newTmpFile() tmpFile {
	file, err := ioutil.TempFile(os.TempDir(), "0644")
	if err != nil {
		panic(err)
	}
	stat, _ := file.Stat()
	return tmpFile{Path: file.Name(), BaseName: stat.Name()}
}

func (f tmpFile) recreate() {
	if f.exists() {
		return
	}
	_, err := os.Create(f.Path)
	if err != nil {
		panic(err)
	}
}

func (f tmpFile) delete() {
	_ = os.RemoveAll(f.Path)
}

func (f tmpFile) exists() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

type tmpDir struct {
	Path     string
	BaseName string
}

func newTmpDir() tmpDir {
	dir, err := ioutil.TempDir(os.TempDir(), "0644")
	if err != nil {
		panic(err)
	}
	stat, _ := os.Stat(dir)
	return tmpDir{Path: dir, BaseName: stat.Name()}
}

func (f tmpDir) recreate() {
	if f.exists() {
		return
	}
	err := os.MkdirAll(f.Path, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (f tmpDir) delete() {
	_ = os.RemoveAll(f.Path)
}

func (f tmpDir) exists() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}
