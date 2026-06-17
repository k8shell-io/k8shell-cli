package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/k8shell-io/k8shell/internal/output"
	"github.com/spf13/cobra"
)

var sessionUser string

var sessionColumns = []output.Column{
	{Header: "SESSION_ID", MaxWidth: 15},
	{Header: "USERNAME", MaxWidth: 20},
	{Header: "WORKSPACE", MaxWidth: 20},
	{Header: "BLUEPRINT", MaxWidth: 20},
	{Header: "CLIENT_IP", MaxWidth: 15},
	{Header: "CHANNELS", MaxWidth: 25},
	{Header: "START", MaxWidth: 16},
	{Header: "END", MaxWidth: 16},
	{Header: "BYTES_IN", MaxWidth: 10},
	{Header: "BYTES_OUT", MaxWidth: 10},
}

var userSessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sessions (defaults to your own via /me)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		sessions, err := client.New(ctx, debug).ListSessions(sessionUser)
		if err != nil {
			return err
		}

		if printer.IsJSON() {
			return printer.JSON(sessions)
		}

		rows := make([][]string, len(sessions))
		for i, s := range sessions {
			rows[i] = []string{
				s.SessionID,
				s.Username,
				s.Workspace,
				s.Blueprint,
				s.ClientIP,
				strings.Join(s.Channels, ","),
				formatTime(s.StartTime),
				formatTime(s.EndTime),
				formatBytes(s.BytesIn),
				formatBytes(s.BytesOut),
			}
		}
		printer.Table(sessionColumns, rows)
		return nil
	},
}

func init() {
	userSessionListCmd.Flags().StringVarP(&sessionUser, "user", "u", "", "username (defaults to the authenticated user)")
}

func formatBytes(b int64) string {
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

func formatTime(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Local().Format("2006-01-02 15:04")
}
