// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var userTokenDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"del"},
	Short:   "Revoke a personal access token",
	Long: "Permanently revoke a personal access token by ID, as found in `token list`, " +
		"for your user or another user's with --username.",
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

		if err := newClient(ctx).DeleteUserToken(cmd.Context(), tokenUsername, id); err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"id": id, "deleted": true})
		}

		printer.Println(fmt.Sprintf("%d: deleted", id))
		return nil
	},
}

func init() {
	addTokenUsernameFlag(userTokenDeleteCmd)
	userTokenCmd.AddCommand(userTokenDeleteCmd)
}
