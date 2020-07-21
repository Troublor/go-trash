package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// TODO generate default setting configuration json file
//go:generate

type Setting struct {
	AutoRemove         bool          `json:"autoRemove"`         // whether to auto remove trashes in the bin
	AutoRemoveInterval time.Duration `json:"autoRemoveInterval"` // only valid when AutoRemove == true, auto remove interval period
}

func GenDefaultSettingJsonFile(filePath string) error {
	// default setting
	setting := Setting{
		AutoRemove:         false,
		AutoRemoveInterval: 0,
	}
	str, err := json.Marshal(&setting)
	if err != nil {
		panic(err)
	}
	return ioutil.WriteFile(filePath, str, os.ModePerm)
}
