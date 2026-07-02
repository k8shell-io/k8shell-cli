// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workspaceDeleteCmd = &cobra.Command{
	Use:               "delete <workspace-name>",
	Aliases:           []string{"del"},
	Short:             "Permanently delete a workspace and its data",
	Long:              "Stop a workspace's pod and permanently delete its data. Use 'workspace shutdown' to stop a workspace without deleting its data.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeWorkspaceNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		if err := newClient(ctx).DeleteWorkspace(cmd.Context(), args[0], true); err != nil {
			return err
		}

		fmt.Printf("workspace %s deleted\n", args[0])
		return nil
	},
}
