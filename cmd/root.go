package cmd

import (
	"fmt"
	"github.com/Troublor/go-trash/service"
	"github.com/Troublor/go-trash/storage"
	"github.com/spf13/cobra"
	"os"
)
var db *storage.Database

var RootCmd = &cobra.Command{
	Use:   "gotrash",
	Short: "Go-trash is a linux command-line trash files management tool",
	Long: `Go-trash is a linux command-line trash files management tool which provides
			 features similar to the Recycle Bin in Windows.
		   Developed by Troublor, 2019`,
	PostRun: func(cmd *cobra.Command, args []string) {
		// after execution, close db
		if db != nil {
			_ = db.Close()
		}
	},
	Version: storage.Version(),
}

func Execute() {
	service.MustEventHappen("onCmdStart")
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		service.MustEventHappen("onCmdExitWithErr")
		os.Exit(-1)
	}
	service.MustEventHappen("onCmdEnd")
}

func init() {
	cobra.OnInitialize(initialize)

	RootCmd.AddCommand(lsCmd)
	RootCmd.AddCommand(rmCmd)
	RootCmd.AddCommand(ssCmd)
	RootCmd.AddCommand(urCmd)
	RootCmd.AddCommand(cleanCmd)
}

func initialize() {
	db = storage.NewDatabase(GetDbPath())
	err := db.Open()
	if err != nil {
		panic(err)
	}
	//if system.IsSudo() {
	//	openPermission := func(event service.Event) {
	//		cmd := exec.Command("sudo chmod 666 -R " + storage.GetDbPath() + " " + storage.GetTrashBinPath())
	//		_, err := cmd.Output()
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//	service.MustSubscribeEvent("onCmdEnd", openPermission)
	//	service.MustSubscribeEvent("onCmdExitWithErr", openPermission)
	//}
}
