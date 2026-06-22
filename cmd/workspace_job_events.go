// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
	"github.com/k8shell-io/k8shell/internal/client"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var workspaceJobEventsCmd = &cobra.Command{
	Use:               "job-events <workspace-name>",
	Short:             "Stream job events for a workspace",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeWorkspaceNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := cfg.ActiveContext()
		if err != nil {
			return err
		}

		c := client.New(ctx, debug, insecure || ctx.Insecure)

		ws, err := c.GetWorkspace(args[0])
		if err != nil {
			return err
		}
		if ws.JobId == "" {
			return fmt.Errorf("workspace %q has no associated job", args[0])
		}

		rc, err := c.MonitorWorkspace("/api/v1/jobs/" + ws.JobId)
		if err != nil {
			return err
		}
		defer rc.Close()

		return printEventStream(rc)
	},
}

// printProgressStream reads an SSE stream and updates a single progress line in place.
// Progress events update the percentage; all other events are suppressed.
// A trailing newline is printed when the stream closes.
func printProgressStream(rc io.Reader, workspace string) error {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		data := line
		if after, ok := strings.CutPrefix(line, "data:"); ok {
			data = strings.TrimPrefix(after, " ")
		}
		var event models.WorkspaceStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}
		if event.Type == models.WorkspaceStreamEventTypeProgress {
			pct := strings.TrimSuffix(event.Message, " complete")
			fmt.Printf("\rCreating workspace %s... %s", workspace, pct)
		}
	}
	fmt.Println()
	return scanner.Err()
}

// printEventStream reads an SSE stream from rc and prints log events to stdout,
// skipping progress events. Lines are trimmed to the terminal width unless -w is set.
func printEventStream(rc io.Reader) error {
	termW := 0
	if !wrap {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
			termW = w
		}
	}

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		data := line
		if after, ok := strings.CutPrefix(line, "data:"); ok {
			data = strings.TrimPrefix(after, " ")
		}
		var event models.WorkspaceStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}
		if event.Type == models.WorkspaceStreamEventTypeProgress {
			continue
		}
		s := event.String()
		if s == "" {
			s = event.Message
		}
		if s == "" {
			continue
		}
		if termW > 0 && len([]rune(s)) > termW {
			s = string([]rune(s)[:termW])
		}
		fmt.Println(s)
	}
	return scanner.Err()
}
