package storage

import (
	"errors"
	"github.com/Troublor/go-trash/errs"
	"github.com/Troublor/go-trash/storage/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"path"
	"time"
)

type Database struct {
	dbPath string
	orm    *gorm.DB
}

func NewDatabase(dbPath string) *Database {
	return &Database{dbPath: dbPath}
}

/**
Open database sqlite3 connection, must close after use.
Schema migration will be performed if necessary.
*/
func (db *Database) Open() error {
	var err error
	// open database sqlite3
	db.orm, err = gorm.Open("sqlite3", db.dbPath)
	if err != nil {
		return err
	}
	// migrate the database schema
	db.orm.AutoMigrate(&model.TrashMetadata{})
	return nil
}

/**
Close database connection
*/
func (db *Database) Close() error {
	err := db.orm.Close()
	return err
}

/**
Add a new trash record
*/
func (db *Database) InsertTrashItem(id, originalPath, trashDir, baseName, itemType, owner string) error {
	if itemType != model.TYPE_FILE && itemType != model.TYPE_DIRECTORY {
		panic(errors.New("invalid value for argument itemType: " + itemType))
	}
	deleteTime := time.Now()

	trashMetadata := model.TrashMetadata{
		ID:           id,
		OriginalPath: originalPath,
		TrashPath:    path.Join(trashDir, id),
		BaseName:     baseName,
		Type:         itemType,
		Owner:        owner,
		CreatedAt:    deleteTime,
	}
	if !db.orm.Where("id = ?", id).First(&model.TrashMetadata{}).RecordNotFound() {
		panic(errors.New("trash id " + id + " exists, this should not happen"))
	}
	if err := db.orm.Create(&trashMetadata).Error; err != nil {
		return err
	}
	return nil
}

/**
List all trash items
*/
func (db *Database) ListTrashItems(user string) model.TrashMetadataList {
	var metadataList []model.TrashMetadata
	db.orm.Where("owner = ?", user).Find(&metadataList)
	return metadataList
}

/**
Delete one trash item by id
*/
func (db *Database) DeleteTrashItem(id, user string) error {
	if id == "" {
		panic(errors.New("delete trash item id is empty"))
	}
	// get the item to be deleted
	var deletedItem model.TrashMetadata
	if db.orm.Where("id = ?", id).First(&deletedItem).RecordNotFound() {
		return errs.ItemNotExistError
	}
	if deletedItem.Owner != user {
		return errs.PermissionError
	}
	return db.orm.Delete(&deletedItem).Error
}

/**
Get trash item by id
*/
func (db *Database) GetTrashItemById(id string, user string) (model.TrashMetadata, error) {
	var metadata model.TrashMetadata
	if db.orm.Where("id = ?", id).First(&metadata).RecordNotFound() {
		return model.TrashMetadata{}, errs.ItemNotExistError
	}
	if metadata.Owner != user {
		return model.TrashMetadata{}, errs.PermissionError
	}
	return metadata, nil
}
