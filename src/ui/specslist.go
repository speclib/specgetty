package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

const (
	specLineCapability = iota
	specLineHeader
	specLineSpacer
)

type specLine struct {
	kind      int
	label     string
	specIndex int
}

func splitSpecName(name string) (group, rest string) {
	if i := strings.IndexByte(name, '/'); i >= 0 {
		return name[:i], name[i+1:]
	}
	return "", name
}

func specLines(names []string) []specLine {
	held := make(map[string]int, len(names))
	for _, name := range names {
		if group, _ := splitSpecName(name); group != "" {
			held[group]++
		}
	}

	var lines []specLine
	group := ""
	for i, name := range names {
		g, rest := splitSpecName(name)
		if g != "" && g != group {
			if len(lines) > 0 {
				lines = append(lines, specLine{kind: specLineSpacer, specIndex: -1})
			}
			lines = append(lines, specLine{
				kind:      specLineHeader,
				label:     fmt.Sprintf("%s (%d)", g, held[g]),
				specIndex: -1,
			})
			group = g
		}
		lines = append(lines, specLine{
			kind:      specLineCapability,
			label:     rest,
			specIndex: i,
		})
	}
	return lines
}

func ownersOfSpecs(lines []specLine) itemLines {
	owners := make(itemLines, len(lines))
	for i, l := range lines {
		owners[i] = lineOwner
		if l.kind == specLineCapability {
			owners[i] = l.specIndex
		}
	}
	return owners
}

func renderSpecList(lines []specLine, cursor, width, height int, dimmed bool) string {
	width = max(1, width)

	var b strings.Builder
	offset := ownersOfSpecs(lines).offsetFor(cursor, height)
	if offset < len(lines) && lines[offset].kind == specLineSpacer {
		offset++
	}
	end := offset + height
	if end > len(lines) {
		end = len(lines)
	}

	for i := offset; i < end; i++ {
		if i > offset {
			b.WriteString("\n")
		}
		l := lines[i]
		if l.kind == specLineSpacer {
			continue
		}
		label := ansi.Truncate(l.label, width, "")
		if l.kind == specLineHeader {
			b.WriteString(sectionHeaderStyle.Render(label))
			continue
		}
		switch {
		case l.specIndex != cursor:
			b.WriteString(normalStyle.Render(label))
		case dimmed:
			b.WriteString(dimSelectedStyle.Width(width).Render(label))
		default:
			b.WriteString(selectedStyle.Width(width).Render(label))
		}
	}
	return b.String()
}
