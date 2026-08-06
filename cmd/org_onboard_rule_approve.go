// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var onboardRuleApproveOrg string

var orgOnboardRuleApproveCmd = &cobra.Command{
	Use:   "approve <id> --org <org>",
	Short: "Approve a pending onboarding rule",
	Long: "Approve a pending (\"waitlist\") onboarding rule, flipping its action to \"allow\", and immediately\n" +
		"onboard the user it names rather than waiting for their next login attempt.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid rule ID %q: %w", args[0], err)
		}

		user, err := newClient(ctx).ApproveOnboardRule(cmd.Context(), onboardRuleApproveOrg, int32(id))
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(user)
		}

		printer.Println(fmt.Sprintf("%d: approved (onboarded as %s)", id, user.Username))
		return nil
	},
}

func init() {
	orgOnboardRuleApproveCmd.Flags().StringVar(&onboardRuleApproveOrg, "org", "", "organization the rule is scoped to (required)")
	_ = orgOnboardRuleApproveCmd.MarkFlagRequired("org")
}
