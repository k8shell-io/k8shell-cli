// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import "github.com/spf13/cobra"

var userCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"usr"},
	Short:   "Manage users",
}

func init() {
	userCmd.AddCommand(userListCmd)
}
