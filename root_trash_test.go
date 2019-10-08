package main

import (
	"github.com/Troublor/trash-go/cmd"
	"github.com/Troublor/trash-go/storage"
	"os"
	"path"
	"testing"
)

func TestInit(t *testing.T) {
	storage.InitStorage()
}

func TestCrossDriverRemove(t *testing.T) {
	filePath := "/var/www/file.txt"
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(filePath)
	}()
	_, _ = file.WriteString("123")
	_ = file.Close()
	id, err := cmd.Remove(filePath, false, false)
	if err != nil {
		t.Fatal("remove failed")
	}
	if _, err = os.Stat(path.Join(storage.GetTrashPath(), id)); err != nil {
		t.Fatal("remove unfinished")
	}
	_, err = cmd.UnRemove(id, true, false)
	if err != nil {
		t.Fatal("un-remove failed")
	}
	if _, err = os.Stat(filePath); err != nil {
		t.Fatal("un-remove unfinished")
	}

	//directory
	dirPath1, dirPath2 := "/var/www/parent", "child"
	filePath1, filePath2 := "file1.txt", "file2.txt"
	err = os.Mkdir(dirPath1, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(dirPath1)
	}()
	err = os.Mkdir(path.Join(dirPath1, dirPath2), os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(path.Join(dirPath1, dirPath2))
	}()
	file, err = os.Create(path.Join(dirPath1, filePath1))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(path.Join(dirPath1, filePath1))
	}()
	err = file.Close()
	if err != nil {
		panic(err)
	}
	file, err = os.Create(path.Join(dirPath1, dirPath2, filePath2))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(path.Join(dirPath1, dirPath2, filePath2))
	}()
	err = file.Close()
	if err != nil {
		panic(err)
	}
	id, err = cmd.Remove(dirPath1, true, true)
	if err != nil {
		t.Fatal("remove failed")
	}
	if _, err = os.Stat(dirPath1); err == nil {
		t.Fatal("remove unfinished")
	}
	_, err = cmd.UnRemove(id, true, false)
	if err != nil {
		t.Fatal("un-remove failed")
	}
	if _, err = os.Stat(dirPath1); err != nil {
		t.Fatal("un-remove unfinished")
	}
}

func TestToRemoveTestFiles(t *testing.T) {
	_ = os.RemoveAll("tmp")
	_ = os.RemoveAll("trash_info.db")
	_ = os.RemoveAll("trash_bin")
}
