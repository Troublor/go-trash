package main

import (
	"github.com/Troublor/trash-go/cmd"
	"github.com/Troublor/trash-go/errs"
	"github.com/Troublor/trash-go/service"
	"github.com/Troublor/trash-go/storage"
	"github.com/Troublor/trash-go/system"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	mockConfigFile := func() {
		cmdDir := system.GetTrashCmdDir()
		testPayload := `{"trashDir":"` + cmdDir + `"}`
		err := ioutil.WriteFile(filepath.Join(system.GetTrashCmdDir(), "gotrash-config.json"), []byte(testPayload), 0666)
		if err != nil {
			panic(err)
		}
	}

	// test context initialization here
	mockConfigFile()
	storage.InitStorage()
	service.MustSubscribeEvent("onTestEnd", func(event service.Event) {
		// delete trash_bin trash_info.db after test finishes
		_ = os.RemoveAll(storage.GetTrashBinPath())
		_ = os.Remove(storage.GetDbPath())
	})
	createTestDir("tmp")
	defer removeTestDir("tmp")
	os.Chdir("tmp")
	defer os.Chdir("..")
	service.MustEventHappen("onTestStart")
	defer service.MustEventHappen("onTestEnd")
	_ = m.Run()
}

func TestNormalFile(t *testing.T) {
	filePath := "abc.txt"
	createTestFileAndClose(filePath, "123")
	defer removeTestFile(filePath)
	_, err := cmd.Remove(filePath, true, false, false)
	if err != errs.IsFileError {
		panic("report wrong error type")
	}
	id, err := cmd.Remove(filePath, false, false, false)
	if err != nil {
		t.Error("remove error: " + err.Error())
	}
	_, err = os.Stat(filePath)
	if err == nil {
		t.Error("didn't remove")
	}
	trashPath := path.Join(storage.GetTrashBinPath(), path.Base(id))
	_, err = os.Stat(trashPath)
	if err != nil {
		t.Error("removed item is not in trash bin")
	}
	infos := storage.DbListAllTrashItems(system.GetUser())
	if len(infos) != 1 {
		t.Error("the length of database record is wrong")
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	if infos[0].Id != id ||
		infos[0].BaseName != path.Base(filePath) ||
		infos[0].OriginalPath != absFilePath ||
		infos[0].TrashPath != trashPath ||
		infos[0].ItemType != storage.TYPE_FILE {
		t.Error("database record error")
	}
	trashInfo, err := cmd.UnRemove(id, true, false, "/original", false)
	if err != nil {
		t.Error("un-remove error")
	}
	_, err = os.Stat(trashInfo.TrashPath)
	if err == nil {
		t.Error("file still in the trash bin")
	}
	_, err = os.Stat(trashInfo.OriginalPath)
	if err != nil {
		t.Error("file is not in the original path")
	}
	infos = storage.DbListAllTrashItems(system.GetUser())
	if len(infos) > 0 {
		t.Error("database record is not deleted")
	}
}

func createTestFile(filePath, content string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	return file
}

func createTestFileAndClose(filePath, content string) {
	file := createTestFile(filePath, content)
	_ = file.Close()
}

func removeTestFile(filePath string) {
	_ = os.Remove(filePath)
}

func createTestDir(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func removeTestDir(dirPath string) {
	_ = os.RemoveAll(dirPath)
}

func TestWrongFilePath(t *testing.T) {
	_, err := cmd.Remove("path/not/exist", false, false, false)
	if err == nil {
		t.Error("don't report file not exist error")
	}
	if err != errs.ItemNotExistError {
		t.Error("report a wrong error type")
	}
	_, err = cmd.UnRemove("non-exist", false, false, "/original", false)
	if err == nil {
		t.Error("don't report file not exist error")
	}
	if err != errs.ItemNotExistError {
		t.Error("report a wrong error type")
	}
}

func TestEmptyDirectory(t *testing.T) {
	dirPath := "test_dir"
	createTestDir(dirPath)
	defer removeTestDir(dirPath)
	_, err := cmd.Remove(dirPath, false, false, false)
	if err == nil {
		t.Error("delete directory when it shouldn't")
	}
	if err != errs.IsDirectoryError {
		t.Error("report wrong error type")
	}
	id, err := cmd.Remove(dirPath, true, false, false)
	if err != nil {
		t.Error("remove directory failed")
	}
	_, err = os.Stat(dirPath)
	if err == nil {
		t.Error("directory is not deleted at all")
	}
	trashPath := path.Join(storage.GetTrashBinPath(), path.Base(id))
	_, err = os.Stat(trashPath)
	if err != nil {
		t.Error("removed item is not in trash bin")
	}
	infos := storage.DbListAllTrashItems(system.GetUser())
	if len(infos) != 1 {
		t.Error("the length of database record is wrong")
	}
	dirPath, err = filepath.Abs(dirPath)
	if infos[0].Id != id ||
		infos[0].BaseName != path.Base(dirPath) ||
		infos[0].OriginalPath != dirPath ||
		infos[0].TrashPath != trashPath ||
		infos[0].ItemType != storage.TYPE_DIRECTORY {
		t.Error("database record error")
	}
	trashInfo, err := cmd.UnRemove(id, true, false, "/original", false)
	if err != nil {
		t.Error("un-remove error")
	}
	_, err = os.Stat(trashInfo.TrashPath)
	if err == nil {
		t.Error("file still in the trash bin")
	}
	_, err = os.Stat(trashInfo.OriginalPath)
	if err != nil {
		t.Error("file is not in the original path")
	}
	infos = storage.DbListAllTrashItems(system.GetUser())
	if len(infos) > 0 {
		t.Error("database record is not deleted")
	}
}

func TestNestedDirectory(t *testing.T) {
	dirPath1, dirPath2 := "parent", "child"
	filePath1, filePath2 := "file1.txt", "file2.txt"
	createTestDir(dirPath1)
	defer removeTestDir(dirPath1)
	createTestDir(filepath.Join(dirPath1, dirPath2))
	defer removeTestDir(filepath.Join(dirPath1, dirPath2))
	createTestFileAndClose(filepath.Join(dirPath1, filePath1), "")
	defer removeTestFile(filepath.Join(dirPath1, filePath1))
	createTestFileAndClose(path.Join(dirPath1, dirPath2, filePath2), "")
	defer removeTestFile(path.Join(dirPath1, dirPath2, filePath2))
	_, err := cmd.Remove(dirPath1, false, false, false)
	if err == nil {
		t.Error("remove dir when it shouldn't")
	}
	_, err = cmd.Remove(dirPath1, true, false, false)
	if err == nil {
		t.Error("remove a non-empty dir when it shouldn't")
	}
	id, err := cmd.Remove(dirPath1, true, true, false)
	if err != nil {
		t.Error("remove dir failed")
	}
	info, err := os.Stat(path.Join(storage.GetTrashBinPath(), id))
	if err != nil {
		t.Error("removed item is not in trash bin")
	}
	if !info.IsDir() {
		t.Error("item type is wrong")
	}
	infos := storage.DbListAllTrashItems(system.GetUser())
	if len(infos) != 1 {
		t.Error("number of records in database is wrong")
	}
	originalPath, _ := filepath.Abs(dirPath1)
	if infos[0].Id != id ||
		infos[0].OriginalPath != originalPath ||
		infos[0].TrashPath != path.Join(storage.GetTrashBinPath(), id) ||
		infos[0].ItemType != storage.TYPE_DIRECTORY ||
		infos[0].BaseName != path.Base(dirPath1) ||
		infos[0].Owner != system.GetUser() {
		t.Error("record information is wrong")
	}
	_, err = cmd.UnRemove(id, true, false, "/original", false)
	if err != nil {
		t.Error("un-remove failed")
	}
	_, err = os.Stat(dirPath1)
	if err != nil {
		t.Error("un-remove not complete")
	}
}
func TestOverride(t *testing.T) {
	filePath := "file.txt"
	createTestFileAndClose(filePath, "abc")
	defer removeTestFile(filePath)
	id, err := cmd.Remove(filePath, false, false, false)
	if err != nil {
		t.Error("remove failed")
	}
	createTestFileAndClose(filePath, "")
	_, err = cmd.UnRemove(id, true, false, "/original", false)
	if err == nil {
		t.Error("override when it shouldn't")
	} else if err != errs.ItemExistError {
		t.Error("report wrong error")
	}
	_, err = cmd.UnRemove(id, true, true, "/original", false)
	if err != nil {
		t.Error("un-remove failed")
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	if string(data) != "abc" {
		t.Error("not the original file")
	}
}

func TestUnremoveRedirect(t *testing.T) {
	filePath := "abc.txt"
	dirPath := "redirect"
	createTestFileAndClose(filePath, "")
	defer removeTestFile(filePath)
	_, err := cmd.Remove(filePath, false, false, false)
	if err != nil {
		t.Error("remove failed")
	}
	_, err = cmd.UnRemove(filePath, false, false, dirPath, true)
	if err != nil {
		t.Error("un-remove failed")
	}
	defer removeTestDir(dirPath)
	if _, err := os.Stat(filePath); err == nil {
		t.Error("file is still at the original path")
	}
	if _, err := os.Stat(filepath.Join(dirPath, filePath)); err != nil {
		t.Error("file is not at the target path")
	}
}

func TestCrossDriverRemove(t *testing.T) {
	filePath := "/tmp/file.txt"
	createTestFileAndClose(filePath, "123")
	defer removeTestFile(filePath)
	id, err := cmd.Remove(filePath, false, false, false)
	if err != nil {
		t.Fatal("remove failed")
	}
	if _, err = os.Stat(path.Join(storage.GetTrashBinPath(), id)); err != nil {
		t.Fatal("remove unfinished")
	}
	_, err = cmd.UnRemove(id, true, false, "/original", false)
	if err != nil {
		t.Fatal("un-remove failed")
	}
	if _, err = os.Stat(filePath); err != nil {
		t.Fatal("un-remove unfinished")
	}

	//directory
	dirPath1, dirPath2 := "~/parent", "child"
	filePath1, filePath2 := "file1.txt", "file2.txt"
	createTestDir(dirPath1)
	defer removeTestDir(dirPath1)
	createTestDir(filepath.Join(dirPath1, dirPath2))
	defer removeTestDir(filepath.Join(dirPath1, dirPath2))
	createTestFileAndClose(path.Join(dirPath1, filePath1), "")
	defer removeTestFile(path.Join(dirPath1, filePath1))
	createTestFileAndClose(path.Join(dirPath1, dirPath2, filePath2), "")
	defer removeTestFile(path.Join(dirPath1, dirPath2, filePath2))
	id, err = cmd.Remove(dirPath1, true, true, false)
	if err != nil {
		t.Fatal("remove failed")
	}
	if _, err = os.Stat(dirPath1); err == nil {
		t.Fatal("remove unfinished")
	}
	_, err = cmd.UnRemove(id, true, false, "/original", false)
	if err != nil {
		t.Fatal("un-remove failed")
	}
	if _, err = os.Stat(dirPath1); err != nil {
		t.Fatal("un-remove unfinished")
	}
}
