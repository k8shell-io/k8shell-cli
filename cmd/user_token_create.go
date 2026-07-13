// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"time"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	createTokenScopes   []string
	createTokenTTL      time.Duration
	createTokenRenew    bool
	createTokenInactive bool
)

var userTokenCreateCmd = &cobra.Command{
	Use:   "create <name> --scopes <scope1,scope2,...> [flags]",
	Short: "Create a new personal access token",
	Long: "Create a new personal access token for your user, or another user's with --username.\n" +
		"The raw token value is returned exactly once and cannot be retrieved again — save it now.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		req := models.AccessTokenCreateRequest{
			Name:   args[0],
			Scopes: createTokenScopes,
			Renew:  createTokenRenew,
		}
		if createTokenTTL > 0 {
			exp := time.Now().Add(createTokenTTL).UTC().Format(time.RFC3339)
			req.ExpiresAt = &exp
		}
		if createTokenInactive {
			active := false
			req.Active = &active
		}

		created, err := newClient(ctx).CreateUserToken(cmd.Context(), tokenUsername, req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(created)
		}

		printer.Println(fmt.Sprintf("%d: created\nToken: %s\n\nSave this token now — it will not be shown again.", created.ID, created.Token))
		return nil
	},
}

func init() {
	addTokenUsernameFlag(userTokenCreateCmd)
	userTokenCreateCmd.Flags().StringSliceVar(&createTokenScopes, "scopes", nil, "granted scopes, comma-separated (required)")
	userTokenCreateCmd.Flags().DurationVar(&createTokenTTL, "ttl", 0, "expire the token after this duration, e.g. 720h (omit for a non-expiring token)")
	userTokenCreateCmd.Flags().BoolVar(&createTokenRenew, "renew", false, "rotate an existing active token with the same name in place, instead of creating a new one")
	userTokenCreateCmd.Flags().BoolVar(&createTokenInactive, "inactive", false, "create the token in a disabled state")
	_ = userTokenCreateCmd.MarkFlagRequired("scopes")
	userTokenCmd.AddCommand(userTokenCreateCmd)
}
