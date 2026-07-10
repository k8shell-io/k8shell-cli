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

		if printer.IsJSON() {
			return printer.JSON(created)
		}

		printer.Println(fmt.Sprintf("%s: credential added (id %d)", created.ServiceName, created.ID))
		return nil
	},
}

func init() {
	addCredentialUsernameFlag(userCredentialAddKubernetesCmd)
	userCredentialAddKubernetesCmd.Flags().StringVar(&addKubernetesScope, "scope", "", "kubernetes namespace to scope the credential to (required)")
	userCredentialAddKubernetesCmd.Flags().StringVar(&addKubernetesSubject, "subject", "", "kubernetes service account name (required)")
	_ = userCredentialAddKubernetesCmd.MarkFlagRequired("scope")
	_ = userCredentialAddKubernetesCmd.MarkFlagRequired("subject")

	userCredentialAddCmd.AddCommand(userCredentialAddKubernetesCmd)
	userCredentialCmd.AddCommand(userCredentialAddCmd)
}
