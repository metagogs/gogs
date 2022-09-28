package cmd

import (
	"fmt"

	"github.com/metagogs/gogs/tools/gogs/docker"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(dockerCmd)
	dockerCmd.Flags().BoolVarP(&proxycn, "cn", "c", false, "use golang proxy for china")
}

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "generate docker file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gen, err := docker.NewDockerGen(proxycn)
		if err != nil {
			fmt.Println(err)
		}
		_ = gen.Generate()
	},
}
