package cmd

import (
	"fmt"
	"github.com/Troublor/trash-go/cmd/setting"
	"github.com/Troublor/trash-go/service"
	"github.com/Troublor/trash-go/storage"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "gotrash",
	Short: "Go-trash is a linux command-line trash files management tool",
	Long: `Go-trash is a linux command-line trash files management tool which provides
			 features similar to the Recycle Bin in Windows.
		   Developed by Troublor, 2019`,
	Version: storage.VersionString(),
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
	RootCmd.AddCommand(setting.RootCmd)
	RootCmd.AddCommand(ssCmd)
	RootCmd.AddCommand(urCmd)
	RootCmd.AddCommand(cleanCmd)
}

func initialize() {
	storage.InitStorage()
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
