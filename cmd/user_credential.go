// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var userCredentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Manage external service credentials",
}

// credentialUsername holds the --username override. Empty means the token's own user.
var credentialUsername string

func addCredentialUsernameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&credentialUsername, "username", "u", "", "act on this user's credentials instead of your own (admin only)")
	_ = cmd.RegisterFlagCompletionFunc("username", completeUsernames)
}

var userCredentialColumns = []table.Col[models.UserCredential]{
	{Header: "ID", MaxWidth: 8, Help: "credential ID, used with `get`", Field: "id"},
	{Header: "USERNAME", MaxWidth: 20, Help: "owning user", Field: "username"},
	{Header: "SERVICE", MaxWidth: 20, Help: "external service name", Field: "serviceName"},
	{Header: "SCOPE", MaxWidth: 20, Help: "granted OAuth scope; for kubernetes, the namespace", Field: "serviceScope"},
	{Header: "SOURCE", MaxWidth: 15, Help: "how the credential was obtained", Field: "credentialSource"},
	{Header: "SUBJECT", MaxWidth: 25, Help: "identity subject on the external service; for kubernetes, the service account", Field: "subject"},
	{Header: "ACTIVE", MaxWidth: 6, Help: "whether the credential is active", Field: "isActive", Fmt: table.FmtBool},
	{Header: "CREATED", MaxWidth: 16, Help: "creation timestamp (local time)", Fn: func(c models.UserCredential) string {
		return c.CreatedAt.Local().Format("2006-01-02 15:04")
	}},
	{Header: "UPDATED", MaxWidth: 16, Help: "last update timestamp (local time)", Fn: func(c models.UserCredential) string {
		return c.UpdatedAt.Local().Format("2006-01-02 15:04")
	}},
}

var userCredentialSortFlag string

var userCredentialListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List external service credentials for your user",
	Long: "List external service credentials stored for your user, or another user's with --username.\n\n" +
		table.ColumnHelp(userCredentialColumns),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		creds, err := newClient(ctx).ListUserCredentials(cmd.Context(), credentialUsername)
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
	Use:   "get <id>",
	Short: "Show a single external service credential by ID",
	Long: "Show a single external service credential by ID, as found in `credential list`.\n\n" +
		table.ColumnHelp(userCredentialColumns),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid credential ID %q: %w", args[0], err)
		}

		cred, err := newClient(ctx).GetUserCredential(cmd.Context(), credentialUsername, uint32(id))
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(cred)
		}

		return table.Detail(printer, userCredentialColumns, *cred)
	},
}

var userCredentialDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"del"},
	Short:   "Delete a single external service credential by ID",
	Long:    "Permanently delete a single external service credential by ID, as found in `credential list`.",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid credential ID %q: %w", args[0], err)
		}

		if err := newClient(ctx).DeleteUserCredential(cmd.Context(), credentialUsername, uint32(id)); err != nil {
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
	userCredentialListCmd.Flags().StringVar(&userCredentialSortFlag, "sort", "", "sort by fields, e.g. serviceName,-createdAt (prefix - for descending)")
	addCredentialUsernameFlag(userCredentialListCmd)
	addCredentialUsernameFlag(userCredentialGetCmd)
	addCredentialUsernameFlag(userCredentialDeleteCmd)
	userCredentialCmd.AddCommand(userCredentialListCmd)
	userCredentialCmd.AddCommand(userCredentialGetCmd)
	userCredentialCmd.AddCommand(userCredentialDeleteCmd)
}
