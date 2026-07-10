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
	setCredentialScope       string
	setCredentialSubject     string
	setCredentialSecret      bool
	setCredentialSecretStdin bool
	setCredentialActivate    bool
	setCredentialDeactivate  bool
)

var userCredentialSetCmd = &cobra.Command{
	Use:   "set <id> [flags]",
	Short: "Update fields on an external service credential",
	Long:  "Update one or more fields on an external service credential, identified by ID as found in `credential list`.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid credential ID %q: %w", args[0], err)
		}

		if setCredentialSecret && setCredentialSecretStdin {
			return fmt.Errorf("--secret cannot be combined with --secret-stdin")
		}
		if setCredentialActivate && setCredentialDeactivate {
			return fmt.Errorf("--activate cannot be combined with --deactivate")
		}

		var req models.UserCredentialUpdateRequest
		var updated []string

		if cmd.Flags().Changed("scope") {
			req.Scope = &setCredentialScope
			updated = append(updated, "scope")
		}
		if cmd.Flags().Changed("subject") {
			req.Subject = &setCredentialSubject
			updated = append(updated, "subject")
		}
		if setCredentialSecret || setCredentialSecretStdin {
			secret, err := readPassword("Secret: ", setCredentialSecretStdin)
			if err != nil {
				return err
			}
			req.Secret = &secret
			updated = append(updated, "secret")
		}
		if setCredentialActivate {
			active := true
			req.Active = &active
			updated = append(updated, "activated")
		}
		if setCredentialDeactivate {
			active := false
			req.Active = &active
			updated = append(updated, "deactivated")
		}

		if len(updated) == 0 {
			return fmt.Errorf("specify at least one field to update (--scope, --subject, --secret, --secret-stdin, --activate, --deactivate)")
		}

		cred, err := newClient(ctx).UpdateUserCredential(cmd.Context(), credentialUsername, uint32(id), req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(cred)
		}

		printer.Println(fmt.Sprintf("%d: updated %s", id, strings.Join(updated, ", ")))
		return nil
	},
}

func init() {
	addCredentialUsernameFlag(userCredentialSetCmd)
	userCredentialSetCmd.Flags().StringVar(&setCredentialScope, "scope", "", "new scope (registry host, git host, or kubernetes namespace)")
	userCredentialSetCmd.Flags().StringVar(&setCredentialSubject, "subject", "", "new subject (username or service account)")
	userCredentialSetCmd.Flags().BoolVar(&setCredentialSecret, "secret", false, "set a new secret, prompted interactively")
	userCredentialSetCmd.Flags().BoolVar(&setCredentialSecretStdin, "secret-stdin", false, "set a new secret by reading it from stdin")
	userCredentialSetCmd.Flags().BoolVar(&setCredentialActivate, "activate", false, "mark the credential active")
	userCredentialSetCmd.Flags().BoolVar(&setCredentialDeactivate, "deactivate", false, "mark the credential inactive")

	userCredentialCmd.AddCommand(userCredentialSetCmd)
}
