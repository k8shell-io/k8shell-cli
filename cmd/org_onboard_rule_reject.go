// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	onboardRuleRejectOrg  string
	onboardRuleRejectNote string
)

var orgOnboardRuleRejectCmd = &cobra.Command{
	Use:   "reject <id> --org <org> [flags]",
	Short: "Reject a pending onboarding rule",
	Long: "Reject a pending (\"waitlist\") onboarding rule, flipping its action to \"reject\" so the user\n" +
		"cannot re-trigger a new waitlist entry by trying again.",
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

		req := models.OnboardRuleRejectRequest{Note: onboardRuleRejectNote}
		rule, err := newClient(ctx).RejectOnboardRule(cmd.Context(), onboardRuleRejectOrg, int32(id), req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(rule)
		}

		printer.Println(fmt.Sprintf("%d: rejected", id))
		return nil
	},
}

func init() {
	orgOnboardRuleRejectCmd.Flags().StringVar(&onboardRuleRejectOrg, "org", "", "organization the rule is scoped to (required)")
	orgOnboardRuleRejectCmd.Flags().StringVar(&onboardRuleRejectNote, "note", "", "rejection reason")
	_ = orgOnboardRuleRejectCmd.MarkFlagRequired("org")
}
