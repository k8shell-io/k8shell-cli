// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var onboardRuleDeleteOrg string

var orgOnboardRuleDeleteCmd = &cobra.Command{
	Use:     "delete <id> --org <org>",
	Aliases: []string{"del"},
	Short:   "Delete an onboarding rule",
	Long:    "Permanently delete an onboarding rule by ID, as found in `onboard-rule list`.",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid rule ID %q: %w", args[0], err)
		}

		if err := newClient(ctx).DeleteOnboardRule(cmd.Context(), onboardRuleDeleteOrg, int32(id)); err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"id": id, "deleted": true})
		}

		printer.Println(fmt.Sprintf("%d: deleted", id))
		return nil
	},
}

func init() {
	orgOnboardRuleDeleteCmd.Flags().StringVar(&onboardRuleDeleteOrg, "org", "", "organization the rule is scoped to (required)")
	_ = orgOnboardRuleDeleteCmd.MarkFlagRequired("org")
}
