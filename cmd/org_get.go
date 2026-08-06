// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var orgGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Show details for a single organization",
	Long:  "Show details for a single organization.\n\n" + table.ColumnHelp(orgColumns),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		org, err := newClient(ctx).GetOrganization(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(org)
		}

		return table.Detail(printer, orgColumns, *org)
	},
}
