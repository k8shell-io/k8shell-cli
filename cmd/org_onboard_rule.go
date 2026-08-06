// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import "github.com/spf13/cobra"

var orgOnboardRuleCmd = &cobra.Command{
	Use:     "onboard-rule",
	Aliases: []string{"onboard-rules"},
	Short:   "Manage onboarding rules",
}

func init() {
	orgOnboardRuleCmd.AddCommand(orgOnboardRuleListCmd)
	orgOnboardRuleCmd.AddCommand(orgOnboardRuleCreateCmd)
	orgOnboardRuleCmd.AddCommand(orgOnboardRuleSetCmd)
	orgOnboardRuleCmd.AddCommand(orgOnboardRuleDeleteCmd)
	orgOnboardRuleCmd.AddCommand(orgOnboardRuleApproveCmd)
	orgOnboardRuleCmd.AddCommand(orgOnboardRuleRejectCmd)
}
