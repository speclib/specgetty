package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// projectRow is one line of the project picker.
type projectRow struct {
	path    string
	display string
	info    scanner.ProjectInfo
	files   []scanner.FileEntry
}

func (r projectRow) key() string        { return r.path }
func (r projectRow) searchName() string { return r.display }

// searchBodies exposes a project to a ':' query. Every spec and artifact is
// searchable by content, and the file paths are searchable as one body, so
// looking for a filename and looking for a phrase both work through one sigil.
func (r projectRow) searchBodies() map[string]string {
	bodies := make(map[string]string)

	var paths strings.Builder
	for _, f := range r.files {
		paths.WriteString(f.Path)
		paths.WriteString("\n")
	}
	bodies["paths"] = paths.String()

	for name, content := range r.info.SpecContents {
		bodies["specs/"+name] = content
	}
	for _, ci := range append(append([]scanner.ChangeInfo{}, r.info.Changes...), r.info.ArchivedChanges...) {
		for file, content := range ci.ArtifactContents {
			bodies[ci.Name+"/"+trimMarkdownSuffix(file)] = content
		}
		for spec, content := range ci.SpecContents {
			bodies[ci.Name+"/specs/"+spec] = content
		}
	}
	return bodies
}

// projectFields are fixed. The change list takes --change-fields because
// columns were asked for there; nobody has asked to configure these.
var projectFields = []fieldDef[projectRow]{
	{
		id: "name", header: "project", width: 0,
		value: func(r projectRow) string { return r.display },
	},
	{
		id: "specs", header: "specs", width: 6,
		value: func(r projectRow) string { return fmt.Sprintf("%d", r.info.SpecCount) },
	},
	{
		id: "changes", header: "active", width: 7,
		value: func(r projectRow) string { return fmt.Sprintf("%d", len(r.info.Changes)) },
	},
	{
		id: "archived", header: "archived", width: 9,
		value: func(r projectRow) string { return fmt.Sprintf("%d", len(r.info.ArchivedChanges)) },
	},
	{
		id: "tasks", header: "tasks", width: 9,
		value: func(r projectRow) string {
			if r.info.TasksTotal == 0 {
				return ""
			}
			return fmt.Sprintf("%d/%d", r.info.TasksDone, r.info.TasksTotal)
		},
	},
}

// pickerLoadedMsg carries the result of discovering projects.
type pickerLoadedMsg struct {
	rows []projectRow
	err  error
}

// buildProjectRows turns a scanned map into sorted rows with disambiguated
// display names.
func buildProjectRows(projects scanner.ProjectMap) []projectRow {
	paths := make([]string, 0, len(projects))
	for p := range projects {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	names := projectDisplayNames(paths)
	rows := make([]projectRow, len(paths))
	for i, p := range paths {
		rows[i] = projectRow{
			path:    p,
			display: names[i],
			info:    projects[p].Info,
			files:   projects[p].Files,
		}
	}
	return rows
}

// loadProjects discovers projects for the picker.
//
// Without refresh it prefers the cache, which holds only paths, and parses
// statistics from disk so the numbers shown are current. With refresh it walks
// the configured directories and rewrites the cache.
func loadProjects(config *scanner.Config, ignoreDirErrors bool, refresh bool) tea.Cmd {
	return func() tea.Msg {
		cachePath, cacheErr := scanner.CachePath()

		if !refresh && cacheErr == nil {
			if paths, ok := scanner.LoadCache(cachePath, config); ok {
				return pickerLoadedMsg{rows: buildProjectRows(scanner.ScanPaths(paths))}
			}
		}

		projects, err := scanner.Scan(config, ignoreDirErrors)
		if err != nil {
			return pickerLoadedMsg{err: err}
		}

		if cacheErr == nil {
			paths := make([]string, 0, len(projects))
			for p := range projects {
				paths = append(paths, p)
			}
			sort.Strings(paths)
			// A cache that cannot be written is not worth failing over: the
			// next run simply walks again.
			_ = scanner.SaveCache(cachePath, config, paths, time.Now())
		}

		return pickerLoadedMsg{rows: buildProjectRows(projects)}
	}
}

// pickerRows applies the picker's own filter.
func (m model) pickerVisibleRows() []filtered[projectRow] {
	return filterRows(m.pickerAll, parseQuery(m.pickerInput.Value()))
}

// pickerSelected returns the project under the picker cursor.
func (m model) pickerSelected() (projectRow, bool) {
	rows := m.pickerVisibleRows()
	if m.pickerCursor < 0 || m.pickerCursor >= len(rows) {
		return projectRow{}, false
	}
	return rows[m.pickerCursor].row, true
}

// pickerRemember records which project the cursor is on, by key.
func (m *model) pickerRemember() {
	rows := m.pickerVisibleRows()
	if m.pickerCursor >= 0 && m.pickerCursor < len(rows) {
		m.pickerKey = rows[m.pickerCursor].row.key()
	}
}

// pickerSync keeps the cursor on the same project when the list changes under
// it, whether from a keystroke or from a refresh.
func (m *model) pickerSync() {
	rows := m.pickerVisibleRows()
	if len(rows) == 0 {
		m.pickerCursor = 0
		return
	}
	if m.pickerKey != "" {
		if i := indexOfKey(rows, m.pickerKey); i >= 0 {
			m.pickerCursor = i
			return
		}
	}
	if m.pickerCursor >= len(rows) {
		m.pickerCursor = len(rows) - 1
	}
	if m.pickerCursor < 0 {
		m.pickerCursor = 0
	}
	m.pickerKey = rows[m.pickerCursor].row.key()
}

// pickerBox returns the overlay dimensions, in total columns and rows
// including the border.
func (m model) pickerBox() (width, height int) {
	width = min(m.width-8, 110)
	if width < 40 {
		width = max(m.width-4, 20)
	}
	// Sized to the number of projects, not to how many currently match, so the
	// box does not jump about while a filter is being typed.
	rows := len(m.pickerAll)
	if rows < 3 {
		rows = 3
	}
	height = min(rows+4, min(m.height-6, 22))
	if height < 6 {
		height = max(m.height-2, 4)
	}
	return width, height
}

// renderPicker draws the overlay: a table of projects with a prompt beneath.
func (m model) renderPicker() string {
	width, height := m.pickerBox()

	// lipgloss v2 counts the border and the padding inside Style.Width, so a box
	// of `width` total columns asks for exactly that, and the content area is
	// four columns narrower: one border and one padding column on each side.
	//
	// Handing the table the wrong number here does not error. It silently
	// overflows, lipgloss wraps the row, and the result shows up as the
	// highlighted row alone looking misaligned, because only that row carries a
	// background colour into the wrapped remainder. The line count is what
	// catches it; equal line widths do not, since the wrap is padded.
	styleWidth := width
	inner := width - 4

	var body string
	switch {
	case m.pickerLoading:
		body = m.spinner.View() + " Looking for OpenSpec projects..."
	case m.pickerErr != "":
		body = "Could not scan: " + m.pickerErr
	case len(m.pickerAll) == 0:
		body = dimStyle.Render("No OpenSpec projects found. Press r to scan again.")
	default:
		rows := m.pickerVisibleRows()
		tableHeight := height - 4
		if tableHeight < 1 {
			tableHeight = 1
		}
		if len(rows) == 0 {
			body = dimStyle.Render(fmt.Sprintf("No projects match %q", m.pickerInput.Value()))
			body = truncateContent(body, tableHeight)
		} else {
			body = renderTable(rows, projectFields, m.pickerCursor, inner, tableHeight)
			body = truncateContent(body, tableHeight)
		}
		body += "\n" + renderSearchPrompt(
			m.pickerInput.Value(), m.pickerFocused, len(rows), len(m.pickerAll))
	}

	title := headerStyle.Render("Projects")
	content := title + "\n" + body

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("2")).
		Padding(0, 1).
		Width(styleWidth).
		Render(content)
}

// choosePickerProject switches to the project under the picker cursor.
//
// It always lands on the project view with the changes tab active. Staying at
// the change level would be meaningless, because the change that was open does
// not exist in the newly selected project.
func (m model) choosePickerProject() (tea.Model, tea.Cmd) {
	r, ok := m.pickerSelected()
	if !ok {
		return m, nil
	}

	m.pickerOpen = false
	m.pickerFocused = false
	m.pickerInput.Blur()

	// Watching follows the selection, and only one project is watched at a
	// time, so the previous watcher is closed before the next one opens.
	m.stopWatcher()

	m.projects = scanner.ProjectMap{r.path: scanner.ProjectStatus{Files: r.files, Info: r.info}}
	m.repoPaths = []string{r.path}
	m.displayNames = []string{r.display}
	m.cursor = 0

	m.level = levelProject
	m.detailTab = tabChanges
	m.focus = focusDetail
	m.resetProjectState()

	var cmds []tea.Cmd
	if cmd := m.startWatcher(r.path); cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
