package cmd

import (
	"fmt"

	"github.com/k8shell-io/k8shell/internal/output"
	"github.com/spf13/cobra"
)

var contextColumns = []output.Column{
	{Header: ""},
	{Header: "NAME", MaxWidth: 20},
	{Header: "SERVER", MaxWidth: 50},
}

var contextListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contexts",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(cfg.Contexts) == 0 {
			fmt.Println("No contexts configured.")
			return nil
		}

		if printer.IsJSON() {
			return printer.JSON(cfg.Contexts)
		}

		rows := make([][]string, len(cfg.Contexts))
		for i, ctx := range cfg.Contexts {
			marker := " "
			name := ctx.Name
			if ctx.Name == cfg.CurrentContext {
				marker = output.Active("*")
				name = output.Active(ctx.Name)
			}
			rows[i] = []string{marker, name, ctx.Server}
		}
		printer.Table(contextColumns, rows)
		return nil
	},
}
