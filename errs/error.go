package errs

import (
	"errors"
)

var ItemNotExistError error
var ItemExistError error
var MultipleItemsError error
var DirectoryNotEmptyError error
var IsDirectoryError error
var IsFileError error

func init() {
	ItemNotExistError = errors.New("item not exists")
	ItemExistError = errors.New("item has already existed")
	MultipleItemsError = errors.New("find multiple items")
	DirectoryNotEmptyError = errors.New("the directory is not empty")
	IsDirectoryError = errors.New("the item is a directory")
	IsFileError = errors.New("the item is a file")
}
