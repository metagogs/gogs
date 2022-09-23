package cmd

import (
	"fmt"
	"os"

	"github.com/metagogs/gogs/tools/gogs/csharp"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(csharpCmd)
	csharpCmd.Flags().StringVarP(&protoFile, "file", "f", "", "proto file")
	csharpCmd.Flags().BoolVarP(&gogsFiles, "gogs", "g", false, "should generate gogs framework code")
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

		fmt.Println(gogsFiles)
		gen, err := csharp.NewCSharpGen(protoFile, gogsFiles)
		if err != nil {
			fmt.Println(err)
		}
		_ = gen.Generate()
	},
}
