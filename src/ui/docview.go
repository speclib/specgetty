package ui

import (
	"fmt"
	"strings"
)

// Documents are identified by where they came from, so that one rule decides
// every scroll-position reset. The separator cannot appear in a path, a change
// name or an artifact name.
const docKeySep = "\x00"

// docRegion returns the size of the scrolling region for whichever document
// pane is on screen.
//
// The width has to match the content width of the panel exactly. If it does
// not, the lipgloss box re-wraps the rows and the viewport's line count, and
// therefore its reported position, stops matching the screen.
func (m model) docRegion() (width, height int) {
	width = m.width - 2 // the panel content width that View passes down
	panelH := m.mainPanelHeight()

	switch {
	case m.level == levelChange:
		// The change name line and the sub-tab row stay put above the document.
		height = panelH - 2

	case m.detailTab == tabSpecs:
		// The specs tab splits its width; the document is the content half.
		_, contentWidth := specsSplit(width)
		width = contentWidth
		height = panelH - 5

	default:
		// The project header takes four rows and the tab bar one, then the
		// config pane's own source line and the blank line under it.
		height = panelH - 5 - 2
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
	if m.activeView == viewLog {
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
		return info.ConfigFile != ""
	case tabSpecs:
		// The specs tab has two halves. The vertical keys belong to the content
		// only while the content holds the keyboard.
		return m.specsFocus == specsFocusContent &&
			m.specCursor < len(info.SpecNames)
	}
	return false
}

// renderChangeArtifact renders the scrolling part of an open change: the
// selected artifact, or the change's specs.
func renderChangeArtifact(r changeRow, artifactTab, width int) string {
	var content strings.Builder

	if artifactTab < len(r.ci.ArtifactFiles) {
		filename := r.ci.ArtifactFiles[artifactTab]
		if filename == "tasks.md" && r.ci.TasksTotal > 0 {
			// Part of the document, so it scrolls away with the tasks it counts.
			content.WriteString(sectionHeaderStyle.Render(
				fmt.Sprintf("Tasks: %d/%d complete", r.ci.TasksDone, r.ci.TasksTotal)))
			content.WriteString("\n\n")
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

	return content.String()
}

// currentDocument returns the identity and the rendered rows of the document
// that should be on screen, if there is one.
func (m model) currentDocument() (key, content string, ok bool) {
	if len(m.repoPaths) == 0 || m.cursor >= len(m.repoPaths) {
		return "", "", false
	}
	project := m.repoPaths[m.cursor]
	width, _ := m.docRegion()

	switch {
	case m.level == levelChange:
		r, found := m.selectedRow()
		if !found {
			return "", "", false
		}
		names := r.artifactTabNames()
		if len(names) == 0 {
			return "", "", false
		}
		tab := m.changeArtifactTab
		if tab >= len(names) {
			tab = len(names) - 1
		}
		key = strings.Join([]string{project, "change", r.key(), names[tab]}, docKeySep)
		return key, renderChangeArtifact(r, tab, width), true

	case m.level == levelProject && m.detailTab == tabSpecs:
		info := m.projects[project].Info
		if m.specCursor >= len(info.SpecNames) {
			return "", "", false
		}
		name := info.SpecNames[m.specCursor]
		// Keyed by spec name, so moving the cursor is moving to a different
		// document and starts at the top, by the same rule as artifact sub-tabs.
		key = strings.Join([]string{project, "spec", name}, docKeySep)
		content, found := info.SpecContents[name]
		if !found || content == "" {
			return key, dimStyle.Render("No spec.md found"), true
		}
		return key, renderMarkdown(content, width), true

	case m.level == levelProject && m.detailTab == tabConfig:
		info := m.projects[project].Info
		if info.ConfigFile == "" {
			return "", "", false
		}
		key = strings.Join([]string{project, "config", info.ConfigFile}, docKeySep)
		if strings.HasSuffix(info.ConfigFile, ".md") {
			return key, renderMarkdown(info.ConfigContent, width), true
		}
		return key, renderYAML(info.ConfigContent, width), true
	}

	return "", "", false
}

// syncDocument loads the document that should be on screen into the viewport.
//
// One rule covers every reset case: a different document starts at the top, the
// same document keeps its position. That is what makes the filesystem watcher
// bearable while reading, since a rewritten file is still the same document and
// the position stays where the eye is.
func (m *model) syncDocument() {
	key, content, ok := m.currentDocument()
	if !ok {
		m.docKey = ""
		return
	}

	width, height := m.docRegion()
	m.docViewport.SetWidth(width)
	m.docViewport.SetHeight(height)

	m.docViewport.SetContent(content)
	if key != m.docKey {
		m.docKey = key
		m.docViewport.GotoTop()
		return
	}
	// Same document, possibly shorter than it was. SetYOffset clamps.
	m.docViewport.SetYOffset(m.docViewport.YOffset())
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
