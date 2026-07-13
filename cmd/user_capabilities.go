// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var userCapabilitiesCmd = &cobra.Command{
	Use:     "capabilities",
	Aliases: []string{"caps"},
	Short:   "Show a user's policy capabilities",
}

var userCapabilitiesColumns = []table.Col[models.Capability]{
	{Header: "ACTION", MaxWidth: 30, Help: "policy action identifier", Field: "action"},
	{Header: "ALLOWED", MaxWidth: 8, Help: "whether the action is permitted (true/false)", Field: "allowed", Fmt: table.FmtAllowed},
	{Header: "OBLIGATIONS", MaxWidth: 60, Help: "policy obligations attached to the action, if allowed", Fn: formatCapabilityObligations},
}

// formatCapabilityObligations renders a capability's obligations, but only
// when the action is allowed — a denied action's obligations don't apply.
// Renders "-" when there's nothing to show.
func formatCapabilityObligations(c models.Capability) string {
	if !c.Allowed {
		return "-"
	}
	if s := table.FmtObligations(c.Obligations); s != "" {
		return s
	}
	return "-"
}

var userCapabilitiesSortFlag string

// capabilitiesUsername holds the --user override. Empty means the token's own user.
var capabilitiesUsername string

// capabilitiesResourceOwner holds the --on override, checking capabilities as
// they'd apply to resources owned by that user instead of the caller's own.
var capabilitiesResourceOwner string

var userCapabilitiesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List a user's policy capabilities",
	Long: "List which actions are allowed or denied for the authenticated token, or for another user's with --user.\n" +
		"Use --on to check capabilities against resources owned by a specific user.\n\n" +
		table.ColumnHelp(userCapabilitiesColumns),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		caps, err := newClient(ctx).GetCapabilities(cmd.Context(), capabilitiesUsername, capabilitiesResourceOwner)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(caps)
		}

		return table.Table(printer, userCapabilitiesColumns, caps, userCapabilitiesSortFlag)
	},
}

func init() {
	userCapabilitiesListCmd.Flags().StringVar(&userCapabilitiesSortFlag, "sort", "", "sort by fields, e.g. action,-allowed (prefix - for descending)")
	userCapabilitiesListCmd.Flags().StringVarP(&capabilitiesUsername, "user", "u", "", "act on this user's capabilities instead of your own (admin only)")
	userCapabilitiesListCmd.Flags().StringVar(&capabilitiesResourceOwner, "on", "", "check capabilities against resources owned by this username")
	_ = userCapabilitiesListCmd.RegisterFlagCompletionFunc("user", completeUsernames)
	_ = userCapabilitiesListCmd.RegisterFlagCompletionFunc("on", completeUsernames)
	userCapabilitiesCmd.AddCommand(userCapabilitiesListCmd)
}
