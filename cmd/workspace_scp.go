// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var wsSCPPath string

var workspaceSCPCmd = &cobra.Command{
	Use:   "scp [flags] [local-destination]",
	Short: "Copy files from a workspace using SCP",
	Long: `Copy files from a workspace using SCP.

One of --pod or --repo is required to identify the workspace.
Use --path to specify the remote path to copy from (default: /).
Use --print to print the SCP string instead of running it.

The optional positional argument sets the local destination (default: .).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ustr, server, _, err := wsConnResolve(false)
		if err != nil {
			return err
		}
		localDest := "."
		if len(args) > 0 {
			localDest = args[0]
		}
		remote := fmt.Sprintf("%s@%s:%s", ustr, server, wsSCPPath)
		scpStr := fmt.Sprintf("scp %s %s", remote, localDest)

		if wsConnPrint || printer.IsJSON() {
			if printer.IsJSON() {
				return printer.JSON(struct {
					SCP        string `json:"scp"`
					Userstring string `json:"userstring"`
					Server     string `json:"server"`
					Path       string `json:"path"`
				}{SCP: scpStr, Userstring: ustr, Server: server, Path: wsSCPPath})
			}
			fmt.Println(scpStr)
			return nil
		}

		c := exec.Command("scp", remote, localDest)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	workspaceSCPCmd.Flags().StringVar(&wsConnPod, "pod", "", "workspace name")
	workspaceSCPCmd.Flags().StringVar(&wsConnRepo, "repo", "", "repo workspace in owner/name format")
	workspaceSCPCmd.Flags().StringVar(&wsConnRef, "ref", "", "git ref (only with --repo)")
	workspaceSCPCmd.Flags().StringVar(&wsConnHost, "host", "", "override SSH server hostname")
	workspaceSCPCmd.Flags().BoolVar(&wsConnB64, "b64", false, "encode userstring as base64")
	workspaceSCPCmd.Flags().BoolVarP(&wsConnPrint, "print", "p", false, "print SCP string instead of running")
	workspaceSCPCmd.Flags().StringVar(&wsSCPPath, "path", "/", "remote path to copy from")
	workspaceSCPCmd.Flags().BoolVar(&wsConnNoLookup, "no-lookup", false, "skip workspace API lookup; uses context username and derives server from context URL")
	_ = workspaceSCPCmd.RegisterFlagCompletionFunc("pod", completeWorkspaceNames)
}
