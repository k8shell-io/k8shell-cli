// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

// userDetailColumns lists every field of models.User (except the write-only password),
// used to render `user get` as a two-column field/value listing.
var userDetailColumns = []table.Col[models.User]{
	{Header: "USERNAME", Field: "username"},
	{Header: "FULLNAME", Field: "fullname"},
	{Header: "EMAIL", Field: "email"},
	{Header: "ORGANIZATION", Field: "organization"},
	{Header: "UID", Field: "uid"},
	{Header: "GID", Field: "gid"},
	{Header: "ROLES", Field: "roles", Fmt: table.FmtRoles},
	{Header: "BLUEPRINTS", Field: "blueprints", Fmt: table.FmtJoin},
	{Header: "SUDO", Field: "sudo", Fmt: table.FmtBool},
	{Header: "LOCKED", Field: "locked", Fmt: table.FmtBool},
	{Header: "VALID", Field: "isValid", Fmt: table.FmtBool},
	{Header: "SOURCE", Field: "source"},
	{Header: "SHELL", Field: "shell"},
	{Header: "AUTHKEYS", Field: "authKeys", Fmt: table.FmtJoin},
	{Header: "EXPIRES", Fn: func(u models.User) string { return u.ExpiresAt.Local().Format("2006-01-02 15:04") }},
}

var userGetCmd = &cobra.Command{
	Use:               "get <username>",
	Short:             "Show details for a single user",
	Long:              "Show details for a single user.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeUsernames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		user, err := newClient(ctx).GetUserProfile(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(user)
		}

		return table.Detail(printer, userDetailColumns, *user)
	},
}
