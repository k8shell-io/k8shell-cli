// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package table

import (
	"strings"

	"github.com/k8shell-io/common/pkg/models"
)

// FmtBool renders a bool field as "yes" or "no".
func FmtBool(v any) string {
	if b, ok := v.(bool); ok && b {
		return "yes"
	}
	return "no"
}

// FmtJoin renders a []string field as a comma-separated string.
func FmtJoin(v any) string {
	ss, _ := v.([]string)
	return strings.Join(ss, ",")
}

// FmtRoles renders a []models.Role field as a comma-separated string of role names.
func FmtRoles(v any) string {
	roles, _ := v.([]models.Role)
	s := make([]string, len(roles))
	for i, r := range roles {
		s[i] = string(r)
	}
	return strings.Join(s, ",")
}
