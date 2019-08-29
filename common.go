package main

import (
	"os/exec"
	"strconv"
)

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
