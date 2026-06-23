// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var shutdownDelete bool

var workspaceShutdownCmd = &cobra.Command{
	Use:               "shutdown <workspace-name>",
	Short:             "Shutdown a workspace",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeWorkspaceNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		if err := newClient(ctx).DeleteWorkspace(cmd.Context(), args[0], shutdownDelete); err != nil {
			return err
		}

		fmt.Printf("workspace %s shutdown\n", args[0])
		return nil
	},
}

func init() {
	workspaceShutdownCmd.Flags().BoolVar(&shutdownDelete, "delete", false, "permanently delete workspace data")
}
