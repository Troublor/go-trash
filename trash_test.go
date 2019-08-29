package main

import (
	"github.com/Troublor/trash-go/operation"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestNormalFile(t *testing.T) {
	filePath := "/home/troublor/workspace/go/trash-go/abc.txt"
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_ = file.Close()
	defer func() {
		_ = os.Remove(filePath)
	}()
	_, err = operation.Remove(filePath, true, false)
	if err != operation.IsFileError {
		panic("report wrong error type")
	}
	id, err := operation.Remove(filePath, false, false)
	if err != nil {
		t.Fatal("remove error: " + err.Error())
	}
	_, err = os.Stat(filePath)
	if err == nil {
		t.Fatal("didn't remove")
	}
	trashPath := path.Join(operation.GetTrashPath(), path.Base(id))
	_, err = os.Stat(trashPath)
	if err != nil {
		t.Fatal("removed item is not in trash bin")
	}
	infos := operation.DbListAllTrashItems()
	if len(infos) != 1 {
		t.Fatal("the length of database record is wrong")
	}
	if infos[0].Id != id ||
		infos[0].BaseName != path.Base(filePath) ||
		infos[0].OriginalPath != filePath ||
		infos[0].TrashPath != trashPath ||
		infos[0].ItemType != operation.TYPE_FILE {
		t.Fatal("database record error")
	}
	trashInfo, err := operation.UnRemove(id, true, false)
	if err != nil {
		t.Fatal("un-remove error")
	}
	_, err = os.Stat(trashInfo.TrashPath)
	if err == nil {
		t.Fatal("file still in the trash bin")
	}
	_, err = os.Stat(trashInfo.OriginalPath)
	if err != nil {
		t.Fatal("file is not in the original path")
	}
	infos = operation.DbListAllTrashItems()
	if len(infos) > 0 {
		t.Fatal("database record is not deleted")
	}
}

func TestWrongFilePath(t *testing.T) {
	_, err := operation.Remove("path/not/exist", false, false)
	if err == nil {
		t.Fatal("don't report file not exist error")
	}
	if err != operation.ItemNotExistError {
		t.Fatal("report a wrong error type")
	}
	_, err = operation.UnRemove("non-exist", false, false)
	if err == nil {
		t.Fatal("don't report file not exist error")
	}
	if err != operation.ItemNotExistError {
		t.Fatal("report a wrong error type")
	}
}

func TestEmptyDirectory(t *testing.T) {

	dirPath := "test_dir"
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(dirPath)
	}()
	_, err = operation.Remove(dirPath, false, false)
	if err == nil {
		t.Fatal("delete directory when it shouldn't")
	}
	if err != operation.IsDirectoryError {
		t.Fatal("report wrong error type")
	}
	id, err := operation.Remove(dirPath, true, false)
	if err != nil {
		t.Fatal("remove directory failed")
	}
	_, err = os.Stat(dirPath)
	if err == nil {
		t.Fatal("directory is not deleted at all")
	}
	trashPath := path.Join(operation.GetTrashPath(), path.Base(id))
	_, err = os.Stat(trashPath)
	if err != nil {
		t.Fatal("removed item is not in trash bin")
	}
	infos := operation.DbListAllTrashItems()
	if len(infos) != 1 {
		t.Fatal("the length of database record is wrong")
	}
	dirPath, err = filepath.Abs(dirPath)
	if infos[0].Id != id ||
		infos[0].BaseName != path.Base(dirPath) ||
		infos[0].OriginalPath != dirPath ||
		infos[0].TrashPath != trashPath ||
		infos[0].ItemType != operation.TYPE_DIRECTORY {
		t.Fatal("database record error")
	}
	trashInfo, err := operation.UnRemove(id, true, false)
	if err != nil {
		t.Fatal("un-remove error")
	}
	_, err = os.Stat(trashInfo.TrashPath)
	if err == nil {
		t.Fatal("file still in the trash bin")
	}
	_, err = os.Stat(trashInfo.OriginalPath)
	if err != nil {
		t.Fatal("file is not in the original path")
	}
	infos = operation.DbListAllTrashItems()
	if len(infos) > 0 {
		t.Fatal("database record is not deleted")
	}
}

func TestNestedDirectory(t *testing.T) {
	dirPath1, dirPath2 := "parent", "child"
	filePath1, filePath2 := "file1.txt", "file2.txt"
	err := os.Mkdir(dirPath1, os.ModePerm)
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
	file, err := os.Create(path.Join(dirPath1, filePath1))
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

	_, err = operation.Remove(dirPath1, false, false)
	if err == nil {
		t.Fatal("remove dir when it shouldn't")
	}
	_, err = operation.Remove(dirPath1, true, false)
	if err == nil {
		t.Fatal("remove a non-empty dir when it shouldn't")
	}
	id, err := operation.Remove(dirPath1, true, true)
	if err != nil {
		t.Fatal("remove dir failed")
	}
	info, err := os.Stat(path.Join(operation.GetTrashPath(), id))
	if err != nil {
		t.Fatal("removed item is not in trash bin")
	}
	if !info.IsDir() {
		t.Fatal("item type is wrong")
	}
	infos := operation.DbListAllTrashItems()
	if len(infos) != 1 {
		t.Fatal("number of records in database is wrong")
	}
	originalPath, _ := filepath.Abs(dirPath1)
	if infos[0].Id != id ||
		infos[0].OriginalPath != originalPath ||
		infos[0].TrashPath != path.Join(operation.GetTrashPath(), id) ||
		infos[0].ItemType != operation.TYPE_DIRECTORY ||
		infos[0].BaseName != dirPath1 {
		t.Fatal("record information is wrong")
	}
	_, err = operation.UnRemove(id, true, false)
	if err != nil {
		t.Fatal("un-remove failed")
	}
	_, err = os.Stat(dirPath1)
	if err != nil {
		t.Fatal("un-remove not complete")
	}
}
func TestOverride(t *testing.T) {
	filePath := "file.txt"
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(filePath)
	}()
	_, err = file.WriteString("abc")
	if err != nil {
		panic(err)
	}
	_ = file.Close()
	id, err := operation.Remove(filePath, false, false)
	if err != nil {
		t.Fatal("remove failed")
	}
	file, err = os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_ = file.Close()
	_, err = operation.UnRemove(id, true, false)
	if err == nil {
		t.Fatal("override when it shouldn't")
	} else if err != operation.ItemExistError {
		t.Fatal("report wrong error")
	}
	_, err = operation.UnRemove(id, true, true)
	if err != nil {
		t.Fatal("un-remove failed")
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	if string(data) != "abc" {
		t.Fatal("not the original file")
	}
}
