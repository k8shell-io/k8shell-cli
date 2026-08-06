// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var (
	setFullname           string
	setShell              string
	setEmail              string
	setOrg                string
	setSudo               string
	setUID                uint32
	setGID                uint32
	setLock               bool
	setUnlock             bool
	setPassword           bool
	setPasswordStdin      bool
	setPreserveWorkspaces bool

	setRoles       []string
	setAddRoles    []string
	setRemoveRoles []string

	setAddKeys        []string
	setAddKeyFiles    []string
	setRemoveKeyIndex []int
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
		if setLock && setUnlock {
			return fmt.Errorf("--lock cannot be combined with --unlock")
		}
		if setPassword && setPasswordStdin {
			return fmt.Errorf("--password cannot be combined with --password-stdin")
		}

		username := args[0]
		c := newClient(ctx)

		addKeys := append([]string{}, setAddKeys...)
		for _, path := range setAddKeyFiles {
			key, err := loadPublicKeyFile(path)
			if err != nil {
				return err
			}
			addKeys = append(addKeys, key)
		}

		var req models.UserUpdateRequest
		var updated []string
		var profileFieldsChanged bool

		if cmd.Flags().Changed("fullname") {
			req.Fullname = &setFullname
			profileFieldsChanged = true
			updated = append(updated, "fullname")
		}
		if cmd.Flags().Changed("shell") {
			req.Shell = &setShell
			profileFieldsChanged = true
			updated = append(updated, "shell")
		}
		if cmd.Flags().Changed("email") {
			req.Email = &setEmail
			profileFieldsChanged = true
			updated = append(updated, "email")
		}
		if cmd.Flags().Changed("uid") {
			req.UID = &setUID
			profileFieldsChanged = true
			updated = append(updated, "uid")
		}
		if cmd.Flags().Changed("gid") {
			req.GID = &setGID
			profileFieldsChanged = true
			updated = append(updated, "gid")
		}
		if cmd.Flags().Changed("org") {
			req.Org = &setOrg
			req.PreserveWorkspaces = setPreserveWorkspaces
			profileFieldsChanged = true
			updated = append(updated, "org")
		}
		if cmd.Flags().Changed("roles") {
			req.Roles = toRoles(setRoles)
			profileFieldsChanged = true
			updated = append(updated, "roles")
		}
		if cmd.Flags().Changed("sudo") {
			sudo, err := parseBool(setSudo)
			if err != nil {
				return fmt.Errorf("--sudo: %w", err)
			}
			req.Sudo = &sudo
			profileFieldsChanged = true
			updated = append(updated, "sudo")
		}
		if setLock {
			locked := true
			req.Locked = &locked
			profileFieldsChanged = true
			updated = append(updated, "locked")
		}
		if setUnlock {
			profile, err := c.GetUserProfile(cmd.Context(), username)
			if err != nil {
				return fmt.Errorf("checking lock state: %w", err)
			}
			if profile.PasswordLocked {
				if err := c.ClearUserPasswordLockout(cmd.Context(), username); err != nil {
					return err
				}
			}
			if profile.AccountLocked {
				locked := false
				req.Locked = &locked
				profileFieldsChanged = true
			}
			updated = append(updated, "unlocked")
		}

		if profileFieldsChanged {
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
		if len(setRemoveKeyIndex) > 0 {
			for _, idx := range setRemoveKeyIndex {
				if err := c.RemoveUserAuthKey(cmd.Context(), username, idx); err != nil {
					return err
				}
			}
			updated = append(updated, "remove-key")
		}
		if len(addKeys) > 0 {
			if err := c.AddUserKeys(cmd.Context(), username, addKeys); err != nil {
				return err
			}
			updated = append(updated, "add-key")
		}
		if setPassword || setPasswordStdin {
			var currentPassword string
			if ctx.Username != "" && username == ctx.Username {
				profile, err := c.GetProfile(cmd.Context())
				if err != nil {
					return fmt.Errorf("checking sudo access: %w", err)
				}
				if !profile.Sudo {
					currentPassword, err = readPassword("Current password: ", setPasswordStdin)
					if err != nil {
						return err
					}
				}
			}
			password, err := readPassword("New password: ", setPasswordStdin)
			if err != nil {
				return err
			}
			if _, err := c.SetUserPassword(cmd.Context(), username, password, currentPassword); err != nil {
				return err
			}
			updated = append(updated, "password")
		}

		if len(updated) == 0 {
			return fmt.Errorf("specify at least one field to update (--fullname, --shell, --email, --org, --uid, " +
				"--gid, --roles, --sudo, --lock, --unlock, --password, --password-stdin, " +
				"--add-role, --remove-role, --add-key, --add-key-file, --remove-key)")
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
	userSetCmd.Flags().BoolVar(&setPreserveWorkspaces, "preserve-workspaces", false, "when combined with --org, keep the user's workspaces instead of deleting them on the org move")
	userSetCmd.Flags().StringVar(&setSudo, "sudo", "", "sudo access (true/false)")
	userSetCmd.Flags().BoolVar(&setLock, "lock", false, "lock the account")
	userSetCmd.Flags().BoolVar(&setUnlock, "unlock", false, "unlock the account")
	userSetCmd.Flags().BoolVar(&setPassword, "password", false, "set the account's local password")
	userSetCmd.Flags().BoolVar(&setPasswordStdin, "password-stdin", false, "set the account's local password by reading it from stdin")

	userSetCmd.Flags().StringSliceVar(&setRoles, "roles", nil, "replace all roles, comma-separated")
	repeatableStringVar(userSetCmd.Flags(), &setAddRoles, "add-role", "grant a role, in addition to existing roles (repeatable)")
	repeatableStringVar(userSetCmd.Flags(), &setRemoveRoles, "remove-role", "revoke a role, leaving others untouched (repeatable)")

	repeatableStringVar(userSetCmd.Flags(), &setAddKeys, "add-key", "add an SSH public key, in addition to existing ones (repeatable)")
	repeatableStringVar(userSetCmd.Flags(), &setAddKeyFiles, "add-key-file", "add an SSH public key read from a file — public or private, PEM or OpenSSH format (repeatable)")
	userSetCmd.Flags().IntSliceVar(&setRemoveKeyIndex, "remove-key", nil, "remove an SSH public key by its index (repeatable)")
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

// stdinPasswords lazily wraps os.Stdin so successive stdin-mode readPassword
// calls each consume one line, letting --password-stdin supply more than one
// secret (e.g. current password followed by new password) in sequence.
var stdinPasswords = bufio.NewReader(os.Stdin)

// readPassword obtains a password either from stdin (one line, trimmed of
// trailing line endings) or by prompting on stderr with echo disabled,
// depending on stdin. It returns an error if the resulting password is empty.
func readPassword(prompt string, stdin bool) (string, error) {
	var password string
	if stdin {
		line, err := stdinPasswords.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", fmt.Errorf("reading password from stdin: %w", err)
		}
		password = strings.TrimRight(line, "\r\n")
	} else {
		fmt.Fprint(os.Stderr, prompt)
		raw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", fmt.Errorf("reading password: %w", err)
		}
		fmt.Fprintln(os.Stderr)
		password = strings.TrimSpace(string(raw))
	}
	if password == "" {
		return "", fmt.Errorf("password must not be empty")
	}
	return password, nil
}

// loadPublicKeyFile reads the file at path and returns its contents as a single
// SSH authorized_keys line, converting it from PEM-encoded public/private key
// format if necessary. It fails locally if the file's contents cannot be turned
// into a valid SSH public key, so a malformed key is never sent to the server.
func loadPublicKeyFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading key file %q: %w", path, err)
	}
	key, err := toAuthorizedKey(data, path)
	if err != nil {
		return "", fmt.Errorf("converting key file %q to SSH public key format: %w", path, err)
	}
	return key, nil
}

// toAuthorizedKey converts raw key material into a single SSH authorized_keys
// line. It accepts input already in authorized_keys format, PEM-encoded (PKIX)
// public keys, and private keys (PEM or OpenSSH format, optionally passphrase
// protected), deriving the public key from the latter. path is used only to
// prompt for a passphrase if the private key is encrypted.
func toAuthorizedKey(data []byte, path string) (string, error) {
	if pub, comment, _, _, err := ssh.ParseAuthorizedKey(data); err == nil {
		line := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pub)))
		if comment != "" {
			line += " " + comment
		}
		return line, nil
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return "", fmt.Errorf("not a valid SSH public key, private key, or PEM-encoded key")
	}

	if strings.Contains(block.Type, "PRIVATE KEY") {
		return authorizedKeyFromPrivateKey(data, path)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parsing PEM public key: %w", err)
	}
	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("unsupported public key type: %w", err)
	}
	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPub))), nil
}

// authorizedKeyFromPrivateKey derives the SSH public key line for a PEM or
// OpenSSH format private key. If the key is passphrase-protected, the
// passphrase is prompted for securely on stderr, the same way --password does.
func authorizedKeyFromPrivateKey(data []byte, path string) (string, error) {
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		var missing *ssh.PassphraseMissingError
		if !errors.As(err, &missing) {
			return "", fmt.Errorf("parsing private key: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Passphrase for %s: ", path)
		raw, perr := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if perr != nil {
			return "", fmt.Errorf("reading passphrase: %w", perr)
		}
		signer, err = ssh.ParsePrivateKeyWithPassphrase(data, raw)
		if err != nil {
			return "", fmt.Errorf("parsing private key: %w", err)
		}
	}
	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(signer.PublicKey()))), nil
}
