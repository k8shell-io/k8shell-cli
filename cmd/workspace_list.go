// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var workspaceColumns = []table.Col[models.WorkspaceDetails]{
	{Header: "NAME", MaxWidth: 25, Help: "workspace name", Field: "name"},
	{Header: "USERNAME", MaxWidth: 20, Help: "owner username", Field: "username"},
	{Header: "STATUS", MaxWidth: 12, Help: "current pod status (Starting, Running, Stopped, ...)", Fn: func(w models.WorkspaceDetails) string { return string(w.Status) }},
	{Header: "ORIGIN", MaxWidth: 30, Help: "repo as owner/name, or blueprint when no repo is set", Fn: func(w models.WorkspaceDetails) string {
		if w.RepoOwner != "" || w.RepoName != "" {
			return w.RepoOwner + "/" + w.RepoName
		}
		return w.Blueprint
	}},
	{Header: "CPU", MaxWidth: 8, Help: "CPU resource limit (e.g. 500m)", Field: "cpu"},
	{Header: "MEMORY", MaxWidth: 10, Help: "memory resource limit (e.g. 512Mi)", Field: "memory"},
	{Header: "VERSION", MaxWidth: 12, Help: "app version running in the workspace", Field: "appVersion"},
	{Header: "IP", MaxWidth: 15, Help: "pod IP address", Field: "podIP"},
	{Header: "CREATED", MaxWidth: 16, Help: "creation timestamp (local time)", Fn: func(w models.WorkspaceDetails) string {
		if w.Created.IsZero() {
			return "-"
		}
		return w.Created.Local().Format("2006-01-02 15:04")
	}},
}

var (
	workspaceSortFlag     string
	workspaceUsernameFlag string
	workspaceAllFlag      bool
)

var workspaceListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List workspaces",
	Long:    "List workspaces visible to the authenticated token.\n\n" + table.ColumnHelp(workspaceColumns),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		workspaces, err := client.New(ctx, debug, insecure || ctx.Insecure).ListWorkspaces(workspaceUsernameFlag, workspaceAllFlag)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(workspaces)
		}

		return table.Table(printer, workspaceColumns, workspaces, workspaceSortFlag)
	},
}

func init() {
	workspaceListCmd.Flags().StringVarP(&workspaceUsernameFlag, "username", "u", "", "filter by owner username")
	workspaceListCmd.Flags().StringVar(&workspaceSortFlag, "sort", "", "sort by fields, e.g. name,-username (prefix - for descending); origin is repo (owner/name) or blueprint")
	workspaceListCmd.Flags().BoolVar(&workspaceAllFlag, "all", false, "include all workspaces")
	_ = workspaceListCmd.RegisterFlagCompletionFunc("username", completeUsernames)
}
