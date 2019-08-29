package main

import (
	"github.com/Troublor/trash-go/operation"
	"path/filepath"
	"testing"
)

func TestAbsPath(t *testing.T) {
	if absPath, _ := filepath.Abs("abc.txt"); absPath != "/home/troublor/workspace/go/trash-go/abc.txt" {
		t.Fatal("abs wrong")
	}
	if operation.GetAbsPath("abc.txt") != "/home/troublor/workspace/go/trash-go/abc.txt" {
		t.Fatal("abs wrong")
	}
}
