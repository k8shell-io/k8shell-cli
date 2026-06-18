package cmd

import (
	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var userColumns = []table.Col[models.User]{
	{Header: "USERNAME",   MaxWidth: 20,  Field: "username"},
	{Header: "FULLNAME",   MaxWidth: 20,  Field: "fullname"},
	{Header: "EMAIL",      MaxWidth: 30,  Field: "email"},
	{Header: "ORG",        MaxWidth: 15,  Field: "organization"},
	{Header: "ROLES",      MaxWidth: 20,  Field: "roles",      Fmt: fmtRoles},
	{Header: "BLUEPRINTS", MaxWidth: 30,  Field: "blueprints", Fmt: fmtJoin},
	{Header: "SUDO",       MaxWidth: 5,   Field: "sudo",       Fmt: fmtBool},
	{Header: "SOURCE",     MaxWidth: 122, Field: "source"},
	{Header: "STATUS",     MaxWidth: 8,   Fn: userStatus},
}

var userSortFlag string

var userListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List users",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		users, err := client.New(ctx, debug, insecure || ctx.Insecure).ListUsers()
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(users)
		}

		return table.Table(printer, userColumns, users, userSortFlag)
	},
}

func init() {
	userListCmd.Flags().StringVar(&userSortFlag, "sort", "", "sort by fields, e.g. username,-email (prefix - for descending)")
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
