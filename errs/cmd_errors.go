package errs

import (
	"errors"
	"fmt"
)

var ItemNotExistError error
var ItemExistError error
var MultipleItemsError error
var DirectoryNotEmptyError error
var IsDirectoryError error
var IsFileError error
var PermissionError error

func init() {
	ItemNotExistError = errors.New("item not exists")
	ItemExistError = errors.New("item has already existed")
	MultipleItemsError = errors.New("find multiple items")
	DirectoryNotEmptyError = errors.New("the directory is not empty")
	IsDirectoryError = errors.New("the item is a directory")
	IsFileError = errors.New("the item is a file")
	PermissionError = errors.New("permission denied")
}

type FileOrDirNotExistError struct {
	path string
}

func NewFileOrDirNotExistError(path string) FileOrDirNotExistError {
	return FileOrDirNotExistError{path: path}
}

func (f FileOrDirNotExistError) Error() string {
	return fmt.Sprintf("no such file or directory: %s", f.path)
}
