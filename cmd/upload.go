package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [filter]",
	Args:  cobra.MinimumNArgs(1),
	Short: "(re)Upload a back catalog of replays specified",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}
