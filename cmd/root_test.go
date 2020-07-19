package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	GOTRASH_PATH = "./tmp/"
	// create GOTRASH_PATH if not exist
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
	Delete()
	Exists() bool
}

type tmpFile struct {
	Path string
}

func newTmpFile() tmpFile {
	file, err := ioutil.TempFile(os.TempDir(), "0644")
	if err != nil {
		panic(err)
	}
	return tmpFile{Path: file.Name()}
}

func (f tmpFile) Delete() {
	_ = os.RemoveAll(f.Path)
}

func (f tmpFile) Exists() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

type tmpDir struct {
	Path string
}

func newTmpDir() tmpDir {
	dir, err := ioutil.TempDir(os.TempDir(), "0644")
	if err != nil {
		panic(err)
	}
	return tmpDir{Path: dir}
}

func (f tmpDir) Delete() {
	_ = os.RemoveAll(f.Path)
}

func (f tmpDir) Exists() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}
