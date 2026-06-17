package cmd

import (
	"strings"

	"github.com/k8shell-io/common/pkg/models"
)

func fmtBool(v any) string {
	if b, ok := v.(bool); ok && b {
		return "yes"
	}
	return "no"
}

func fmtJoin(v any) string {
	ss, _ := v.([]string)
	return strings.Join(ss, ",")
}

func fmtRoles(v any) string {
	roles, _ := v.([]models.Role)
	s := make([]string, len(roles))
	for i, r := range roles {
		s[i] = string(r)
	}
	return strings.Join(s, ",")
}
