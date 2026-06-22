// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strings"

	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/spf13/cobra"
)

var (
	createUsername  string
	createBlueprint string
	createRepo      string
	createRef       string
	createEvents    bool
)

var workspaceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintSet := createBlueprint != ""
		repoSet := createRepo != ""

		switch {
		case blueprintSet && repoSet:
			return fmt.Errorf("specify either --blueprint or --repo, not both")
		case !blueprintSet && !repoSet:
			return fmt.Errorf("specify either --blueprint or --repo")
		}

		var repoOwner, repoName string
		if repoSet {
			if i := strings.IndexByte(createRepo, '/'); i >= 0 {
				repoOwner = createRepo[:i]
				repoName = createRepo[i+1:]
			} else {
				repoOwner = createUsername
				repoName = createRepo
			}
		}

		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		username := createUsername
		if username == "" {
			username = ctx.Username
		}
		if username == "" {
			return fmt.Errorf("no username specified and none saved in context; use --username or log in again")
		}

		c := client.New(ctx, debug, insecure || ctx.Insecure)

		resp, err := c.CreateWorkspace(client.WorkspaceCreateRequest{
			Username:  username,
			Blueprint: createBlueprint,
			RepoOwner: repoOwner,
			RepoName:  repoName,
			RepoRef:   createRef,
		})
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(resp)
		}

		if createEvents {
			fmt.Printf("Creating workspace %s (job %s)\n", resp.Workspace, resp.JobID)
		} else {
			fmt.Printf("Creating workspace %s...", resp.Workspace)
		}

		rc, err := c.MonitorWorkspace(resp.MonitorURL)
		if err != nil {
			if !createEvents {
				fmt.Println()
			}
			return fmt.Errorf("monitoring workspace: %w", err)
		}
		defer rc.Close()

		if createEvents {
			return printEventStream(rc)
		}

		// Default mode: update a single progress line in place.
		return printProgressStream(rc, resp.Workspace)
	},
}

func init() {
	workspaceCreateCmd.Flags().StringVarP(&createUsername, "username", "u", "", "owner username (defaults to the logged-in user)")
	workspaceCreateCmd.Flags().StringVar(&createBlueprint, "blueprint", "", "blueprint name")
	workspaceCreateCmd.Flags().StringVar(&createRepo, "repo", "", "repository as owner/name (owner defaults to username when omitted)")
	workspaceCreateCmd.Flags().StringVar(&createRef, "ref", "", "repository ref — branch, tag, or commit (optional, used with --repo)")
	workspaceCreateCmd.Flags().BoolVar(&createEvents, "events", false, "show log events instead of progress percentage")
	_ = workspaceCreateCmd.RegisterFlagCompletionFunc("username", completeUsernames)
}
