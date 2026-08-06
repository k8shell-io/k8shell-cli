// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var orgCreateDescription string

var orgCreateCmd = &cobra.Command{
	Use:   "create <name> [flags]",
	Short: "Create a new organization",
	Long:  "Register a new organization.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		req := models.OrganizationCreateRequest{
			Name:        args[0],
			Description: orgCreateDescription,
		}

		org, err := newClient(ctx).CreateOrganization(cmd.Context(), req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(org)
		}

		printer.Println(fmt.Sprintf("%s: created", org.Name))
		return nil
	},
}

func init() {
	orgCreateCmd.Flags().StringVar(&orgCreateDescription, "description", "", "organization description")
}
