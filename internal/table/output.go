// Copyright 2026 The k8shell Authors.
// SPDX-License-Identifier: AGPL-3.0-only

package table

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"golang.org/x/term"
)

var (
	Bold   = color.New(color.Bold).SprintFunc()
	Active = color.New(color.FgGreen, color.Bold).SprintFunc()
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// visibleLen returns the number of visible runes in s, excluding ANSI escape sequences.
func visibleLen(s string) int {
	return utf8.RuneCountInString(ansiEscape.ReplaceAllString(s, ""))
}

// truncate shortens s to at most max visible runes, appending "…" if truncation occurs.
func truncate(s string, max int) string {
	if max <= 0 || visibleLen(s) <= max {
		return s
	}
	runes := []rune(ansiEscape.ReplaceAllString(s, ""))
	if max <= 1 {
		return string(runes[:max])
	}
	return string(runes[:max-1]) + "…"
}

// truncateLine hard-truncates s to maxWidth runes after stripping ANSI codes, with no ellipsis.
func truncateLine(s string, maxWidth int) string {
	plain := ansiEscape.ReplaceAllString(s, "")
	if utf8.RuneCountInString(plain) <= maxWidth {
		return s
	}
	return string([]rune(plain)[:maxWidth])
}

// termWidth returns the current terminal width, or 0 if it cannot be determined.
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 0
	}
	return w
}

// Column is the internal row/width representation used by the table renderer.
type Column struct {
	Header   string
	MaxWidth int
}

// Col defines a typed table column. Exactly one of Field or Fn must be set.
//
//   - Field: JSON or YAML tag name resolved via reflection; Fmt is applied to
//     the raw value if provided, otherwise fmt.Sprint is used.
//   - Fn: used for computed columns whose value cannot be derived from a single
//     struct field (e.g. status derived from multiple fields, or values that
//     depend on external state).
type Col[T any] struct {
	Header   string
	MaxWidth int
	Help     string           // short description shown in --help output
	Field    string           // json/yaml tag name
	Fmt      func(any) string // optional; applied to the raw field value
	Fn       func(T) string   // for computed columns; takes precedence over Field
}

// ColumnHelp returns a "Columns:\n  HEADER  description\n..." block built from
// the Help strings of cols. Columns with an empty Help are omitted.
func ColumnHelp[T any](cols []Col[T]) string {
	maxLen := 0
	for _, c := range cols {
		if c.Help != "" && len(c.Header) > maxLen {
			maxLen = len(c.Header)
		}
	}
	var sb strings.Builder
	sb.WriteString("Columns:\n")
	for _, c := range cols {
		if c.Help == "" {
			continue
		}
		fmt.Fprintf(&sb, "  %-*s  %s\n", maxLen, c.Header, c.Help)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// SortSpec holds a single sort criterion parsed from the --sort flag.
type SortSpec struct {
	Field string
	Desc  bool
}

// ParseSort parses a comma-separated sort string into SortSpecs.
// Prefix a field name with - for descending order, e.g. "username,-email".
func ParseSort(s string) ([]SortSpec, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	specs := make([]SortSpec, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		spec := SortSpec{Field: p}
		if strings.HasPrefix(p, "-") {
			spec.Field = p[1:]
			spec.Desc = true
		}
		if spec.Field == "" {
			return nil, fmt.Errorf("invalid sort field %q", p)
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

// validateSortFields checks that every sort spec refers to a column Field.
// Fn-only columns are excluded since they have no single backing struct field.
func validateSortFields[T any](specs []SortSpec, cols []Col[T]) error {
	valid := make(map[string]struct{}, len(cols))
	for _, c := range cols {
		if c.Field != "" {
			valid[c.Field] = struct{}{}
		}
	}
	var unknown []string
	for _, s := range specs {
		if _, ok := valid[s.Field]; !ok {
			unknown = append(unknown, s.Field)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	keys := make([]string, 0, len(valid))
	for k := range valid {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return fmt.Errorf("unknown sort field(s) %v; valid fields: %s", unknown, strings.Join(keys, ", "))
}

// resolveField returns the value of the struct field whose json or yaml tag
// matches tag, searching in priority order: json → yaml.
func resolveField(v reflect.Value, tag string) (any, bool) {
	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		for _, key := range []string{"json", "yaml"} {
			name := strings.SplitN(f.Tag.Get(key), ",", 2)[0]
			if name == tag {
				return v.Field(i).Interface(), true
			}
		}
	}
	return nil, false
}

// compareAny compares two values of the same underlying type.
// Falls back to string comparison for unknown types.
func compareAny(a, b any, desc bool) int {
	var result int
	switch av := a.(type) {
	case string:
		result = cmp.Compare(av, b.(string))
	case int:
		result = cmp.Compare(av, b.(int))
	case int32:
		result = cmp.Compare(av, b.(int32))
	case int64:
		result = cmp.Compare(av, b.(int64))
	case uint32:
		result = cmp.Compare(av, b.(uint32))
	case uint64:
		result = cmp.Compare(av, b.(uint64))
	case float64:
		result = cmp.Compare(av, b.(float64))
	case bool:
		// false < true
		result = cmp.Compare(boolInt(av), boolInt(b.(bool)))
	case *time.Time:
		bv, _ := b.(*time.Time)
		result = compareTimes(av, bv)
	case time.Time:
		bv, _ := b.(time.Time)
		result = av.Compare(bv)
	default:
		result = cmp.Compare(fmt.Sprint(a), fmt.Sprint(b))
	}
	if desc {
		return -result
	}
	return result
}

// boolInt converts a bool to 0 or 1 for comparison purposes.
func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// compareTimes compares two *time.Time values, treating nil as less than any non-nil value.
func compareTimes(a, b *time.Time) int {
	switch {
	case a == nil && b == nil:
		return 0
	case a == nil:
		return -1
	case b == nil:
		return 1
	default:
		return a.Compare(*b)
	}
}

// sortItems returns a sorted copy of items according to specs.
func sortItems[T any](items []T, specs []SortSpec) []T {
	if len(specs) == 0 {
		return items
	}
	sorted := slices.Clone(items)
	slices.SortStableFunc(sorted, func(a, b T) int {
		ra, rb := reflect.ValueOf(a), reflect.ValueOf(b)
		for _, spec := range specs {
			av, _ := resolveField(ra, spec.Field)
			bv, _ := resolveField(rb, spec.Field)
			if c := compareAny(av, bv, spec.Desc); c != 0 {
				return c
			}
		}
		return 0
	})
	return sorted
}

// Table renders a typed slice as a formatted table. For each column, Fn takes
// precedence; otherwise the field named by Field is resolved via reflection.
// sort is the raw --sort flag value (e.g. "username,-email"); ignored when empty.
// Only fields declared in the column definitions (via Field) are valid sort keys.
func Table[T any](p *Printer, cols []Col[T], items []T, sort string) error {
	if sort != "" {
		specs, err := ParseSort(sort)
		if err != nil {
			return err
		}
		if err := validateSortFields(specs, cols); err != nil {
			return err
		}
		items = sortItems(items, specs)
	}

	baseCols := make([]Column, len(cols))
	for i, c := range cols {
		baseCols[i] = Column{Header: c.Header, MaxWidth: c.MaxWidth}
	}

	rows := make([][]string, len(items))
	for i, item := range items {
		rv := reflect.ValueOf(item)
		row := make([]string, len(cols))
		for j, col := range cols {
			if col.Fn != nil {
				row[j] = col.Fn(item)
				continue
			}
			raw, _ := resolveField(rv, col.Field)
			if col.Fmt != nil {
				row[j] = col.Fmt(raw)
			} else {
				row[j] = fmt.Sprint(raw)
			}
		}
		rows[i] = row
	}
	p.table(baseCols, rows)
	return nil
}

// Printer renders output as either a formatted table or indented JSON.
type Printer struct {
	jsonMode bool
	wrap     bool
}

// New creates a Printer. When noANSI is true, color output is globally disabled.
func New(jsonMode, noANSI, wrap bool) *Printer {
	color.NoColor = noANSI
	return &Printer{jsonMode: jsonMode, wrap: wrap}
}

// table writes a header row followed by data rows to stdout, padding columns to a consistent width.
// Lines are clipped to the terminal width unless wrap is enabled.
func (p *Printer) table(cols []Column, rows [][]string) {
	lineMax := 0
	if !p.wrap {
		lineMax = termWidth()
	}

	ncols := len(cols)
	widths := make([]int, ncols)
	for i, c := range cols {
		widths[i] = len(c.Header)
	}
	for _, row := range rows {
		for i := 0; i < ncols && i < len(row); i++ {
			cell := truncate(row[i], cols[i].MaxWidth)
			if w := visibleLen(cell); w > widths[i] {
				widths[i] = w
			}
		}
	}

	const gap = "  "

	buildRow := func(cells []string, isHeader bool) string {
		var sb strings.Builder
		for i := 0; i < ncols && i < len(cells); i++ {
			cell := cells[i]
			if !isHeader {
				cell = truncate(cell, cols[i].MaxWidth)
			}
			if i > 0 {
				sb.WriteString(gap)
			}
			sb.WriteString(cell)
			if i < ncols-1 {
				if pad := widths[i] - visibleLen(cell); pad > 0 {
					sb.WriteString(strings.Repeat(" ", pad))
				}
			}
		}
		line := sb.String()
		if lineMax > 0 && visibleLen(line) > lineMax {
			line = truncateLine(line, lineMax)
		}
		return line
	}

	headers := make([]string, ncols)
	for i, c := range cols {
		headers[i] = c.Header
	}
	fmt.Fprintln(os.Stdout, buildRow(headers, true))
	for _, row := range rows {
		fmt.Fprintln(os.Stdout, buildRow(row, false))
	}
}

// JSON encodes v as indented JSON and writes it to stdout.
func (p *Printer) JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// IsJSON reports whether the printer is in JSON output mode.
func (p *Printer) IsJSON() bool {
	return p.jsonMode
}
