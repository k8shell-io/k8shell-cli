// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	k8shell "github.com/k8shell-io/k8shell-go"

	"github.com/k8shell-io/k8shell/internal/config"
)

// newClient builds a k8shell API client from the given context, applying the
// global debug and insecure flags together with any context-level insecure setting.
func newClient(ctx *config.Context) *k8shell.Client {
	var opts []k8shell.Option
	if debug {
		opts = append(opts, k8shell.WithDebug())
	}
	if insecure || ctx.Insecure {
		opts = append(opts, k8shell.WithInsecure())
	}
	return k8shell.New(ctx.Server, ctx.Token, opts...)
}
