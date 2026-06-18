// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import "github.com/spf13/cobra"

var sessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"ses"},
	Short:   "Manage sessions",
}

func init() {
	sessionCmd.AddCommand(sessionListCmd)
}
