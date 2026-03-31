package ui

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mipmip/specgetty/src/scanner"
)

const (
	viewProjects = 0
	viewDetail   = 1
	viewLog      = 2
)

const (
	tabOverview = 0
	tabSpecs    = 1
	tabChanges  = 2
	tabConfig   = 3
	tabSearch   = 4
)

var tabNames = []string{"overview", "specs", "changes", "config", "search"}

// Message types

type scanMsg struct {
	projects scanner.ProjectMap
	err      error
}

type logMsg string

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
	config          *scanner.Config
	ignoreDirErrors bool
	projects        scanner.ProjectMap
	repoPaths       []string
	displayNames    []string
	cursor          int
	activeView      int
	scanning        bool
	err             error
	spinner         spinner.Model
	detailViewport  viewport.Model
	logViewport     viewport.Model
	logContent      string
	detailTab       int
	fileCursor      int
	filePaths       []string
	logVisible      bool
	logShownOnce    bool
	pendingKey      string
	width           int
	height          int
	program         *tea.Program
	version         string
}

func newModel(config *scanner.Config, ignoreDirErrors bool, version string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		config:          config,
		ignoreDirErrors: ignoreDirErrors,
		scanning:        true,
		version:         version,
		spinner:         s,
		detailViewport:  viewport.New(0, 0),
		logViewport:     viewport.New(0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.doScan(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		key := msg.String()
		if m.pendingKey == "g" {
			m.pendingKey = ""
			if key == "g" {
				switch m.activeView {
				case viewProjects:
					m.cursor = 0
					m.updateFileList()
				case viewDetail:
					m.fileCursor = 0
				case viewLog:
					m.logViewport.GotoTop()
				}
				return m, nil
			}
		}

		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "s":
			m.scanning = true
			cmds = append(cmds, m.doScan())
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
					m.activeView = viewProjects
				}
				m.recalcLayout()
			}
		case "tab":
			if m.logVisible {
				m.activeView = (m.activeView + 1) % 3
			} else {
				if m.activeView == viewProjects {
					m.activeView = viewDetail
				} else {
					m.activeView = viewProjects
				}
			}
		case "g":
			m.pendingKey = "g"
		case "G":
			switch m.activeView {
			case viewProjects:
				if len(m.repoPaths) > 0 {
					m.cursor = len(m.repoPaths) - 1
					m.updateFileList()
				}
			case viewDetail:
				if len(m.filePaths) > 0 {
					m.fileCursor = len(m.filePaths) - 1
				}
			case viewLog:
				m.logViewport.GotoBottom()
			}
		case "pgdown", "ctrl+f":
			half := m.halfPage()
			switch m.activeView {
			case viewProjects:
				m.cursor = min(m.cursor+half, len(m.repoPaths)-1)
				m.updateFileList()
			case viewDetail:
				if len(m.filePaths) > 0 {
					m.fileCursor = min(m.fileCursor+half, len(m.filePaths)-1)
				}
			case viewLog:
				m.logViewport.LineDown(half)
			}
		case "pgup", "ctrl+b":
			half := m.halfPage()
			switch m.activeView {
			case viewProjects:
				m.cursor = max(m.cursor-half, 0)
				m.updateFileList()
			case viewDetail:
				if len(m.filePaths) > 0 {
					m.fileCursor = max(m.fileCursor-half, 0)
				}
			case viewLog:
				m.logViewport.LineUp(half)
			}
		case "left":
			if m.activeView == viewDetail {
				if m.detailTab > 0 {
					m.detailTab--
				}
			}
		case "right":
			if m.activeView == viewDetail {
				if m.detailTab < len(tabNames)-1 {
					m.detailTab++
				}
			}
		case "1":
			if m.activeView == viewDetail {
				m.detailTab = tabOverview
			}
		case "2":
			if m.activeView == viewDetail {
				m.detailTab = tabSpecs
			}
		case "3":
			if m.activeView == viewDetail {
				m.detailTab = tabChanges
			}
		case "4":
			if m.activeView == viewDetail {
				m.detailTab = tabConfig
			}
		case "5":
			if m.activeView == viewDetail {
				m.detailTab = tabSearch
			}
		case "up", "k":
			if m.activeView == viewProjects {
				if m.cursor > 0 {
					m.cursor--
					m.updateFileList()
				}
			} else if m.activeView == viewDetail {
				if len(m.filePaths) > 0 && m.fileCursor > 0 {
					m.fileCursor--
				}
			} else if m.activeView == viewLog {
				m.logViewport.LineUp(1)
			}
		case "down", "j":
			if m.activeView == viewProjects {
				if m.cursor < len(m.repoPaths)-1 {
					m.cursor++
					m.updateFileList()
				}
			} else if m.activeView == viewDetail {
				if len(m.filePaths) > 0 && m.fileCursor < len(m.filePaths)-1 {
					m.fileCursor++
				}
			} else if m.activeView == viewLog {
				m.logViewport.LineDown(1)
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
			m.updateFileList()
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

func (m model) leftPanelWidth() int {
	w := m.width * 3 / 10
	if w < 20 {
		w = 20
	}
	if w > 40 {
		w = 40
	}
	return w
}

func (m model) rightPanelWidth() int {
	return m.width - m.leftPanelWidth() - 1 // 1 for gap between panels
}

func (m *model) recalcLayout() {
	if m.width == 0 || m.height == 0 {
		return
	}

	rightInner := m.rightPanelWidth() - 2 // border
	panelH := m.mainPanelHeight()

	m.detailViewport.Width = rightInner
	m.detailViewport.Height = panelH

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

func (m *model) updateFileList() {
	if len(m.repoPaths) == 0 {
		m.filePaths = nil
		m.fileCursor = 0
		return
	}
	currentProject := m.repoPaths[m.cursor]
	st, ok := m.projects[currentProject]
	if !ok || len(st.Files) == 0 {
		m.filePaths = nil
		m.fileCursor = 0
		return
	}

	m.filePaths = make([]string, 0, len(st.Files))
	for _, f := range st.Files {
		m.filePaths = append(m.filePaths, f.Path)
	}
	m.fileCursor = 0
}

func (m model) doScan() tea.Cmd {
	config := m.config
	ignoreDirErrors := m.ignoreDirErrors
	return func() tea.Msg {
		projects, err := scanner.Scan(config, ignoreDirErrors)
		return scanMsg{projects: projects, err: err}
	}
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	if m.width < 60 || m.height < 20 {
		return "Terminal too small. Need at least 60x20."
	}

	leftW := m.leftPanelWidth() - 2 // inner width (minus border)
	rightW := m.rightPanelWidth() - 2
	panelH := m.mainPanelHeight()

	// Left panel: project list
	projectContent := m.renderProjectList(leftW, panelH)
	leftPanel := m.renderPanel(viewProjects, m.leftPanelWidth()-2, panelH, projectContent)

	// Right panel: detail
	detailContent := m.renderDetailPanel(rightW, panelH)
	rightPanel := m.renderPanel(viewDetail, m.rightPanelWidth()-2, panelH, detailContent)

	// Join panels horizontally
	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Nav bar
	navBar := m.renderNavBar()

	var view string
	if m.logVisible {
		logPanel := m.renderPanel(viewLog, m.width-2, m.logPanelHeight(), m.logViewport.View())
		view = lipgloss.JoinVertical(lipgloss.Left, mainRow, logPanel, navBar)
	} else {
		view = lipgloss.JoinVertical(lipgloss.Left, mainRow, navBar)
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

	return padToHeight(view, m.height)
}

func (m model) renderProjectList(width int, height int) string {
	if len(m.repoPaths) == 0 {
		return "No OpenSpec projects found."
	}

	var b strings.Builder
	offset := 0
	if m.cursor >= height {
		offset = m.cursor - height + 1
	}

	end := offset + height
	if end > len(m.repoPaths) {
		end = len(m.repoPaths)
	}

	for i := offset; i < end; i++ {
		if i > offset {
			b.WriteString("\n")
		}
		name := m.displayNames[i]
		if i == m.cursor {
			b.WriteString(selectedStyle.Width(width).Render(name))
		} else {
			b.WriteString(normalStyle.Render(name))
		}
	}
	return b.String()
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
		return "No project selected."
	}

	var b strings.Builder

	// Tab header
	b.WriteString(m.renderTabHeader(width))
	b.WriteString("\n")

	// Tab content (height minus tab header line)
	contentHeight := height - 1
	if contentHeight < 1 {
		contentHeight = 1
	}

	switch m.detailTab {
	case tabOverview:
		b.WriteString(m.renderOverview(width, contentHeight))
	case tabConfig:
		b.WriteString(m.renderConfigTab(width, contentHeight))
	default:
		b.WriteString(m.renderNotImplemented(width, contentHeight))
	}

	return b.String()
}

func (m model) renderOverview(width int, height int) string {
	currentProject := m.repoPaths[m.cursor]
	st := m.projects[currentProject]
	info := st.Info

	var b strings.Builder

	// Project path
	b.WriteString(headerStyle.Render(currentProject))
	b.WriteString("\n")

	// Stats line
	statsLine := fmt.Sprintf("Specs: %d  Changes: %d active  Archived: %d",
		info.SpecCount, len(info.ActiveChanges), len(info.ArchivedChanges))
	b.WriteString(dimStyle.Render(statsLine))
	b.WriteString("\n")

	// Active changes
	b.WriteString("\n")
	b.WriteString(sectionHeaderStyle.Render("Active Changes"))
	b.WriteString("\n")
	if len(info.ActiveChanges) == 0 {
		b.WriteString(dimStyle.Render("  None"))
	} else {
		for _, c := range info.ActiveChanges {
			b.WriteString("  - " + c + "\n")
		}
	}

	// Recently archived
	b.WriteString("\n")
	b.WriteString(sectionHeaderStyle.Render("Recently Archived"))
	b.WriteString("\n")
	if len(info.ArchivedChanges) == 0 {
		b.WriteString(dimStyle.Render("  None"))
	} else {
		for _, a := range info.ArchivedChanges {
			date := a.Date.Format("2006-01-02")
			b.WriteString(fmt.Sprintf("  - %s (%s)\n", a.Name, date))
		}
	}

	return b.String()
}

func (m model) renderConfigTab(width int, height int) string {
	if len(m.repoPaths) == 0 {
		return "No project selected."
	}

	currentProject := m.repoPaths[m.cursor]
	info := m.projects[currentProject].Info

	if info.ConfigFile == "" {
		return dimStyle.Render("No project configuration found")
	}

	var b strings.Builder

	// File source indicator
	b.WriteString(dimStyle.Render("openspec/" + info.ConfigFile))
	b.WriteString("\n\n")

	// Render based on file type
	if strings.HasSuffix(info.ConfigFile, ".md") {
		b.WriteString(renderMarkdown(info.ConfigContent, width))
	} else {
		b.WriteString(renderYAML(info.ConfigContent, width))
	}

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

func renderMarkdown(content string, width int) string {
	lines := strings.Split(content, "\n")
	var b strings.Builder

	for i, line := range lines {
		if i > 0 {
			b.WriteString("\n")
		}

		trimmed := strings.TrimSpace(line)

		// Headers
		if strings.HasPrefix(trimmed, "#") {
			b.WriteString(mdHeaderStyle.Render(trimmed))
			continue
		}

		// List items
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			b.WriteString("  " + renderInlineMarkdown(trimmed))
			continue
		}

		// Regular text with inline formatting
		b.WriteString(renderInlineMarkdown(line))
	}

	return b.String()
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

func renderYAML(content string, width int) string {
	lines := strings.Split(content, "\n")
	var b strings.Builder

	for i, line := range lines {
		if i > 0 {
			b.WriteString("\n")
		}

		trimmed := strings.TrimSpace(line)

		// Comment lines
		if strings.HasPrefix(trimmed, "#") {
			b.WriteString(yamlCommentStyle.Render(line))
			continue
		}

		// Key: value lines
		colonIdx := strings.Index(line, ":")
		if colonIdx > 0 {
			// Check for inline comment
			key := line[:colonIdx]
			rest := line[colonIdx:]

			commentIdx := strings.Index(rest, " #")
			if commentIdx > 0 {
				value := rest[:commentIdx]
				comment := rest[commentIdx:]
				b.WriteString(yamlKeyStyle.Render(key))
				b.WriteString(yamlValueStyle.Render(value))
				b.WriteString(yamlCommentStyle.Render(comment))
			} else {
				b.WriteString(yamlKeyStyle.Render(key))
				b.WriteString(yamlValueStyle.Render(rest))
			}
			continue
		}

		// Plain lines (list items, etc)
		b.WriteString(line)
	}

	return b.String()
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
	case viewProjects:
		title = " Projects "
	case viewDetail:
		title = " Detail "
	case viewLog:
		title = " Log "
	}

	borderColor := lipgloss.Color("240")
	if m.activeView == view {
		borderColor = lipgloss.Color("2")
	}

	border := lipgloss.RoundedBorder()
	titleStyled := lipgloss.NewStyle().Foreground(borderColor).Bold(true).Render(title)
	topBorder := border.TopLeft +
		strings.Repeat(border.Top, 1) +
		titleStyled +
		strings.Repeat(border.Top, max(0, width-lipgloss.Width(title)-2)) +
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
	keys := []struct{ key, action string }{
		{"q", "quit"},
		{"s", "scan"},
		{"tab", "switch"},
		{"jk/\u2191\u2193", "navigate"},
		{"\u2190\u2192/1-5", "tabs"},
		{"l", "log"},
		{"gg/G", "jump"},
	}

	var left strings.Builder
	for i, k := range keys {
		if i > 0 {
			left.WriteString(navBarStyle.Render("  "))
		}
		left.WriteString(navBarKeyStyle.Render(k.key))
		left.WriteString(navBarStyle.Render(" " + k.action))
	}

	right := navBarStyle.Render("specgetty " + m.version)

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

func Run(config *scanner.Config, ignoreDirErrors bool, version string) error {
	m := newModel(config, ignoreDirErrors, version)
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
