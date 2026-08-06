// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	roleCreateOrg         string
	roleCreateDescription string
)

var roleCreateCmd = &cobra.Command{
	Use:   "create <name> --org <org> [flags]",
	Short: "Create a new role",
	Long:  "Create a new assignable role scoped to --org. Global roles cannot be created.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		req := models.RoleCreateRequest{
			Name:        args[0],
			Description: roleCreateDescription,
		}

		role, err := newClient(ctx).CreateRole(cmd.Context(), roleCreateOrg, req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(role)
		}

		printer.Println(fmt.Sprintf("%s: created", role.Name))
		return nil
	},
}

func init() {
	roleCreateCmd.Flags().StringVar(&roleCreateOrg, "org", "", "organization to scope the role to (required)")
	roleCreateCmd.Flags().StringVar(&roleCreateDescription, "description", "", "role description")
	_ = roleCreateCmd.MarkFlagRequired("org")
}
