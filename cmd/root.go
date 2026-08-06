// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"errors"
	"fmt"
	"os"

	k8shell "github.com/k8shell-io/k8shell-go"
	"github.com/k8shell-io/k8shell/internal/config"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	contextName string
	jsonMode    bool
	noANSI      bool
	wrap        bool
	debug       bool
	curl        bool
	verbose     bool
	location    bool
	insecure    bool
	cfg         *config.Config
	printer     *table.Printer
)

var rootCmd = &cobra.Command{
	Use:   "k8shell",
	Short: "Manage users, workspaces, SSH sessions, and server contexts for k8shell",
	Long: `k8shell connects to a k8shell server and provides commands for managing
its resources.

Run 'k8shell --help' on any command to see its available options.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug && curl {
			return fmt.Errorf("--curl cannot be used with --debug")
		}
		if verbose && !curl {
			return fmt.Errorf("--verbose requires --curl")
		}
		if location && !curl {
			return fmt.Errorf("--location requires --curl")
		}
		printer = table.New(jsonMode, noANSI, wrap)

		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		if contextName != "" {
			cfg.CurrentContext = contextName
		}
		return nil
	},
}

// Execute runs the root cobra command and exits with a non-zero status on error.
// k8shell.ErrDryRun (from --curl) is treated as success: the curl command has
// already been printed and there is no request result to act on or report.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, k8shell.ErrDryRun) {
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.config/k8shell/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&contextName, "context", "c", "", "override the active context")
	rootCmd.PersistentFlags().BoolVar(&jsonMode, "json", false, "output as JSON")
	rootCmd.PersistentFlags().BoolVar(&noANSI, "no-ansi", false, "disable ANSI color output")
	rootCmd.PersistentFlags().BoolVarP(&wrap, "wrap", "w", false, "allow lines to wrap beyond terminal width")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print request and response headers to stderr")
	rootCmd.PersistentFlags().BoolVar(&curl, "curl", false, "print an equivalent curl command for each request, including the bearer token (cannot be combined with --debug)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "add curl's own -v flag to the command printed by --curl (requires --curl)")
	rootCmd.PersistentFlags().BoolVarP(&location, "location", "L", false, "add curl's own -L flag to the command printed by --curl, to follow redirects (requires --curl)")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "skip TLS certificate verification")

	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(roleCmd)
	rootCmd.AddCommand(orgCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(loginCmd)
}
