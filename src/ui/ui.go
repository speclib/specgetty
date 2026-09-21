package ui

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
	"github.com/mipmip/specgetty/src/watcher"
)

// There is one panel, so the title it carries has no alternative to choose
// between. The constant survives as the argument renderPanel takes, which keeps
// its signature honest about drawing one named thing.
const viewDetail = 0

// Navigation depth. enter descends, esc ascends. levelProject is the floor:
// the project list is an overlay now, not a level above it.
const (
	levelProject = 0 // one project, with its tab bar
	levelChange  = 1 // one change, with its artifact sub-tabs
	levelSpec    = 2 // one spec, as an outline beside a card
	// levelChangeSpec is a change's spec deltas, as one outline beside a card.
	// Reached from a change rather than from the project, which is why level is
	// a set of views with explicit enter and esc mappings and not a depth.
	levelChangeSpec = 3
)

// specLevel reports whether an outline and a card own the panel.
//
// The two spec levels differ in what they read and what a card may show, and in
// nothing else: the outline, the cursor, the paging and the borders are one
// implementation serving both.
func (m model) specLevel() bool {
	return m.level == levelSpec || m.level == levelChangeSpec
}

// Where the keyboard is. One position, one value.
//
// This used to be two fields, activeView and specsFocus, that every handler had
// to move together. They could describe a state that means nothing, such as the
// the list holding the keyboard while the content also holds it, and the only
// thing preventing it was a paired assignment in five places. Harmless while
// nothing read the pair, which stopped being true once each border is drawn lit
// or dim by asking where the keyboard is.
//
// viewDetail survives above, because which panel gets which title is a
// different question from where the keyboard is.
const (
	focusDetail      = iota // the active tab's content, or an open change's artifact
	focusListPane           // the list half of a split tab
	focusContentPane        // the document half of a split tab
)

// Changes lead, because that is what the tool is usually opened to look at.
const (
	tabChanges    = 0
	tabSpecs      = 1
	tabProperties = 2
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

// The export asks where before it asks whether. The prompt is the confirmation:
// `enter` exports, `esc` cancels, and what would have been confirmed is instead
// editable. Replacing an existing file is the one thing still worth a yes or no.
const (
	exportIdle      = 0
	exportPrompting = 1
	exportReplacing = 2
	exportRunning   = 3
	exportResult    = 4
)

var tabNames = []string{"changes", "specs", "properties"}

// Message types

type scanMsg struct {
	projects scanner.ProjectMap
	err      error
}

// fsChangeMsg is sent when the filesystem watcher detects changes in the openspec directory.
type fsChangeMsg struct{}

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

	// lipgloss v2 counts border and padding inside Style.Width, where v1 added
	// them outside. modalStyle spends two columns on its border and four on its
	// padding, so a modal that wants N columns of text has to ask for N+6 or
	// its content wraps where it used to fit.
	modalChrome = 6

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

	// The selected row of a list that no longer holds the keyboard.
	dimSelectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("238")).
				Foreground(lipgloss.Color("252"))

	// The header's store mark. A chip rather than a glyph: it has to read the
	// same in every terminal font, which a geometric symbol does not.
	storeMarkStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("5")). // magenta
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 1)

	// A store declaration that could not be followed, shown where a tab would
	// otherwise report the project as empty.
	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")) // red
)

type model struct {
	config            *scanner.Config
	ignoreDirErrors   bool
	projects          scanner.ProjectMap
	repoPaths         []string
	displayNames      []string
	cursor            int
	focus             int // focusDetail, focusListPane, focusContentPane
	scanning          bool
	err               error
	spinner           spinner.Model
	docViewport       viewport.Model
	docKey            string
	docLines          []sourceLine // source to screen mapping, when the document has a cursor
	docCursor         int          // index into docLines
	docPath           string       // the file a toggle writes back to
	detailTab         int
	specCursor        int
	changeCursor      int
	changeArtifactTab int
	level             int    // which view is open; see the level constants
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
	fields        []string
	searchInput   textinput.Model
	searchFocused bool
	selectedKey   string // identifies the selected change across re-filter and rescan
	pendingKey    string

	// statusMsg is a transient one-line report shown in place of the nav bar.
	// Every other result in specgetty is a modal that must be dismissed, which
	// is the wrong weight for something instant and harmless.
	statusMsg         string
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
	exportDirInput    textinput.Model
	exportProblem     string
	exportDirName     string
	exportResultMsg   string
	exportResultOk    bool
	exportIsArchived  bool
	width             int
	height            int
	program           *tea.Program
	version           string
	watcher           *watcher.Watcher
	watchedRoot       string
	watchedDirs       []string

	// Which of a comparable node's three views its card is showing. Reset to
	// the difference on every move, a choice being about one node.
	cardView int

	// The spec open at levelSpec, parsed once when it is opened and dropped on
	// the way out. Reparsing on every render would cost a parse per keystroke.
	//
	// The level has two states. A file that fits the grammar gives a tree and
	// no problems; one that does not gives problems and no tree, and the view
	// reports them. specName holds the spec either way, because a report has to
	// name the file it is about.
	specTree     specTree
	specProblems []specProblem
	specName     string
	specNode     int
	specNodePath string

	// Which row of the properties tab is selected, and what is known about the
	// schemas of the project it belongs to.
	propSection int
	schemas     map[string]schemaState
	schemaFor   string
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

	// The text input accepts suggestions and binds tab to them, so completing a
	// path is a matter of handing it what is on disk.
	ei := textinput.New()
	ei.Prompt = ""
	ei.Placeholder = ""
	ei.ShowSuggestions = true

	return model{
		config:          config,
		ignoreDirErrors: ignoreDirErrors,
		version:         version,
		spinner:         s,
		docViewport:     viewport.New(),
		fields:          append([]string(nil), defaultFields...),
		searchInput:     ti,
		pickerInput:     pi,
		exportDirInput:  ei,
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

	case tea.KeyPressMsg:
		// The status line is cleared by the next keystroke rather than by a
		// timer: no tea.Tick, no re-render loop, and it stays exactly as long
		// as the user is still looking at it. Handlers below may set it again.
		m.statusMsg = ""

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
				cmds = append(cmds, doArchiveChange(m.currentRoot(), m.archiveChangeName))
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

		// The export prompt intercepts all keys. Every printable one belongs to
		// the field, so `a`, `d`, `e` and `q` type rather than act: a path
		// called `~/archive` must not archive a change.
		if m.exportState == exportPrompting {
			switch msg.String() {
			case "esc":
				m.exportState = exportIdle
				m.exportDirInput.Blur()
			case "enter":
				m, cmds = m.submitExport(cmds)
			case "tab":
				if s := scanner.CompleteDir(m.exportDirInput.Value()); len(s) > 0 {
					m.exportDirInput.SetSuggestions(s)
				}
				var cmd tea.Cmd
				m.exportDirInput, cmd = m.exportDirInput.Update(msg)
				cmds = append(cmds, cmd)
			default:
				m.exportProblem = ""
				var cmd tea.Cmd
				m.exportDirInput, cmd = m.exportDirInput.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}
		if m.exportState == exportReplacing {
			switch msg.String() {
			case "y":
				m.exportState = exportRunning
				cmds = append(cmds, doExportChange(m.currentRoot(), m.exportDirName,
					m.exportChangeName, m.exportIsArchived, m.exportDir()))
			case "n", "esc":
				// Back to the prompt with the directory still in it, so another
				// can be chosen rather than the whole export lost.
				m.exportState = exportPrompting
				m.exportDirInput.Focus()
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
				cmds = append(cmds, doDiscardChange(m.currentRoot(), m.discardChangeName))
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
			case "pgdown", "ctrl+f":
				m.movePickerCursor(m.pickerPage())
			case "pgup", "ctrl+b":
				m.movePickerCursor(-m.pickerPage())
			case "ctrl+d":
				m.movePickerCursor(max(1, m.pickerPage()/2))
			case "ctrl+u":
				m.movePickerCursor(-max(1, m.pickerPage()/2))
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
				case m.docActive():
					m.docViewport.GotoTop()
				default:
					m.gotoListEnd(false)
				}
				return m, nil
			}
		}

		switch key {
		case "q", "ctrl+c":
			m.stopWatcher()
			return m, tea.Quit

		case "enter":
			if m.specLevel() {
				// Nothing below a spec to descend into.
			} else if m.level == levelChange && m.onChangeSpecsTab() {
				m = m.openChangeSpecs()
			} else if m.level == levelProject && m.detailTab == tabSpecs {
				m = m.openSelectedSpec()
			} else if m.level == levelProject && m.detailTab == tabChanges {
				if _, ok := m.selectedRow(); ok {
					m.rememberSelection()
					m.level = levelChange
					m.changeArtifactTab = 0
				}
			}

		case "esc":
			// levelProject is the floor. The project list is an overlay now,
			// so there is nothing above it to escape to.
			if m.specLevel() {
				if m.level == levelChangeSpec {
					m.level = levelChange
					m.changeArtifactTab = m.changeSpecsTabIndex()
				} else {
					m.level = levelProject
					m.detailTab = tabSpecs
				}
				m.focus = focusListPane
				m.specTree = specTree{}
				m.specProblems = nil
				m.specName = ""
				m.specNode = 0
				m.specNodePath = ""
				m.cardView = viewDiff
			} else if m.level == levelChange {
				m.level = levelProject
				m.syncCursor()
			}

		case "p":
			// The cursor opens on the open project's own row. A project is
			// listed under the directory it was resolved from, which is what
			// the key is, so this matches without translation. A store opened
			// by path has no row: the lookup misses and the cursor stays put,
			// rather than landing on some unrelated project.
			if key := m.currentKey(); key != "" {
				m.pickerKey = key
			}
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

		case "s":
			if len(m.repoPaths) > 0 {
				m.scanning = true
				cmds = append(cmds, m.rescanCurrent())
			}

		case "E":
			var cmd tea.Cmd
			m, cmd = m.openInEditor()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

		case "tab":
			// A split tab has two halves, so tab moves the keyboard between
			// them. Both the specs tab and the properties tab are splits;
			// elsewhere there is one panel and nowhere else for tab to go.
			if m.splitTab() {
				if m.focus == focusListPane {
					m.focus = focusContentPane
				} else {
					m.focus = focusListPane
				}
			}

		case "g":
			m.pendingKey = "g"

		case "G":
			switch {
			case m.docActive():
				m.docViewport.GotoBottom()
			default:
				m.gotoListEnd(true)
			}

		case "pgdown", "ctrl+f":
			// A full page in a document, vim style. Lists keep halfPage(),
			// which is how they have always moved.
			switch {
			case m.docActive():
				m.docViewport.PageDown()
			default:
				m.moveListCursor(m.listPage())
			}

		case "pgup", "ctrl+b":
			switch {
			case m.docActive():
				m.docViewport.PageUp()
			default:
				m.moveListCursor(-m.listPage())
			}

		// 5.3: half page, in a document only.
		case "ctrl+d":
			switch {
			case m.docActive():
				m.docViewport.HalfPageDown()
			default:
				m.moveListCursor(max(1, m.listPage()/2))
			}

		case "ctrl+u":
			switch {
			case m.docActive():
				m.docViewport.HalfPageUp()
			default:
				m.moveListCursor(-max(1, m.listPage()/2))
			}

		case "left":
			// Sub-tabs belong to an open change, the tab bar to the project.
			// Neither spills into the other.
			//
			// Switched on the level rather than chained, so a level added later
			// has to say what these keys do there. The chain this replaced gave
			// every level it did not name the project tab bar by omission,
			// which is how the spec view came to move a tab bar it does not
			// show.
			switch m.level {
			case levelSpec:
				// One spec, no sibling to move to. Its two halves are reached
				// with tab.
			case levelChangeSpec:
				m.moveCardView(-1)
			case levelChange:
				if m.changeArtifactTab > 0 {
					m.changeArtifactTab--
				}
			case levelProject:
				if m.detailTab > 0 {
					m.detailTab--
					m.focus = m.defaultFocus()
					cmds = append(cmds, m.enterTab()...)
				}
			}

		case "right":
			switch m.level {
			case levelSpec:
				// As for left.
			case levelChangeSpec:
				m.moveCardView(1)
			case levelChange:
				if m.changeArtifactTab < m.changeArtifactTabCount()-1 {
					m.changeArtifactTab++
				}
			case levelProject:
				if m.detailTab < len(tabNames)-1 {
					m.detailTab++
					m.focus = m.defaultFocus()
					cmds = append(cmds, m.enterTab()...)
				}
			}

		case "1", "2", "3":
			// Number keys address the project tab bar, so they are inert while
			// a change is open.
			if m.level != levelChange {
				m.detailTab = int(key[0] - '1')
				m.focus = m.defaultFocus()
				cmds = append(cmds, m.enterTab()...)
			}

		case "up", "k":
			if m.docHasCursor() {
				m.moveDocCursor(-1)
			} else if m.docActive() {
				m.docViewport.ScrollUp(1)
			} else if m.specLevel() {
				if m.specNode > 0 {
					m.specNode--
					m.rememberSpecNode()
				}
			} else if m.level != levelChange {
				switch m.detailTab {
				case tabSpecs:
					if m.specCursor > 0 {
						m.specCursor--
					}
				case tabProperties:
					if m.propSection > 0 {
						m.propSection--
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
			if m.docHasCursor() {
				m.moveDocCursor(1)
			} else if m.docActive() {
				m.docViewport.ScrollDown(1)
			} else if m.specLevel() {
				if m.specNode < len(m.specTree.nodes)-1 {
					m.specNode++
					m.rememberSpecNode()
				}
			} else if m.level != levelChange {
				switch m.detailTab {
				case tabSpecs:
					if m.specCursor < len(m.currentSpecNames())-1 {
						m.specCursor++
					}
				case tabProperties:
					if m.propSection < len(m.currentSections())-1 {
						m.propSection++
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

		case "space":
			if m.docHasCursor() {
				msg, cmd := m.toggleSelectedTask()
				m.statusMsg = msg
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}

		case "y":
			if r, ok := m.selectedRow(); ok && m.detailTab == tabChanges {
				m.statusMsg = copyToClipboard("name", r.ci.Name)
			}

		case "Y":
			if r, ok := m.selectedRow(); ok && m.detailTab == tabChanges {
				m.statusMsg = copyToClipboard("path",
					changeDirPath(m.currentRoot(), r))
			}

		case "e":
			if r, ok := m.selectedRow(); ok && m.detailTab == tabChanges {
				m.exportChangeName = r.ci.Name
				m.exportDirName = r.ci.DirName
				m.exportIsArchived = r.archived
				m.exportState = exportPrompting
				m.exportProblem = ""
				m.exportDirInput.SetValue(scanner.DefaultExportDir(m.config))
				m.exportDirInput.SetSuggestions(nil)
				m.exportDirInput.Focus()
				m.exportDirInput.CursorEnd()
			}
		}

	case archiveMsg:
		m.archiveResultOk = msg.ok
		m.archiveResultMsg = msg.output
		m.archiveState = archiveResult
		if msg.ok && len(m.repoPaths) > 0 {
			cmds = append(cmds, m.rescanCurrent())
		}

	case discardMsg:
		m.discardResultOk = msg.ok
		m.discardResultMsg = msg.output
		m.discardState = discardResult
		if msg.ok && len(m.repoPaths) > 0 {
			cmds = append(cmds, m.rescanCurrent())
		}

	case exportMsg:
		m.exportResultOk = msg.ok
		m.exportResultMsg = msg.output
		m.exportState = exportResult

	case fsChangeMsg:
		if m.level >= levelProject && len(m.repoPaths) > 0 {
			m.scanning = true
			cmds = append(cmds, m.rescanCurrent())
			if m.watcher != nil {
				cmds = append(cmds, waitForFsChange(m.watcher))
			}
		}

	case editorDoneMsg:
		var cmd tea.Cmd
		m, cmd = m.editorFinished(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
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
			m.displayNames = projectDisplayNamesFor(m.repoPaths, m.projects)
			if m.cursor >= len(m.repoPaths) {
				m.cursor = max(0, len(m.repoPaths)-1)
			}
			m.recalcLayout()
			// A rescan fires on every file save while watching, so the filter
			// and the selection are preserved across it rather than reset.
			m.syncCursor()
			m.reparseOpenSpec()

			// Start watching the open project after its first scan. A rescan
			// that resolved to a different root, which is what editing a
			// `store:` key does, moves the watch with it rather than leaving
			// it on a tree nothing is read from any more.
			if key := m.currentKey(); key != "" {
				switch {
				case m.watcher == nil:
					if cmd := m.startWatcher(key); cmd != nil {
						cmds = append(cmds, cmd)
					}
				case m.watchedRoot != key || !sameDirs(m.watchedDirs, m.watchDirs(key)):
					// A rescan that resolved elsewhere, which is what editing a
					// `store:` key does, moves the watch with it rather than
					// leaving it on a tree nothing is read from any more.
					m.stopWatcher()
					if cmd := m.startWatcher(key); cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
			}
		}

	case schemaMsg:
		// An answer for a project that is no longer open is dropped rather than
		// shown against the wrong one.
		if msg.project == m.schemaFor {
			if m.schemas == nil {
				m.schemas = make(map[string]schemaState)
			}
			m.schemas[msg.name] = schemaState{details: msg.details, problem: msg.problem}
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
	return disambiguate(paths, nil)
}

// projectDisplayNamesFor names each row the way the user should read it.
//
// A repo is named by its own directory, never by the store it reads from: the
// directory is where the work happens and what the user will look for. Only a
// store opened on its own account, by path rather than from the list, is named
// by the id it declares, because there is no repo to name it after.
func projectDisplayNamesFor(paths []string, projects scanner.ProjectMap) []string {
	preferred := make([]string, len(paths))
	for i, p := range paths {
		st, ok := projects[p]
		if !ok || st.Info.StoreID == "" {
			continue
		}
		if st.Info.ResolvedElsewhere() {
			continue
		}
		preferred[i] = st.Info.StoreID
	}
	return disambiguate(paths, preferred)
}

// disambiguate names each path, preferring the given name where one is set and
// falling back to the basename, then qualifies any name that is not unique
// with its parent directory. A store id colliding with a project name is
// settled the same way two equal basenames are.
func disambiguate(paths []string, preferred []string) []string {
	names := make([]string, len(paths))
	counts := make(map[string]int)

	nameOf := func(i int) string {
		if preferred != nil && preferred[i] != "" {
			return preferred[i]
		}
		return filepath.Base(paths[i])
	}

	for i := range paths {
		counts[nameOf(i)]++
	}

	for i, p := range paths {
		base := nameOf(i)
		if counts[base] > 1 {
			names[i] = base + " (" + filepath.Base(filepath.Dir(p)) + ")"
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
	m.docViewport.SetWidth(docW)
	m.docViewport.SetHeight(docH)

}

// panelContentWidth is how many columns a view drawn inside the main panel has
// to work with.
//
// Every consumer of this number must get the same one. The panel is drawn at
// m.width and spends a column on each border and a column of inset inside each
// border, so its content is four narrower. If a caller derives that separately and drifts, the lipgloss box
// re-wraps the rows and the viewport's line count, and therefore its reported
// position, stops matching the screen. That warning used to sit on docRegion
// and on specsSplit, written twice and enforced nowhere, while three files
// spelled the arithmetic out for themselves.
func (m model) panelContentWidth() int {
	w := m.width - 4
	if w < 1 {
		w = 1
	}
	return w
}

// The content of a tab sits in a border of its own. These are what that border
// costs: a column of border and a column of inset on each side, and a row of
// border top and bottom.
const (
	boxChrome = 4
	boxRows   = 2
)

// contentBoxWidth is how many columns a tab's content has inside its own
// border. panelContentWidth is what the box itself is drawn at.
func (m model) contentBoxWidth() int {
	w := m.panelContentWidth() - boxChrome
	if w < 1 {
		w = 1
	}
	return w
}

// contentBox draws the border around a tab's content. No title: the tab bar
// directly above it already names what is inside.
//
// lit follows containment: a border is drawn in the active colour when the
// region holding the keyboard lies inside it. Borders nest, so the panel's and
// this one can both be lit, and the innermost lit border is the region actually
// receiving the keys.
func contentBox(width, height int, lit bool, content string) string {
	colour := lipgloss.Color("240")
	if lit {
		colour = lipgloss.Color("2")
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colour).
		Padding(0, 1).
		Width(width).
		Height(height).
		MaxHeight(height).
		Render(content)
}

// mainPanelHeight is the panel's height: the terminal, less the nav bar, less
// its own border. Nothing sits between the panel and the nav bar.
func (m model) mainPanelHeight() int {
	remaining := m.height - 1 // -1 for nav bar
	if remaining < 5 {
		return 3
	}
	return remaining - 2 // -border
}

func (m model) halfPage() int {
	return max(1, m.mainPanelHeight()/2)
}

// defaultFocus is where the keyboard lands when a tab becomes active. A split
// tab starts on its list, because that is what chooses what the content shows.
func (m model) defaultFocus() int {
	if m.splitTab() {
		return focusListPane
	}
	return focusDetail
}

func (m model) currentSpecNames() []string {
	if len(m.repoPaths) == 0 {
		return nil
	}
	return m.projects[m.repoPaths[m.cursor]].Info.SpecNames
}

// currentGroups is the change list's two groups for the open project, before
// any search narrows them.
func (m model) currentGroups() []changeGroup {
	key := m.currentKey()
	if key == "" {
		return nil
	}
	return buildGroups(m.projects[key].Info)
}

// allRows is every change in selection order, which is what the cursor indexes.
// The groups decide that order; nothing downstream reorders it.
func (m model) allRows() []changeRow {
	return allRowsOf(m.currentGroups())
}

// currentRows is the single seam every consumer of the change list goes
// through: the renderer, the actions, the confirm modals and the nav bar. The
// grouping order and the search filter both apply here, so nothing downstream
// has to know either exists.
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
	m.propSection = 0
	m.selectedKey = ""
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

// doScanSingle reads one project, resolving startDir to the root its content
// lives in.
//
// The argument is where to start looking, not where to read. For a project
// that reads from a store those are two different directories, and starting
// from the origin every time is what lets an edited `store:` key be picked up
// by the next rescan.
func (m model) doScanSingle(startDir string) tea.Cmd {
	return func() tea.Msg {
		root, st, err := scanner.ScanResolved(startDir)
		if err != nil {
			return scanMsg{err: err}
		}
		if root == "" {
			return scanMsg{projects: scanner.ProjectMap{}}
		}
		return scanMsg{projects: scanner.ProjectMap{root: st}}
	}
}

// rescanCurrent rereads the open project from where its resolution started, so
// that the store declaration is followed again rather than assumed.
func (m model) rescanCurrent() tea.Cmd {
	key := m.currentKey()
	if key == "" {
		return nil
	}
	return m.doScanSingle(m.startDirOf(key))
}

// currentRoot is the directory every filesystem operation on the open project
// targets.
//
// It is the root the content was resolved to, never the repo the user started
// in. The projects map is keyed by root for exactly this reason: archiving,
// discarding and exporting all reach for the key, and a store-backed project
// whose actions reached for the origin instead would create a discarded/
// directory under a repo that holds no changes at all.
func (m model) currentRoot() string {
	key := m.currentKey()
	if key == "" {
		return ""
	}
	if st, ok := m.projects[key]; ok && st.Info.Root != "" {
		return st.Info.Root
	}
	return key
}

// sameDirs reports whether two watch sets are the same list in the same order,
// which is how a re-resolution that moved the store is noticed.
func sameDirs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// currentKey identifies the open project: the directory it was resolved from,
// which is what the picker lists and what the header names. It is the root only
// when no store declaration was followed.
func (m model) currentKey() string {
	if len(m.repoPaths) == 0 || m.cursor >= len(m.repoPaths) {
		return ""
	}
	return m.repoPaths[m.cursor]
}

// watchDirs is every tree the project under key depends on.
//
// One for an ordinary project. Two when the content comes from a store: the
// store's, where specs and changes move, and the originating repo's, which
// holds a single file whose `store:` key decides which store is read at all.
// Watching only the first would show stale content with no sign that anything
// had happened.
func (m model) watchDirs(key string) []string {
	root := key
	if st, ok := m.projects[key]; ok && st.Info.Root != "" {
		root = st.Info.Root
	}
	dirs := []string{filepath.Join(root, "openspec")}
	if key != root {
		dirs = append(dirs, filepath.Join(key, "openspec"))
	}
	return dirs
}

// startDirOf returns the directory a project's resolution began at. That is the
// key itself now that projects are filed under where the reading started, and
// the lookup survives only for a map that predates a scan.
func (m model) startDirOf(key string) string {
	if st, ok := m.projects[key]; ok && st.Info.Origin != "" {
		return st.Info.Origin
	}
	return key
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

// startWatcher watches every tree the open project depends on.
//
// That is one tree for an ordinary project, and two when the content comes
// from a store: the store's, where specs and changes move, and the originating
// repo's, which holds the declaration that decides which store is read at all.
func (m *model) startWatcher(key string) tea.Cmd {
	w, err := watcher.New(m.watchDirs(key)...)
	if err != nil {
		log.Printf("watcher: failed to start: %v", err)
		return nil
	}
	m.watcher = w
	m.watchedRoot = key
	m.watchedDirs = m.watchDirs(key)
	return waitForFsChange(w)
}

func (m *model) stopWatcher() {
	if m.watcher != nil {
		m.watcher.Close()
		m.watcher = nil
	}
	m.watchedRoot = ""
	m.watchedDirs = nil
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

// doArchiveChange runs the OpenSpec CLI in the resolved root.
//
// For a store-backed project that is the store, not the repo the user started
// in. The CLI would resolve the declaration for itself either way, but running
// it where the content is keeps this command honest about what it touches.
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

// exportFileName is the name of the zip, which is generated and never typed.
// It strips an archived change's date prefix and stamps the export date, and
// both do work that retyping would lose.
func exportFileName(semanticName string) string {
	return semanticName + "-" + time.Now().Format("2006-01-02") + ".zip"
}

// exportDir is the directory the prompt currently holds, expanded.
func (m model) exportDir() string {
	return scanner.ExpandPath(m.exportDirInput.Value())
}

// submitExport answers `enter` at the prompt: refuse what cannot be used, ask
// before replacing, otherwise write.
func (m model) submitExport(cmds []tea.Cmd) (model, []tea.Cmd) {
	dir := m.exportDir()
	if p := scanner.CheckExportDir(dir); p != nil {
		// The typed text stays in the field to be corrected.
		m.exportProblem = p.Detail
		return m, cmds
	}
	if _, err := os.Stat(filepath.Join(dir, exportFileName(m.exportChangeName))); err == nil {
		m.exportState = exportReplacing
		m.exportDirInput.Blur()
		return m, cmds
	}
	m.exportState = exportRunning
	m.exportDirInput.Blur()
	return m, append(cmds, doExportChange(m.currentRoot(), m.exportDirName,
		m.exportChangeName, m.exportIsArchived, dir))
}

// doExportChange zips a change directory.
//
// dirName is the directory on disk, which for an archived change still carries
// its date prefix. semanticName is the display name, which does not, and which
// names both the zip's root folder and the zip file itself. Passing the display
// name as the source path is the defect this signature exists to prevent.
func doExportChange(projectPath string, dirName string, semanticName string, isArchived bool, destDir string) tea.Cmd {
	return func() tea.Msg {
		var srcDir string
		if isArchived {
			srcDir = filepath.Join(projectPath, "openspec", "changes", "archive", dirName)
		} else {
			srcDir = filepath.Join(projectPath, "openspec", "changes", dirName)
		}

		if _, err := os.Stat(srcDir); err != nil {
			return exportMsg{ok: false, output: fmt.Sprintf("Source not found: %s", srcDir)}
		}

		destPath := filepath.Join(destDir, exportFileName(semanticName))

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

// View returns the frame plus the terminal features it needs.
//
// In v2 the alternate screen is a property of the view rather than a program
// option, so it is declared here on every render. Getting this wrong means
// specgetty draws over the scrollback instead of taking its own screen, and no
// test can see that: what is asserted below is that the flag is set, not that
// the terminal honoured it.
func (m model) View() tea.View {
	v := tea.NewView(m.renderFrame())
	v.AltScreen = true
	return v
}

// renderFrame builds the whole screen as a styled string.
func (m model) renderFrame() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	if m.width < 60 || m.height < 20 {
		return "Terminal too small. Need at least 60x20."
	}

	panelH := m.mainPanelHeight()

	detailContent := m.renderDetailPanel(m.panelContentWidth(), panelH)
	mainRow := m.renderPanel(viewDetail, m.width, panelH, detailContent)

	// Nav bar
	navBar := m.renderNavBar()

	view := lipgloss.JoinVertical(lipgloss.Left, mainRow, navBar)

	// Each of these replaces the frame rather than being drawn over it, so the
	// order below is precedence, not layering: whichever runs last is what the
	// user sees. Only one can be up at a time in practice, because whatever
	// holds the keyboard refuses the keys that would raise another.
	if m.pickerOpen {
		view = modalFrame(m.width, m.height, m.renderPicker())
	}

	if m.askOpenPicker {
		modal := modalStyle.Width(modalWidth(m, 56)).Render(
			"No OpenSpec project here.\n\nOpen the project picker? (y/n)")
		view = modalFrame(m.width, m.height, modal)
	}

	// Modal overlays
	if m.scanning {
		modal := modalStyle.Width(modalWidth(m, 40)).Render(m.spinner.View() + " Scanning for OpenSpec sources...")
		view = modalFrame(m.width, m.height, modal)
	}
	if m.err != nil {
		errText := fmt.Sprintf("Error: %v", m.err)
		modal := modalStyle.Width(modalWidth(m, m.width*3/4)).Render(errText)
		view = modalFrame(m.width, m.height, modal)
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
		modal := modalStyle.Width(modalWidth(m, 50)).Render(content)
		view = modalFrame(m.width, m.height, modal)
	case archiveRunning:
		modal := modalStyle.Width(modalWidth(m, 40)).Render(m.spinner.View() + " Archiving...")
		view = modalFrame(m.width, m.height, modal)
	case archiveResult:
		var prefix string
		if m.archiveResultOk {
			prefix = "✓ "
		} else {
			prefix = "✗ "
		}
		content := prefix + m.archiveResultMsg + "\n\nPress any key to dismiss."
		modal := modalStyle.Width(modalWidth(m, m.width*3/4)).Render(content)
		view = modalFrame(m.width, m.height, modal)
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
		modal := modalStyle.Width(modalWidth(m, 50)).Render(content)
		view = modalFrame(m.width, m.height, modal)
	case discardRunning:
		modal := modalStyle.Width(modalWidth(m, 40)).Render(m.spinner.View() + " Discarding...")
		view = modalFrame(m.width, m.height, modal)
	case discardResult:
		var prefix string
		if m.discardResultOk {
			prefix = "✓ "
		} else {
			prefix = "✗ "
		}
		content := prefix + m.discardResultMsg + "\n\nPress any key to dismiss."
		modal := modalStyle.Width(modalWidth(m, m.width*3/4)).Render(content)
		view = modalFrame(m.width, m.height, modal)
	}

	// Export modals
	switch m.exportState {
	case exportPrompting:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("Export %q", m.exportChangeName))
		b.WriteString("\n\n→ " + m.exportDirInput.View() + "/")
		b.WriteString("\n  " + dimStyle.Render(exportFileName(m.exportChangeName)))
		if m.exportProblem != "" {
			b.WriteString("\n\n" + warnStyle.Render(m.exportProblem))
		}
		// A field implies no keys, so the modal names them.
		b.WriteString("\n\n" + dimStyle.Render("tab completes   ⏎ export   esc cancel"))
		modal := modalStyle.Width(modalWidth(m, 60)).Render(b.String())
		view = modalFrame(m.width, m.height, modal)
	case exportReplacing:
		content := fmt.Sprintf("%s already exists.\n\nReplace it? (y/n)",
			filepath.Join(m.exportDir(), exportFileName(m.exportChangeName)))
		modal := modalStyle.Width(modalWidth(m, 60)).Render(content)
		view = modalFrame(m.width, m.height, modal)
	case exportRunning:
		modal := modalStyle.Width(modalWidth(m, 40)).Render(m.spinner.View() + " Exporting...")
		view = modalFrame(m.width, m.height, modal)
	case exportResult:
		var prefix string
		if m.exportResultOk {
			prefix = "✓ "
		} else {
			prefix = "✗ "
		}
		content := prefix + m.exportResultMsg + "\n\nPress any key to dismiss."
		modal := modalStyle.Width(modalWidth(m, m.width*3/4)).Render(content)
		view = modalFrame(m.width, m.height, modal)
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
		return "\n" + dimStyle.Render("No project selected. Press p to pick one.")
	}

	// A spec fills the panel on its own, the way an open change does, whether
	// it opened as an outline or as a report.
	if m.specLevel() {
		return m.renderSpecDetail(width, height)
	}

	// An open change fills the panel on its own: no project header, no tab bar,
	// so its artifact sub-tabs own the full width and their own key axis.
	if m.level == levelChange {
		if r, ok := m.selectedRow(); ok {
			return m.renderChangeDetail(r, m.changeArtifactTab, width, height)
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
	headerContent := headerStyle.Render(headerName(info, currentProject)) +
		storeMark(info) + "\n" + dimStyle.Render(statsLine)
	// Vertical only. The horizontal half is the panel's job now, and keeping
	// both would put the header two columns in while everything under it sits
	// at one.
	b.WriteString(lipgloss.NewStyle().Padding(1, 0).Render(headerContent))
	b.WriteString("\n")

	b.WriteString(m.renderTabHeader(width))
	b.WriteString("\n")

	// The region below the tab bar, which is what the box is drawn at. The
	// header takes four rows (two of content, two of padding) and the tab bar
	// one.
	boxHeight := height - 5

	if boxHeight < boxRows+1 {
		boxHeight = boxRows + 1
	}
	inner := boxHeight - boxRows
	w := m.contentBoxWidth()

	// The panel always holds the keyboard, so a single box is always lit. A
	// split tab has two and decides for itself which one is.
	lit := true

	switch m.detailTab {
	case tabSpecs:
		b.WriteString(m.renderSpecsTab(width, boxHeight))
	case tabProperties:
		b.WriteString(m.renderPropertiesTab(width, boxHeight))
	case tabChanges:
		b.WriteString(contentBox(width, boxHeight, lit, m.renderChangesTab(w, inner)))
	default:
		b.WriteString(contentBox(width, boxHeight, lit, m.renderNotImplemented(w, inner)))
	}

	return b.String()
}

// renderChangesTab draws the full-width change table, plus the search prompt
// whenever a filter is active or being typed.
func (m model) renderChangesTab(width int, height int) string {
	// A store declaration that could not be followed leaves no changes to
	// list, which is indistinguishable from a project that has none. The
	// reason goes where the emptiness shows.
	if len(m.repoPaths) > 0 && m.cursor < len(m.repoPaths) {
		if line := storeProblemLine(m.projects[m.repoPaths[m.cursor]].Info); line != "" {
			return warnStyle.Render(line) + "\n\n" +
				dimStyle.Render("See the config tab's store details.")
		}
	}

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
	case total == 0 && len(m.currentGroups()) == 0:
		body = dimStyle.Render("No project selected.")
	default:
		body = renderGroupedTable(m.currentGroups(), rows, changeFieldDefs(m.fields),
			m.changeCursor, width, tableHeight)
		// The groups stay, each counting zero, so the shape of the list does
		// not change under a query. The message says why they are empty, which
		// the counts alone do not.
		if len(rows) == 0 && m.searchInput.Value() != "" {
			body += "\n" + noMatchMessage("changes", m.searchInput.Value())
		}
	}

	if !showPrompt {
		return body
	}
	return truncateContent(body, tableHeight) + "\n" +
		renderSearchPrompt(m.searchInput.Value(), m.searchFocused, len(rows), total, width)
}

// specsSplit divides the specs tab into its list half and its content half.
//
// renderSpecsTab and docRegion both need this, and they must agree: if the
// document is wrapped to a different width than the pane it is drawn in, the
// box re-wraps the rows and the reported position stops matching the screen.
func specsSplit(width int) (listWidth, contentWidth int) {
	listWidth = width * 3 / 10
	if listWidth < 15 {
		listWidth = 15
	}
	return listWidth, width - listWidth - 1 // 1 for the gap
}

// renderSpecsTab draws the two halves, each in its own border. width and height
// are the whole region below the tab bar, because this function owns the chrome
// rather than being handed the space inside it.
func (m model) renderSpecsTab(width int, height int) string {
	const lit = true
	if len(m.repoPaths) == 0 {
		return contentBox(width, height, lit, "No project selected.")
	}

	info := m.projects[m.repoPaths[m.cursor]].Info

	if len(info.SpecNames) == 0 {
		// Nothing to split, so nothing to tell apart: one box. An unfollowed
		// store declaration says why there is nothing, rather than letting the
		// emptiness stand as the whole report.
		body := dimStyle.Render("No specs found")
		if line := storeProblemLine(info); line != "" {
			body = warnStyle.Render(line) + "\n\n" +
				dimStyle.Render("See the config tab's store details.")
		}
		return contentBox(width, height, lit, body)
	}

	listOuter, contentOuter := specsSplit(width)
	// The document is already wrapped to the right box's inside by docRegion,
	// which is why only the list needs its own width here.
	listWidth := listOuter - boxChrome
	rows := height - boxRows
	if rows < 1 {
		rows = 1
	}

	// Render spec list
	var listB strings.Builder
	offset := 0
	if m.specCursor >= rows {
		offset = m.specCursor - rows + 1
	}
	end := offset + rows
	if end > len(info.SpecNames) {
		end = len(info.SpecNames)
	}

	for i := offset; i < end; i++ {
		if i > offset {
			listB.WriteString("\n")
		}
		name := info.SpecNames[i]
		switch {
		case i == m.specCursor && m.focus == focusListPane:
			listB.WriteString(selectedStyle.Width(listWidth).Render(name))
		case i == m.specCursor && m.focus == focusContentPane:
			// Still the selected spec, but the keys belong to the content now.
			listB.WriteString(dimSelectedStyle.Width(listWidth).Render(name))
		default:
			listB.WriteString(normalStyle.Render(name))
		}
	}

	specList := truncateContent(listB.String(), rows)
	specContentStr := truncateContent(m.docViewport.View(), rows)

	// Each half says for itself whether the keyboard is in it. Exactly one of
	// them always does, there being nowhere else for it to be.
	leftBox := contentBox(listOuter, height, m.focus == focusListPane, specList)
	rightBox := contentBox(contentOuter, height, m.focus == focusContentPane, specContentStr)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
}

// enterTab is what a tab costs when it is selected.
//
// Only the properties tab costs anything: its schema rows need a definition
// located, which is a subprocess. Nothing runs until the tab is asked for, so a
// session that never opens it never pays.
func (m *model) enterTab() []tea.Cmd {
	if m.level != levelProject || m.detailTab != tabProperties {
		return nil
	}
	return m.ensureSchemasLoaded()
}

// splitTab reports whether the active tab draws a list beside its content.
//
// The specs tab and the properties tab are both splits, and everything about
// focus, borders and the vertical keys asks this rather than naming a tab.
func (m model) splitTab() bool {
	if m.specLevel() {
		// A report is one panel, so tab has nowhere to go in that state.
		return m.specStructured()
	}
	if m.level != levelProject {
		return false
	}
	switch m.detailTab {
	case tabSpecs:
		return len(m.currentSpecNames()) > 0
	case tabProperties:
		return len(m.currentSections()) > 0
	}
	return false
}

// currentSections returns the open project's properties rows.
func (m model) currentSections() []propSection {
	key := m.currentKey()
	if key == "" {
		return nil
	}
	return propSections(m.projects[key].Info)
}

// sectionIndex clamps the selected row to the set this project has, so
// switching from a project with three schemas to one with a single schema
// cannot land past the end.
func (m model) sectionIndex(sections []propSection) int {
	if len(sections) == 0 {
		return 0
	}
	if m.propSection < 0 {
		return 0
	}
	if m.propSection >= len(sections) {
		return len(sections) - 1
	}
	return m.propSection
}

// renderPropertiesTab draws the row list and the content beside it, each in its
// own border, by the same split the specs tab uses.
func (m model) renderPropertiesTab(width int, height int) string {
	const lit = true
	sections := m.currentSections()
	if len(sections) == 0 {
		return contentBox(width, height, lit, dimStyle.Render("No project selected."))
	}

	listOuter, contentOuter := m.propertiesSplit(width)
	listWidth := listOuter - boxChrome
	rows := height - boxRows
	if rows < 1 {
		rows = 1
	}

	active := m.sectionIndex(sections)
	var listB strings.Builder
	offset := 0
	if active >= rows {
		offset = active - rows + 1
	}
	end := offset + rows
	if end > len(sections) {
		end = len(sections)
	}
	for i := offset; i < end; i++ {
		if i > offset {
			listB.WriteString("\n")
		}
		label := sections[i].label
		switch {
		case i == active && m.focus == focusListPane:
			listB.WriteString(selectedStyle.Width(listWidth).Render(label))
		case i == active:
			// Still the selected row, but the keys belong to the content now.
			listB.WriteString(dimSelectedStyle.Width(listWidth).Render(label))
		default:
			listB.WriteString(normalStyle.Render(label))
		}
	}

	list := truncateContent(listB.String(), rows)
	content := truncateContent(m.docViewport.View(), rows)

	leftBox := contentBox(listOuter, height, m.focus == focusListPane, list)
	rightBox := contentBox(contentOuter, height, m.focus == focusContentPane, content)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
}

// propertiesSplit sizes the row list to its labels rather than to a share of
// the panel.
//
// The specs tab gives its list thirty percent because capability names are
// long. These labels are short words and a schema name, and the content beside
// them carries absolute paths, which is the thing that suffers from a narrow
// column.
func (m model) propertiesSplit(width int) (listOuter, contentOuter int) {
	widest := 0
	for _, s := range m.currentSections() {
		if n := len([]rune(s.label)); n > widest {
			widest = n
		}
	}
	listOuter = widest + boxChrome
	if min := 9 + boxChrome; listOuter < min {
		listOuter = min
	}
	// Never let the list crowd out the content it exists to label.
	if max := width / 3; listOuter > max {
		listOuter = max
	}
	return listOuter, width - listOuter - 1
}

var (
	mdHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6")) // cyan

	mdBoldStyle = lipgloss.NewStyle().Bold(true)

	mdItalicStyle = lipgloss.NewStyle().Italic(true)

	// A backticked span. Specs name capabilities, keys and paths this way.
	mdCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan

	// A scenario's keyword, which opens every clause of a card.
	specKeywordStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("3")) // yellow

	// What a change does to a requirement, and what one of its scenarios does
	// to the one it restates. Green, yellow and red are what a diff has always
	// used for these three, so the marks read before the legend is found.
	opAddedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	opModifiedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	opRemovedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red

	// A scenario a modified requirement restates without changing it. Drawn
	// back rather than drawn out: "you have read this before" is what the
	// reader needs from it, and half of a modified requirement is this.
	unchangedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	// The capability a delta belongs to, rooting its requirements.
	capabilityStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6")) // cyan

	// A word the change adds, and one it drops. Struck through rather than
	// only coloured, so the two read apart without relying on colour.
	addedWordStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	removedWordStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")). // red
				Strikethrough(true)

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

	// Task checkboxes, drawn as boxes rather than as their punctuation. The
	// prefixes matched here are exactly the ones the scanner counts, so the
	// glyphs and the totals can never disagree about what a task is.
	if rest, ok := strings.CutPrefix(line, taskUncheckedPrefix); ok {
		return "  " + checkboxUnchecked + " " + renderInlineMarkdown(rest)
	}
	if rest, ok := strings.CutPrefix(line, taskCheckedPrefix); ok {
		return "  " + checkboxChecked + " " + renderInlineMarkdown(rest)
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
// sourceLine records where one line of the source ended up on screen.
//
// The renderer wraps as it goes, so a source line becomes one to four screen
// rows and the connection between them is otherwise lost. Two things need it
// back: a highlight has to cover every row a line produced, and a toggle has to
// know the exact text of the line it is about to rewrite.
type sourceLine struct {
	index    int    // position in the source, counting from zero
	text     string // the source line verbatim, which is what a save matches on
	rowStart int    // first screen row it produced
	rowEnd   int    // last screen row it produced, inclusive
}

// renderMarkdownLines renders content and reports where each source line landed.
func renderMarkdownLines(content string, width int) (string, []sourceLine) {
	var rows []string
	var lines []sourceLine

	for i, line := range strings.Split(content, "\n") {
		start := len(rows)
		rows = append(rows, wrapStyled(styleMarkdownLine(line), width)...)
		lines = append(lines, sourceLine{
			index:    i,
			text:     line,
			rowStart: start,
			rowEnd:   len(rows) - 1,
		})
	}
	return strings.Join(rows, "\n"), lines
}

// renderMarkdown renders content without reporting the mapping, for the panes
// that have no cursor and do not need it.
func renderMarkdown(content string, width int) string {
	out, _ := renderMarkdownLines(content, width)
	return out
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

	// Code: `text`. Specs name capabilities, keys and files in backticks
	// constantly, and leaving the marks in makes them the loudest punctuation
	// on the line. An unclosed backtick is left alone: the rest of the line is
	// prose and deserves to survive.
	for {
		start := strings.Index(result, "`")
		if start == -1 {
			break
		}
		end := strings.Index(result[start+1:], "`")
		if end == -1 {
			break
		}
		end += start + 1
		code := result[start+1 : end]
		result = result[:start] + mdCodeStyle.Render(code) + result[end+1:]
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

// renderPanel draws a titled box.
//
// width is the TOTAL width of the box including its borders. lipgloss v2
// changed this: Style.Width now covers border and padding, where v1 added them
// outside. Every caller therefore passes the terminal width rather than the
// terminal width minus two, and the content handed in must already be sized to
// width-2.
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
	}

	// The panel always holds the keyboard: there is nowhere else for it to be.
	borderColor := lipgloss.Color("2")

	border := lipgloss.RoundedBorder()
	titleStyled := lipgloss.NewStyle().Foreground(borderColor).Bold(true).Render(title)
	// The box renders as exactly `width` columns. This line has to match it,
	// and spends three of them on the two corners and the segment before the
	// title.
	topBorder := border.TopLeft +
		strings.Repeat(border.Top, 1) +
		titleStyled +
		strings.Repeat(border.Top, max(0, width-lipgloss.Width(title)-3)) +
		border.TopRight

	boxStyle := lipgloss.NewStyle().
		Border(border).
		BorderTop(false).
		BorderForeground(borderColor).
		// One column of air inside each border, so no view has to indent
		// itself and none of them disagree about how far. In lipgloss v2
		// Width is the total, so the padding comes out of the content.
		Padding(0, 1).
		Width(width).
		Height(height).
		MaxHeight(height + 2) // +2 for border lines

	rendered := topBorder + "\n" + boxStyle.Render(content)
	// Ensure the panel is exactly height+2 lines (title + border top + content area + border bottom)
	return padToHeight(rendered, height+2)
}

func (m model) renderNavBar() string {
	// A status message takes the nav bar's row. The keys it would have listed
	// are all still bound; the message is gone on the next keystroke.
	if m.statusMsg != "" {
		return lipgloss.PlaceHorizontal(
			m.width,
			lipgloss.Left,
			navBarStyle.Render(" "+m.statusMsg+" "),
			lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(lipgloss.Color("236"))),
		)
	}

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
		case levelSpec, levelChangeSpec:
			back := "back to specs"
			if m.level == levelChangeSpec {
				back = "back to change"
			}
			if !m.specStructured() {
				// A report offers what a report can do: read it, open the file,
				// leave. None of the navigation keys apply.
				keys = []struct{ key, action string }{
					{"q", "quit"},
					{"esc", back},
					{"jk/\u2191\u2193", "scroll"},
					{"^f^b", "page"},
					{"p", "projects"},
					{"s", "scan"},
				}
				break
			}
			keys = []struct{ key, action string }{
				{"q", "quit"},
				{"esc", back},
				{"jk/\u2191\u2193", "navigate"},
			}
			// Advertised only where it does something, which is a node with an
			// original to show beside itself.
			if m.cardViewRowShown() {
				keys = append(keys,
					struct{ key, action string }{"\u2190\u2192", "diff/old/new"})
			}
			if m.focus == focusContentPane {
				keys = append(keys,
					struct{ key, action string }{"tab", "focus outline"},
					struct{ key, action string }{"^f^b", "page"},
					struct{ key, action string }{"^d^u", "half"},
					struct{ key, action string }{"gg/G", "ends"})
			} else {
				keys = append(keys,
					struct{ key, action string }{"tab", "focus card"})
			}
			keys = append(keys,
				struct{ key, action string }{"p", "projects"},
				struct{ key, action string }{"s", "scan"})
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
			}
		case levelProject:
			keys = []struct{ key, action string }{
				{"q", "quit"},
				{"jk/\u2191\u2193", "navigate"},
				{"\u2190\u2192/1-3", "tabs"},
			}
			if m.detailTab == tabProperties && m.splitTab() {
				if m.focus == focusContentPane {
					keys = append(keys,
						struct{ key, action string }{"tab", "focus list"},
						struct{ key, action string }{"^f^b", "page"},
						struct{ key, action string }{"^d^u", "half"},
						struct{ key, action string }{"gg/G", "ends"})
				} else {
					keys = append(keys,
						struct{ key, action string }{"tab", "focus content"})
				}
			}
			if m.detailTab == tabSpecs && len(m.currentSpecNames()) > 0 {
				if m.focus == focusContentPane {
					keys = append(keys,
						struct{ key, action string }{"tab", "focus list"},
						struct{ key, action string }{"^f^b", "page"},
						struct{ key, action string }{"^d^u", "half"},
						struct{ key, action string }{"gg/G", "ends"})
				} else {
					keys = append(keys,
						struct{ key, action string }{"tab", "focus content"})
				}
			}
			if m.detailTab == tabChanges {
				keys = append(keys,
					struct{ key, action string }{"\u23ce", "view"},
					struct{ key, action string }{"/", "search"},
				)
				if r, ok := m.selectedRow(); ok {
					if !r.archived {
						keys = append(keys,
							struct{ key, action string }{"a", "archive"},
							struct{ key, action string }{"d", "discard"})
					}
					keys = append(keys,
						struct{ key, action string }{"e", "export"},
						struct{ key, action string }{"y/Y", "copy name/path"})
				}
			}
			keys = append(keys,
				struct{ key, action string }{"p", "projects"},
				struct{ key, action string }{"s", "scan"},
				struct{ key, action string }{"gg/G", "jump"},
			)
		}

		// One rule for every level: the key is listed exactly where the pane is
		// showing a file, which is the same fact the key itself asks. A pane
		// showing several files or an assembled report has no path and gets no
		// hint.
		if m.docPath != "" {
			keys = append(keys, struct{ key, action string }{"E", "edit"})
		}
	}

	right := navBarStyle.Render("specgetty " + m.version)

	// Hints are dropped from the end rather than allowed to run underneath the
	// version on the right. A collided nav bar is worse than a short one.
	// Two for the gutters at each end, two for the gap before the version.
	budget := m.width - lipgloss.Width(right) - 4

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

	// A gutter at each end to match the panel, and the fill between the hints
	// and the version rendered through navBarStyle. Bare spaces here left an
	// unpainted hole in the middle of a strip whose whole job is to mark the
	// bottom edge of the screen.
	gutter := navBarStyle.Render(" ")
	fill := max(0, m.width-2-lipgloss.Width(left.String())-lipgloss.Width(right))
	bar := lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		gutter+left.String()+navBarStyle.Render(strings.Repeat(" ", fill))+right+gutter,
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(lipgloss.Color("236"))),
	)

	return bar
}

// modalFrame builds a whole frame holding the modal and nothing else, centred
// on blank space. It does not draw over the view it replaces, and it never took
// a background to draw over: it used to accept one and discard it, which read
// as compositing to everyone who saw the call.
// modalWidth is the width a modal asks for, clamped to what the terminal has.
//
// `modal-presentation` requires a modal never to push the frame out of shape,
// and a fixed width does exactly that on a narrow terminal: the export prompt
// asked for 66 columns and the startup question for 62, both wider than the
// 60-column minimum.
func modalWidth(m model, want int) int {
	max := m.width - 2
	if want+modalChrome > max {
		if max < 20 {
			max = 20
		}
		return max
	}
	return want + modalChrome
}

func modalFrame(width, height int, modal string) string {
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle()),
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
	// The alternate screen is declared by View now, not here.
	p := tea.NewProgram(m)

	m.program = p
	// Discarded rather than left alone. The standard logger writes to stderr by
	// default, which under a full-screen interface means writing over the
	// frame, and the scanner logs a line per project on every scan. `--debug`
	// never reaches here, so it keeps the output it has always printed.
	log.SetOutput(io.Discard)

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
