// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var userCredentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Manage a user's external service credentials",
}

var userCredentialColumns = []table.Col[models.UserCredential]{
	{Header: "SERVICE", MaxWidth: 20, Help: "external service name", Field: "serviceName"},
	{Header: "SCOPE", MaxWidth: 20, Help: "granted OAuth scope", Field: "serviceScope"},
	{Header: "SOURCE", MaxWidth: 15, Help: "how the credential was obtained", Field: "credentialSource"},
	{Header: "SUBJECT", MaxWidth: 25, Help: "identity subject on the external service", Field: "subject"},
	{Header: "ACTIVE", MaxWidth: 6, Help: "whether the credential is active", Field: "isActive", Fmt: table.FmtBool},
	{Header: "CREATED", MaxWidth: 16, Help: "creation timestamp (local time)", Fn: func(c models.UserCredential) string {
		return c.CreatedAt.Local().Format("2006-01-02 15:04")
	}},
	{Header: "EXPIRES", MaxWidth: 16, Help: "expiry timestamp, or - if it does not expire", Fn: func(c models.UserCredential) string {
		if c.ExpiresAt == nil {
			return "-"
		}
		return c.ExpiresAt.Local().Format("2006-01-02 15:04")
	}},
}

var userCredentialSortFlag string

var userCredentialListCmd = &cobra.Command{
	Use:               "list <username>",
	Aliases:           []string{"ls"},
	Short:             "List a user's external service credentials",
	Long:              "List external service credentials stored for a user.\n\n" + table.ColumnHelp(userCredentialColumns),
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeUsernames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		creds, err := newClient(ctx).ListUserCredentials(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(creds)
		}

		return table.Table(printer, userCredentialColumns, creds, userCredentialSortFlag)
	},
}

var userCredentialGetCmd = &cobra.Command{
	Use:   "get <username> <service-name>",
	Short: "Show a user's credential for a single external service",
	Long:  "Show a user's credential for a single external service.\n\n" + table.ColumnHelp(userCredentialColumns),
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeUsernames(cmd, args, toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		cred, err := newClient(ctx).GetUserCredential(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(cred)
		}

		return table.Table(printer, userCredentialColumns, []models.UserCredential{*cred}, "")
	},
}

func init() {
	userCredentialListCmd.Flags().StringVar(&userCredentialSortFlag, "sort", "", "sort by fields, e.g. serviceName,-createdAt (prefix - for descending)")
	userCredentialCmd.AddCommand(userCredentialListCmd)
	userCredentialCmd.AddCommand(userCredentialGetCmd)
}
