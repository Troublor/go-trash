package cmd

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	GOTRASH_PATH = "./tmp/"
	// create GOTRASH_PATH if not exist
	if _, err := os.Stat(GOTRASH_PATH); os.IsNotExist(err) {
		err := os.MkdirAll(GOTRASH_PATH, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	defer func() {
		// remove the GOTRASH_PATH
		err := os.RemoveAll(GetDbPath())
		if err != nil {
			panic(err)
		}
	}()

	// do initialization work
	initialize()
	defer func() {
		// remove the db
		err := os.RemoveAll(GetDbPath())
		if err != nil {
			panic(err)
		}
	}()
	defer func() {
		// close db
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_ = os.MkdirAll(GOTRASH_PATH, os.ModePerm)
	defer func() {
		_ = os.RemoveAll(GOTRASH_PATH)
	}()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_ = m.Run()
}
