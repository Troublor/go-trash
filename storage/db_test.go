package storage

import (
	"github.com/Troublor/go-trash/errs"
	"github.com/Troublor/go-trash/storage/model"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

/**
Test open a new database that does not exist before
*/
func TestDatabase_Open(t *testing.T) {
	// a temp database file
	tmp := path.Join("/tmp", time.Now().String())
	defer func() {
		// delete tmp file after test
		_ = os.Remove(tmp)
	}()
	db := NewDatabase(tmp)
	err := db.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	// check if the database is created
	if _, err := os.Stat(tmp); os.IsNotExist(err) {
		t.Fatal("db not exist after open")
	}
}

/**
Test insert item into database and check if item exist
*/
func TestDatabase_InsertGetDeleteListTrashItem(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatal(err)
		}
	}()
	tmp := path.Join("/tmp", time.Now().String())
	defer func() {
		// delete tmp file after test
		_ = os.Remove(tmp)
	}()
	db := NewDatabase(tmp)
	err := db.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	// before insert, should be empty
	list := db.ListTrashItems("user")
	if len(list) > 0 {
		t.Fatal("database should not have data")
	}
	// insert
	err = db.InsertTrashItem("1", "1", "1", "1", TYPE_FILE, "user")
	if err != nil {
		t.Fatal(err)
	}
	// after insert, should have one
	list = db.ListTrashItems("user")
	if len(list) != 1 {
		t.Fatal("database should only contain one item")
	}
	// test if item in database
	item, err := db.GetTrashItemById("1", "user")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(item, model.TrashMetadata{
		ID:           "1",
		OriginalPath: "1",
		TrashPath:    "1/1",
		BaseName:     "1",
		Type:         "F",
		Owner:        "user",
		CreatedAt:    item.CreatedAt,
	}) {
		t.Fatal("inserted item is not as expected")
	}
	// test for permission control
	_, err = db.GetTrashItemById("1", "guest")
	if err != errs.PermissionError {
		t.Fatal("permission control not work")
	}
	// test for item not found err
	_, err = db.GetTrashItemById("2", "user")
	if err != errs.ItemNotExistError {
		t.Fatal("item should not exist")
	}

	// test delete
	err = db.DeleteTrashItem("1", "guest")
	if err != errs.PermissionError {
		t.Fatal("permission control not work")
	}
	err = db.DeleteTrashItem("2", "user")
	if err != errs.ItemNotExistError {
		t.Fatal("item should not exist")
	}
	err = db.DeleteTrashItem("1", "user")
	if err != nil {
		t.Fatal(err)
	}
	// after delete, should be empty
	list = db.ListTrashItems("user")
	if len(list) > 0 {
		t.Fatal("database should be empty")
	}
}
