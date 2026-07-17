// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"strings"

	"github.com/spf13/pflag"
)

// repeatableStringValue implements pflag.Value with the same behavior as
// pflag's built-in stringArray type (each flag occurrence is appended verbatim,
// without splitting on commas), but reports its type as "string" rather than
// "stringArray" in --help output, since users pass one plain string value per
// occurrence rather than an array literal.
type repeatableStringValue struct {
	value   *[]string
	changed bool
}

func (v *repeatableStringValue) String() string {
	if v.value == nil {
		return ""
	}
	return strings.Join(*v.value, ",")
}

func (v *repeatableStringValue) Set(s string) error {
	if !v.changed {
		*v.value = []string{s}
		v.changed = true
	} else {
		*v.value = append(*v.value, s)
	}
	return nil
}

func (v *repeatableStringValue) Type() string {
	return "string"
}

// repeatableStringVar registers a flag on fs that can be passed multiple
// times to build up p, in usage and behavior identical to StringArrayVar.
func repeatableStringVar(fs *pflag.FlagSet, p *[]string, name string, usage string) {
	*p = nil
	fs.Var(&repeatableStringValue{value: p}, name, usage)
}
