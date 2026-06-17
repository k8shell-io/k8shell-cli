package output

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"golang.org/x/term"
)

var (
	Bold   = color.New(color.Bold).SprintFunc()
	Active = color.New(color.FgGreen, color.Bold).SprintFunc()
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleLen(s string) int {
	return utf8.RuneCountInString(ansiEscape.ReplaceAllString(s, ""))
}

// truncate shortens s to max visible characters, appending "…" if cut.
// If s contains ANSI codes, they are stripped before truncating.
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

// truncateLine truncates a full row string to maxWidth visible characters.
// ANSI codes are stripped when the line needs to be cut.
func truncateLine(s string, maxWidth int) string {
	plain := ansiEscape.ReplaceAllString(s, "")
	if utf8.RuneCountInString(plain) <= maxWidth {
		return s
	}
	return string([]rune(plain)[:maxWidth])
}

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 0
	}
	return w
}

// Column defines a table column header and its optional maximum display width.
type Column struct {
	Header   string
	MaxWidth int // 0 = no limit
}

type Printer struct {
	jsonMode bool
	wrap     bool
}

func New(jsonMode, noANSI, wrap bool) *Printer {
	color.NoColor = noANSI
	return &Printer{jsonMode: jsonMode, wrap: wrap}
}

func (p *Printer) Table(cols []Column, rows [][]string) {
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

func (p *Printer) JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (p *Printer) IsJSON() bool {
	return p.jsonMode
}
