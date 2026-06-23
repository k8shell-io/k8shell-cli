// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import "github.com/spf13/cobra"

var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage workspaces",
}

func init() {
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceCreateCmd)
	workspaceCmd.AddCommand(workspaceJobEventsCmd)
	workspaceCmd.AddCommand(workspaceShutdownCmd)
	workspaceCmd.AddCommand(workspaceSSHCmd)
	workspaceCmd.AddCommand(workspaceSCPCmd)
	workspaceCmd.AddCommand(workspaceCodeCmd)
}
