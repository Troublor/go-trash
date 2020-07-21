package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Troublor/go-trash/storage/model"
	"github.com/Troublor/go-trash/system"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
)

// GOTRASH_PATH is set at compilation using go build -X options, this variable should not be directly accessed.
// GoTrashPath should be accessed by GetGoTrashPath()
var GOTRASH_PATH string

func GetGoTrashPath() string {
	p := os.Getenv("GOTRASH_PATH")
	if p != "" {
		// override GOTRASH_PATH if environment is set globally
		return p
	}
	if GOTRASH_PATH != "" {
		// return GOTRASH_PATH if set at compilation time
		return GOTRASH_PATH
	}
	// default GOTRASH_PATH
	if system.GetUser() == "root" {
		return "/etc/gotrash"
	} else {
		return path.Join(os.Getenv("HOME"), ".gotrash")
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

var envCmd = &cobra.Command{
	Use:   "env items...",
	Short: "Show gotrash environment variables",
	Long:  `Show gotrash environment variables`,
	Run: func(cmd *cobra.Command, args []string) {
		var envVars = map[string]string{
			"GOTRASH_PATH": GetGoTrashPath(),
		}
		if len(args) == 0 {
			// show all environment variables
			for name, value := range envVars {
				fmt.Printf("%s=%s\n", name, value)
			}
		} else {
			for _, arg := range args {
				variable, ok := envVars[arg]
				if !ok {
					fmt.Println("environment variable", arg, "not found")
				} else {
					fmt.Printf("%s=%s\n", arg, variable)
				}
			}
		}
	},
}
