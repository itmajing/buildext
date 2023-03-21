package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(dockerCmd)
}

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Simplified docker cli",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}
