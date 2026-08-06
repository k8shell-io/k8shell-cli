// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import "github.com/spf13/cobra"

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage roles",
}

func init() {
	roleCmd.AddCommand(roleListCmd)
	roleCmd.AddCommand(roleGetCmd)
	roleCmd.AddCommand(roleCreateCmd)
	roleCmd.AddCommand(roleSetCmd)
	roleCmd.AddCommand(roleDeleteCmd)
}
