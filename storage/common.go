package storage

import (
	"github.com/creamdog/gonfig"
	"github.com/otiai10/copy"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
)

var config gonfig.Gonfig

func GetTrashCmdDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func GetConfig() gonfig.Gonfig {
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

func IsSudo() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		panic(err)
	}
	if i == 0 {
		return true
	} else {
		return false
	}
}

func GetDbPath() string {
	if IsSudo() {
		r, err := GetConfig().GetString("trashDir", "wrong path")
		if err != nil {
			panic(err)
		}
		return path.Join(r, "root_trash_info.db")
	} else {
		r, err := GetConfig().GetString("trashDir", "wrong path")
		if err != nil {
			panic(err)
		}
		return path.Join(r, "trash_info.db")
	}

}

func GetTrashPath() string {
	if IsSudo() {
		r, err := GetConfig().GetString("trashDir", "wrong path")
		if err != nil {
			panic(err)
		}
		return path.Join(r, "root_trash_bin")
	} else {
		r, err := GetConfig().GetString("trashDir", "wrong path")
		if err != nil {
			panic(err)
		}
		return path.Join(r, "trash_bin")
	}
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
