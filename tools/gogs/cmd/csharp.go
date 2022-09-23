package cmd

import (
	"fmt"
	"os"

	"github.com/metagogs/gogs/tools/gogs/gen"
	"github.com/metagogs/gogs/utils/gomod"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(csharpCmd)
	csharpCmd.Flags().StringVarP(&protoFile, "file", "f", "", "proto文件")
}

var csharpCmd = &cobra.Command{
	Use:   "csharp",
	Short: "generate csharp code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(protoFile) == 0 {
			pterm.Error.Printfln("proto file is empty")
			os.Exit(1)
		}
		if _, err := os.Stat("go.mod"); err != nil {
			pterm.Error.Printfln("go.mod not found")
			os.Exit(1)
		}
		goModule, err := gomod.GetMod()
		if err != nil {
			pterm.Error.Printfln("get go mod error: " + err.Error())
			os.Exit(1)
		}
		if !goModule.IsInGoMod() {
			pterm.Error.Printfln("not in go mod mode")
			os.Exit(1)
		}

		gen, err := gen.NewGen(protoFile, goModule.Path)
		if err != nil {
			fmt.Println(err)
		}
		_ = gen.Generate()
	},
}
