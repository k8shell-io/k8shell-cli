// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	k8shell "github.com/k8shell-io/k8shell-go"
	"github.com/spf13/cobra"
)

var (
	onboardRuleSetOrg      string
	onboardRuleSetAction   string
	onboardRuleSetPriority int32
	onboardRuleSetRoles    []string
	onboardRuleSetSudo     bool
	onboardRuleSetNote     string
)

var orgOnboardRuleSetCmd = &cobra.Command{
	Use:   "set <id> --org <org> [flags]",
	Short: "Update fields on an onboarding rule",
	Long: "Update an onboarding rule's action, priority, roles, sudo, and/or note, identified by ID as found\n" +
		"in `onboard-rule list`. idp/pattern/org are immutable — delete and recreate the rule to change them.\n" +
		"Fields left unspecified keep their current value.",
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

		c := newClient(ctx)

		rules, err := c.ListOnboardRules(cmd.Context(), onboardRuleSetOrg, k8shell.OnboardRuleFilter{})
		if err != nil {
			return err
		}
		var current *models.OnboardRule
		for i := range rules {
			if int64(rules[i].ID) == id {
				current = &rules[i]
				break
			}
		}
		if current == nil {
			return fmt.Errorf("rule %d not found in org %q", id, onboardRuleSetOrg)
		}

		req := models.OnboardRuleUpdateRequest{
			Action:   string(current.Action),
			Priority: current.Priority,
			Roles:    current.Roles,
			Sudo:     current.Sudo,
			Note:     current.Note,
		}
		var updated []string
		if cmd.Flags().Changed("action") {
			req.Action = onboardRuleSetAction
			updated = append(updated, "action")
		}
		if cmd.Flags().Changed("priority") {
			req.Priority = onboardRuleSetPriority
			updated = append(updated, "priority")
		}
		if cmd.Flags().Changed("roles") {
			req.Roles = onboardRuleSetRoles
			updated = append(updated, "roles")
		}
		if cmd.Flags().Changed("sudo") {
			req.Sudo = onboardRuleSetSudo
			updated = append(updated, "sudo")
		}
		if cmd.Flags().Changed("note") {
			req.Note = onboardRuleSetNote
			updated = append(updated, "note")
		}

		if len(updated) == 0 {
			return fmt.Errorf("specify at least one field to update (--action, --priority, --roles, --sudo, --note)")
		}

		rule, err := c.UpdateOnboardRule(cmd.Context(), onboardRuleSetOrg, int32(id), req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(rule)
		}

		printer.Println(fmt.Sprintf("%d: updated %s", id, strings.Join(updated, ", ")))
		return nil
	},
}

func init() {
	orgOnboardRuleSetCmd.Flags().StringVar(&onboardRuleSetOrg, "org", "", "organization the rule is scoped to (required)")
	orgOnboardRuleSetCmd.Flags().StringVar(&onboardRuleSetAction, "action", "", "allow, reject, or waitlist")
	orgOnboardRuleSetCmd.Flags().Int32Var(&onboardRuleSetPriority, "priority", 0, "lower wins among matching rules of the same specificity")
	orgOnboardRuleSetCmd.Flags().StringSliceVar(&onboardRuleSetRoles, "roles", nil, "replace the roles granted when this rule resolves to allow, comma-separated")
	orgOnboardRuleSetCmd.Flags().BoolVar(&onboardRuleSetSudo, "sudo", false, "grant sudo access when this rule resolves to allow")
	orgOnboardRuleSetCmd.Flags().StringVar(&onboardRuleSetNote, "note", "", "admin comment")
	_ = orgOnboardRuleSetCmd.MarkFlagRequired("org")
}
