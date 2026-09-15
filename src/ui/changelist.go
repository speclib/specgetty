package ui

import (
	"fmt"
	"strings"
)

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
