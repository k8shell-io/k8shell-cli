// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var roleDeleteOrg string

var roleDeleteCmd = &cobra.Command{
	Use:     "delete <name> --org <org>",
	Aliases: []string{"del"},
	Short:   "Delete a role",
	Long:    "Permanently delete a role scoped to --org. Fails if any user still holds the role.",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		name := args[0]
		if err := newClient(ctx).DeleteRole(cmd.Context(), roleDeleteOrg, name); err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"name": name, "deleted": true})
		}

		printer.Println(fmt.Sprintf("%s: deleted", name))
		return nil
	},
}

func init() {
	roleDeleteCmd.Flags().StringVar(&roleDeleteOrg, "org", "", "organization the role is scoped to (required)")
	_ = roleDeleteCmd.MarkFlagRequired("org")
}
