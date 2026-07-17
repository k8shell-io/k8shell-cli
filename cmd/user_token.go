// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var userTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage personal access tokens",
}

// tokenUsername holds the --username override. Empty means the token's own user.
var tokenUsername string

func addTokenUsernameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&tokenUsername, "username", "u", "", "act on this user's tokens instead of your own (admin only)")
	_ = cmd.RegisterFlagCompletionFunc("username", completeUsernames)
}

var userTokenColumns = []table.Col[models.AccessToken]{
	{Header: "ID", MaxWidth: 8, Help: "token ID", Fn: func(t models.AccessToken) string { return fmt.Sprint(t.ID) }},
	{Header: "USERNAME", MaxWidth: 20, Help: "owning user", Fn: func(t models.AccessToken) string { return t.Username }},
	{Header: "NAME", MaxWidth: 20, Help: "token name", Fn: func(t models.AccessToken) string { return t.Name }},
	{Header: "SCOPES", MaxWidth: 30, Help: "granted scopes (comma-separated)", Wrap: true, Fn: func(t models.AccessToken) string {
		return strings.Join(t.Scopes, ", ")
	}},
	{Header: "PREVIEW", MaxWidth: 12, Help: "truncated token preview", Fn: func(t models.AccessToken) string {
		if t.TokenPreview == nil {
			return "-"
		}
		return *t.TokenPreview
	}},
	{Header: "ACTIVE", MaxWidth: 6, Help: "whether the token is active", Fn: func(t models.AccessToken) string { return table.FmtBool(t.IsActive) }},
	{Header: "EXPIRES", MaxWidth: 16, Help: "expiry timestamp (local time), or - if it never expires", Fn: func(t models.AccessToken) string { return fmtTime(t.ExpiresAt) }},
	{Header: "LAST_USED", MaxWidth: 16, Help: "last used timestamp (local time), or - if never used", Fn: func(t models.AccessToken) string { return fmtTime(t.LastUsedAt) }},
	{Header: "CREATED", MaxWidth: 16, Help: "creation timestamp (local time)", Fn: func(t models.AccessToken) string {
		return t.CreatedAt.Local().Format("2006-01-02 15:04")
	}},
}

var userTokenSortFlag string

var userTokenListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List personal access tokens for your user",
	Long: "List personal access tokens issued for your user, or another user's with --username.\n\n" +
		table.ColumnHelp(userTokenColumns),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		tokens, err := newClient(ctx).ListUserTokens(cmd.Context(), tokenUsername)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(tokens)
		}

		return table.Table(printer, userTokenColumns, tokens, userTokenSortFlag)
	},
}

var userTokenGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Show details for a single personal access token",
	Long: "Show details for a single personal access token, as found in `token list`.\n\n" +
		table.ColumnHelp(userTokenColumns),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid token ID %q: %w", args[0], err)
		}

		token, err := newClient(ctx).GetUserToken(cmd.Context(), tokenUsername, id)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(token)
		}

		return table.Detail(printer, userTokenColumns, *token)
	},
}

func init() {
	userTokenListCmd.Flags().StringVar(&userTokenSortFlag, "sort", "", "sort by fields, e.g. name,-createdAt (prefix - for descending)")
	addTokenUsernameFlag(userTokenListCmd)
	addTokenUsernameFlag(userTokenGetCmd)
	userTokenCmd.AddCommand(userTokenListCmd)
	userTokenCmd.AddCommand(userTokenGetCmd)
}
