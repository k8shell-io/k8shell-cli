// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	roleSetOrg              string
	roleSetDescription      string
	roleSetAddBlueprints    []string
	roleSetRemoveBlueprints []string
)

var roleSetCmd = &cobra.Command{
	Use:   "set <name> --org <org> [flags]",
	Short: "Update fields on a role",
	Long:  "Update a role's description and/or the blueprints it grants.\nEvery user holding the role gains or loses access to a blueprint as it's added or removed.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		name := args[0]
		c := newClient(ctx)
		var updated []string

		if cmd.Flags().Changed("description") {
			req := models.RoleUpdateRequest{Description: &roleSetDescription}
			if _, err := c.UpdateRole(cmd.Context(), roleSetOrg, name, req); err != nil {
				return err
			}
			updated = append(updated, "description")
		}
		if len(roleSetRemoveBlueprints) > 0 {
			if err := c.RemoveRoleBlueprints(cmd.Context(), roleSetOrg, name, roleSetRemoveBlueprints); err != nil {
				return err
			}
			updated = append(updated, "remove-blueprint")
		}
		if len(roleSetAddBlueprints) > 0 {
			if err := c.AddRoleBlueprints(cmd.Context(), roleSetOrg, name, roleSetAddBlueprints); err != nil {
				return err
			}
			updated = append(updated, "add-blueprint")
		}

		if len(updated) == 0 {
			return fmt.Errorf("specify at least one field to update (--description, --add-blueprint, --remove-blueprint)")
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"name": name, "updated": updated})
		}

		printer.Println(fmt.Sprintf("%s: updated %s", name, strings.Join(updated, ", ")))
		return nil
	},
}

func init() {
	roleSetCmd.Flags().StringVar(&roleSetOrg, "org", "", "organization the role is scoped to (required)")
	roleSetCmd.Flags().StringVar(&roleSetDescription, "description", "", "role description")
	repeatableStringVar(roleSetCmd.Flags(), &roleSetAddBlueprints, "add-blueprint", "grant a blueprint, in addition to existing ones (repeatable)")
	repeatableStringVar(roleSetCmd.Flags(), &roleSetRemoveBlueprints, "remove-blueprint", "revoke a blueprint, leaving others untouched (repeatable)")
	_ = roleSetCmd.MarkFlagRequired("org")
}
