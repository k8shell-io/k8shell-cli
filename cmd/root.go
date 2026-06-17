package cmd

import (
	"fmt"
	"os"

	"github.com/k8shell-io/k8shell/internal/config"
	"github.com/k8shell-io/k8shell/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	contextName string
	jsonMode    bool
	noANSI      bool
	wrap        bool
	debug       bool
	cfg         *config.Config
	printer     *output.Printer
)

var rootCmd = &cobra.Command{
	Use:          "k8shell",
	Short:        "k8shell — CLI for k8shell resources",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		printer = output.New(jsonMode, noANSI, wrap)

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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
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

	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(contextCmd)
}
