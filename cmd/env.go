package cmd

import (
	"encoding/json"
	"github.com/Troublor/go-trash/storage/model"
	"github.com/Troublor/go-trash/system"
	"github.com/creamdog/gonfig"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// GOTRASH_PATH is set at compilation using go build -X options
var GOTRASH_PATH string

var config gonfig.Gonfig

var settingConfigFile = "trash-config.json"
var dbFile = "trash-bin.db"

/**
Get the instance of Setting object

At present, we read setting from json file
*/
func GetSetting() (*model.Setting, error) {
	file := path.Join(GOTRASH_PATH, settingConfigFile)
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

func GetConfig() gonfig.Gonfig {
	if config != nil {
		return config
	}
	file, err := os.Open(filepath.Join(system.GetTrashCmdDir(), GetConfigPath()))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	config, err = gonfig.FromJson(file)
	if err != nil {
		panic(err)
	}
	return config
}

func GetConfigPath() string {
	return path.Join(GOTRASH_PATH, settingConfigFile)
}

func GetDbPath() string {
	return path.Join(GOTRASH_PATH, dbFile)
}

func GetTrashBinPath() string {
	return path.Join(GOTRASH_PATH, "trash_bin")
}
