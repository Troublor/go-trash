package operation

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

func Remove(itemPath string, isDirectory bool, recursive bool) (string, error) {
	trashDir := GetTrashPath()
	var err error
	if !filepath.IsAbs(itemPath) {
		itemPath, err = filepath.Abs(itemPath)
		if err != nil {
			panic(err)
		}
	}
	fileInfo, err := os.Stat(itemPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ItemNotExistError
		} else {
			panic(err)
		}
	}

	if !isDirectory {
		if fileInfo.IsDir() {
			return "", IsDirectoryError
		}
		// add information in database
		id := DbInsertTrashItem(itemPath, trashDir, fileInfo.Name(), TYPE_FILE)
		// move the item into trash directory
		err := SafeRename(itemPath, path.Join(trashDir, id))
		if err != nil {
			panic(err)
		}
		return id, nil
	} else {
		if !fileInfo.IsDir() {
			return "", IsFileError
		}
		isEmpty, err := DirectoryIsEmpty(itemPath)
		if err != nil {
			panic(err)
		}
		if !isEmpty && !recursive {
			return "", DirectoryNotEmptyError
		}
		// add information in database
		id := DbInsertTrashItem(itemPath, trashDir, fileInfo.Name(), TYPE_DIRECTORY)
		// move the item into trash directory
		err = SafeRename(itemPath, path.Join(trashDir, id))
		if err != nil {
			panic(err)
		}
		return id, nil
	}
}

func DirectoryIsEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
