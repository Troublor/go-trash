package main

import (
	"fmt"
	"github.com/Troublor/go-trash/cmd"
)

func main() {
	fmt.Println(cmd.GOTRASH_PATH)
	cmd.Execute()
}
