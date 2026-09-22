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
	case m.specLevel() && !m.specStructured():
		// A report fills the panel on its own: there is no outline beside it.
		width = m.panelContentWidth() - boxChrome
		height = panelH - 1 - boxRows

	case m.specLevel():
		// The name stays put above the two halves, and the card gets what is
		// inside the right-hand border. A change's card gives one more row to
		// the chooser above it.
		_, cardOuter := specDetailSplit(m.panelContentWidth())
		width = cardOuter - boxChrome
		height = panelH - 1 - boxRows
		if m.cardViewRowShown() {
			height--
		}

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

	case m.detailTab == tabProperties:
		// The same split, sized to its own labels rather than to a share.
		_, contentOuter := m.propertiesSplit(m.panelContentWidth())
		width = contentOuter - boxChrome
		height = panelH - 5 - boxRows

	default:
		// The project header takes four rows and the tab bar one, all above
		// the box.
		height = panelH - 5 - boxRows
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
	if m.level == levelChange {
		return true
	}
	if m.specLevel() {
		if !m.specStructured() {
			// One panel, so the vertical keys scroll the report, the way they
			// scroll an open change.
			return len(m.specProblems) > 0
		}
		// Two halves, so the vertical keys belong to the card only while the
		// card holds the keyboard.
		return m.focus == focusContentPane
	}
	if m.level != levelProject {
		return false
	}

	info := m.projects[m.repoPaths[m.cursor]].Info
	switch m.detailTab {
	case tabProperties:
		// The properties tab has two halves. The vertical keys belong to the
		// content only while the content holds the keyboard, exactly as on the
		// specs tab.
		return m.focus == focusContentPane && len(propSections(info)) > 0
	case tabSpecs:
		// The specs tab has two halves. The vertical keys belong to the content
		// only while the content holds the keyboard.
		return m.focus == focusContentPane &&
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

	// The file this artifact came from. Every artifact has one, whether or not
	// anything writes to it: the path says where the content on screen lives,
	// which is what an editor needs and what a toggle happens to need too.
	dir := "changes"
	if r.archived {
		dir = filepath.Join("changes", "archive")
	}
	path := filepath.Join(projectPath, "openspec", dir, r.ci.DirName, filename)

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
		// Only tasks.md maps its rows back to source lines, which is what gives
		// a pane its cursor. The others are read, not edited in place.
		return content.String(), nil, path
	}

	// The prefix occupies rows the source knows nothing about, so every mapped
	// range shifts down by it.
	for i := range lines {
		lines[i].rowStart += prefixRows
		lines[i].rowEnd += prefixRows
	}

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
	case m.specLevel() && !m.specStructured():
		if len(m.specProblems) == 0 {
			return document{}, false
		}
		// The report scrolls like any other document, so a file with many
		// faults is reachable rather than clipped.
		key := strings.Join([]string{project, "spec-report", m.specName}, docKeySep)
		// The path is the file the report is about, so `E` opens it from here
		// and the nav bar offers the key, both by the rule every other pane
		// already follows: a document that came from one file names it.
		return document{key: key,
			content: renderSpecReport(m.specName, m.specProblems, width),
			path:    filepath.Join(m.currentRoot(), "openspec", "specs", m.specName, "spec.md"),
		}, true

	case m.level == levelSpec:
		n := m.specTree.nodes[m.selectedSpecNode()]
		// Keyed by the node, so moving to another one is moving to another
		// document and starts at the top.
		key := strings.Join([]string{project, "spec", m.specTree.name, n.path}, docKeySep)
		return document{key: key, content: renderSpecCard(n, width)}, true

	case m.level == levelChangeSpec:
		n := m.specTree.nodes[m.selectedSpecNode()]
		// Keyed by the node and the view, so moving between the difference, the
		// original and the new text starts each at its own top.
		key := strings.Join([]string{project, "change-spec", m.specTree.name,
			n.path, cardViewNames[m.cardView]}, docKeySep)
		// The delta this node came from, which is the one file the pane is
		// showing. The specs sub-tab one level up shows several and so offers
		// no key.
		path := ""
		if r, ok := m.selectedRow(); ok && n.capability != "" {
			dir := "changes"
			if r.archived {
				dir = filepath.Join("changes", "archive")
			}
			path = filepath.Join(m.currentRoot(), "openspec", dir, r.ci.DirName,
				"specs", n.capability, "spec.md")
		}
		return document{key: key, path: path,
			content: renderChangeCard(n, m.cardView, width)}, true

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
		// The resolved root, not the directory the user started in: in a
		// store-backed project the change is in the store, and a path built
		// from the starting directory points at nothing.
		d.content, d.lines, d.path = renderChangeArtifact(m.currentRoot(), r, tab, width)
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
		path := filepath.Join(m.currentRoot(), "openspec", "specs", name, "spec.md")
		content, found := info.SpecContents[name]
		if !found || content == "" {
			// No path: there is no file to hand to an editor, which is what the
			// message says.
			return document{key: key, content: dimStyle.Render("No spec.md found")}, true
		}
		return document{key: key, content: renderMarkdown(content, width), path: path}, true

	case m.level == levelProject && m.detailTab == tabProperties:
		sections := m.currentSections()
		if len(sections) == 0 {
			return document{}, false
		}
		info := m.projects[project].Info
		sec := sections[m.sectionIndex(sections)]
		// Keyed by the row, so selecting a different one is moving to a
		// different document and starts at the top, by the same rule the spec
		// list and the artifact sub-tabs already follow.
		key := strings.Join([]string{project, "properties", sec.label}, docKeySep)

		switch sec.kind {
		case sectionConfig:
			if info.ConfigFile == "" {
				return document{key: key, content: dimStyle.Render("No project configuration found")}, true
			}
			// Where the configuration came from is part of what this row
			// reports, not a label on it: for a store-backed project the file
			// is the store's, and nothing else on screen would say so. The
			// list beside it is the chrome that names the row.
			head := dimStyle.Render("# "+sec.source) + "\n\n"
			path := filepath.Join(m.currentRoot(), "openspec", info.ConfigFile)
			if sec.md {
				return document{key: key, content: head + renderMarkdown(info.ConfigContent, width), path: path}, true
			}
			return document{key: key, content: head + renderYAML(info.ConfigContent, width), path: path}, true
		case sectionSchema:
			body := renderSchemaSection(info, sec.schema, m.schemaStateOf(sec.schema))
			return document{key: key, content: renderYAML(body, width)}, true
		default:
			return document{key: key, content: renderYAML(renderStoreSection(info), width)}, true
		}
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
		m.docTasks = taskItems{}
		m.docPath = ""
		return
	}

	width, height := m.docRegion()
	m.docViewport.SetWidth(width)
	m.docViewport.SetHeight(height)

	fresh := d.key != m.docKey
	m.docPath = d.path
	m.docTasks = taskItemsIn(d.lines, strings.Count(d.content, "\n")+1)

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
	m.docCursor = clampIndex(m.docCursor, m.docTasks.count())
}

// selectedTask returns the task under the cursor, if the document has one.
func (m model) selectedTask() (taskItem, bool) {
	return m.docTasks.at(m.docCursor)
}

// selectedTaskRows returns the first and last screen row of the selected task,
// which is every row of its checkbox line and of its continuation lines.
func (m model) selectedTaskRows() (first, last int, ok bool) {
	if m.docTasks.count() == 0 {
		return 0, 0, false
	}
	first, last = m.docTasks.rows.span(m.docCursor)
	return first, last, first >= 0
}

// highlightCursorLine marks every row the selected task produced.
//
// Every row, not the checkbox line's rows: a task wraps onto continuation lines
// and a band that covered the first of them and stopped read as a task half
// selected, which is what the span from itemLines fixes.
//
// The rows are already styled, and layering a background over them would leave
// the inner colours showing through in patches. Stripping first and restyling
// gives one flat highlighted band, which is how the change list draws its
// selected row too.
func (m model) highlightCursorLine(d document) string {
	if !d.cursored() {
		return d.content
	}
	first, last, ok := m.selectedTaskRows()
	if !ok {
		return d.content
	}

	rows := strings.Split(d.content, "\n")
	width, _ := m.docRegion()
	for i := first; i <= last && i < len(rows); i++ {
		if i < 0 {
			continue
		}
		rows[i] = selectedStyle.Width(width).Render(ansi.Strip(rows[i]))
	}
	return strings.Join(rows, "\n")
}

// scrollCursorIntoView moves the viewport so the whole selected task is shown.
//
// A task is several rows tall, so bringing the first row into view is not
// enough: the last row has to fit as well, and when the task is taller than the
// pane the top wins, a task scrolled to its last row reading as a fragment with
// no beginning.
//
// The scroll is the smallest one that shows the task, rather than itemLines'
// offsetFor, which answers a different question: where the pane would sit if it
// were being laid out from scratch. The outline can ask that because it is laid
// out from scratch every frame. A document holds a reading position across
// rescans and resizes, and moving it further than the cursor needed would throw
// that position away on every keystroke.
func (m *model) scrollCursorIntoView() {
	first, last, ok := m.selectedTaskRows()
	if !ok {
		return
	}
	top := m.docViewport.YOffset()
	height := m.docViewport.Height()
	if height < 1 {
		return
	}

	if last >= top+height {
		top = last - height + 1
	}
	if first < top {
		top = first
	}
	m.docViewport.SetYOffset(top)
}

// moveDocCursor steps the cursor by delta tasks and brings it into view.
func (m *model) moveDocCursor(delta int) {
	if m.docTasks.count() == 0 {
		return
	}
	m.docCursor += delta
	m.clampDocCursor()
	m.scrollCursorIntoView()
}

// gotoDocEnd sends the cursor to the first or last task.
func (m *model) gotoDocEnd(last bool) {
	n := m.docTasks.count()
	if n == 0 {
		return
	}
	to := 0
	if last {
		to = n - 1
	}
	m.docCursor = clampIndex(to, n)
	m.scrollCursorIntoView()
}

// docHasCursor reports whether the document on screen is one the cursor drives.
//
// A tasks file with no checkbox at column zero has no tasks to select, so it
// answers no and is scrolled by rows like any other document. That falls out of
// the count rather than being a case of its own.
func (m model) docHasCursor() bool {
	return m.docActive() && m.docTasks.count() > 0 && m.docPath != ""
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
