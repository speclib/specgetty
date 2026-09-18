package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Documents are identified by where they came from, so that one rule decides
// every scroll-position reset. The separator cannot appear in a path, a change
// name or an artifact name.
const docKeySep = "\x00"

// docRegion returns the size of the scrolling region for whichever document
// pane is on screen.
//
// The width comes from panelContentWidth rather than being worked out here, so
// that it cannot drift from the width the panel is actually drawn at.
func (m model) docRegion() (width, height int) {
	// Inside the tab's own border, not the panel's: the document is drawn in
	// the box, so it has to be wrapped to the box.
	width = m.contentBoxWidth()
	panelH := m.mainPanelHeight()

	switch {
	case m.level == levelChange:
		// The change name line and the sub-tab row stay put above the box.
		height = panelH - 2 - boxRows

	case m.detailTab == tabSpecs:
		// The specs tab splits the region below the tab bar into two boxes, so
		// the split is over the whole region and the document gets what is
		// inside the right-hand one.
		_, contentOuter := specsSplit(m.panelContentWidth())
		width = contentOuter - boxChrome
		height = panelH - 5 - boxRows

	default:
		// The project header takes four rows and the tab bar one, then the
		// config pane's source line and the blank line under it, all above the
		// box.
		height = panelH - 5 - 2 - boxRows
	}

	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return width, height
}

// docActive reports whether a document viewer currently owns the vertical axis.
//
// An overlay that takes the keyboard takes it from the document too, which is
// why the picker, the startup prompt and the change list search prompt all
// answer no here.
func (m model) docActive() bool {
	if len(m.repoPaths) == 0 || m.cursor >= len(m.repoPaths) {
		return false
	}
	if m.pickerOpen || m.askOpenPicker || m.searchFocused {
		return false
	}
	if m.focus == focusLog {
		return false
	}
	if m.level == levelChange {
		return true
	}
	if m.level != levelProject {
		return false
	}

	info := m.projects[m.repoPaths[m.cursor]].Info
	switch m.detailTab {
	case tabConfig:
		return len(configPanes(info)) > 0
	case tabSpecs:
		// The specs tab has two halves. The vertical keys belong to the content
		// only while the content holds the keyboard.
		return m.focus == focusSpecsContent &&
			m.specCursor < len(info.SpecNames)
	}
	return false
}

// renderChangeArtifact renders the scrolling part of an open change: the
// selected artifact, or the change's specs.
//
// For the tasks artifact it also reports where each source line landed and the
// file those lines came from, which is what the cursor highlights and what a
// toggle writes back to. Every other artifact returns neither, and so has no
// cursor.
func renderChangeArtifact(projectPath string, r changeRow, artifactTab, width int) (string, []sourceLine, string) {
	if artifactTab >= len(r.ci.ArtifactFiles) {
		return renderChangeSpecs(r, width), nil, ""
	}

	filename := r.ci.ArtifactFiles[artifactTab]
	var content strings.Builder
	prefixRows := 0

	if filename == "tasks.md" && r.ci.TasksTotal > 0 {
		// Part of the document, so it scrolls away with the tasks it counts.
		content.WriteString(sectionHeaderStyle.Render(
			fmt.Sprintf("Tasks: %d/%d complete", r.ci.TasksDone, r.ci.TasksTotal)))
		content.WriteString("\n\n")
		prefixRows = 2
	}

	body, lines := renderMarkdownLines(r.ci.ArtifactContents[filename], width)
	content.WriteString(body)

	if filename != "tasks.md" {
		return content.String(), nil, ""
	}

	// The prefix occupies rows the source knows nothing about, so every mapped
	// range shifts down by it.
	for i := range lines {
		lines[i].rowStart += prefixRows
		lines[i].rowEnd += prefixRows
	}

	dir := "changes"
	name := r.ci.DirName
	if r.archived {
		dir = filepath.Join("changes", "archive")
	}
	path := filepath.Join(projectPath, "openspec", dir, name, filename)

	return content.String(), lines, path
}

// renderChangeSpecs renders the specs sub-tab of an open change.
func renderChangeSpecs(r changeRow, width int) string {
	if len(r.ci.SpecNames) == 0 {
		return dimStyle.Render("No specs in this change")
	}
	var content strings.Builder
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
	return content.String()
}

// document is what the viewer needs to show one file.
//
// lines and path are set only for a document with a cursor, which today means a
// change's tasks artifact. Everything else leaves them empty and behaves as it
// always has.
type document struct {
	key     string
	content string
	lines   []sourceLine // source line to screen row mapping, for the cursor
	path    string       // the file on disk, for a toggle to write back
}

// cursored reports whether this document supports a line cursor.
func (d document) cursored() bool { return len(d.lines) > 0 && d.path != "" }

// currentDocument returns the identity and the rendered rows of the document
// that should be on screen, if there is one.
func (m model) currentDocument() (key, content string, ok bool) {
	d, found := m.currentDoc()
	return d.key, d.content, found
}

func (m model) currentDoc() (document, bool) {
	if len(m.repoPaths) == 0 || m.cursor >= len(m.repoPaths) {
		return document{}, false
	}
	project := m.repoPaths[m.cursor]
	width, _ := m.docRegion()

	switch {
	case m.level == levelChange:
		r, found := m.selectedRow()
		if !found {
			return document{}, false
		}
		names := r.artifactTabNames()
		if len(names) == 0 {
			return document{}, false
		}
		tab := m.changeArtifactTab
		if tab >= len(names) {
			tab = len(names) - 1
		}
		d := document{key: strings.Join([]string{project, "change", r.key(), names[tab]}, docKeySep)}
		d.content, d.lines, d.path = renderChangeArtifact(project, r, tab, width)
		return d, true

	case m.level == levelProject && m.detailTab == tabSpecs:
		info := m.projects[project].Info
		if m.specCursor >= len(info.SpecNames) {
			return document{}, false
		}
		name := info.SpecNames[m.specCursor]
		// Keyed by spec name, so moving the cursor is moving to a different
		// document and starts at the top, by the same rule as artifact sub-tabs.
		key := strings.Join([]string{project, "spec", name}, docKeySep)
		content, found := info.SpecContents[name]
		if !found || content == "" {
			return document{key: key, content: dimStyle.Render("No spec.md found")}, true
		}
		return document{key: key, content: renderMarkdown(content, width)}, true

	case m.level == levelProject && m.detailTab == tabConfig:
		panes := m.currentConfigPanes()
		if len(panes) == 0 {
			return document{}, false
		}
		i := m.configPaneIndex(panes)
		pane := panes[i]
		// Keyed by the pane, so selecting a different sub-tab is moving to a
		// different document and starts at the top, by the same rule the spec
		// list and the artifact sub-tabs already follow.
		key := strings.Join([]string{project, "config", pane.label, pane.source}, docKeySep)
		if pane.md {
			return document{key: key, content: renderMarkdown(pane.content, width)}, true
		}
		return document{key: key, content: renderYAML(pane.content, width)}, true
	}

	return document{}, false
}

// syncDocument loads the document that should be on screen into the viewport.
//
// One rule covers every reset case: a different document starts at the top, the
// same document keeps its position. That is what makes the filesystem watcher
// bearable while reading, since a rewritten file is still the same document and
// the position stays where the eye is.
func (m *model) syncDocument() {
	d, ok := m.currentDoc()
	if !ok {
		m.docKey = ""
		m.docLines = nil
		m.docPath = ""
		return
	}

	width, height := m.docRegion()
	m.docViewport.SetWidth(width)
	m.docViewport.SetHeight(height)

	fresh := d.key != m.docKey
	m.docLines = d.lines
	m.docPath = d.path

	if fresh {
		// A different document starts at the top, cursor included.
		m.docCursor = 0
	}
	m.clampDocCursor()

	m.docViewport.SetContent(m.highlightCursorLine(d))

	if fresh {
		m.docKey = d.key
		m.docViewport.GotoTop()
		m.scrollCursorIntoView()
		return
	}
	// Same document, possibly shorter than it was. SetYOffset clamps.
	m.docViewport.SetYOffset(m.docViewport.YOffset())
}

// clampDocCursor keeps the cursor inside the document.
func (m *model) clampDocCursor() {
	if len(m.docLines) == 0 {
		m.docCursor = 0
		return
	}
	if m.docCursor < 0 {
		m.docCursor = 0
	}
	if m.docCursor >= len(m.docLines) {
		m.docCursor = len(m.docLines) - 1
	}
}

// selectedSourceLine returns the line under the cursor, if the document has one.
func (m model) selectedSourceLine() (sourceLine, bool) {
	if len(m.docLines) == 0 || m.docCursor < 0 || m.docCursor >= len(m.docLines) {
		return sourceLine{}, false
	}
	return m.docLines[m.docCursor], true
}

// highlightCursorLine marks every row the selected source line produced.
//
// The rows are already styled, and layering a background over them would leave
// the inner colours showing through in patches. Stripping first and restyling
// gives one flat highlighted band, which is how the change list draws its
// selected row too.
func (m model) highlightCursorLine(d document) string {
	if !d.cursored() {
		return d.content
	}
	sel, ok := m.selectedSourceLine()
	if !ok {
		return d.content
	}

	rows := strings.Split(d.content, "\n")
	width, _ := m.docRegion()
	for i := sel.rowStart; i <= sel.rowEnd && i < len(rows); i++ {
		if i < 0 {
			continue
		}
		rows[i] = selectedStyle.Width(width).Render(ansi.Strip(rows[i]))
	}
	return strings.Join(rows, "\n")
}

// scrollCursorIntoView moves the viewport so the whole selected line is shown.
//
// A source line can be several rows tall, so bringing the first row into view
// is not enough: the last row has to fit as well, and when the line is taller
// than the pane the top wins.
func (m *model) scrollCursorIntoView() {
	sel, ok := m.selectedSourceLine()
	if !ok {
		return
	}
	top := m.docViewport.YOffset()
	height := m.docViewport.Height()
	if height < 1 {
		return
	}

	if sel.rowEnd >= top+height {
		top = sel.rowEnd - height + 1
	}
	if sel.rowStart < top {
		top = sel.rowStart
	}
	m.docViewport.SetYOffset(top)
}

// moveDocCursor steps the cursor and brings it into view.
func (m *model) moveDocCursor(delta int) {
	if len(m.docLines) == 0 {
		return
	}
	m.docCursor += delta
	m.clampDocCursor()
	m.scrollCursorIntoView()
}

// docHasCursor reports whether the document on screen is one the cursor drives.
func (m model) docHasCursor() bool {
	return m.docActive() && len(m.docLines) > 0 && m.docPath != ""
}

// docScrollPercent reports how far down the document the viewport sits, or -1
// when the whole document fits and there is nothing to report.
func (m model) docScrollPercent() int {
	if !m.docActive() {
		return -1
	}
	if m.docViewport.TotalLineCount() <= m.docViewport.Height() {
		return -1
	}
	p := int(m.docViewport.ScrollPercent() * 100)
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}
	return p
}
