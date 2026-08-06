// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var orgSetDescription string

var orgSetCmd = &cobra.Command{
	Use:   "set <name> [flags]",
	Short: "Update fields on an organization",
	Long:  "Update an organization's description. The name is immutable.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		if !cmd.Flags().Changed("description") {
			return fmt.Errorf("specify at least one field to update (--description)")
		}

		name := args[0]
		req := models.OrganizationUpdateRequest{Description: &orgSetDescription}
		if _, err := newClient(ctx).UpdateOrganization(cmd.Context(), name, req); err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"name": name, "updated": []string{"description"}})
		}

		printer.Println(fmt.Sprintf("%s: updated description", name))
		return nil
	},
}

func init() {
	orgSetCmd.Flags().StringVar(&orgSetDescription, "description", "", "organization description")
}
