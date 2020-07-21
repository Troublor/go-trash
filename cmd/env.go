package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Troublor/go-trash/storage/model"
	"github.com/Troublor/go-trash/system"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

// GOTRASH_PATH is set at compilation using go build -X options, this variable should not be directly accessed.
// GoTrashPath should be accessed by GetGoTrashPath()
var GOTRASH_PATH string

func GetGoTrashPath() string {
	if GOTRASH_PATH != "" {
		return GOTRASH_PATH
	} else if system.GetUser() == "root" {
		return "/etc/gotrash"
	} else {
		u, err := user.Current()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		return path.Join(u.HomeDir, ".gotrash")
	}
}

var settingConfigFile = "trash-config.json"
var dbFile = "trash-bin.db"

/**
Get the instance of Setting object

At present, we read setting from json file
*/
func GetSetting() (*model.Setting, error) {
	file := path.Join(GetGoTrashPath(), settingConfigFile)
	// check if setting file exists, create default file if not exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("Creating default setting file at", file)
		err := model.GenDefaultSettingJsonFile(file)
		if err != nil {
			return nil, err
		}
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var setting model.Setting
	err = json.Unmarshal(content, &setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func GetDbPath() string {
	return path.Join(GetGoTrashPath(), dbFile)
}

func GetTrashBinPath() string {
	return path.Join(GetGoTrashPath(), "trash_bin")
}
