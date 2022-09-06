package main

import (
	"fmt"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/tools/gogs/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Long:  `show the gogs version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(gogs.Version)
	},
}

func main() {
	cmd.Execute()
}
