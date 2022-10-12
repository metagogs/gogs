package cmd

import (
	"github.com/metagogs/gogs/tools/gogs/gen"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	gomodpath string
	gopackage string
)

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&gomodpath, "pakcage", "p", "", "the go mod package")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init gogs project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(gomodpath) == 0 {
			gomodpath = "gogs"
			pterm.Warning.Printfln("You didn't set the package name. I'll use the default package name - gogs")
		}
		initGen := gen.NewInit(gomodpath)
		err := initGen.Generate()
		if err != nil {
			pterm.Error.Println("generate error: " + err.Error())
		}
	},
}
