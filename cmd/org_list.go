// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var orgColumns = []table.Col[models.Organization]{
	{Header: "NAME", MaxWidth: 20, Help: "organization name", Field: "name"},
	{Header: "DESCRIPTION", MaxWidth: 40, Help: "organization description", Field: "description"},
	{Header: "USERS", MaxWidth: 6, Help: "number of users in this organization", Field: "userCount"},
	{Header: "ADMINS", MaxWidth: 30, Help: "usernames holding the org-admin role (comma-separated)", Field: "adminUsernames", Fmt: table.FmtJoin},
	{Header: "READONLY", MaxWidth: 8, Help: "built-in organization that cannot be updated or deleted", Field: "readOnly", Fmt: table.FmtBool},
	{Header: "CREATED", MaxWidth: 16, Help: "creation timestamp (local time)", Fn: func(o models.Organization) string {
		return o.CreatedAt.Local().Format("2006-01-02 15:04")
	}},
}

var orgSortFlag string

var orgListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List organizations",
	Long:    "List the registered organizations.\n\n" + table.ColumnHelp(orgColumns),
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		orgs, err := newClient(ctx).ListOrganizations(cmd.Context())
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(orgs)
		}

		return table.Table(printer, orgColumns, orgs, orgSortFlag)
	},
}

func init() {
	orgListCmd.Flags().StringVar(&orgSortFlag, "sort", "", "sort by fields, e.g. name,-userCount (prefix - for descending)")
}
