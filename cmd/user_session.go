package cmd

import "github.com/spf13/cobra"

var userSessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"ses"},
	Short:   "Manage user sessions",
}

func init() {
	userSessionCmd.AddCommand(userSessionListCmd)
}
