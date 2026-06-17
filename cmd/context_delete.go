package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contextDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := cfg.DeleteContext(name); err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Context %q deleted.\n", name)
		if cfg.CurrentContext == "" {
			fmt.Println("No active context — use `k8shell context use <name>` to set one.")
		}
		return nil
	},
}
