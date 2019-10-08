package storage

import (
	"github.com/Troublor/trash-go/errs"
	"github.com/Troublor/trash-go/system"
	"github.com/creamdog/gonfig"
	"github.com/otiai10/copy"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var config gonfig.Gonfig

func InitStorage() {
	initDB()
}

func GetTrashCmdDir() string {
	if system.IsTesting() {
		return ""
	} else {
		r, _ := os.Executable()
		dir, err := filepath.Abs(filepath.Dir(r))
		if err != nil {
			panic(err)
		}
		return dir
	}

}

func GetConfig() gonfig.Gonfig {
	if system.IsTesting() {
		currentDir, _ := os.Getwd()
		testPayload := `{"trashDir":"` + currentDir + `"}`
		s := strings.NewReader(testPayload)
		config, err := gonfig.FromJson(s)
		if err != nil {
			panic(err)
		}
		return config
	} else {
		if config != nil {
			return config
		}
		file, err := os.Open(filepath.Join(GetTrashCmdDir(), "gotrash-config.json"))
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

}

func GetDbPath() string {
	r, err := GetConfig().GetString("trashDir", "wrong path")
	if err != nil {
		panic(err)
	}
	return path.Join(r, "trash_info.db")
}

func GetTrashPath() string {
	r, err := GetConfig().GetString("trashDir", "wrong path")
	if err != nil {
		panic(err)
	}
	return path.Join(r, "trash_bin")
}

func GetAbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

func SafeRename(originalPath, targetPath string) error {
	if _, err := os.Stat(originalPath); err != nil {
		return errs.ItemNotExistError
	}
	if _, err := os.Stat(targetPath); err == nil {
		return errs.ItemExistError
	}
	err := os.Rename(originalPath, targetPath)
	if err != nil {
		// rename failed, try copy and delete
		err = copy.Copy(originalPath, targetPath)
		if err != nil {
			return err
		}
		err = os.RemoveAll(originalPath)
		if err != nil {
			return err
		}
	}
	return nil
}
