package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// PROGRAM is the human readable product name
	PROGRAM = "set-by-main"

	// VERSION is the human readable product version
	VERSION = "set-by-main"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the program's version and exits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.Version)
	},
}
