package cmd

import "github.com/spf13/cobra"

var userCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"usr"},
	Short:   "Manage users",
}

func init() {
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userSessionCmd)
}
