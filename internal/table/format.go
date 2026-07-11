// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package table

import (
	"fmt"
	"sort"
	"strings"

	"github.com/k8shell-io/common/pkg/models"
)

// FmtBool renders a bool field as "true" or "false".
func FmtBool(v any) string {
	if b, ok := v.(bool); ok && b {
		return "true"
	}
	return "false"
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

// FmtObligations renders a map[string]string field as a comma-separated string
// of "key=value" pairs, sorted by key for stable output.
func FmtObligations(v any) string {
	m, _ := v.(map[string]string)
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%s=%s", k, m[k])
	}
	return strings.Join(parts, ",")
}
