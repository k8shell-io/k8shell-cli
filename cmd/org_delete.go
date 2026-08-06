// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	orgDeleteMoveUsersTo        string
	orgDeletePreserveWorkspaces bool
)

var orgDeleteCmd = &cobra.Command{
	Use:     "delete <name> [flags]",
	Aliases: []string{"del"},
	Short:   "Delete an organization",
	Long: "Permanently delete an organization. By default this fails if any user still belongs to it;\n" +
		"pass --move-users-to to re-home them into another organization instead of deleting them.\n" +
		"Deleting (or moving) a user deletes their workspaces too, unless --preserve-workspaces is set.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		name := args[0]
		req := models.OrganizationDeleteRequest{
			MoveUsersToOrg:     orgDeleteMoveUsersTo,
			PreserveWorkspaces: orgDeletePreserveWorkspaces,
		}
		if err := newClient(ctx).DeleteOrganization(cmd.Context(), name, req); err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"name": name, "deleted": true})
		}

		printer.Println(fmt.Sprintf("%s: deleted", name))
		return nil
	},
}

func init() {
	orgDeleteCmd.Flags().StringVar(&orgDeleteMoveUsersTo, "move-users-to", "", "move the organization's users into this organization instead of deleting them")
	orgDeleteCmd.Flags().BoolVar(&orgDeletePreserveWorkspaces, "preserve-workspaces", false, "keep workspaces belonging to deleted or moved users instead of deleting them")
}
