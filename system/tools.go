package system

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetUser() string {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(output)
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
		return false
	} else {
		return false
	}
}

func IsTesting() bool {
	return strings.Contains(os.Args[0], "_test")
}
