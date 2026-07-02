// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var workspaceGetCmd = &cobra.Command{
	Use:               "get <workspace-name>",
	Short:             "Show details for a single workspace",
	Long:              "Show details for a single workspace.\n\n" + table.ColumnHelp(workspaceColumns),
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeWorkspaceNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		ws, err := newClient(ctx).GetWorkspace(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(ws)
		}

		return table.Table(printer, workspaceColumns, []models.WorkspaceDetails{*ws}, "")
	},
}
