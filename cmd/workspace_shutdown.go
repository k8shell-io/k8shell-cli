// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workspaceShutdownCmd = &cobra.Command{
	Use:               "shutdown <workspace-name>",
	Short:             "Shutdown a workspace, preserving its data",
	Long:              "Stop a workspace's pod without deleting its data. Use 'workspace delete' to permanently remove a workspace and its data.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeWorkspaceNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		if err := newClient(ctx).DeleteWorkspace(cmd.Context(), args[0], false); err != nil {
			return err
		}

		fmt.Printf("workspace %s shutdown\n", args[0])
		return nil
	},
}
