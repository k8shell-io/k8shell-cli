// Copyright 2026 The k8shell CLI Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package cmd

import (
	"strings"

	"github.com/k8shell-io/common/pkg/models"
)

// fmtBool renders a bool field as "yes" or "no".
func fmtBool(v any) string {
	if b, ok := v.(bool); ok && b {
		return "yes"
	}
	return "no"
}

// fmtJoin renders a []string field as a comma-separated string.
func fmtJoin(v any) string {
	ss, _ := v.([]string)
	return strings.Join(ss, ",")
}

// fmtRoles renders a []models.Role field as a comma-separated string of role names.
func fmtRoles(v any) string {
	roles, _ := v.([]models.Role)
	s := make([]string, len(roles))
	for i, r := range roles {
		s[i] = string(r)
	}
	return strings.Join(s, ",")
}
