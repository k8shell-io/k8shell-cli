// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userDeleteCmd = &cobra.Command{
	Use:               "delete <username>",
	Aliases:           []string{"del", "rm"},
	Short:             "Delete a user",
	Long:              "Permanently delete a user.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeUsernames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		username := args[0]
		if err := newClient(ctx).DeleteUser(cmd.Context(), username); err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"username": username, "deleted": true})
		}

		printer.Println(fmt.Sprintf("%s: deleted", username))
		return nil
	},
}
