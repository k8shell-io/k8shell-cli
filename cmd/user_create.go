// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	createUserOrg           string
	createUserFullname      string
	createUserEmail         string
	createUserShell         string
	createUserUID           uint32
	createUserGID           uint32
	createUserSudo          bool
	createUserLocked        bool
	createUserPassword      bool
	createUserPasswordStdin bool
	createUserRoles         []string
	createUserBlueprints    []string
)

var userCreateCmd = &cobra.Command{
	Use:   "create <username> --org <org> [flags]",
	Short: "Create a new user",
	Long:  "Create a new local user with no backing identity provider.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		if createUserPassword && createUserPasswordStdin {
			return fmt.Errorf("--password cannot be combined with --password-stdin")
		}

		username := args[0]
		req := models.UserCreateRequest{
			Username:   username,
			Org:        createUserOrg,
			Fullname:   createUserFullname,
			Email:      createUserEmail,
			Shell:      createUserShell,
			Sudo:       createUserSudo,
			Locked:     createUserLocked,
			UID:        createUserUID,
			GID:        createUserGID,
			Roles:      toRoles(createUserRoles),
			Blueprints: createUserBlueprints,
		}

		if createUserPassword || createUserPasswordStdin {
			password, err := readPassword("Password: ", createUserPasswordStdin)
			if err != nil {
				return err
			}
			req.Password = password
		}

		user, err := newClient(ctx).CreateUser(cmd.Context(), req)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(user)
		}

		printer.Println(fmt.Sprintf("%s: created", user.Username))
		return nil
	},
}

func init() {
	userCreateCmd.Flags().StringVar(&createUserOrg, "org", "", "organization (required)")
	userCreateCmd.Flags().StringVar(&createUserFullname, "fullname", "", "display name")
	userCreateCmd.Flags().StringVar(&createUserEmail, "email", "", "email address")
	userCreateCmd.Flags().StringVar(&createUserShell, "shell", "", "login shell (e.g. /bin/bash)")
	userCreateCmd.Flags().Uint32Var(&createUserUID, "uid", 0, "numeric user ID")
	userCreateCmd.Flags().Uint32Var(&createUserGID, "gid", 0, "numeric group ID")
	userCreateCmd.Flags().BoolVar(&createUserSudo, "sudo", false, "grant sudo access")
	userCreateCmd.Flags().BoolVar(&createUserLocked, "locked", false, "create the account in a locked state")
	userCreateCmd.Flags().BoolVar(&createUserPassword, "password", false, "set the account's local password")
	userCreateCmd.Flags().BoolVar(&createUserPasswordStdin, "password-stdin", false, "set the account's local password by reading it from stdin")
	userCreateCmd.Flags().StringSliceVar(&createUserRoles, "roles", nil, "roles to grant, comma-separated")
	userCreateCmd.Flags().StringSliceVar(&createUserBlueprints, "blueprints", nil, "allowed blueprints, comma-separated")
	_ = userCreateCmd.MarkFlagRequired("org")
}
