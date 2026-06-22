// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/k8shell-io/k8shell/internal/config"
	"github.com/k8shell-io/k8shell/internal/table"
	"github.com/spf13/cobra"
)

var contextSortFlag string

var contextListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all contexts",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(cfg.Contexts) == 0 {
			fmt.Println("No contexts configured.")
			return nil
		}

		if printer.IsJSON() {
			return printer.JSON(cfg.Contexts)
		}

		cols := []table.Col[config.Context]{
			{Header: "", MaxWidth: 0, Fn: func(ctx config.Context) string {
				if ctx.Name == cfg.CurrentContext {
					return table.Active("*")
				}
				return " "
			}},
			{Header: "NAME", MaxWidth: 20, Fn: func(ctx config.Context) string {
				if ctx.Name == cfg.CurrentContext {
					return table.Active(ctx.Name)
				}
				return ctx.Name
			}},
			{Header: "SERVER", MaxWidth: 50, Field: "server"},
		}

		return table.Table(printer, cols, cfg.Contexts, contextSortFlag)
	},
}

func init() {
	contextListCmd.Flags().StringVar(&contextSortFlag, "sort", "", "sort by fields, e.g. name,-server (prefix - for descending)")
}
