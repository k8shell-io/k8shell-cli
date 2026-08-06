// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userDeletePreserveWorkspaces bool

var userDeleteCmd = &cobra.Command{
	Use:               "delete <username>",
	Aliases:           []string{"del"},
	Short:             "Delete a user",
	Long:              "Permanently delete a user. Deleting a user also deletes their workspaces, unless --preserve-workspaces is set.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeUsernames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		username := args[0]
		if err := newClient(ctx).DeleteUser(cmd.Context(), username, userDeletePreserveWorkspaces); err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"username": username, "deleted": true})
		}

		printer.Println(fmt.Sprintf("%s: deleted", username))
		return nil
	},
}

func init() {
	userDeleteCmd.Flags().BoolVar(&userDeletePreserveWorkspaces, "preserve-workspaces", false, "keep the user's workspaces instead of deleting them along with the account")
}
