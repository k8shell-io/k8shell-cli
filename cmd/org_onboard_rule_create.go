// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	onboardRuleCreateOrg      string
	onboardRuleCreateIDP      string
	onboardRuleCreateAction   string
	onboardRuleCreatePriority int32
	onboardRuleCreateRoles    []string
	onboardRuleCreateSudo     bool
	onboardRuleCreateNote     string
)

var orgOnboardRuleCreateCmd = &cobra.Command{
	Use:   "create <pattern> --org <org> --idp <idp> --action allow|reject|waitlist [flags]",
	Short: "Create a new onboarding rule",
	Long: "Register a new onboarding rule scoped to --org: a standing pattern policy (pattern contains '*'),\n" +
		"or a one-off decision for a specific username. On \"allow\", the listed --roles (and --sudo) are\n" +
		"granted to the user.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		req := models.OnboardRuleCreateRequest{
			IDP:             onboardRuleCreateIDP,
			UsernamePattern: args[0],
			Action:          onboardRuleCreateAction,
			Priority:        onboardRuleCreatePriority,
			Roles:           onboardRuleCreateRoles,
			Sudo:            onboardRuleCreateSudo,
			Note:            onboardRuleCreateNote,
		}

		rule, err := newClient(ctx).CreateOnboardRule(cmd.Context(), onboardRuleCreateOrg, req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(rule)
		}

		printer.Println(fmt.Sprintf("%d: created", rule.ID))
		return nil
	},
}

func init() {
	orgOnboardRuleCreateCmd.Flags().StringVar(&onboardRuleCreateOrg, "org", "", "destination organization for users this rule matches (required)")
	orgOnboardRuleCreateCmd.Flags().StringVar(&onboardRuleCreateIDP, "idp", "", "identity provider, \"local\", or \"*\" for any (required)")
	orgOnboardRuleCreateCmd.Flags().StringVar(&onboardRuleCreateAction, "action", "", "allow, reject, or waitlist (required)")
	orgOnboardRuleCreateCmd.Flags().Int32Var(&onboardRuleCreatePriority, "priority", 0, "lower wins among matching rules of the same specificity")
	orgOnboardRuleCreateCmd.Flags().StringSliceVar(&onboardRuleCreateRoles, "roles", nil, "roles granted to the user when this rule resolves to allow, comma-separated")
	orgOnboardRuleCreateCmd.Flags().BoolVar(&onboardRuleCreateSudo, "sudo", false, "grant sudo access when this rule resolves to allow")
	orgOnboardRuleCreateCmd.Flags().StringVar(&onboardRuleCreateNote, "note", "", "admin comment")
	_ = orgOnboardRuleCreateCmd.MarkFlagRequired("org")
	_ = orgOnboardRuleCreateCmd.MarkFlagRequired("idp")
	_ = orgOnboardRuleCreateCmd.MarkFlagRequired("action")
}
