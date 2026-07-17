// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workspaceStartCmd = &cobra.Command{
	Use:               "start <workspace-name>",
	Short:             "Start a stopped workspace",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeWorkspaceNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		if err := newClient(ctx).StartWorkspace(cmd.Context(), args[0]); err != nil {
			return err
		}

		fmt.Printf("workspace %s started\n", args[0])
		return nil
	},
}
