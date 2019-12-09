package cmd

import (
	"fmt"
	"github.com/Troublor/trash-go/service"
	"github.com/Troublor/trash-go/storage"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gotrash",
	Short: "Go-trash is a linux command-line trash files management tool",
	Long: `Go-trash is a linux command-line trash files management tool which provides
			 features similar to the Recycle Bin in Windows.
		   Developed by Troublor, 2019`,
}

func Execute() {
	service.MustEventHappen("onCmdStart")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		service.MustEventHappen("onCmdExitWithErr")
		os.Exit(-1)
	}
	service.MustEventHappen("onCmdEnd")
}

func init() {
	cobra.OnInitialize(initialize)

	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(settingCmd)
	rootCmd.AddCommand(ssCmd)
	rootCmd.AddCommand(urCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cleanCmd)
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
