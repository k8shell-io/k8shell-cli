// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"time"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var sessionColumns = []table.Col[models.SSHSession]{
	{Header: "SESSION_ID", MaxWidth: 15, Help: "unique session identifier", Field: "sessionID"},
	{Header: "USERNAME", MaxWidth: 20, Help: "session owner", Field: "username"},
	{Header: "WORKSPACE", MaxWidth: 20, Help: "workspace the session is attached to", Field: "workspace"},
	{Header: "CLIENT_IP", MaxWidth: 15, Help: "IP address of the connecting client", Field: "clientIP"},
	{Header: "OPERATIONS", MaxWidth: 25, Help: "SSH operations (comma-separated)", Field: "operations", Fmt: table.FmtJoin},
	{Header: "START", MaxWidth: 16, Help: "session start time (local time)", Field: "startTime", Fmt: fmtTime},
	{Header: "END", MaxWidth: 16, Help: "session end time, or - if still active", Field: "endTime", Fmt: fmtTime},
	{Header: "BYTES_IN", MaxWidth: 10, Help: "bytes received from the client", Field: "bytesIn", Fmt: fmtBytes},
	{Header: "BYTES_OUT", MaxWidth: 10, Help: "bytes sent to the client", Field: "bytesOut", Fmt: fmtBytes},
}

var (
	sessionSortFlag      string
	sessionUsernameFlag  string
	sessionWorkspaceFlag string
	sessionLastFlag      int
	sessionAllFlag       bool
)

var sessionListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List sessions",
	Long: "List SSH sessions for your own user. Pass --username to look up another user " +
		"(admin only), or --all to see every session visible to the authenticated token.\n\n" + table.ColumnHelp(sessionColumns),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		username := sessionUsernameFlag
		if !cmd.Flags().Changed("username") && !sessionAllFlag {
			username = ctx.Username
		}

		sessions, err := newClient(ctx).ListSessions(cmd.Context(), username, sessionWorkspaceFlag, sessionLastFlag, sessionAllFlag)
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
	sessionListCmd.Flags().StringVarP(&sessionUsernameFlag, "username", "u", "", "look up another user's sessions instead of your own (admin only)")
	sessionListCmd.Flags().StringVarP(&sessionWorkspaceFlag, "workspace", "k", "", "filter by workspace")
	sessionListCmd.Flags().StringVar(&sessionSortFlag, "sort", "",
		"sort by fields, e.g. startTime,-bytesIn (prefix - for descending)")
	sessionListCmd.Flags().IntVarP(&sessionLastFlag, "last", "n", 0, "show only the last N sessions")
	sessionListCmd.Flags().BoolVar(&sessionAllFlag, "all", false, "include all sessions visible to the token, instead of just your own")
	_ = sessionListCmd.RegisterFlagCompletionFunc("username", completeUsernames)
}

// fmtBytes renders a byte count as a human-readable string (B, KB, MB, GB).
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

// fmtTime renders a *time.Time field as a local "YYYY-MM-DD HH:MM" string, or "-" if nil.
func fmtTime(v any) string {
	t, _ := v.(*time.Time)
	if t == nil {
		return "-"
	}
	return t.Local().Format("2006-01-02 15:04")
}
