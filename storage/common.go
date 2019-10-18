package storage

import (
	"github.com/Troublor/trash-go/system"
	"github.com/creamdog/gonfig"
	"os"
	"path"
	"path/filepath"
)

var config gonfig.Gonfig

func GetConfig() gonfig.Gonfig {
	if config != nil {
		return config
	}
	file, err := os.Open(filepath.Join(system.GetTrashCmdDir(), "gotrash-config.json"))
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

func GetDbPath() string {
	r, err := GetConfig().GetString("trashDir", "wrong path")
	if err != nil {
		panic(err)
	}
	return path.Join(r, "trash_info.db")
}

func GetTrashBinPath() string {
	r, err := GetConfig().GetString("trashDir", "wrong path")
	if err != nil {
		panic(err)
	}
	return path.Join(r, "trash_bin")
}
