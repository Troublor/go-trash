package operation

import (
	"github.com/creamdog/gonfig"
	"github.com/otiai10/copy"
	"os"
	"path"
	"path/filepath"
)

var Config gonfig.Gonfig

func init() {
	file, err := os.Open("../config.json")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	Config, err = gonfig.FromJson(file)
	if err != nil {
		panic(err)
	}
}

func GetDbPath() string {
	r, err := Config.GetString("trashDir", "wrong path")
	if err != nil {
		panic(err)
	}
	return path.Join(r, "trash_info.db")
}

func GetTrashPath() string {
	r, err := Config.GetString("trashDir", "wrong path")
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
	return err
}
