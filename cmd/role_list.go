// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var roleColumns = []table.Col[models.RoleInfo]{
	{Header: "NAME", MaxWidth: 20, Help: "role name", Field: "name"},
	{Header: "ORG", MaxWidth: 15, Help: "organization the role is scoped to (empty for global roles)", Field: "org"},
	{Header: "DESCRIPTION", MaxWidth: 40, Help: "role description", Field: "description"},
	{Header: "BLUEPRINTS", MaxWidth: 30, Help: "blueprints granted by this role (comma-separated)", Field: "blueprints", Fmt: table.FmtJoin},
	{Header: "USERS", MaxWidth: 6, Help: "number of users holding this role", Field: "userCount"},
	{Header: "CREATED", MaxWidth: 16, Help: "creation timestamp (local time)", Fn: func(r models.RoleInfo) string {
		return r.CreatedAt.Local().Format("2006-01-02 15:04")
	}},
}

var (
	roleListOrg  string
	roleSortFlag string
)

var roleListCmd = &cobra.Command{
	Use:     "list --org <org>",
	Aliases: []string{"ls"},
	Short:   "List roles",
	Long: "List roles scoped to --org, plus all global roles (there is no endpoint to list global roles on their own).\n\n" +
		table.ColumnHelp(roleColumns),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		roles, err := newClient(ctx).ListRoles(cmd.Context(), roleListOrg)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(roles)
		}

		return table.Table(printer, roleColumns, roles, roleSortFlag)
	},
}

func init() {
	roleListCmd.Flags().StringVar(&roleListOrg, "org", "", "organization to list roles for, plus all global roles (required)")
	roleListCmd.Flags().StringVar(&roleSortFlag, "sort", "", "sort by fields, e.g. name,-userCount (prefix - for descending)")
	_ = roleListCmd.MarkFlagRequired("org")
}
