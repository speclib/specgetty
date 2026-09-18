package ui

import (
	"fmt"
	"strings"
)

// renderChangeDetail draws a single open change at full panel width: its name,
// its artifact sub-tabs, and the artifact itself as a scrolling document.
//
// The name line and the sub-tab row stay put; only the document scrolls. The
// content comes from the viewport rather than being truncated, which is what
// makes the rows below the fold reachable at all.
func (m model) renderChangeDetail(r changeRow, artifactTab, width, height int) string {
	var b strings.Builder

	state := "active"
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

	// The change name and the sub-tab row are chrome; the artifact is content.
	boxHeight := height - 2
	if boxHeight < boxRows+1 {
		boxHeight = boxRows + 1
	}
	b.WriteString(contentBox(width, boxHeight, m.focus != focusLog, m.docViewport.View()))

	return b.String()
}
