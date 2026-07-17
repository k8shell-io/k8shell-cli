// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	setTokenScopes     []string
	setTokenActivate   bool
	setTokenDeactivate bool
)

var userTokenSetCmd = &cobra.Command{
	Use:   "set <id> [flags]",
	Short: "Update a personal access token's active state and/or scopes",
	Long: "Update a personal access token's active state and/or scopes, identified by ID as found in `token list`.\n" +
		"Name and expiry are immutable after creation.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid token ID %q: %w", args[0], err)
		}

		if setTokenActivate && setTokenDeactivate {
			return fmt.Errorf("--activate cannot be combined with --deactivate")
		}

		var req models.AccessTokenUpdateRequest
		var updated []string

		if cmd.Flags().Changed("scopes") {
			req.Scopes = &setTokenScopes
			updated = append(updated, "scopes")
		}
		if setTokenActivate {
			active := true
			req.Active = &active
			updated = append(updated, "activated")
		}
		if setTokenDeactivate {
			active := false
			req.Active = &active
			updated = append(updated, "deactivated")
		}

		if len(updated) == 0 {
			return fmt.Errorf("specify at least one field to update (--scopes, --activate, --deactivate)")
		}

		token, err := newClient(ctx).UpdateUserToken(cmd.Context(), tokenUsername, id, req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(token)
		}

		printer.Println(fmt.Sprintf("%d: updated %s", id, strings.Join(updated, ", ")))
		return nil
	},
}

func init() {
	addTokenUsernameFlag(userTokenSetCmd)
	userTokenSetCmd.Flags().StringSliceVar(&setTokenScopes, "scopes", nil, "replace granted scopes, comma-separated")
	userTokenSetCmd.Flags().BoolVar(&setTokenActivate, "activate", false, "mark the token active")
	userTokenSetCmd.Flags().BoolVar(&setTokenDeactivate, "deactivate", false, "mark the token inactive")
	userTokenCmd.AddCommand(userTokenSetCmd)
}
