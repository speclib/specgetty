package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderChangeTable draws the change list across the full panel width: a header
// row of column names, then one row per change.
func renderChangeTable(rows []changeRow, fields []string, cursor int, width, height int) string {
	// A body search needs room to say which files matched, otherwise the
	// flexible name column absorbs the whole width and the hint never shows.
	// The column only appears while such a search is active.
	hintWidth := 0
	for _, r := range rows {
		if len(r.matchedFiles) > 0 {
			hintWidth = min(28, width/3)
			break
		}
	}
	tableWidth := width - hintWidth
	if tableWidth < 1 {
		tableWidth = width
		hintWidth = 0
	}

	kept, widths := layoutFields(fields, tableWidth)

	var b strings.Builder

	// Column header.
	headerCells := make([]string, len(kept))
	for i, f := range kept {
		headerCells[i] = fitCell(knownFields[f].header, widths[i])
	}
	header := fitCell(strings.Join(headerCells, " "), tableWidth)
	if hintWidth > 0 {
		header += fitCell("matched", hintWidth)
	}
	b.WriteString(dimStyle.Render(header))

	bodyHeight := height - 1
	if bodyHeight < 1 {
		return b.String()
	}

	// Keep the cursor on screen.
	offset := 0
	if cursor >= bodyHeight {
		offset = cursor - bodyHeight + 1
	}
	end := offset + bodyHeight
	if end > len(rows) {
		end = len(rows)
	}

	for i := offset; i < end; i++ {
		b.WriteString("\n")
		r := rows[i]

		cells := make([]string, len(kept))
		for j, f := range kept {
			cells[j] = fitCell(knownFields[f].value(r), widths[j])
		}
		base := fitCell(strings.Join(cells, " "), tableWidth)

		hint := ""
		if hintWidth > 0 {
			hint = fitCell(strings.Join(r.matchedFiles, ", "), hintWidth)
		}

		if i == cursor {
			b.WriteString(selectedStyle.Width(width).Render(fitCell(base+hint, width)))
		} else {
			b.WriteString(normalStyle.Render(base))
			if hint != "" {
				b.WriteString(dimStyle.Render(hint))
			}
		}
	}

	return b.String()
}

// renderChangeDetail draws a single open change at full panel width: its name,
// its artifact sub-tabs, and the selected artifact's content.
func renderChangeDetail(r changeRow, artifactTab int, width, height int) string {
	var b strings.Builder

	state := "open"
	if r.archived {
		state = "archived"
	}
	b.WriteString(headerStyle.Render(r.ci.Name))
	b.WriteString(dimStyle.Render("  (" + state + ")"))
	if r.ci.TasksTotal > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  tasks %d/%d", r.ci.TasksDone, r.ci.TasksTotal)))
	}
	b.WriteString("\n")

	names := r.artifactTabNames()
	for i, name := range names {
		if i > 0 {
			b.WriteString(" ")
		}
		if i == artifactTab {
			b.WriteString(activeTabStyle.Render(name))
		} else {
			b.WriteString(inactiveTabStyle.Render(name))
		}
	}
	b.WriteString("\n")

	contentHeight := height - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	var content strings.Builder
	if artifactTab < len(r.ci.ArtifactFiles) {
		filename := r.ci.ArtifactFiles[artifactTab]
		if filename == "tasks.md" && r.ci.TasksTotal > 0 {
			content.WriteString(sectionHeaderStyle.Render(
				fmt.Sprintf("Tasks: %d/%d complete\n\n", r.ci.TasksDone, r.ci.TasksTotal)))
		}
		content.WriteString(renderMarkdown(r.ci.ArtifactContents[filename], width))
	} else if len(r.ci.SpecNames) > 0 {
		for i, name := range r.ci.SpecNames {
			if i > 0 {
				content.WriteString("\n")
			}
			content.WriteString(sectionHeaderStyle.Render(name))
			content.WriteString("\n")
			if c, ok := r.ci.SpecContents[name]; ok {
				content.WriteString(renderMarkdown(c, width))
			} else {
				content.WriteString(dimStyle.Render("  No spec.md found"))
			}
			content.WriteString("\n")
		}
	} else {
		content.WriteString(dimStyle.Render("No specs in this change"))
	}

	b.WriteString(truncateContent(content.String(), contentHeight))
	return b.String()
}

// renderSearchPrompt draws the filter line under the list. It stays visible for
// as long as a query is applied, so a narrowed list always shows why.
func renderSearchPrompt(input string, focused bool, shown, total int) string {
	var b strings.Builder
	b.WriteString(navBarKeyStyle.Render("/"))
	b.WriteString(navBarStyle.Render(input))
	if focused {
		b.WriteString(navBarStyle.Render("_"))
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("   %d of %d shown", shown, total)))
	return lipgloss.NewStyle().Render(b.String())
}
