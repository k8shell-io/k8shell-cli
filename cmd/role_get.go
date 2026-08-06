// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var roleGetOrg string

var roleGetCmd = &cobra.Command{
	Use:   "get <name> --org <org>",
	Short: "Show details for a single role",
	Long:  "Show details for a single role scoped to --org (or a global role visible within it).\n\n" + table.ColumnHelp(roleColumns),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		role, err := newClient(ctx).GetRole(cmd.Context(), roleGetOrg, args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(role)
		}

		return table.Detail(printer, roleColumns, *role)
	},
}

func init() {
	roleGetCmd.Flags().StringVar(&roleGetOrg, "org", "", "organization the role is scoped to (required)")
	_ = roleGetCmd.MarkFlagRequired("org")
}
