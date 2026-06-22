// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/spf13/cobra"
)

// completeWorkspaceNames returns workspace names for shell completion.
// The status is appended as a tab-separated description shown by the shell.
func completeWorkspaceNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if cfg == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ctx, err := cfg.ActiveContext()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	workspaces, err := client.New(ctx, false, ctx.Insecure).ListWorkspaces("", false)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := make([]string, len(workspaces))
	for i, ws := range workspaces {
		origin := ws.Blueprint
		if ws.RepoOwner != "" || ws.RepoName != "" {
			origin = ws.RepoOwner + "/" + ws.RepoName
		}
		desc := string(ws.Status)
		if origin != "" {
			desc += " [" + origin + "]"
		}
		names[i] = ws.Name + "\t" + desc
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// completeUsernames returns usernames for shell flag completion.
func completeUsernames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if cfg == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ctx, err := cfg.ActiveContext()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	users, err := client.New(ctx, false, ctx.Insecure).ListUsers()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := make([]string, len(users))
	for i, u := range users {
		names[i] = u.Username
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
