package cmd

import (
	"fmt"
	"time"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var sessionUser string

var sessionColumns = []table.Col[models.SSHSession]{
	{Header: "SESSION_ID", MaxWidth: 15, Field: "sessionID"},
	{Header: "USERNAME", MaxWidth: 20, Field: "username"},
	{Header: "WORKSPACE", MaxWidth: 20, Field: "workspace"},
	{Header: "BLUEPRINT", MaxWidth: 20, Field: "blueprint"},
	{Header: "CLIENT_IP", MaxWidth: 15, Field: "clientIP"},
	{Header: "CHANNELS", MaxWidth: 25, Field: "channels", Fmt: fmtJoin},
	{Header: "START", MaxWidth: 16, Field: "startTime", Fmt: fmtTime},
	{Header: "END", MaxWidth: 16, Field: "endTime", Fmt: fmtTime},
	{Header: "BYTES_IN", MaxWidth: 10, Field: "bytesIn", Fmt: fmtBytes},
	{Header: "BYTES_OUT", MaxWidth: 10, Field: "bytesOut", Fmt: fmtBytes},
}

var sessionSortFlag string

var userSessionListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List sessions (defaults to your own via /me)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		sessions, err := client.New(ctx, debug, insecure).ListSessions(sessionUser)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(sessions)
		}

		return table.Table(printer, sessionColumns, sessions, sessionSortFlag)
	},
}

func init() {
	userSessionListCmd.Flags().StringVarP(&sessionUser, "user", "u", "",
		"username (defaults to the authenticated user)")
	userSessionListCmd.Flags().StringVar(&sessionSortFlag, "sort", "",
		"sort by fields, e.g. startTime,-bytesIn (prefix - for descending)")
}

func fmtBytes(v any) string {
	b, _ := v.(int64)
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func fmtTime(v any) string {
	t, _ := v.(*time.Time)
	if t == nil {
		return "-"
	}
	return t.Local().Format("2006-01-02 15:04")
}
