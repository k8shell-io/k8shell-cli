// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var userColumns = []table.Col[models.User]{
	{Header: "USERNAME",   MaxWidth: 20,  Help: "login username",                                  Field: "username"},
	{Header: "FULLNAME",   MaxWidth: 20,  Help: "display name",                                    Field: "fullname"},
	{Header: "EMAIL",      MaxWidth: 30,  Help: "email address",                                   Field: "email"},
	{Header: "ORG",        MaxWidth: 15,  Help: "organization",                                    Field: "organization"},
	{Header: "ROLES",      MaxWidth: 20,  Help: "assigned roles (comma-separated)",                Field: "roles",      Fmt: table.FmtRoles},
	{Header: "BLUEPRINTS", MaxWidth: 30,  Help: "allowed blueprints (comma-separated)",            Field: "blueprints", Fmt: table.FmtJoin},
	{Header: "SUDO",       MaxWidth: 5,   Help: "sudo access (true/false)",                        Field: "sudo",       Fmt: table.FmtBool},
	{Header: "SOURCE",     MaxWidth: 122, Help: "identity source (e.g. github, google)",           Field: "source"},
	{Header: "STATUS",     MaxWidth: 8,   Help: "active, locked, or invalid",                      Fn: userStatus},
}

var userSortFlag string

var userListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short: "List users",
	Long:  "List users visible to the authenticated token.\n\n" + table.ColumnHelp(userColumns),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		users, err := newClient(ctx).ListUsers(cmd.Context())
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(users)
		}

		return table.Table(printer, userColumns, users, userSortFlag)
	},
}

func init() {
	userListCmd.Flags().StringVar(&userSortFlag, "sort", "", "sort by fields, e.g. username,-email (prefix - for descending)")
}

// userStatus derives a display status string from the user's locked and validity fields.
func userStatus(u models.User) string {
	if u.Locked {
		return "locked"
	}
	if !u.IsValid {
		return "invalid"
	}
	return "active"
}
