// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var userCredentialAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new external service credential for a user",
}

func printAddedCredential(created *models.UserCredential) error {
	if printer.IsJSON() {
		return printer.JSON(created)
	}

	printer.Println(fmt.Sprintf("%s: credential added (id %d)", created.ServiceName, created.ID))
	return nil
}

var (
	addKubernetesScope   string
	addKubernetesSubject string
)

var userCredentialAddKubernetesCmd = &cobra.Command{
	Use:   "kubernetes --scope <namespace> --subject <service-account>",
	Short: "Add a Kubernetes service account credential",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		created, err := newClient(ctx).AddKubernetesUserCredential(cmd.Context(), credentialUsername, models.UserKubernetesCredentialRequest{
			Scope:   addKubernetesScope,
			Subject: addKubernetesSubject,
		})
		if err != nil {
			return err
		}

		return printAddedCredential(created)
	},
}

var (
	addGitScope       string
	addGitSubject     string
	addGitSecret      bool
	addGitSecretStdin bool
)

var userCredentialAddGitCmd = &cobra.Command{
	Use:   "git --scope <host> --subject <username> --secret",
	Short: "Add a Git credential",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if addGitSecret && addGitSecretStdin {
			return fmt.Errorf("--secret cannot be combined with --secret-stdin")
		}
		if !addGitSecret && !addGitSecretStdin {
			return fmt.Errorf("one of --secret or --secret-stdin is required")
		}

		secret, err := readPassword("Secret: ", addGitSecretStdin)
		if err != nil {
			return err
		}

		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		created, err := newClient(ctx).AddGitUserCredential(cmd.Context(), credentialUsername, models.UserGitCredentialRequest{
			Scope:   addGitScope,
			Subject: addGitSubject,
			Secret:  secret,
		})
		if err != nil {
			return err
		}

		return printAddedCredential(created)
	},
}

var (
	addRegistryScope       string
	addRegistrySubject     string
	addRegistrySecret      bool
	addRegistrySecretStdin bool
)

var userCredentialAddRegistryCmd = &cobra.Command{
	Use:   "registry --scope <registry-host> --subject <username> --secret",
	Short: "Add a container registry credential",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if addRegistrySecret && addRegistrySecretStdin {
			return fmt.Errorf("--secret cannot be combined with --secret-stdin")
		}
		if !addRegistrySecret && !addRegistrySecretStdin {
			return fmt.Errorf("one of --secret or --secret-stdin is required")
		}

		secret, err := readPassword("Secret: ", addRegistrySecretStdin)
		if err != nil {
			return err
		}

		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		created, err := newClient(ctx).AddRegistryUserCredential(cmd.Context(), credentialUsername, models.UserRegistryCredentialRequest{
			Scope:   addRegistryScope,
			Subject: addRegistrySubject,
			Secret:  secret,
		})
		if err != nil {
			return err
		}

		return printAddedCredential(created)
	},
}

func init() {
	addCredentialUsernameFlag(userCredentialAddKubernetesCmd)
	userCredentialAddKubernetesCmd.Flags().StringVar(&addKubernetesScope, "scope", "", "kubernetes namespace to scope the credential to (required)")
	userCredentialAddKubernetesCmd.Flags().StringVar(&addKubernetesSubject, "subject", "", "kubernetes service account name (required)")
	_ = userCredentialAddKubernetesCmd.MarkFlagRequired("scope")
	_ = userCredentialAddKubernetesCmd.MarkFlagRequired("subject")

	addCredentialUsernameFlag(userCredentialAddGitCmd)
	userCredentialAddGitCmd.Flags().StringVar(&addGitScope, "scope", "", "git host to scope the credential to, e.g. github.com (required)")
	userCredentialAddGitCmd.Flags().StringVar(&addGitSubject, "subject", "", "git username (required)")
	userCredentialAddGitCmd.Flags().BoolVar(&addGitSecret, "secret", false, "set the git password/token, prompted interactively")
	userCredentialAddGitCmd.Flags().BoolVar(&addGitSecretStdin, "secret-stdin", false, "set the git password/token by reading it from stdin")
	_ = userCredentialAddGitCmd.MarkFlagRequired("scope")
	_ = userCredentialAddGitCmd.MarkFlagRequired("subject")

	addCredentialUsernameFlag(userCredentialAddRegistryCmd)
	userCredentialAddRegistryCmd.Flags().StringVar(&addRegistryScope, "scope", "", "registry host (required)")
	userCredentialAddRegistryCmd.Flags().StringVar(&addRegistrySubject, "subject", "", "registry username (required)")
	userCredentialAddRegistryCmd.Flags().BoolVar(&addRegistrySecret, "secret", false, "set the registry password/token, prompted interactively")
	userCredentialAddRegistryCmd.Flags().BoolVar(&addRegistrySecretStdin, "secret-stdin", false, "set the registry password/token by reading it from stdin")
	_ = userCredentialAddRegistryCmd.MarkFlagRequired("scope")
	_ = userCredentialAddRegistryCmd.MarkFlagRequired("subject")

	userCredentialAddCmd.AddCommand(userCredentialAddKubernetesCmd)
	userCredentialAddCmd.AddCommand(userCredentialAddGitCmd)
	userCredentialAddCmd.AddCommand(userCredentialAddRegistryCmd)
	userCredentialCmd.AddCommand(userCredentialAddCmd)
}
