// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import "github.com/spf13/cobra"

var orgCmd = &cobra.Command{
	Use:     "org",
	Aliases: []string{"organization"},
	Short:   "Manage organizations",
}

func init() {
	orgCmd.AddCommand(orgListCmd)
	orgCmd.AddCommand(orgGetCmd)
	orgCmd.AddCommand(orgCreateCmd)
	orgCmd.AddCommand(orgSetCmd)
	orgCmd.AddCommand(orgDeleteCmd)
	orgCmd.AddCommand(orgOnboardRuleCmd)
}
