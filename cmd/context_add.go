// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/k8shell-io/k8shell/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	addServer   string
	addToken    string
	addInsecure bool
)

var contextAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := addToken
		if token == "" {
			fmt.Fprint(os.Stderr, "Token: ")
			raw, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("reading token: %w", err)
			}
			fmt.Fprintln(os.Stderr)
			token = strings.TrimSpace(string(raw))
		}
		if token == "" {
			return fmt.Errorf("token must not be empty")
		}

		tmpCtx := &config.Context{Server: addServer, Token: token, Insecure: addInsecure}
		profile, err := newClient(tmpCtx).GetProfile(cmd.Context())
		if err != nil {
			return fmt.Errorf("verifying token: %w", err)
		}

		ctx := config.Context{
			Name:     args[0],
			Server:   addServer,
			Token:    token,
			Username: profile.Username,
			Insecure: addInsecure,
		}
		ctx.SetIntegrity()
		if err := cfg.AddContext(ctx); err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Context %q added (username: %s).\n", args[0], profile.Username)
		return nil
	},
}

func init() {
	contextAddCmd.Flags().StringVar(&addServer, "server", "", "API server URL (required)")
	contextAddCmd.Flags().StringVar(&addToken, "token", "", "PAT token (prompted securely if omitted)")
	contextAddCmd.Flags().BoolVar(&addInsecure, "insecure", false, "skip TLS certificate verification for this context")
	_ = contextAddCmd.MarkFlagRequired("server")
}
