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
	createServer   string
	createToken    string
	createInsecure bool
)

var contextCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := createToken
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

		tmpCtx := &config.Context{Server: createServer, Token: token, Insecure: createInsecure}
		profile, err := newClient(tmpCtx).GetProfile(cmd.Context())
		if err != nil {
			return fmt.Errorf("verifying token: %w", err)
		}

		ctx := config.Context{
			Name:     args[0],
			Server:   createServer,
			Token:    token,
			Username: profile.Username,
			Insecure: createInsecure,
		}
		ctx.SetIntegrity()
		if err := cfg.AddContext(ctx); err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Context %q created (username: %s).\n", args[0], profile.Username)
		return nil
	},
}

func init() {
	contextCreateCmd.Flags().StringVar(&createServer, "server", "", "API server URL (required)")
	contextCreateCmd.Flags().StringVar(&createToken, "token", "", "PAT token (prompted securely if omitted)")
	contextCreateCmd.Flags().BoolVar(&createInsecure, "insecure", false, "skip TLS certificate verification for this context")
	_ = contextCreateCmd.MarkFlagRequired("server")
}
