package cmd

import (
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/k8shell-io/k8shell/internal/output"
	"github.com/spf13/cobra"
)

var userColumns = []output.Column{
	{Header: "USERNAME", MaxWidth: 20},
	{Header: "FULLNAME", MaxWidth: 20},
	{Header: "EMAIL", MaxWidth: 30},
	{Header: "ORG", MaxWidth: 15},
	{Header: "ROLES", MaxWidth: 20},
	{Header: "BLUEPRINTS", MaxWidth: 30},
	{Header: "SUDO", MaxWidth: 5},
	{Header: "SOURCE", MaxWidth: 122},
	{Header: "STATUS", MaxWidth: 8},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		users, err := client.New(ctx, debug).ListUsers()
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(users)
		}

		rows := make([][]string, len(users))
		for i, u := range users {
			rows[i] = []string{
				u.Username,
				u.Fullname,
				u.Email,
				u.Organization,
				formatRoles(u.Roles),
				strings.Join(u.Blueprints, ","),
				boolVal(u.Sudo),
				u.Source,
				userStatus(u),
			}
		}
		printer.Table(userColumns, rows)
		return nil
	},
}

func formatRoles(roles []models.Role) string {
	s := make([]string, len(roles))
	for i, r := range roles {
		s[i] = string(r)
	}
	return strings.Join(s, ",")
}

func boolVal(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func userStatus(u models.User) string {
	if u.Locked {
		return "locked"
	}
	if !u.IsValid {
		return "invalid"
	}
	return "active"
}
