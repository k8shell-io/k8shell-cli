// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startEvents bool

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

		c := newClient(ctx)

		resp, err := c.StartWorkspace(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(resp)
		}

		if startEvents {
			fmt.Printf("Starting workspace %s (job %s)\n", resp.Workspace, resp.JobID)
		} else {
			fmt.Printf("Starting workspace %s...", resp.Workspace)
		}

		rc, err := c.MonitorWorkspace(cmd.Context(), resp.MonitorURL)
		if err != nil {
			if !startEvents {
				fmt.Println()
			}
			return fmt.Errorf("monitoring workspace: %w", err)
		}
		defer rc.Close()

		if startEvents {
			return printEventStream(rc)
		}

		// Default mode: update a single progress line in place.
		return printProgressStream(rc, "Starting", resp.Workspace)
	},
}

func init() {
	workspaceStartCmd.Flags().BoolVar(&startEvents, "events", false, "show log events instead of progress percentage")
}
