// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var workspaceSSHCmd = &cobra.Command{
	Use:   "ssh [flags] [-- ssh-args...]",
	Short: "Open an SSH session to a workspace",
	Long: `Open an SSH session to a workspace.

One of --pod or --repo is required to identify the workspace.
Use --print to print the connection string instead of connecting.

Any arguments after -- are passed directly to ssh.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ustr, server, _, err := wsConnResolve(false)
		if err != nil {
			return err
		}
		target := ustr + "@" + server
		sshStr := "ssh " + target

		if wsConnPrint || printer.IsJSON() {
			if printer.IsJSON() {
				return printer.JSON(struct {
					SSH        string `json:"ssh"`
					Userstring string `json:"userstring"`
					Server     string `json:"server"`
				}{SSH: sshStr, Userstring: ustr, Server: server})
			}
			fmt.Println(sshStr)
			return nil
		}

		c := exec.Command("ssh", append([]string{target}, args...)...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	workspaceSSHCmd.Flags().StringVar(&wsConnPod, "pod", "", "workspace name")
	workspaceSSHCmd.Flags().StringVar(&wsConnRepo, "repo", "", "repo workspace in owner/name format")
	workspaceSSHCmd.Flags().StringVar(&wsConnRef, "ref", "", "git ref (only with --repo)")
	workspaceSSHCmd.Flags().StringVar(&wsConnHost, "host", "", "override SSH server hostname")
	workspaceSSHCmd.Flags().BoolVar(&wsConnB64, "b64", false, "encode userstring as base64")
	workspaceSSHCmd.Flags().BoolVarP(&wsConnPrint, "print", "p", false, "print connection string instead of connecting")
	workspaceSSHCmd.Flags().BoolVar(&wsConnNoLookup, "no-lookup", false, "skip workspace API lookup; uses context username and derives server from context URL")
	_ = workspaceSSHCmd.RegisterFlagCompletionFunc("pod", completeWorkspaceNames)
}
