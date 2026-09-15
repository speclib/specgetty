package ui

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
	"github.com/mipmip/specgetty/src/watcher"
)

// Focus. There is one main panel, so the only question is whether the optional
// log panel has the keyboard.
const (
	viewDetail = 0
	viewLog    = 1
)

// Navigation depth. enter descends, esc ascends. levelProject is the floor:
// the project list is an overlay now, not a level above it.
const (
	levelProject = 0 // one project, with its tab bar
	levelChange  = 1 // one change, with its artifact sub-tabs
)

// Changes lead, because that is what the tool is usually opened to look at.
const (
	tabChanges = 0
	tabSpecs   = 1
	tabConfig  = 2
)

const (
	archiveIdle       = 0
	archiveConfirming = 1
	archiveRunning    = 2
	archiveResult     = 3
)

const (
	discardIdle       = 0
	discardConfirming = 1
	discardRunning    = 2
	discardResult     = 3
)

const (
	exportIdle       = 0
	exportConfirming = 1
	exportRunning    = 2
	exportResult     = 3
)

var tabNames = []string{"changes", "specs", "config"}

// Message types

type scanMsg struct {
	projects scanner.ProjectMap
	err      error
}

type logMsg string

// fsChangeMsg is sent when the filesystem watcher detects changes in the openspec directory.
type fsChangeMsg struct{}

// logWriter sends log output as tea messages to the program.
type logWriter struct {
	program *tea.Program
}

func (w logWriter) Write(p []byte) (n int, err error) {
	w.program.Send(logMsg(string(p)))
	return len(p), nil
}

// Styles

var (
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("2")). // green
			Foreground(lipgloss.Color("0")). // black
			Width(0)                         // set dynamically

	normalStyle = lipgloss.NewStyle()

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("1")). // red
			Padding(1, 2).
			Align(lipgloss.Center)

	navBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252"))

	navBarKeyStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("2")).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("4")) // blue

	activeTabStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("2")).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("250")).
				Padding(0, 1)

	sectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("3")) // yellow

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))
)

type model struct {
	config            *scanner.Config
	ignoreDirErrors   bool
	projects          scanner.ProjectMap
	repoPaths         []string
	displayNames      []string
	cursor            int
	activeView        int
	scanning          bool
	err               error
	spinner           spinner.Model
	docViewport       viewport.Model
	docKey            string
	logViewport       viewport.Model
	logContent        string
	detailTab         int
	specCursor        int
	changeCursor      int
	changeArtifactTab int
	level             int    // navigation depth: levelProject, levelChange
	startupPath       string // project resolved at startup, scanned on its own

	// Project picker state. The picker is an overlay, not a level: it opens
	// from anywhere and always returns to the project view.
	pickerOpen    bool
	pickerAll     []projectRow
	pickerCursor  int
	pickerKey     string
	pickerInput   textinput.Model
	pickerFocused bool
	pickerLoading bool
	pickerLoaded  bool
	pickerErr     string

	startView string // "single" or "all"

	// Set when startup found no project, so the user is asked whether to pick
	// one rather than being dropped into an empty view with no explanation.
	askOpenPicker bool

	// Change list state.
	listMode          int // modeOpen, modeArchived, modeBoth
	fields            []string
	searchInput       textinput.Model
	searchFocused     bool
	selectedKey       string // identifies the selected change across re-filter and rescan
	logVisible        bool
	logShownOnce      bool
	pendingKey        string
	archiveState      int
	archiveChangeName string
	archiveResultMsg  string
	archiveResultOk   bool
	discardState      int
	discardChangeName string
	discardResultMsg  string
	discardResultOk   bool
	exportState       int
	exportChangeName  string
	exportResultMsg   string
	exportResultOk    bool
	exportIsArchived  bool
	width             int
	height            int
	program           *tea.Program
	version           string
	watcher           *watcher.Watcher
}

func newModel(config *scanner.Config, ignoreDirErrors bool, version string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = ""

	pi := textinput.New()
	pi.Prompt = ""
	pi.Placeholder = ""

	return model{
		config:          config,
		ignoreDirErrors: ignoreDirErrors,
		version:         version,
		spinner:         s,
		docViewport:     viewport.New(0, 0),
		logViewport:     viewport.New(0, 0),
		listMode:        modeOpen,
		fields:          append([]string(nil), defaultFields...),
		searchInput:     ti,
		pickerInput:     pi,
	}
}

func (m model) Init() tea.Cmd {
	// Startup never walks the configured directories. Either a project was
	// resolved from the working directory, which costs about a millisecond, or
	// the picker is asked for explicitly and pays for discovery itself.
	if m.startupPath != "" {
		return tea.Batch(m.spinner.Tick, m.doScanSingle(m.startupPath))
	}
	if m.pickerOpen {
		return tea.Batch(m.spinner.Tick, loadProjects(m.config, m.ignoreDirErrors, false))
	}
	return m.spinner.Tick
}

// Update wraps the message handling so that the document viewer is
// resynchronised on every path, including the early returns taken by the
// overlays. Re-rendering costs a wrap of one document and keeps the viewport
// from ever holding stale rows.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.update(msg)
	next := updated.(model)
	next.syncDocument()
	return next, cmd
}

func (m model) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()

	case tea.KeyMsg:
		if m.scanning {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

		if m.err != nil {
			m.err = nil
			return m, nil
		}

		// Archive modal intercepts all keys
		if m.archiveState == archiveConfirming {
			switch msg.String() {
			case "y":
				m.archiveState = archiveRunning
				projectPath := m.repoPaths[m.cursor]
				cmds = append(cmds, doArchiveChange(projectPath, m.archiveChangeName))
			case "n", "esc":
				m.archiveState = archiveIdle
			}
			return m, tea.Batch(cmds...)
		}
		if m.archiveState == archiveResult {
			m.archiveState = archiveIdle
			return m, nil
		}
		if m.archiveState == archiveRunning {
			return m, nil
		}

		// Export modal intercepts all keys
		if m.exportState == exportConfirming {
			switch msg.String() {
			case "y":
				m.exportState = exportRunning
				projectPath := m.repoPaths[m.cursor]
				cmds = append(cmds, doExportChange(projectPath, m.exportChangeName, m.exportIsArchived))
			case "n", "esc":
				m.exportState = exportIdle
			}
			return m, tea.Batch(cmds...)
		}
		if m.exportState == exportResult {
			m.exportState = exportIdle
			return m, nil
		}
		if m.exportState == exportRunning {
			return m, nil
		}

		// Discard modal intercepts all keys
		if m.discardState == discardConfirming {
			switch msg.String() {
			case "y":
				m.discardState = discardRunning
				projectPath := m.repoPaths[m.cursor]
				cmds = append(cmds, doDiscardChange(projectPath, m.discardChangeName))
			case "n", "esc":
				m.discardState = discardIdle
			}
			return m, tea.Batch(cmds...)
		}
		if m.discardState == discardResult {
			m.discardState = discardIdle
			return m, nil
		}
		if m.discardState == discardRunning {
			return m, nil
		}

		key := msg.String()

		// The startup prompt is a question that must be answered before
		// anything else makes sense.
		if m.askOpenPicker {
			switch key {
			case "y", "Y", "enter":
				m.askOpenPicker = false
				m.pickerOpen = true
				if !m.pickerLoaded {
					m.pickerLoading = true
					cmds = append(cmds, loadProjects(m.config, m.ignoreDirErrors, false))
				}
			case "n", "N", "esc":
				m.askOpenPicker = false
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, tea.Batch(cmds...)
		}

		// The picker is an overlay drawn over the current view, so its keys
		// outrank the view beneath it. A confirmation modal still outranks the
		// picker, which is why those checks come first.
		if m.pickerOpen {
			if m.pickerFocused {
				switch key {
				case "esc":
					m.pickerFocused = false
					m.pickerInput.SetValue("")
					m.pickerInput.Blur()
					m.pickerSync()
					return m, nil
				case "enter":
					return m.choosePickerProject()
				case "up", "ctrl+p":
					if m.pickerCursor > 0 {
						m.pickerCursor--
						m.pickerRemember()
					}
					return m, nil
				case "down", "ctrl+n":
					if m.pickerCursor < len(m.pickerVisibleRows())-1 {
						m.pickerCursor++
						m.pickerRemember()
					}
					return m, nil
				case "ctrl+c":
					m.stopWatcher()
					return m, tea.Quit
				}
				var cmd tea.Cmd
				m.pickerInput, cmd = m.pickerInput.Update(msg)
				m.pickerSync()
				return m, cmd
			}

			switch key {
			case "esc", "p":
				m.pickerOpen = false
			case "q", "ctrl+c":
				m.stopWatcher()
				return m, tea.Quit
			case "enter":
				return m.choosePickerProject()
			case "/":
				m.pickerFocused = true
				m.pickerInput.Focus()
			case "r":
				m.pickerLoading = true
				cmds = append(cmds, loadProjects(m.config, m.ignoreDirErrors, true))
			case "up", "k", "ctrl+p":
				if m.pickerCursor > 0 {
					m.pickerCursor--
					m.pickerRemember()
				}
			case "down", "j", "ctrl+n":
				if m.pickerCursor < len(m.pickerVisibleRows())-1 {
					m.pickerCursor++
					m.pickerRemember()
				}
			case "g":
				m.pickerCursor = 0
				m.pickerRemember()
			case "G":
				if n := len(m.pickerVisibleRows()); n > 0 {
					m.pickerCursor = n - 1
					m.pickerRemember()
				}
			}
			return m, tea.Batch(cmds...)
		}

		// While the search prompt has focus every rune belongs to it, because
		// a, d, e, s, l and q are all actions. Only the keys below escape it.
		if m.searchFocused {
			switch key {
			case "esc":
				m.searchFocused = false
				m.searchInput.SetValue("")
				m.searchInput.Blur()
				m.syncCursor()
				return m, nil
			case "enter":
				if _, ok := m.selectedRow(); ok {
					m.rememberSelection()
					m.level = levelChange
					m.changeArtifactTab = 0
				}
				return m, nil
			case "up", "ctrl+p":
				if m.changeCursor > 0 {
					m.changeCursor--
					m.rememberSelection()
				}
				return m, nil
			case "down", "ctrl+n":
				if m.changeCursor < len(m.currentRows())-1 {
					m.changeCursor++
					m.rememberSelection()
				}
				return m, nil
			case "ctrl+c":
				m.stopWatcher()
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			m.syncCursor()
			return m, cmd
		}

		if m.pendingKey == "g" {
			m.pendingKey = ""
			if key == "g" {
				switch {
				case m.activeView == viewLog:
					m.logViewport.GotoTop()
				case m.docActive():
					m.docViewport.GotoTop()
				}
				return m, nil
			}
		}

		switch key {
		case "q", "ctrl+c":
			m.stopWatcher()
			return m, tea.Quit

		case "enter":
			if m.level == levelProject && m.detailTab == tabChanges {
				if _, ok := m.selectedRow(); ok {
					m.rememberSelection()
					m.level = levelChange
					m.changeArtifactTab = 0
				}
			}

		case "esc":
			// levelProject is the floor. The project list is an overlay now,
			// so there is nothing above it to escape to.
			if m.level == levelChange {
				m.level = levelProject
				m.syncCursor()
			}

		case "p":
			m.pickerOpen = true
			m.pickerSync()
			if !m.pickerLoaded && !m.pickerLoading {
				m.pickerLoading = true
				cmds = append(cmds, loadProjects(m.config, m.ignoreDirErrors, false))
			}

		case "/":
			if m.level == levelProject && m.detailTab == tabChanges {
				m.searchFocused = true
				m.searchInput.Focus()
			}

		case "f":
			if m.level == levelProject && m.detailTab == tabChanges {
				m.listMode = (m.listMode + 1) % 3
				m.syncCursor()
			}

		case "s":
			if len(m.repoPaths) > 0 {
				m.scanning = true
				cmds = append(cmds, m.doScanSingle(m.repoPaths[m.cursor]))
			}

		case "l":
			m.logVisible = !m.logVisible
			if m.logVisible {
				m.recalcLayout()
				if !m.logShownOnce {
					m.logShownOnce = true
					m.logViewport.GotoBottom()
				}
			} else {
				if m.activeView == viewLog {
					m.activeView = viewDetail
				}
				m.recalcLayout()
			}

		case "tab":
			// One main panel, so tab only has somewhere to go when the log
			// panel is open.
			if m.logVisible {
				if m.activeView == viewDetail {
					m.activeView = viewLog
				} else {
					m.activeView = viewDetail
				}
			}

		case "g":
			m.pendingKey = "g"

		case "G":
			switch {
			case m.activeView == viewLog:
				m.logViewport.GotoBottom()
			case m.docActive():
				m.docViewport.GotoBottom()
			}

		case "pgdown", "ctrl+f":
			// A full page in a document, vim style. Lists keep halfPage(),
			// which is how they have always moved.
			switch {
			case m.activeView == viewLog:
				m.logViewport.LineDown(m.halfPage())
			case m.docActive():
				m.docViewport.ViewDown()
			}

		case "pgup", "ctrl+b":
			switch {
			case m.activeView == viewLog:
				m.logViewport.LineUp(m.halfPage())
			case m.docActive():
				m.docViewport.ViewUp()
			}

		// 5.3: half page, in a document only.
		case "ctrl+d":
			if m.docActive() {
				m.docViewport.HalfPageDown()
			}

		case "ctrl+u":
			if m.docActive() {
				m.docViewport.HalfPageUp()
			}

		case "left":
			// Sub-tabs belong to an open change, the tab bar to the project.
			// Neither spills into the other.
			if m.level == levelChange {
				if m.changeArtifactTab > 0 {
					m.changeArtifactTab--
				}
			} else if m.activeView == viewDetail && m.detailTab > 0 {
				m.detailTab--
			}

		case "right":
			if m.level == levelChange {
				if m.changeArtifactTab < m.changeArtifactTabCount()-1 {
					m.changeArtifactTab++
				}
			} else if m.activeView == viewDetail && m.detailTab < len(tabNames)-1 {
				m.detailTab++
			}

		case "1", "2", "3":
			// Number keys address the project tab bar, so they are inert while
			// a change is open.
			if m.level != levelChange && m.activeView == viewDetail {
				m.detailTab = int(key[0] - '1')
			}

		case "up", "k":
			if m.activeView == viewLog {
				m.logViewport.LineUp(1)
			} else if m.docActive() {
				m.docViewport.LineUp(1)
			} else if m.level != levelChange {
				switch m.detailTab {
				case tabSpecs:
					if m.specCursor > 0 {
						m.specCursor--
					}
				case tabChanges:
					if m.changeCursor > 0 {
						m.changeCursor--
						m.changeArtifactTab = 0
						m.rememberSelection()
					}
				}
			}

		case "down", "j":
			if m.activeView == viewLog {
				m.logViewport.LineDown(1)
			} else if m.docActive() {
				m.docViewport.LineDown(1)
			} else if m.level != levelChange {
				switch m.detailTab {
				case tabSpecs:
					if m.specCursor < len(m.currentSpecNames())-1 {
						m.specCursor++
					}
				case tabChanges:
					if m.changeCursor < len(m.currentRows())-1 {
						m.changeCursor++
						m.changeArtifactTab = 0
						m.rememberSelection()
					}
				}
			}

		case "a":
			// Archiving an already archived change is a no-op.
			if r, ok := m.selectedRow(); ok && m.detailTab == tabChanges && !r.archived {
				m.archiveChangeName = r.ci.Name
				m.archiveState = archiveConfirming
			}

		case "d":
			if r, ok := m.selectedRow(); ok && m.detailTab == tabChanges && !r.archived {
				m.discardChangeName = r.ci.Name
				m.discardState = discardConfirming
			}

		case "e":
			if r, ok := m.selectedRow(); ok && m.detailTab == tabChanges {
				m.exportChangeName = r.ci.Name
				m.exportIsArchived = r.archived
				m.exportState = exportConfirming
			}
		}

	case archiveMsg:
		m.archiveResultOk = msg.ok
		m.archiveResultMsg = msg.output
		m.archiveState = archiveResult
		if msg.ok && len(m.repoPaths) > 0 {
			cmds = append(cmds, m.doScanSingle(m.repoPaths[m.cursor]))
		}

	case discardMsg:
		m.discardResultOk = msg.ok
		m.discardResultMsg = msg.output
		m.discardState = discardResult
		if msg.ok && len(m.repoPaths) > 0 {
			cmds = append(cmds, m.doScanSingle(m.repoPaths[m.cursor]))
		}

	case exportMsg:
		m.exportResultOk = msg.ok
		m.exportResultMsg = msg.output
		m.exportState = exportResult

	case fsChangeMsg:
		if m.level >= levelProject && len(m.repoPaths) > 0 {
			m.scanning = true
			cmds = append(cmds, m.doScanSingle(m.repoPaths[m.cursor]))
			if m.watcher != nil {
				cmds = append(cmds, waitForFsChange(m.watcher))
			}
		}

	case scanMsg:
		m.scanning = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.projects = msg.projects
			m.repoPaths = make([]string, 0, len(m.projects))
			for r := range m.projects {
				m.repoPaths = append(m.repoPaths, r)
			}
			sort.Strings(m.repoPaths)
			m.displayNames = projectDisplayNames(m.repoPaths)
			if m.cursor >= len(m.repoPaths) {
				m.cursor = max(0, len(m.repoPaths)-1)
			}
			m.recalcLayout()
			// A rescan fires on every file save while watching, so the filter
			// and the selection are preserved across it rather than reset.
			m.syncCursor()

			// Start watching the open project after its first scan.
			if m.watcher == nil && len(m.repoPaths) > 0 {
				if cmd := m.startWatcher(m.repoPaths[m.cursor]); cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}

	case pickerLoadedMsg:
		m.pickerLoading = false
		if msg.err != nil {
			m.pickerErr = msg.err.Error()
		} else {
			m.pickerErr = ""
			m.pickerAll = msg.rows
			m.pickerLoaded = true
			m.pickerSync()
		}

	case logMsg:
		m.logContent += string(msg)
		m.logViewport.SetContent(m.logContent)
		m.logViewport.GotoBottom()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// projectDisplayNames returns basenames for each path, disambiguating duplicates
// by appending the parent directory name.
func projectDisplayNames(paths []string) []string {
	names := make([]string, len(paths))
	baseCounts := make(map[string]int)

	for _, p := range paths {
		base := filepath.Base(p)
		baseCounts[base]++
	}

	for i, p := range paths {
		base := filepath.Base(p)
		if baseCounts[base] > 1 {
			parent := filepath.Base(filepath.Dir(p))
			names[i] = base + " (" + parent + ")"
		} else {
			names[i] = base
		}
	}
	return names
}

func (m *model) recalcLayout() {
	if m.width == 0 || m.height == 0 {
		return
	}

	docW, docH := m.docRegion()
	m.docViewport.Width = docW
	m.docViewport.Height = docH

	logHeight := m.logPanelHeight()
	if logHeight > 0 {
		m.logViewport.Width = m.width - 2
		m.logViewport.Height = logHeight
	}
}

func (m model) mainPanelHeight() int {
	logH := m.logPanelHeight()
	if logH > 0 {
		logH += 2 // border
	}
	remaining := m.height - logH - 1 // -1 for nav bar
	if remaining < 5 {
		return 3
	}
	return remaining - 2 // -border
}

func (m model) logPanelHeight() int {
	if !m.logVisible {
		return 0
	}
	return min(10, (m.height-6)/3)
}

func (m model) halfPage() int {
	switch m.activeView {
	case viewDetail:
		return max(1, m.mainPanelHeight()/2)
	case viewLog:
		return max(1, m.logViewport.Height/2)
	default:
		return max(1, m.mainPanelHeight()/2)
	}
}

func (m model) currentSpecNames() []string {
	if len(m.repoPaths) == 0 {
		return nil
	}
	return m.projects[m.repoPaths[m.cursor]].Info.SpecNames
}

// allRows builds the merged, mode-filtered change list, before any search
// narrows it.
func (m model) allRows() []changeRow {
	if len(m.repoPaths) == 0 {
		return nil
	}
	return buildRows(m.projects[m.repoPaths[m.cursor]].Info, m.listMode)
}

// currentRows is the single seam every consumer of the change list goes
// through: the renderer, the actions, the confirm modals and the nav bar. The
// active/archived merge, the mode filter and the search filter all apply here,
// so nothing downstream has to know a filter exists.
func (m model) currentRows() []filtered[changeRow] {
	return filterRows(m.allRows(), parseQuery(m.searchInput.Value()))
}

// selectedRow returns the change under the cursor, if there is one.
func (m model) selectedRow() (changeRow, bool) {
	rows := m.currentRows()
	if m.changeCursor < 0 || m.changeCursor >= len(rows) {
		return changeRow{}, false
	}
	return rows[m.changeCursor].row, true
}

// rememberSelection records which change the cursor is on, by key.
func (m *model) rememberSelection() {
	rows := m.currentRows()
	if m.changeCursor >= 0 && m.changeCursor < len(rows) {
		m.selectedKey = rows[m.changeCursor].row.key()
	}
}

// syncCursor keeps the selection on the same change when the list changes
// underneath it, whether from a keystroke narrowing the filter or from a
// rescan. Clamping by index alone would silently move the selection onto a
// different change, which matters because archive, discard and export all act
// on whatever is under the cursor.
func (m *model) syncCursor() {
	rows := m.currentRows()
	if len(rows) == 0 {
		m.changeCursor = 0
		return
	}
	if m.selectedKey != "" {
		if i := indexOfKey(rows, m.selectedKey); i >= 0 {
			m.changeCursor = i
			return
		}
	}
	if m.changeCursor >= len(rows) {
		m.changeCursor = len(rows) - 1
	}
	if m.changeCursor < 0 {
		m.changeCursor = 0
	}
	m.selectedKey = rows[m.changeCursor].row.key()
}

func (m model) changeArtifactTabCount() int {
	r, ok := m.selectedRow()
	if !ok {
		return 0
	}
	return len(r.artifactTabNames())
}

// resetProjectState clears everything scoped to a single project: the cursors,
// the search filter and the list mode.
//
// It runs when the selected project changes, NOT when the same project is
// rescanned. The filesystem watcher fires a rescan on every file save, and
// losing an active filter mid-edit would be maddening.
func (m *model) resetProjectState() {
	m.specCursor = 0
	m.changeCursor = 0
	m.changeArtifactTab = 0
	m.selectedKey = ""
	m.listMode = modeOpen
	m.searchFocused = false
	m.searchInput.SetValue("")
	m.searchInput.Blur()
}

func (m model) doScan() tea.Cmd {
	config := m.config
	ignoreDirErrors := m.ignoreDirErrors
	return func() tea.Msg {
		projects, err := scanner.Scan(config, ignoreDirErrors)
		return scanMsg{projects: projects, err: err}
	}
}

func (m model) doScanSingle(projectPath string) tea.Cmd {
	return func() tea.Msg {
		files, err := scanner.ListOpenSpecContents(projectPath)
		if err != nil {
			return scanMsg{err: err}
		}
		info := scanner.ParseProjectInfo(projectPath)
		projects := scanner.ProjectMap{
			projectPath: scanner.ProjectStatus{
				Files: files,
				Info:  info,
			},
		}
		return scanMsg{projects: projects}
	}
}

func waitForFsChange(w *watcher.Watcher) tea.Cmd {
	return func() tea.Msg {
		_, ok := <-w.Events()
		if !ok {
			return nil // watcher closed
		}
		return fsChangeMsg{}
	}
}

func (m *model) startWatcher(projectPath string) tea.Cmd {
	openspecDir := filepath.Join(projectPath, "openspec")
	w, err := watcher.New(openspecDir)
	if err != nil {
		log.Printf("watcher: failed to start: %v", err)
		return nil
	}
	m.watcher = w
	return waitForFsChange(w)
}

func (m *model) stopWatcher() {
	if m.watcher != nil {
		m.watcher.Close()
		m.watcher = nil
	}
}

type archiveMsg struct {
	ok     bool
	output string
}

type discardMsg struct {
	ok     bool
	output string
}

type exportMsg struct {
	ok     bool
	output string
}

func doArchiveChange(projectPath string, changeName string) tea.Cmd {
	return func() tea.Msg {
		if _, err := exec.LookPath("openspec"); err != nil {
			return archiveMsg{ok: false, output: "openspec CLI not found. Install it to enable archiving."}
		}
		cmd := exec.Command("openspec", "archive", changeName, "-y")
		cmd.Dir = projectPath
		out, err := cmd.CombinedOutput()
		if err != nil {
			return archiveMsg{ok: false, output: strings.TrimSpace(string(out))}
		}
		return archiveMsg{ok: true, output: strings.TrimSpace(string(out))}
	}
}

func doDiscardChange(projectPath string, changeName string) tea.Cmd {
	return func() tea.Msg {
		changesDir := filepath.Join(projectPath, "openspec", "changes")
		discardedDir := filepath.Join(changesDir, "discarded")
		if err := os.MkdirAll(discardedDir, 0755); err != nil {
			return discardMsg{ok: false, output: fmt.Sprintf("Failed to create discarded directory: %v", err)}
		}
		datePrefix := time.Now().Format("2006-01-02")
		target := filepath.Join(discardedDir, datePrefix+"-"+changeName)
		if _, err := os.Stat(target); err == nil {
			return discardMsg{ok: false, output: fmt.Sprintf("Target already exists: %s", datePrefix+"-"+changeName)}
		}
		src := filepath.Join(changesDir, changeName)
		if err := os.Rename(src, target); err != nil {
			return discardMsg{ok: false, output: fmt.Sprintf("Failed to move: %v", err)}
		}
		return discardMsg{ok: true, output: fmt.Sprintf("Discarded \"%s\"", changeName)}
	}
}

var archiveDatePrefix = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-`)

func exportSemanticName(changeName string, isArchived bool) string {
	if isArchived {
		return archiveDatePrefix.ReplaceAllString(changeName, "")
	}
	return changeName
}

func exportDestPath(semanticName string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dateStr := time.Now().Format("2006-01-02")
	return filepath.Join(home, semanticName+"-"+dateStr+".zip")
}

func doExportChange(projectPath string, changeName string, isArchived bool) tea.Cmd {
	return func() tea.Msg {
		var srcDir string
		if isArchived {
			srcDir = filepath.Join(projectPath, "openspec", "changes", "archive", changeName)
		} else {
			srcDir = filepath.Join(projectPath, "openspec", "changes", changeName)
		}

		if _, err := os.Stat(srcDir); err != nil {
			return exportMsg{ok: false, output: fmt.Sprintf("Source not found: %s", srcDir)}
		}

		semanticName := exportSemanticName(changeName, isArchived)
		destPath := exportDestPath(semanticName)

		zipFile, err := os.Create(destPath)
		if err != nil {
			return exportMsg{ok: false, output: fmt.Sprintf("Failed to create zip: %v", err)}
		}
		defer zipFile.Close()

		w := zip.NewWriter(zipFile)
		defer w.Close()

		err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}
			zipPath := filepath.Join(semanticName, relPath)

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = zipPath
			header.Method = zip.Deflate

			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})

		if err != nil {
			os.Remove(destPath)
			return exportMsg{ok: false, output: fmt.Sprintf("Failed to create zip: %v", err)}
		}

		return exportMsg{ok: true, output: fmt.Sprintf("Exported to %s", destPath)}
	}
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	if m.width < 60 || m.height < 20 {
		return "Terminal too small. Need at least 60x20."
	}

	panelH := m.mainPanelHeight()

	fullW := m.width - 2
	detailContent := m.renderDetailPanel(fullW, panelH)
	mainRow := m.renderPanel(viewDetail, fullW, panelH, detailContent)

	// Nav bar
	navBar := m.renderNavBar()

	var view string
	if m.logVisible {
		logPanel := m.renderPanel(viewLog, m.width-2, m.logPanelHeight(), m.logViewport.View())
		view = lipgloss.JoinVertical(lipgloss.Left, mainRow, logPanel, navBar)
	} else {
		view = lipgloss.JoinVertical(lipgloss.Left, mainRow, navBar)
	}

	// The picker is an overlay over the current view. Confirmation modals are
	// drawn after it, so they sit on top.
	if m.pickerOpen {
		view = placeOverlay(m.width, m.height, m.renderPicker(), view)
	}

	if m.askOpenPicker {
		modal := modalStyle.Width(56).Render(
			"No OpenSpec project here.\n\nOpen the project picker? (y/n)")
		view = placeOverlay(m.width, m.height, modal, view)
	}

	// Modal overlays
	if m.scanning {
		modal := modalStyle.Width(40).Render(m.spinner.View() + " Scanning for OpenSpec sources...")
		view = placeOverlay(m.width, m.height, modal, view)
	}
	if m.err != nil {
		errText := fmt.Sprintf("Error: %v", m.err)
		modal := modalStyle.Width(m.width * 3 / 4).Render(errText)
		view = placeOverlay(m.width, m.height, modal, view)
	}

	// Archive modals
	switch m.archiveState {
	case archiveConfirming:
		var content string
		if r, ok := m.selectedRow(); ok {
			if r.ci.TasksTotal > 0 && r.ci.TasksDone < r.ci.TasksTotal {
				incomplete := r.ci.TasksTotal - r.ci.TasksDone
				content = fmt.Sprintf("⚠ %d incomplete task(s) in \"%s\"\n\nArchive anyway? (y/n)", incomplete, m.archiveChangeName)
			} else {
				content = fmt.Sprintf("Archive \"%s\"? (y/n)", m.archiveChangeName)
			}
		}
		modal := modalStyle.Width(50).Render(content)
		view = placeOverlay(m.width, m.height, modal, view)
	case archiveRunning:
		modal := modalStyle.Width(40).Render(m.spinner.View() + " Archiving...")
		view = placeOverlay(m.width, m.height, modal, view)
	case archiveResult:
		var prefix string
		if m.archiveResultOk {
			prefix = "✓ "
		} else {
			prefix = "✗ "
		}
		content := prefix + m.archiveResultMsg + "\n\nPress any key to dismiss."
		modal := modalStyle.Width(m.width * 3 / 4).Render(content)
		view = placeOverlay(m.width, m.height, modal, view)
	}

	// Discard modals
	switch m.discardState {
	case discardConfirming:
		var content string
		if r, ok := m.selectedRow(); ok {
			if r.ci.TasksTotal > 0 && r.ci.TasksDone < r.ci.TasksTotal {
				incomplete := r.ci.TasksTotal - r.ci.TasksDone
				content = fmt.Sprintf("⚠ %d incomplete task(s) in \"%s\"\n\nDiscard anyway? (y/n)", incomplete, m.discardChangeName)
			} else {
				content = fmt.Sprintf("Discard \"%s\"? (y/n)", m.discardChangeName)
			}
		}
		modal := modalStyle.Width(50).Render(content)
		view = placeOverlay(m.width, m.height, modal, view)
	case discardRunning:
		modal := modalStyle.Width(40).Render(m.spinner.View() + " Discarding...")
		view = placeOverlay(m.width, m.height, modal, view)
	case discardResult:
		var prefix string
		if m.discardResultOk {
			prefix = "✓ "
		} else {
			prefix = "✗ "
		}
		content := prefix + m.discardResultMsg + "\n\nPress any key to dismiss."
		modal := modalStyle.Width(m.width * 3 / 4).Render(content)
		view = placeOverlay(m.width, m.height, modal, view)
	}

	// Export modals
	switch m.exportState {
	case exportConfirming:
		semanticName := exportSemanticName(m.exportChangeName, m.exportIsArchived)
		destPath := exportDestPath(semanticName)
		content := fmt.Sprintf("Export \"%s\"?\n\n→ %s\n\n(y/n)", semanticName, destPath)
		modal := modalStyle.Width(60).Render(content)
		view = placeOverlay(m.width, m.height, modal, view)
	case exportRunning:
		modal := modalStyle.Width(40).Render(m.spinner.View() + " Exporting...")
		view = placeOverlay(m.width, m.height, modal, view)
	case exportResult:
		var prefix string
		if m.exportResultOk {
			prefix = "✓ "
		} else {
			prefix = "✗ "
		}
		content := prefix + m.exportResultMsg + "\n\nPress any key to dismiss."
		modal := modalStyle.Width(m.width * 3 / 4).Render(content)
		view = placeOverlay(m.width, m.height, modal, view)
	}

	return padToHeight(view, m.height)
}

func (m model) renderTabHeader(width int) string {
	var b strings.Builder
	for i, name := range tabNames {
		if i > 0 {
			b.WriteString(" ")
		}
		if i == m.detailTab {
			b.WriteString(activeTabStyle.Render(name))
		} else {
			b.WriteString(inactiveTabStyle.Render(name))
		}
	}
	return b.String()
}

func (m model) renderDetailPanel(width int, height int) string {
	if len(m.repoPaths) == 0 {
		return "\n  " + dimStyle.Render("No project selected. Press p to pick one.")
	}

	// An open change fills the panel on its own: no project header, no tab bar,
	// so its artifact sub-tabs own the full width and their own key axis.
	if m.level == levelChange {
		if r, ok := m.selectedRow(); ok {
			return m.renderChangeDetail(r, m.changeArtifactTab)
		}
	}

	currentProject := m.repoPaths[m.cursor]
	info := m.projects[currentProject].Info

	var b strings.Builder

	// Persistent header: project path + stats (1 char padding all sides)
	statsLine := fmt.Sprintf("Specs: %d  Changes: %d active  Archived: %d",
		info.SpecCount, len(info.ActiveChanges), len(info.ArchivedChanges))
	if info.TasksTotal > 0 {
		statsLine += fmt.Sprintf("  Tasks: %d/%d", info.TasksDone, info.TasksTotal)
	}
	headerContent := headerStyle.Render(currentProject) + "\n" + dimStyle.Render(statsLine)
	b.WriteString(lipgloss.NewStyle().Padding(1, 1).Render(headerContent))
	b.WriteString("\n")

	b.WriteString(m.renderTabHeader(width))
	b.WriteString("\n")

	// Tab content (height minus header 4 lines (2 content + 2 padding) and tab header 1 line)
	contentHeight := height - 5
	if contentHeight < 1 {
		contentHeight = 1
	}

	switch m.detailTab {
	case tabSpecs:
		b.WriteString(m.renderSpecsTab(width, contentHeight))
	case tabChanges:
		b.WriteString(m.renderChangesTab(width, contentHeight))
	case tabConfig:
		b.WriteString(m.renderConfigTab(width, contentHeight))
	default:
		b.WriteString(m.renderNotImplemented(width, contentHeight))
	}

	return b.String()
}

// renderChangesTab draws the full-width change table, plus the search prompt
// whenever a filter is active or being typed.
func (m model) renderChangesTab(width int, height int) string {
	rows := m.currentRows()
	total := len(m.allRows())

	showPrompt := m.searchFocused || m.searchInput.Value() != ""
	tableHeight := height
	if showPrompt {
		tableHeight--
	}
	if tableHeight < 1 {
		tableHeight = 1
	}

	var body string
	switch {
	case total == 0:
		body = dimStyle.Render(emptyListMessage(m.listMode))
	case len(rows) == 0:
		body = dimStyle.Render(fmt.Sprintf("No changes match %q", m.searchInput.Value()))
	default:
		body = renderTable(rows, changeFieldDefs(m.fields), m.changeCursor, width, tableHeight)
	}

	if !showPrompt {
		return body
	}
	return truncateContent(body, tableHeight) + "\n" +
		renderSearchPrompt(m.searchInput.Value(), m.searchFocused, len(rows), total)
}

func (m model) renderSpecsTab(width int, height int) string {
	if len(m.repoPaths) == 0 {
		return "No project selected."
	}

	info := m.projects[m.repoPaths[m.cursor]].Info

	if len(info.SpecNames) == 0 {
		return dimStyle.Render("No specs found")
	}

	// Split: ~30% for spec list, ~70% for content
	listWidth := width * 3 / 10
	if listWidth < 15 {
		listWidth = 15
	}
	contentWidth := width - listWidth - 1 // 1 for gap

	// Render spec list
	var listB strings.Builder
	offset := 0
	if m.specCursor >= height {
		offset = m.specCursor - height + 1
	}
	end := offset + height
	if end > len(info.SpecNames) {
		end = len(info.SpecNames)
	}

	isActive := m.activeView == viewDetail && m.detailTab == tabSpecs
	for i := offset; i < end; i++ {
		if i > offset {
			listB.WriteString("\n")
		}
		name := info.SpecNames[i]
		if isActive && i == m.specCursor {
			listB.WriteString(selectedStyle.Width(listWidth).Render(name))
		} else {
			listB.WriteString(normalStyle.Render(name))
		}
	}

	// Render spec content
	var contentB strings.Builder
	selectedSpec := info.SpecNames[m.specCursor]
	specContent, ok := info.SpecContents[selectedSpec]
	if !ok || specContent == "" {
		contentB.WriteString(dimStyle.Render("No spec.md found"))
	} else {
		contentB.WriteString(renderMarkdown(specContent, contentWidth))
	}

	specList := truncateContent(listB.String(), height)
	specContentStr := truncateContent(contentB.String(), height)

	// Force fixed dimensions on both sides to prevent wrapping/flickering
	leftBox := lipgloss.NewStyle().Width(listWidth).Height(height).MaxHeight(height).Render(specList)
	rightBox := lipgloss.NewStyle().Width(contentWidth).Height(height).MaxHeight(height).Render(specContentStr)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
}

func (m model) renderConfigTab(width int, height int) string {
	if len(m.repoPaths) == 0 {
		return "No project selected."
	}

	info := m.projects[m.repoPaths[m.cursor]].Info
	if info.ConfigFile == "" {
		return dimStyle.Render("No project configuration found")
	}

	// The source line stays put above the scrolling region.
	var b strings.Builder
	b.WriteString(dimStyle.Render("openspec/" + info.ConfigFile))
	b.WriteString("\n\n")
	b.WriteString(m.docViewport.View())
	return b.String()
}

var (
	mdHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6")) // cyan

	mdBoldStyle = lipgloss.NewStyle().Bold(true)

	mdItalicStyle = lipgloss.NewStyle().Italic(true)

	yamlKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")) // cyan

	yamlCommentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("250")) // light gray

	yamlValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")) // green
)

// wrapBreakpoints are the characters wrapping may break on, besides spaces and
// the hyphen that ansi always treats as one. Slashes matter here: file paths
// and store paths are common in these documents and are otherwise unbreakable.
const wrapBreakpoints = "/,;:"

// activeSGRAfter returns the SGR sequence still in effect at the end of row,
// starting from the sequence that was active before it. A reset clears it, any
// other sequence replaces it.
//
// This renderer emits simple open/close pairs, one style at a time, so tracking
// the most recent sequence is enough. It is not a general SGR state machine.
func activeSGRAfter(active, row string) string {
	for i := 0; i < len(row); {
		j := strings.Index(row[i:], "\x1b[")
		if j < 0 {
			break
		}
		start := i + j
		end := strings.Index(row[start:], "m")
		if end < 0 {
			break
		}
		seq := row[start : start+end+1]
		if seq == "\x1b[0m" || seq == "\x1b[m" {
			active = ""
		} else {
			active = seq
		}
		i = start + end + 1
	}
	return active
}

// reopenStyles makes every row independently styled.
//
// Wrapping a styled span leaves the opening sequence on the first row only, so
// the continuation rows rely on terminal state carrying across the newline.
// That holds when the rows are printed together and breaks the moment a
// viewport slices them: scroll until a continuation row is at the top and it
// renders unstyled.
//
// Re-opening costs no display width, because escape sequences are not counted
// as cells.
func reopenStyles(rows []string) []string {
	var active string
	out := make([]string, len(rows))
	for i, row := range rows {
		prefixed := active + row
		active = activeSGRAfter(active, row)
		if active != "" {
			prefixed += "\x1b[0m"
		}
		out[i] = prefixed
	}
	return out
}

// wrapStyled breaks one already-styled line to the given width.
//
// ansi.Wrap is used rather than ansi.Wordwrap because it breaks a word that is
// longer than the limit instead of letting it overflow. An overflowing row
// would be re-wrapped by the lipgloss box afterwards, and the viewport's row
// count, and therefore its reported position, would no longer match the screen.
//
// It counts display cells and preserves escape sequences, so a bold span that
// straddles a boundary keeps its styling on both rows.
func wrapStyled(styled string, width int) []string {
	if width < 1 {
		width = 1
	}
	return reopenStyles(strings.Split(ansi.Wrap(styled, width, wrapBreakpoints), "\n"))
}

// styleMarkdownLine applies the styling for a single source line.
func styleMarkdownLine(line string) string {
	trimmed := strings.TrimSpace(line)

	// Headers
	if strings.HasPrefix(trimmed, "#") {
		return mdHeaderStyle.Render(trimmed)
	}

	// List items
	if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
		return "  " + renderInlineMarkdown(trimmed)
	}

	// Regular text with inline formatting
	return renderInlineMarkdown(line)
}

// renderMarkdown renders content as styled rows, each no wider than width.
//
// It returns the rows the terminal will actually show. That is the whole point:
// a viewport slicing rows that are already final can report a position that
// matches the screen. The previous version ignored width entirely, so the
// lipgloss box wrapped afterwards and rows were lost inside the visible height.
func renderMarkdown(content string, width int) string {
	var rows []string
	for _, line := range strings.Split(content, "\n") {
		rows = append(rows, wrapStyled(styleMarkdownLine(line), width)...)
	}
	return strings.Join(rows, "\n")
}

func renderInlineMarkdown(line string) string {
	result := line

	// Bold: **text**
	for {
		start := strings.Index(result, "**")
		if start == -1 {
			break
		}
		end := strings.Index(result[start+2:], "**")
		if end == -1 {
			break
		}
		end += start + 2
		bold := result[start+2 : end]
		result = result[:start] + mdBoldStyle.Render(bold) + result[end+2:]
	}

	// Italic: _text_
	for {
		start := strings.Index(result, "_")
		if start == -1 {
			break
		}
		end := strings.Index(result[start+1:], "_")
		if end == -1 {
			break
		}
		end += start + 1
		italic := result[start+1 : end]
		result = result[:start] + mdItalicStyle.Render(italic) + result[end+1:]
	}

	return result
}

// styleYAMLLine applies the styling for a single line of YAML.
func styleYAMLLine(line string) string {
	trimmed := strings.TrimSpace(line)

	// Comment lines
	if strings.HasPrefix(trimmed, "#") {
		return yamlCommentStyle.Render(line)
	}

	// Key: value lines
	if colonIdx := strings.Index(line, ":"); colonIdx > 0 {
		key := line[:colonIdx]
		rest := line[colonIdx:]

		// Inline comment
		if commentIdx := strings.Index(rest, " #"); commentIdx > 0 {
			return yamlKeyStyle.Render(key) +
				yamlValueStyle.Render(rest[:commentIdx]) +
				yamlCommentStyle.Render(rest[commentIdx:])
		}
		return yamlKeyStyle.Render(key) + yamlValueStyle.Render(rest)
	}

	// Plain lines (list items, etc)
	return line
}

// renderYAML renders content as styled rows, each no wider than width.
func renderYAML(content string, width int) string {
	var rows []string
	for _, line := range strings.Split(content, "\n") {
		rows = append(rows, wrapStyled(styleYAMLLine(line), width)...)
	}
	return strings.Join(rows, "\n")
}

func (m model) renderNotImplemented(width int, height int) string {
	return "\n\n" + dimStyle.Render("  Not yet implemented")
}

// truncateContent ensures content is exactly maxLines tall — truncates or pads.
func truncateContent(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	for len(lines) < maxLines {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

// padToHeight ensures the final rendered view is exactly the given height.
func padToHeight(view string, height int) string {
	lines := strings.Split(view, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func (m model) renderPanel(view int, width int, height int, content string) string {
	content = truncateContent(content, height)
	var title string
	switch view {
	case viewDetail:
		if len(m.displayNames) > 0 && m.cursor < len(m.displayNames) {
			title = " " + m.displayNames[m.cursor] + " "
		} else {
			title = " specgetty "
		}
		// Absence of a percentage is itself the signal that nothing is below.
		if pct := m.docScrollPercent(); pct >= 0 {
			indicator := fmt.Sprintf("%d%% ", pct)
			// A title plus indicator wider than the box would make the border
			// run negative, so the name gives way to the position.
			if lipgloss.Width(title)+lipgloss.Width(indicator)+3 > width {
				title = " "
			}
			title += indicator
		}
	case viewLog:
		title = " Log "
	}

	borderColor := lipgloss.Color("240")
	if m.activeView == view {
		borderColor = lipgloss.Color("2")
	}

	border := lipgloss.RoundedBorder()
	titleStyled := lipgloss.NewStyle().Foreground(borderColor).Bold(true).Render(title)
	// The box below renders as width+2 columns: Width(width) plus a border on
	// each side. This line has to match it. It spends three columns on the two
	// corners and the segment before the title, so the trailing run is
	// width-len(title)-1, not -2.
	topBorder := border.TopLeft +
		strings.Repeat(border.Top, 1) +
		titleStyled +
		strings.Repeat(border.Top, max(0, width-lipgloss.Width(title)-1)) +
		border.TopRight

	boxStyle := lipgloss.NewStyle().
		Border(border).
		BorderTop(false).
		BorderForeground(borderColor).
		Width(width).
		Height(height).
		MaxHeight(height + 2) // +2 for border lines

	rendered := topBorder + "\n" + boxStyle.Render(content)
	// Ensure the panel is exactly height+2 lines (title + border top + content area + border bottom)
	return padToHeight(rendered, height+2)
}

func (m model) renderNavBar() string {
	var keys []struct{ key, action string }

	// While an overlay has the keyboard the ordinary keys are unavailable, so
	// showing them would be a lie. Show what the overlay itself accepts.
	if m.askOpenPicker {
		keys = []struct{ key, action string }{
			{"y", "open picker"},
			{"n", "not now"},
			{"q", "quit"},
		}
	} else if m.pickerOpen && m.pickerFocused {
		keys = []struct{ key, action string }{
			{"esc", "clear filter"},
			{"\u2191\u2193/^p^n", "navigate"},
			{"\u23ce", "open project"},
		}
	} else if m.pickerOpen {
		keys = []struct{ key, action string }{
			{"esc", "close"},
			{"\u23ce", "open project"},
			{"/", "search"},
			{"r", "refresh"},
			{"jk/\u2191\u2193", "navigate"},
			{"q", "quit"},
		}
	} else if m.searchFocused {
		keys = []struct{ key, action string }{
			{"esc", "clear filter"},
			{"\u2191\u2193/^p^n", "navigate"},
			{"\u23ce", "open change"},
		}
	} else {
		switch m.level {
		case levelChange:
			keys = []struct{ key, action string }{
				{"q", "quit"},
				{"esc", "back to list"},
				{"\u2190\u2192", "artifact"},
				{"jk/\u2191\u2193", "scroll"},
				{"^f^b", "page"},
				{"^d^u", "half"},
				{"gg/G", "ends"},
				{"p", "projects"},
				{"s", "scan"},
				{"l", "log"},
			}
		case levelProject:
			keys = []struct{ key, action string }{
				{"q", "quit"},
				{"jk/\u2191\u2193", "navigate"},
				{"\u2190\u2192/1-3", "tabs"},
			}
			if m.detailTab == tabChanges {
				keys = append(keys,
					struct{ key, action string }{"\u23ce", "view"},
					struct{ key, action string }{"/", "search"},
					struct{ key, action string }{"f", "mode:" + listModeNames[m.listMode]},
				)
				if r, ok := m.selectedRow(); ok {
					if !r.archived {
						keys = append(keys,
							struct{ key, action string }{"a", "archive"},
							struct{ key, action string }{"d", "discard"})
					}
					keys = append(keys, struct{ key, action string }{"e", "export"})
				}
			}
			keys = append(keys,
				struct{ key, action string }{"p", "projects"},
				struct{ key, action string }{"s", "scan"},
				struct{ key, action string }{"l", "log"},
				struct{ key, action string }{"gg/G", "jump"},
			)
		}
	}

	right := navBarStyle.Render("specgetty " + m.version)

	// Hints are dropped from the end rather than allowed to run underneath the
	// version on the right. A collided nav bar is worse than a short one.
	budget := m.width - lipgloss.Width(right) - 2

	var left strings.Builder
	used := 0
	for i, k := range keys {
		seg := ""
		if i > 0 {
			seg = "  "
		}
		seg += k.key + " " + k.action
		if used+lipgloss.Width(seg) > budget {
			break
		}
		used += lipgloss.Width(seg)
		if i > 0 {
			left.WriteString(navBarStyle.Render("  "))
		}
		left.WriteString(navBarKeyStyle.Render(k.key))
		left.WriteString(navBarStyle.Render(" " + k.action))
	}

	bar := lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		left.String()+strings.Repeat(" ", max(0, m.width-lipgloss.Width(left.String())-lipgloss.Width(right)))+right,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("236")),
	)

	return bar
}

func placeOverlay(width, height int, modal, background string) string {
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceBackground(lipgloss.NoColor{}),
	)
}

func Run(config *scanner.Config, ignoreDirErrors bool, version string, startupPath string, startView string, fields []string) error {
	m := newModel(config, ignoreDirErrors, version)
	if len(fields) > 0 {
		m.fields = fields
	}
	m.startupPath = startupPath
	m.startView = startView

	switch {
	case startView == "all":
		m.pickerOpen = true
		m.pickerLoading = true
	case startupPath != "":
		m.scanning = true
	default:
		// Nothing to show and nothing asked for. Offer the picker rather than
		// opening an empty view with no explanation.
		m.askOpenPicker = true
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	m.program = p
	log.SetOutput(logWriter{program: p})

	_, err := p.Run()
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
