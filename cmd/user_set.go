// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
)

var (
	setFullname string
	setShell    string
	setEmail    string
	setOrg      string
	setSudo     string
	setUID      uint32
	setGID      uint32
	setLock     bool
	setUnlock   bool

	setRoles       []string
	setAddRoles    []string
	setRemoveRoles []string

	setBlueprints       []string
	setAddBlueprints    []string
	setRemoveBlueprints []string

	setAddKeys    []string
	setRemoveKeys []string
)

var userSetCmd = &cobra.Command{
	Use:               "set <username> [flags]",
	Short:             "Update fields on a user",
	Long:              "Update one or more fields on a user.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeUsernames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		if cmd.Flags().Changed("roles") && (cmd.Flags().Changed("add-role") || cmd.Flags().Changed("remove-role")) {
			return fmt.Errorf("--roles cannot be combined with --add-role/--remove-role")
		}
		if cmd.Flags().Changed("blueprints") && (cmd.Flags().Changed("add-blueprint") || cmd.Flags().Changed("remove-blueprint")) {
			return fmt.Errorf("--blueprints cannot be combined with --add-blueprint/--remove-blueprint")
		}
		if setLock && setUnlock {
			return fmt.Errorf("--lock cannot be combined with --unlock")
		}

		username := args[0]
		c := newClient(ctx)

		var req models.UserUpdateRequest
		var updated []string

		if cmd.Flags().Changed("fullname") {
			req.Fullname = &setFullname
			updated = append(updated, "fullname")
		}
		if cmd.Flags().Changed("shell") {
			req.Shell = &setShell
			updated = append(updated, "shell")
		}
		if cmd.Flags().Changed("email") {
			req.Email = &setEmail
			updated = append(updated, "email")
		}
		if cmd.Flags().Changed("uid") {
			req.UID = &setUID
			updated = append(updated, "uid")
		}
		if cmd.Flags().Changed("gid") {
			req.GID = &setGID
			updated = append(updated, "gid")
		}
		if cmd.Flags().Changed("org") {
			req.Org = &setOrg
			updated = append(updated, "org")
		}
		if cmd.Flags().Changed("roles") {
			req.Roles = toRoles(setRoles)
			updated = append(updated, "roles")
		}
		if cmd.Flags().Changed("sudo") {
			sudo, err := parseBool(setSudo)
			if err != nil {
				return fmt.Errorf("--sudo: %w", err)
			}
			req.Sudo = &sudo
			updated = append(updated, "sudo")
		}
		if cmd.Flags().Changed("blueprints") {
			req.Blueprints = setBlueprints
			updated = append(updated, "blueprints")
		}
		if setLock {
			locked := true
			req.Locked = &locked
			updated = append(updated, "locked")
		}
		if setUnlock {
			locked := false
			req.Locked = &locked
			updated = append(updated, "unlocked")
		}

		if len(updated) > 0 {
			if _, err := c.UpdateUserProfile(cmd.Context(), username, req); err != nil {
				return err
			}
		}

		if len(setRemoveRoles) > 0 {
			if err := c.RemoveUserRoles(cmd.Context(), username, toRoles(setRemoveRoles)); err != nil {
				return err
			}
			updated = append(updated, "remove-role")
		}
		if len(setAddRoles) > 0 {
			if err := c.AddUserRoles(cmd.Context(), username, toRoles(setAddRoles)); err != nil {
				return err
			}
			updated = append(updated, "add-role")
		}
		if len(setRemoveBlueprints) > 0 {
			if err := c.RemoveUserBlueprints(cmd.Context(), username, setRemoveBlueprints); err != nil {
				return err
			}
			updated = append(updated, "remove-blueprint")
		}
		if len(setAddBlueprints) > 0 {
			if err := c.AddUserBlueprints(cmd.Context(), username, setAddBlueprints); err != nil {
				return err
			}
			updated = append(updated, "add-blueprint")
		}
		if len(setRemoveKeys) > 0 {
			if err := c.RemoveUserKeys(cmd.Context(), username, setRemoveKeys); err != nil {
				return err
			}
			updated = append(updated, "remove-key")
		}
		if len(setAddKeys) > 0 {
			if err := c.AddUserKeys(cmd.Context(), username, setAddKeys); err != nil {
				return err
			}
			updated = append(updated, "add-key")
		}

		if len(updated) == 0 {
			return fmt.Errorf("specify at least one field to update (--fullname, --shell, --email, --org, --uid, --gid, --roles, --sudo, --blueprints, --lock, --unlock, " +
				"--add-role, --remove-role, --add-blueprint, --remove-blueprint, --add-key, --remove-key)")
		}

		if printer.IsJSON() {
			return printer.JSON(map[string]any{"username": username, "updated": updated})
		}

		printer.Println(fmt.Sprintf("%s: updated %s", username, strings.Join(updated, ", ")))
		return nil
	},
}

func init() {
	userSetCmd.Flags().StringVar(&setFullname, "fullname", "", "display name")
	userSetCmd.Flags().StringVar(&setShell, "shell", "", "login shell (e.g. /bin/bash)")
	userSetCmd.Flags().StringVar(&setEmail, "email", "", "email address")
	userSetCmd.Flags().Uint32Var(&setUID, "uid", 0, "numeric user ID")
	userSetCmd.Flags().Uint32Var(&setGID, "gid", 0, "numeric group ID")
	userSetCmd.Flags().StringVar(&setOrg, "org", "", "organization")
	userSetCmd.Flags().StringVar(&setSudo, "sudo", "", "sudo access (true/false)")
	userSetCmd.Flags().BoolVar(&setLock, "lock", false, "lock the account")
	userSetCmd.Flags().BoolVar(&setUnlock, "unlock", false, "unlock the account")

	userSetCmd.Flags().StringSliceVar(&setRoles, "roles", nil, "replace all roles, comma-separated (e.g. admin,workspace-user)")
	userSetCmd.Flags().StringArrayVar(&setAddRoles, "add-role", nil, "grant a role, in addition to existing roles (repeatable)")
	userSetCmd.Flags().StringArrayVar(&setRemoveRoles, "remove-role", nil, "revoke a role, leaving others untouched (repeatable)")

	userSetCmd.Flags().StringSliceVar(&setBlueprints, "blueprints", nil, "replace all allowed blueprints, comma-separated")
	userSetCmd.Flags().StringArrayVar(&setAddBlueprints, "add-blueprint", nil, "grant a blueprint, in addition to existing ones (repeatable)")
	userSetCmd.Flags().StringArrayVar(&setRemoveBlueprints, "remove-blueprint", nil, "revoke a blueprint, leaving others untouched (repeatable)")

	userSetCmd.Flags().StringArrayVar(&setAddKeys, "add-key", nil, "add an SSH public key, in addition to existing ones (repeatable)")
	userSetCmd.Flags().StringArrayVar(&setRemoveKeys, "remove-key", nil, "remove an SSH public key, leaving others untouched (repeatable)")
}

// toRoles converts role name strings to models.Role values.
func toRoles(names []string) []models.Role {
	roles := make([]models.Role, len(names))
	for i, r := range names {
		roles[i] = models.Role(r)
	}
	return roles
}

// parseBool parses a "true"/"false" flag value into a bool.
func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("must be %q or %q", "true", "false")
	}
}
