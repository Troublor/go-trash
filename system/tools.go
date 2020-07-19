package system

import (
	"github.com/Troublor/go-trash/errs"
	"github.com/otiai10/copy"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetUser() string {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(output))
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

func IsTesting() bool {
	return strings.Contains(os.Args[0], "_test")
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

func GetTrashCmdDir() string {
	r, _ := os.Executable()
	dir, err := filepath.Abs(filepath.Dir(r))
	if err != nil {
		panic(err)
	}
	return dir
}
