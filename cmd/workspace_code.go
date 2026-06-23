// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var wsCodePath string

var workspaceCodeCmd = &cobra.Command{
	Use:   "code [flags]",
	Short: "Open a workspace in VS Code via the remote SSH extension",
	Long: `Open a workspace in VS Code via the remote SSH extension.

One of --pod or --repo is required to identify the workspace.
Use --path to set the folder to open in VS Code (default: /).
Use --print to print the connection string instead of launching VS Code.

The userstring is always base64-encoded because VS Code does not accept
raw userstrings in remote URIs.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ustr, server, username, err := wsConnResolve(cmd.Context(), true) // always b64
		if err != nil {
			return err
		}
		path := wsCodePath
		if wsConnRepo != "" && !cmd.Flags().Changed("path") {
			parts := strings.SplitN(wsConnRepo, "/", 2)
			path = "/home/" + username + "/" + parts[1]
		}
		uri := fmt.Sprintf("vscode-remote://ssh-remote+%s@%s%s", ustr, server, path)
		codeStr := fmt.Sprintf("code --folder-uri %q", uri)

		if wsConnPrint || printer.IsJSON() {
			if printer.IsJSON() {
				return printer.JSON(struct {
					Code       string `json:"code"`
					URI        string `json:"uri"`
					Userstring string `json:"userstring"`
					Server     string `json:"server"`
					Path       string `json:"path"`
				}{Code: codeStr, URI: uri, Userstring: ustr, Server: server, Path: path})
			}
			fmt.Println(codeStr)
			return nil
		}

		c := exec.Command("code", "--folder-uri", uri)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	workspaceCodeCmd.Flags().StringVar(&wsConnPod, "pod", "", "workspace name")
	workspaceCodeCmd.Flags().StringVar(&wsConnRepo, "repo", "", "repo workspace in owner/name format")
	workspaceCodeCmd.Flags().StringVar(&wsConnRef, "ref", "", "git ref (only with --repo)")
	workspaceCodeCmd.Flags().StringVar(&wsConnHost, "host", "", "override SSH server hostname")
	workspaceCodeCmd.Flags().BoolVarP(&wsConnPrint, "print", "p", false, "print connection string instead of launching VS Code")
	workspaceCodeCmd.Flags().StringVar(&wsCodePath, "path", "/", "folder path to open in VS Code")
	workspaceCodeCmd.Flags().BoolVar(&wsConnNoLookup, "no-lookup", false, "skip workspace API lookup; uses context username and derives server from context URL")
	_ = workspaceCodeCmd.RegisterFlagCompletionFunc("pod", completeWorkspaceNames)
}
