// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	k8shell "github.com/k8shell-io/k8shell-go"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var onboardRuleColumns = []table.Col[models.OnboardRule]{
	{Header: "ID", MaxWidth: 6, Help: "rule ID", Field: "id"},
	{Header: "ORG", MaxWidth: 15, Help: "destination organization for a matching user", Field: "org"},
	{Header: "IDP", MaxWidth: 12, Help: "identity provider, \"local\", or \"*\" (any)", Field: "idp"},
	{Header: "PATTERN", MaxWidth: 20, Help: "exact username, or a pattern containing '*'", Field: "usernamePattern"},
	{Header: "ACTION", MaxWidth: 10, Help: "allow, reject, or waitlist", Field: "action"},
	{Header: "STATUS", MaxWidth: 12, Help: "outcome of an actual onboarding attempt against this rule", Field: "status"},
	{Header: "PRIORITY", MaxWidth: 8, Help: "lower wins among matching rules of the same specificity", Field: "priority"},
	{Header: "ROLES", MaxWidth: 20, Help: "roles granted on allow (comma-separated)", Field: "roles", Fmt: table.FmtJoin},
	{Header: "SUDO", MaxWidth: 5, Help: "sudo access granted on allow", Field: "sudo", Fmt: table.FmtBool},
	{Header: "FULLNAME", MaxWidth: 20, Help: "requester's display name (system-inserted rows)", Field: "fullname"},
	{Header: "EMAIL", MaxWidth: 25, Help: "requester's email (system-inserted rows)", Field: "email"},
	{Header: "REQUESTED", MaxWidth: 16, Help: "when a system-inserted row was requested", Field: "requestedAt", Fmt: fmtTime},
	{Header: "DECIDED", MaxWidth: 16, Help: "when the rule's action left \"waitlist\"", Field: "decidedAt", Fmt: fmtTime},
	{Header: "DECIDED_BY", MaxWidth: 15, Help: "admin username who approved/rejected", Field: "decidedBy"},
	{Header: "NOTE", MaxWidth: 30, Help: "admin comment or rejection reason", Field: "note"},
	{Header: "CREATED", MaxWidth: 16, Help: "creation timestamp (local time)", Fn: func(r models.OnboardRule) string {
		return r.CreatedAt.Local().Format("2006-01-02 15:04")
	}},
}

var (
	onboardRuleListOrg     string
	onboardRuleListIDP     string
	onboardRuleListStatus  string
	onboardRuleListAction  string
	onboardRuleListPending bool
	onboardRuleSortFlag    string
)

var orgOnboardRuleListCmd = &cobra.Command{
	Use:     "list --org <org> [flags]",
	Aliases: []string{"ls"},
	Short:   "List onboarding rules",
	Long: "List the onboarding rules scoped to --org. Use --pending as a shortcut for --action waitlist,\n" +
		"the rules awaiting an admin decision.\n\n" + table.ColumnHelp(onboardRuleColumns),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if onboardRuleListPending && onboardRuleListAction != "" {
			return fmt.Errorf("--pending cannot be combined with --action")
		}

		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		filter := k8shell.OnboardRuleFilter{
			IDP:    onboardRuleListIDP,
			Status: onboardRuleListStatus,
			Action: onboardRuleListAction,
		}
		if onboardRuleListPending {
			filter.Action = "waitlist"
		}

		rules, err := newClient(ctx).ListOnboardRules(cmd.Context(), onboardRuleListOrg, filter)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(rules)
		}

		return table.Table(printer, onboardRuleColumns, rules, onboardRuleSortFlag)
	},
}

func init() {
	orgOnboardRuleListCmd.Flags().StringVar(&onboardRuleListOrg, "org", "", "organization the rules are scoped to (required)")
	orgOnboardRuleListCmd.Flags().StringVar(&onboardRuleListIDP, "idp", "", "filter by identity provider")
	orgOnboardRuleListCmd.Flags().StringVar(&onboardRuleListStatus, "status", "", "filter by status (none, pending, rejected, onboarded)")
	orgOnboardRuleListCmd.Flags().StringVar(&onboardRuleListAction, "action", "", "filter by action (allow, reject, waitlist)")
	orgOnboardRuleListCmd.Flags().BoolVar(&onboardRuleListPending, "pending", false, "shortcut for --action waitlist: rules awaiting an admin decision")
	orgOnboardRuleListCmd.Flags().StringVar(&onboardRuleSortFlag, "sort", "", "sort by fields, e.g. priority,-createdAt (prefix - for descending)")
	_ = orgOnboardRuleListCmd.MarkFlagRequired("org")
}
