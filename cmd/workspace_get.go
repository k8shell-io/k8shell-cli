// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

// workspaceDetailColumns lists every field of models.WorkspaceDetails, used to
// render `workspace get` as a two-column field/value listing.
var workspaceDetailColumns = []table.Col[models.WorkspaceDetails]{
	{Header: "NAME", Field: "name"},
	{Header: "USERNAME", Field: "username"},
	{Header: "STATUS", Fn: func(w models.WorkspaceDetails) string { return string(w.Status) }},
	{Header: "MESSAGE", Fn: func(w models.WorkspaceDetails) string { return w.Message }},
	{Header: "RESTARTS", Fn: func(w models.WorkspaceDetails) string { return fmt.Sprint(w.Restarts) }},
	{Header: "CREATED", Fn: func(w models.WorkspaceDetails) string {
		if w.Created.IsZero() {
			return "-"
		}
		return w.Created.Local().Format("2006-01-02 15:04")
	}},
	{Header: "APPVERSION", Field: "appVersion"},
	{Header: "REPO", Fn: func(w models.WorkspaceDetails) string {
		if w.RepoOwner == "" && w.RepoName == "" {
			return ""
		}
		repo := w.RepoOwner + "/" + w.RepoName
		if w.RepoRef != "" {
			repo += "@" + w.RepoRef
		}
		return repo
	}},
	{Header: "BLUEPRINT", Field: "blueprint"},
	{Header: "ORGANIZATION", Field: "organization"},
	{Header: "CPU", Field: "cpu"},
	{Header: "MEMORY", Field: "memory"},
	{Header: "IP", Field: "podIP"},
	{Header: "NAMESPACE", Field: "namespace"},
}

var workspaceGetCmd = &cobra.Command{
	Use:               "get <workspace-name>",
	Short:             "Show details for a single workspace",
	Long:              "Show details for a single workspace.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeWorkspaceNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		ws, err := newClient(ctx).GetWorkspace(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(ws)
		}

		return table.Detail(printer, workspaceDetailColumns, *ws)
	},
}
