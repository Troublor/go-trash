package storage

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/Troublor/trash-go/errs"
	"github.com/Troublor/trash-go/service"
	"github.com/Troublor/trash-go/system"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/exec"
	"path"
	"time"
)

var database *sql.DB

func initDB() {
	initializeSettings := func() {
		for key, value := range DefaultSettings {
			err := AddSetting(key, value)
			if err != nil && err != errs.ItemExistError {
				panic(err)
			}
		}
	}
	var err error
	//create trash bin directory
	if _, err = os.Stat(GetTrashPath()); err != nil {
		err = os.MkdirAll(GetTrashPath(), os.ModePerm)
		if err != nil {
			panic(err)
		}
		if system.IsTesting() {
			_ = service.SubscribeEvent("onTestEnd", func(event service.Event) {
				_ = os.RemoveAll(GetTrashPath())
			})
		}
	}
	database, err = sql.Open("sqlite3", GetDbPath())
	if err != nil {
		panic(err.Error())
	}
	if system.IsTesting() {
		_ = service.SubscribeEvent("onTestEnd", func(event service.Event) {
			_ = os.Remove(GetDbPath())
		})
	}
	defer func() {
		if system.IsSudo() {
			cmd := exec.Command("sudo chmod 777 -R " + GetDbPath())
			_, err := cmd.Output()
			if err != nil {
				panic(err)
			}
		}
	}()
	// Create Settings table if it does not exist
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS settings (setting_key TEXT PRIMARY KEY, setting_value TEXT)")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	_, err = statement.Exec()
	if err != nil {
		panic(err.Error())
	}
	initializeSettings()
	//Create trash_info table if it does not exist
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS trash_info (id TEXT PRIMARY KEY, original_path TEXT, trash_path TEXT, base_name TEXT, item_type TEXT, owner TEXT, delete_time TIMESTAMP )")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	_, err = statement.Exec()
	if err != nil {
		panic(err.Error())
	}
}

func GetSetting(key string) (value string, err error) {
	statement, e := database.Prepare("SELECT setting_value FROM settings WHERE setting_key=?")
	if e != nil {
		panic(e.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	e = statement.QueryRow(key).Scan(&value)
	if e != nil {
		if e == sql.ErrNoRows {
			err = e
		} else {
			panic(e.Error())
		}
	}
	return
}

func AddSetting(key string, value string) error {
	statement, err := database.Prepare("INSERT INTO settings (setting_key, setting_value) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	_, err = statement.Exec(key, value)
	if err != nil {
		if err, ok := err.(sqlite3.Error); ok {
			//fmt.Println(int(err.Code))
			//fmt.Println(int(err.ExtendedCode))
			if err.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return errs.ItemExistError
			}
		}
		panic(err.Error())
	}
	return nil
}

func UpdateSetting(key string, value string) error {
	statement, err := database.Prepare("UPDATE TABLE settings SET setting_value = ? WHERE setting_key = ?")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	result, err := statement.Exec(value, key)
	if err != nil {
		panic(err.Error())
	}
	n, _ := result.RowsAffected()
	if n < 1 {
		return errors.New("the setting '" + key + "' does not exist")
	} else {
		return nil
	}
}

func ListAllSettings() SettingList {
	rows, err := database.Query("SELECT * FROM settings")
	if err != nil {
		panic(err.Error())
	}
	defer func(r *sql.Rows) {
		err := r.Close()
		if err != nil {
			panic(err.Error())
		}
	}(rows)
	var results []Setting
	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, *NewSetting(key, value))
	}
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}
	return results
}

func DbInsertTrashItem(originalPath, trashDir, baseName, itemType string, owner string) string {
	if itemType != TYPE_FILE && itemType != TYPE_DIRECTORY {
		panic(errors.New("invalid value for argument itemType: " + itemType))
	}
	deleteTime := time.Now().Format("2006-01-02 15:04:05")
	hasher := md5.New()
	hasher.Write([]byte(deleteTime))
	id := hex.EncodeToString(hasher.Sum(nil))[:6]
	statement, err := database.Prepare("INSERT INTO trash_info (id, original_path, trash_path, base_name, item_type, owner, delete_time) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	trashPath := path.Join(trashDir, id)
	_, err = statement.Exec(id, originalPath, trashPath, baseName, itemType, owner, deleteTime)
	if err != nil {
		panic(err.Error())
	}
	return id
}

func DbListAllTrashItems(user string) TrashInfoList {
	rows, err := database.Query("SELECT * FROM trash_info")
	if err != nil {
		panic(err.Error())
	}
	defer func(r *sql.Rows) {
		err := r.Close()
		if err != nil {
			panic(err.Error())
		}
	}(rows)
	var results []TrashInfo
	for rows.Next() {
		var id, originalPath, trashPath, baseName, itemType, owner string
		var deleteTime time.Time
		err = rows.Scan(&id, &originalPath, &trashPath, &baseName, &itemType, &owner, &deleteTime)
		if err != nil {
			panic(err.Error())
		}
		if user == "root" || user == owner {
			results = append(results, TrashInfo{Id: id, OriginalPath: originalPath, TrashPath: trashPath, BaseName: baseName, ItemType: itemType, Owner: owner, DeleteTime: deleteTime})
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}
	return results
}

func DbDeleteTrashItem(id string, user string) error {
	statement, err := database.Prepare("SELECT owner FROM trash_info WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	var owner string
	err = statement.QueryRow(id).Scan(&owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ItemNotExistError
		} else {
			return err
		}
	}
	if owner != user {
		return errs.PermissionError
	}
	statement, err = database.Prepare("DELETE FROM trash_info WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	r, err := statement.Exec(id)
	if err != nil {
		panic(err.Error())
	}
	n, err := r.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	if n < 1 {
		return errs.ItemNotExistError
	}
	return nil
}

func DbGetTrashItemById(id string, user string) (*TrashInfo, error) {
	statement, err := database.Prepare("SELECT original_path, trash_path, base_name, item_type, owner, delete_time FROM trash_info WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer func(s *sql.Stmt) {
		err := s.Close()
		if err != nil {
			panic(err.Error())
		}
	}(statement)
	var originalPath, trashPath, baseName, itemType, owner string
	var deleteTime time.Time
	err = statement.QueryRow(id).Scan(&originalPath, &trashPath, &baseName, &itemType, &owner, &deleteTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return &TrashInfo{}, errs.ItemNotExistError
		}
	}
	if user != owner {
		return &TrashInfo{}, errs.PermissionError
	}
	return &TrashInfo{Id: id, OriginalPath: originalPath, TrashPath: trashPath, BaseName: baseName, ItemType: itemType, Owner: owner, DeleteTime: deleteTime}, nil
}
