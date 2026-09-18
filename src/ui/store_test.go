package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// storeBackedModel is a project opened from a repo that declares a store: the
// content is the store's, the origin is the repo.
func storeBackedModel() model {
	root := "/stores/nivis-tunnel"
	origin := "/work/nivis-tunnel"
	return model{
		width:     100,
		height:    40,
		repoPaths: []string{root},
		cursor:    0,
		detailTab: tabChanges,
		projects: scanner.ProjectMap{
			root: scanner.ProjectStatus{
				Info: scanner.ProjectInfo{
					Root:                root,
					Origin:              origin,
					StoreID:             "nivis-tunnel",
					SpecCount:           6,
					SpecNames:           []string{"tunnel-relay"},
					SpecContents:        map[string]string{"tunnel-relay": "# tunnel-relay\n"},
					ConfigFile:          "config.yaml",
					ConfigContent:       "schema: spec-driven\n# shared by both repos\n",
					OriginConfigFile:    "config.yaml",
					OriginConfigContent: "schema: spec-driven\nstore: nivis-tunnel\ncontext: |\n  repo only\n",
					Store: &scanner.StoreInfo{
						ID: "nivis-tunnel", Root: root, Origin: origin,
						Remote: "git@example.com:n/stores.git",
						Git: &scanner.StoreGit{
							IsRepo: true, OriginURL: "git@example.com:n/stores.git",
							DirtyKnown: true, Dirty: false,
							TrackingKnown: true, Ahead: 2, Behind: 0,
						},
					},
				},
			},
		},
		displayNames: []string{"nivis-tunnel"},
	}
}

// plainModel is an ordinary project that holds its own content.
func plainModel() model {
	root := "/work/specgetty"
	return model{
		width:     100,
		height:    40,
		repoPaths: []string{root},
		cursor:    0,
		detailTab: tabChanges,
		projects: scanner.ProjectMap{
			root: scanner.ProjectStatus{
				Info: scanner.ProjectInfo{
					Root: root, Origin: root,
					SpecCount:     22,
					ConfigFile:    "config.yaml",
					ConfigContent: "schema: spec-driven\n",
				},
			},
		},
		displayNames: []string{"specgetty"},
	}
}

// storeItselfModel is a store opened on its own account from the picker.
func storeItselfModel() model {
	root := "/stores/nivis-tunnel"
	return model{
		width:     100,
		height:    40,
		repoPaths: []string{root},
		cursor:    0,
		detailTab: tabConfig,
		projects: scanner.ProjectMap{
			root: scanner.ProjectStatus{
				Info: scanner.ProjectInfo{
					Root: root, Origin: root,
					StoreID:       "nivis-tunnel",
					ConfigFile:    "config.yaml",
					ConfigContent: "schema: spec-driven\n",
					Store: &scanner.StoreInfo{
						ID: "nivis-tunnel", Root: root,
						Git: &scanner.StoreGit{IsRepo: false},
					},
				},
			},
		},
		displayNames: []string{"nivis-tunnel"},
	}
}

func plainText(m model) string {
	return ansi.Strip(m.renderFrame())
}

// --- 4.1 and 4.2 the header mark ---

func TestHeaderNamesTheRepoAndMarksTheStore(t *testing.T) {
	m := storeBackedModel()
	out := plainText(m)

	if !strings.Contains(out, "/work/nivis-tunnel") {
		t.Error("the header must name the repo the user is standing in")
	}
	if !strings.Contains(out, "store") {
		t.Error("the header must mark that the content came from a store")
	}
}

func TestHeaderNamesAStoreOpenedDirectlyByItsID(t *testing.T) {
	m := storeItselfModel()
	if got := headerName(m.projects[m.currentRoot()].Info, m.currentRoot()); got != "nivis-tunnel" {
		t.Errorf("got %q, want the store id", got)
	}
	if storeMark(m.projects[m.currentRoot()].Info) == "" {
		t.Error("a store carries the mark too")
	}
}

func TestHeaderOfAPlainProjectIsUnchanged(t *testing.T) {
	m := plainModel()
	info := m.projects[m.currentRoot()].Info

	if got := headerName(info, m.currentRoot()); got != "/work/specgetty" {
		t.Errorf("got %q, want the project path", got)
	}
	if got := storeMark(info); got != "" {
		t.Errorf("got %q, want no mark on a plain project", got)
	}
}

func TestTheHeaderMarkCarriesNothingElse(t *testing.T) {
	// The mark says one thing. The id, the root path, the remote and the git
	// state are answers to a question asked once per project, and they live on
	// the config tab instead. This is what keeps the header from growing a row.
	info := storeBackedModel().projects["/stores/nivis-tunnel"].Info
	mark := ansi.Strip(storeMark(info))

	for _, forbidden := range []string{
		"nivis-tunnel", "/stores/", "git@example.com", "ahead", "2",
	} {
		if strings.Contains(mark, forbidden) {
			t.Errorf("the mark %q must not carry %q", mark, forbidden)
		}
	}
	if strings.TrimSpace(mark) != "store" {
		t.Errorf("got %q, want just the word store", strings.TrimSpace(mark))
	}
}

// --- 5.1 to 5.5 the config panes ---

func TestConfigPanesForAStoreBackedProject(t *testing.T) {
	info := storeBackedModel().projects["/stores/nivis-tunnel"].Info
	panes := configPanes(info)

	if len(panes) != 3 {
		t.Fatalf("got %d panes, want repo, store and store details", len(panes))
	}
	want := []string{"repo", "store", "store details"}
	for i, w := range want {
		if panes[i].label != w {
			t.Errorf("pane %d: got %q, want %q", i, panes[i].label, w)
		}
	}
	if !strings.Contains(panes[0].content, "repo only") {
		t.Error("the repo pane must show the repo's own configuration")
	}
	if !strings.Contains(panes[1].content, "shared by both repos") {
		t.Error("the store pane must show the store's configuration")
	}
}

func TestConfigPanesForAPlainProject(t *testing.T) {
	panes := configPanes(plainModel().projects["/work/specgetty"].Info)
	if len(panes) != 1 {
		t.Fatalf("got %d panes, want one", len(panes))
	}
	if panes[0].source != "openspec/config.yaml" {
		t.Errorf("got %q, want the filename", panes[0].source)
	}
}

func TestConfigPanesForAStoreOpenedDirectly(t *testing.T) {
	// No repo to name, so no repo pane is offered.
	panes := configPanes(storeItselfModel().projects["/stores/nivis-tunnel"].Info)
	for _, p := range panes {
		if p.label == "repo" {
			t.Error("a store opened directly has no originating repo")
		}
	}
	if len(panes) != 2 {
		t.Fatalf("got %d panes, want the store's configuration and its details", len(panes))
	}
}

func TestPlainProjectConfigTabDrawsNoSubTabRow(t *testing.T) {
	m := plainModel()
	m.detailTab = tabConfig
	m.syncDocument()
	out := plainText(m)

	if !strings.Contains(out, "openspec/config.yaml") {
		t.Error("a single configuration keeps its dimmed filename")
	}
	if strings.Contains(out, "store details") {
		t.Error("a plain project has no sub-tabs")
	}
}

// --- 5.6 to 5.8 the store details ---

func TestStoreDetailsReportLocalFactsOnly(t *testing.T) {
	info := storeBackedModel().projects["/stores/nivis-tunnel"].Info
	out := storeDetails(info)

	for _, want := range []string{
		"store: nivis-tunnel",
		"root: /stores/nivis-tunnel",
		"origin: /work/nivis-tunnel",
		"registered_remote: git@example.com:n/stores.git",
		"uncommitted_changes: no",
		"ahead: 2",
		"behind: 0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("store details missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "without fetching") {
		t.Error("the report must say the comparison did not fetch")
	}
}

func TestStoreDetailsWithoutAGitWorkingCopy(t *testing.T) {
	out := storeDetails(storeItselfModel().projects["/stores/nivis-tunnel"].Info)
	if !strings.Contains(out, "not a git working copy") {
		t.Errorf("want a plain report rather than an error:\n%s", out)
	}
	if strings.Contains(out, "ahead:") {
		t.Error("no git state to report")
	}
}

func TestStoreDetailsWithNoUpstream(t *testing.T) {
	m := storeBackedModel()
	info := m.projects["/stores/nivis-tunnel"].Info
	info.Store.Git = &scanner.StoreGit{IsRepo: true, DirtyKnown: true, Dirty: true}
	out := storeDetails(info)

	if !strings.Contains(out, "uncommitted_changes: yes") {
		t.Error("a dirty working copy must say so")
	}
	if !strings.Contains(out, "no upstream to compare against") {
		t.Errorf("want the absence of an upstream stated:\n%s", out)
	}
}

// --- 5.9 an unresolved declaration is reported, not shown as emptiness ---

func problemModel() model {
	root := "/work/nivis-tunnel"
	return model{
		width:     100,
		height:    40,
		repoPaths: []string{root},
		cursor:    0,
		detailTab: tabChanges,
		projects: scanner.ProjectMap{
			root: scanner.ProjectStatus{
				Info: scanner.ProjectInfo{
					Root: root, Origin: root,
					ConfigFile:    "config.yaml",
					ConfigContent: "store: nivis-tunnel\n",
					StoreProblem: &scanner.StoreProblem{
						Code:   scanner.ProblemUnknownStore,
						ID:     "nivis-tunnel",
						File:   "/work/nivis-tunnel/openspec/config.yaml",
						Detail: "it is not among the registered stores: other",
					},
				},
			},
		},
		displayNames: []string{"nivis-tunnel"},
	}
}

func TestAnUnfollowedDeclarationIsReportedOnTheChangesTab(t *testing.T) {
	// Zero specs and zero changes reads as a project with nothing in it. That
	// is the exact failure this change exists to fix, so the emptiness never
	// stands as the whole report.
	m := problemModel()
	out := plainText(m)

	if !strings.Contains(out, "nivis-tunnel") {
		t.Error("the report must name the declared store")
	}
	if !strings.Contains(out, "not among the registered stores") {
		t.Errorf("the report must say why:\n%s", out)
	}
	if strings.Contains(out, "No active changes") {
		t.Error("the emptiness must not be the whole report")
	}
}

func TestAnUnfollowedDeclarationIsReportedOnTheSpecsTab(t *testing.T) {
	m := problemModel()
	m.detailTab = tabSpecs
	out := plainText(m)

	if !strings.Contains(out, "not among the registered stores") {
		t.Errorf("the specs tab must say why it is empty:\n%s", out)
	}
}

func TestAnUnfollowedDeclarationGetsAStoreDetailsPane(t *testing.T) {
	info := problemModel().projects["/work/nivis-tunnel"].Info
	panes := configPanes(info)

	var details string
	for _, p := range panes {
		if p.kind == paneDetails {
			details = p.content
		}
	}
	if details == "" {
		t.Fatal("an unresolved declaration still gets a details pane")
	}
	for _, want := range []string{"nivis-tunnel", "openspec/config.yaml", "not among the registered stores"} {
		if !strings.Contains(details, want) {
			t.Errorf("details missing %q:\n%s", want, details)
		}
	}
}

func TestAPlainProjectHasNoProblemLine(t *testing.T) {
	if got := storeProblemLine(plainModel().projects["/work/specgetty"].Info); got != "" {
		t.Errorf("got %q, want nothing", got)
	}
}

// --- 5.10 each pane is its own document ---

func TestSwitchingConfigPaneChangesTheDocument(t *testing.T) {
	m := storeBackedModel()
	m.detailTab = tabConfig
	m.recalcLayout()
	m.syncDocument()

	first, firstContent, ok := m.currentDocument()
	if !ok {
		t.Fatal("expected a document")
	}
	if !strings.Contains(ansi.Strip(firstContent), "repo only") {
		t.Error("the first pane is the repo's configuration")
	}

	m.configPane = 1
	m.syncDocument()
	second, secondContent, _ := m.currentDocument()

	if first == second {
		t.Error("a different pane is a different document")
	}
	if !strings.Contains(ansi.Strip(secondContent), "shared by both repos") {
		t.Error("the second pane is the store's configuration")
	}
}

func TestConfigPaneIndexIsClampedToTheProject(t *testing.T) {
	// Switching from a store-backed project to a plain one must not land past
	// the end of a shorter set of panes.
	m := plainModel()
	m.configPane = 2
	panes := configPanes(m.projects["/work/specgetty"].Info)
	if got := m.configPaneIndex(panes); got != 0 {
		t.Errorf("got %d, want the only pane", got)
	}
}

func TestTabCyclesTheConfigPanes(t *testing.T) {
	m := storeBackedModel()
	m.detailTab = tabConfig
	m.focus = focusDetail
	m.recalcLayout()

	for _, want := range []int{1, 2, 0} {
		updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		m = updated.(model)
		if m.configPane != want {
			t.Fatalf("got pane %d, want %d", m.configPane, want)
		}
	}
}

func TestTabDoesNothingToASingleConfigPane(t *testing.T) {
	m := plainModel()
	m.detailTab = tabConfig
	m.focus = focusDetail
	m.recalcLayout()

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if updated.(model).configPane != 0 {
		t.Error("one pane has nothing to cycle through")
	}
}

// --- 3.2 to 3.4 the actions target the root ---

func TestActionsTargetTheRootNotTheOrigin(t *testing.T) {
	m := storeBackedModel()
	if got := m.currentRoot(); got != "/stores/nivis-tunnel" {
		t.Errorf("got %q, want the store", got)
	}
	if got := m.startDirOf(m.currentRoot()); got != "/work/nivis-tunnel" {
		t.Errorf("start dir: got %q, want the repo", got)
	}
}

func TestDiscardMovesInsideTheStoreAndTouchesNothingElse(t *testing.T) {
	// The defect this fixes: discard built its path from the directory the
	// user stood in, so a store-backed project grew a discarded/ folder in a
	// repo that holds no changes at all.
	store := t.TempDir()
	origin := t.TempDir()
	changeDir := filepath.Join(store, "openspec", "changes", "my-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [x] done\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	msg := doDiscardChange(store, "my-change")()
	res, ok := msg.(discardMsg)
	if !ok || !res.ok {
		t.Fatalf("discard failed: %+v", msg)
	}

	entries, err := os.ReadDir(filepath.Join(store, "openspec", "changes", "discarded"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("the change must land in the store: %v %v", entries, err)
	}
	if _, err := os.Stat(filepath.Join(origin, "openspec")); !os.IsNotExist(err) {
		t.Error("nothing may be created under the repo the user started in")
	}
}

func TestExportReadsFromTheRoot(t *testing.T) {
	store := t.TempDir()
	changeDir := filepath.Join(store, "openspec", "changes", "my-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("# why\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", t.TempDir())

	msg := doExportChange(store, "my-change", "my-change", false)()
	res, ok := msg.(exportMsg)
	if !ok || !res.ok {
		t.Fatalf("export failed: %+v", msg)
	}
}

// --- 4.3 and 4.4 the picker names a store by its id ---

func TestPickerNamesAStoreByItsDeclaredID(t *testing.T) {
	projects := scanner.ProjectMap{
		"/stores/tunnel-dir": {Info: scanner.ProjectInfo{StoreID: "nivis-tunnel"}},
		"/work/specgetty":    {Info: scanner.ProjectInfo{}},
	}
	rows := buildProjectRows(projects)

	byPath := map[string]string{}
	for _, r := range rows {
		byPath[r.path] = r.display
	}
	if byPath["/stores/tunnel-dir"] != "nivis-tunnel" {
		t.Errorf("got %q, want the declared id rather than the folder name", byPath["/stores/tunnel-dir"])
	}
	if byPath["/work/specgetty"] != "specgetty" {
		t.Errorf("got %q, want the basename", byPath["/work/specgetty"])
	}
}

func TestPickerDisambiguatesAStoreIDColliding(t *testing.T) {
	projects := scanner.ProjectMap{
		"/stores/whatever": {Info: scanner.ProjectInfo{StoreID: "nivis"}},
		"/work/nivis":      {Info: scanner.ProjectInfo{}},
	}
	rows := buildProjectRows(projects)
	for _, r := range rows {
		if r.display == "nivis" {
			t.Errorf("a colliding name must be qualified: %q", r.display)
		}
		if !strings.Contains(r.display, "(") {
			t.Errorf("got %q, want a parent directory qualifier", r.display)
		}
	}
}

func TestPickerMarksAStoreRow(t *testing.T) {
	var kind *fieldDef[projectRow]
	for i := range projectFields {
		if projectFields[i].id == "kind" {
			kind = &projectFields[i]
		}
	}
	if kind == nil {
		t.Fatal("the picker has no column saying what a row is")
	}
	if got := kind.value(projectRow{info: scanner.ProjectInfo{StoreID: "alpha"}}); got != "store" {
		t.Errorf("got %q, want store", got)
	}
	if got := kind.value(projectRow{info: scanner.ProjectInfo{}}); got != "" {
		t.Errorf("got %q, want nothing for a plain project", got)
	}
}

// --- 4.5 and 4.6 the picker cursor ---

func TestPickerOpensOnTheRowHoldingTheOpenContent(t *testing.T) {
	m := storeBackedModel()
	m.focus = focusDetail
	m.pickerAll = []projectRow{
		{path: "/work/specgetty", display: "specgetty"},
		{path: "/stores/nivis-tunnel", display: "nivis-tunnel"},
	}
	m.pickerLoaded = true
	m.recalcLayout()

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	um := updated.(model)

	if !um.pickerOpen {
		t.Fatal("the picker must open")
	}
	if um.pickerKey != "/stores/nivis-tunnel" {
		t.Errorf("got %q, want the store's row", um.pickerKey)
	}
	rows := um.pickerVisibleRows()
	if um.pickerCursor >= len(rows) || rows[um.pickerCursor].row.path != "/stores/nivis-tunnel" {
		t.Errorf("cursor at %d, want the store's row", um.pickerCursor)
	}
}

func TestDismissingThePickerLeavesTheProjectUntouched(t *testing.T) {
	m := storeBackedModel()
	m.focus = focusDetail
	m.pickerAll = []projectRow{{path: "/stores/nivis-tunnel", display: "nivis-tunnel"}}
	m.pickerLoaded = true
	m.recalcLayout()

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	updated, _ = updated.(model).Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	um := updated.(model)

	if um.pickerOpen {
		t.Error("esc closes the picker")
	}
	info := um.projects[um.currentRoot()].Info
	if info.Origin != "/work/nivis-tunnel" || info.StoreID != "nivis-tunnel" {
		t.Errorf("the project underneath changed: %+v", info)
	}
}

// --- 6.1 the watched trees ---

func TestStartDirOfAPlainProjectIsItself(t *testing.T) {
	m := plainModel()
	if got := m.startDirOf("/work/specgetty"); got != "/work/specgetty" {
		t.Errorf("got %q, want the project itself", got)
	}
}

func TestWatchDirsCoversBothTreesForAStoreBackedProject(t *testing.T) {
	m := storeBackedModel()
	dirs := m.watchDirs(m.currentRoot())

	if len(dirs) != 2 {
		t.Fatalf("got %v, want the store's tree and the repo's", dirs)
	}
	want := map[string]bool{
		filepath.Join("/stores/nivis-tunnel", "openspec"): true,
		filepath.Join("/work/nivis-tunnel", "openspec"):   true,
	}
	for _, d := range dirs {
		if !want[d] {
			t.Errorf("unexpected watched tree %q", d)
		}
		delete(want, d)
	}
	if len(want) != 0 {
		t.Errorf("not watched: %v", want)
	}
}

func TestWatchDirsIsOneTreeForAPlainProject(t *testing.T) {
	m := plainModel()
	dirs := m.watchDirs(m.currentRoot())
	if len(dirs) != 1 {
		t.Fatalf("got %v, want one tree", dirs)
	}
	if dirs[0] != filepath.Join("/work/specgetty", "openspec") {
		t.Errorf("got %q", dirs[0])
	}
}

func TestWatchDirsIsOneTreeForAStoreOpenedDirectly(t *testing.T) {
	m := storeItselfModel()
	if got := m.watchDirs(m.currentRoot()); len(got) != 1 {
		t.Errorf("got %v, want one tree: there is no separate origin", got)
	}
}

// --- 6.3 a project switch leaves no watcher behind ---

func TestStopWatcherForgetsTheWatchedRoot(t *testing.T) {
	m := storeBackedModel()
	m.watchedRoot = "/stores/nivis-tunnel"
	m.stopWatcher()
	if m.watchedRoot != "" {
		t.Errorf("got %q, want the watched root cleared", m.watchedRoot)
	}
}

// --- 5.1 and 5.2 the config tab's geometry ---

// frameSize reports the rendered frame's rows and its widest column count.
func frameSize(m model) (rows, cols int) {
	lines := strings.Split(m.renderFrame(), "\n")
	for _, l := range lines {
		if w := ansi.StringWidth(l); w > cols {
			cols = w
		}
	}
	return len(lines), cols
}

func TestConfigTabGeometryIsUnchangedByTheSubTabRow(t *testing.T) {
	// A nested row is exactly the shape of change a width assertion misses:
	// lipgloss pads a wrapped remainder back out to full width, so only the
	// line count catches an overflow. Both shapes are measured at the same
	// three widths, including the 60-column minimum.
	for _, size := range []struct{ w, h int }{{60, 20}, {92, 30}, {120, 50}} {
		for _, tc := range []struct {
			name string
			m    model
		}{
			{"plain", plainModel()},
			{"store-backed", storeBackedModel()},
		} {
			t.Run(fmt.Sprintf("%s %dx%d", tc.name, size.w, size.h), func(t *testing.T) {
				m := tc.m
				m.width, m.height = size.w, size.h
				m.detailTab = tabConfig
				m.recalcLayout()
				m.syncDocument()

				rows, cols := frameSize(m)
				if rows != size.h {
					t.Errorf("rows: got %d, want %d", rows, size.h)
				}
				if cols > size.w {
					t.Errorf("columns: got %d, want at most %d", cols, size.w)
				}
			})
		}
	}
}

func TestPlainConfigTabKeepsItsFilenameAboveTheBorder(t *testing.T) {
	m := plainModel()
	m.detailTab = tabConfig
	m.recalcLayout()
	m.syncDocument()

	lines := strings.Split(ansi.Strip(m.renderFrame()), "\n")
	filenameRow := rowContaining(lines, 0, "openspec/config.yaml")
	if filenameRow < 0 {
		t.Fatal("the filename must be shown")
	}
	// The panel draws its own border at the top of the frame, so the box that
	// matters is the first one below the naming row.
	if rowContaining(lines, filenameRow+1, "╭") < 0 {
		t.Error("the filename names the content, so the content box must open below it")
	}
}

// rowContaining returns the first index at or after from whose line holds sub.
func rowContaining(lines []string, from int, sub string) int {
	for i := from; i < len(lines); i++ {
		if strings.Contains(lines[i], sub) {
			return i
		}
	}
	return -1
}

func TestSubTabRowSitsAboveTheBorder(t *testing.T) {
	m := storeBackedModel()
	m.detailTab = tabConfig
	m.recalcLayout()
	m.syncDocument()

	lines := strings.Split(ansi.Strip(m.renderFrame()), "\n")
	subTabRow := rowContaining(lines, 0, "store details")
	if subTabRow < 0 {
		t.Fatal("the sub-tabs must be shown")
	}
	if rowContaining(lines, subTabRow+1, "╭") < 0 {
		t.Error("sub-tabs name the content, so the content box must open below them")
	}
}

// runCmd executes a command, flattening the batch bubbletea wraps them in, and
// returns every message produced.
func runCmd(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		var out []tea.Msg
		for _, c := range batch {
			out = append(out, runCmd(c)...)
		}
		return out
	}
	return []tea.Msg{msg}
}

// TestDiscardFromTheUITargetsTheStore drives the keys rather than calling the
// command directly, so it fails if the dispatch hands over the origin instead
// of the root. That substitution is the defect this change fixes, and calling
// doDiscardChange with an explicit path would not catch it.
func TestDiscardFromTheUITargetsTheStore(t *testing.T) {
	store := t.TempDir()
	origin := t.TempDir()
	if err := os.MkdirAll(filepath.Join(store, "openspec", "changes", "my-change"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(origin, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}

	m := model{
		width: 100, height: 40,
		repoPaths: []string{store},
		detailTab: tabChanges,
		focus:     focusDetail,
		projects: scanner.ProjectMap{
			store: scanner.ProjectStatus{
				Info: scanner.ProjectInfo{
					Root: store, Origin: origin, StoreID: "alpha",
					ActiveChanges: []string{"my-change"},
					Changes:       []scanner.ChangeInfo{{Name: "my-change", DirName: "my-change"}},
				},
			},
		},
		displayNames: []string{"alpha"},
	}
	m.recalcLayout()

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'd', Text: "d"})
	m = updated.(model)
	if m.discardState != discardConfirming {
		t.Fatalf("discardState = %d, want confirming", m.discardState)
	}

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	m = updated.(model)

	var result *discardMsg
	for _, msg := range runCmd(cmd) {
		if d, ok := msg.(discardMsg); ok {
			result = &d
		}
	}
	if result == nil || !result.ok {
		t.Fatalf("discard did not succeed: %+v", result)
	}

	if _, err := os.Stat(filepath.Join(store, "openspec", "changes", "discarded")); err != nil {
		t.Errorf("the change must be moved inside the store: %v", err)
	}
	if _, err := os.Stat(filepath.Join(origin, "openspec", "changes")); !os.IsNotExist(err) {
		t.Error("nothing may be created under the repo the user started in")
	}
}

// TestExportFromTheUIReadsTheStore is the same shape for export: the source
// directory comes from the root, so a dispatch handing over the origin finds
// nothing and the export fails.
func TestExportFromTheUIReadsTheStore(t *testing.T) {
	store := t.TempDir()
	origin := t.TempDir()
	changeDir := filepath.Join(store, "openspec", "changes", "my-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("# why\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", t.TempDir())

	m := model{
		width: 100, height: 40,
		repoPaths: []string{store},
		detailTab: tabChanges,
		focus:     focusDetail,
		projects: scanner.ProjectMap{
			store: scanner.ProjectStatus{
				Info: scanner.ProjectInfo{
					Root: store, Origin: origin, StoreID: "alpha",
					ActiveChanges: []string{"my-change"},
					Changes:       []scanner.ChangeInfo{{Name: "my-change", DirName: "my-change"}},
				},
			},
		},
		displayNames: []string{"alpha"},
	}
	m.recalcLayout()

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	m = updated.(model)
	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})

	var result *exportMsg
	for _, msg := range runCmd(cmd) {
		if e, ok := msg.(exportMsg); ok {
			result = &e
		}
	}
	if result == nil || !result.ok {
		t.Fatalf("export must read from the store: %+v", result)
	}
}

func TestConfigTabKeepsItsPositionAcrossATabSwitch(t *testing.T) {
	// The retention the single-configuration tab already had, kept for the
	// active sub-tab: leaving the config tab and coming back does not rewind.
	m := storeBackedModel()
	m.detailTab = tabConfig
	m.configPane = 1
	m.recalcLayout()
	m.syncDocument()

	before := m.docKey
	m.detailTab = tabChanges
	m.syncDocument()
	m.detailTab = tabConfig
	m.syncDocument()

	if m.docKey != before {
		t.Errorf("the same sub-tab is the same document: got %q, want %q", m.docKey, before)
	}
	if m.configPane != 1 {
		t.Errorf("the active sub-tab survives a tab switch: got %d", m.configPane)
	}
}

func TestSelectingAnotherConfigPaneStartsAtTheTop(t *testing.T) {
	m := storeBackedModel()
	m.detailTab = tabConfig
	m.recalcLayout()
	m.syncDocument()
	first := m.docKey

	m.configPane = 1
	m.syncDocument()

	if m.docKey == first {
		t.Fatal("a different sub-tab is a different document")
	}
	if m.docViewport.YOffset() != 0 {
		t.Errorf("a different document starts at the top, got offset %d", m.docViewport.YOffset())
	}
}

func TestSwitchingProjectResetsTheConfigPane(t *testing.T) {
	m := storeBackedModel()
	m.detailTab = tabConfig
	m.configPane = 2
	m.resetProjectState()
	if m.configPane != 0 {
		t.Errorf("got %d, want the first sub-tab on a new project", m.configPane)
	}
}
