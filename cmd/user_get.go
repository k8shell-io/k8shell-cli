// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

// userDetailColumns lists every field of models.User (except the write-only password),
// used to render `user get` as a two-column field/value listing. The SSH key digests
// are fetched separately (they are no longer part of the profile) and appended as an
// extra KEYS row when rendering.
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
	{Header: "EXPIRES", Fn: func(u models.User) string { return u.ExpiresAt.Local().Format("2006-01-02 15:04") }},
}

// formatAuthKeys renders SSH key digests as "index:digest (source)" entries, one
// per line, so the index can be passed to `user set --remove-key`.
func formatAuthKeys(keys []models.UserAuthKey) string {
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%d:%s (%s)", k.Index, k.Key, k.Source)
	}
	return strings.Join(parts, "\n")
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

		c := newClient(ctx)

		user, err := c.GetUserProfile(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		keys, err := c.ListUserAuthKeys(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			out := struct {
				models.User
				Keys []models.UserAuthKey `json:"keys"`
			}{User: *user, Keys: keys}
			return printer.JSON(out)
		}

		cols := append(append([]table.Col[models.User]{}, userDetailColumns...), table.Col[models.User]{
			Header: "AUTHKEYS",
			Fn:     func(models.User) string { return formatAuthKeys(keys) },
		})

		return table.Detail(printer, cols, *user)
	},
}
