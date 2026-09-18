package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// numberedDoc builds a document whose rows are individually identifiable, so a
// test can say exactly which rows are on screen.
func numberedDoc(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "row-%03d\n", i)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// makeDocModel opens a change on an artifact holding the given content.
func makeDocModel(t *testing.T, content string, termHeight int) model {
	t.Helper()
	m := makeListModel()
	m.width, m.height = 80, termHeight
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Changes: []scanner.ChangeInfo{{
			Name:             "alpha",
			ArtifactFiles:    []string{"proposal.md", "design.md"},
			ArtifactContents: map[string]string{"proposal.md": content, "design.md": "short"},
		}},
	}}}
	m.level = levelChange
	m.recalcLayout()
	m.syncDocument()
	return m
}

func visibleRows(m model) []string {
	return strings.Split(m.docViewport.View(), "\n")
}

// --- 8.1 wrapping ---

func TestRenderMarkdownWrapsToWidthInCells(t *testing.T) {
	long := strings.Repeat("word ", 60)
	for _, width := range []int{20, 40, 79} {
		for _, row := range strings.Split(renderMarkdown(long, width), "\n") {
			if got := ansi.StringWidth(row); got > width {
				t.Errorf("width %d: row is %d cells: %q", width, got, row)
			}
		}
	}
}

func TestRenderMarkdownCountsCellsNotBytes(t *testing.T) {
	// Double-width runes occupy two cells each, so ten of them fill a pane of
	// twenty. Counting bytes or runes would let the row overflow.
	wide := strings.Repeat("日", 40)
	for _, row := range strings.Split(renderMarkdown(wide, 20), "\n") {
		if got := ansi.StringWidth(row); got > 20 {
			t.Errorf("row is %d cells, want at most 20: %q", got, row)
		}
	}
}

func TestRenderMarkdownBreaksUnbreakableTokens(t *testing.T) {
	// A store path has no spaces. Word wrapping alone would let it overflow,
	// and an overflowing row is re-wrapped by the panel box, which is what
	// makes the viewport's row count disagree with the screen.
	path := "/nix/store/phgx7fgi29vpfydrdd65rgmi4sgi5r3w-openspec-1.10.0/lib/openspec/schemas/spec-driven"
	rows := strings.Split(renderMarkdown(path, 30), "\n")
	if len(rows) < 2 {
		t.Fatalf("expected the path to wrap, got %d row(s)", len(rows))
	}
	for _, row := range rows {
		if got := ansi.StringWidth(row); got > 30 {
			t.Errorf("row is %d cells, want at most 30: %q", got, row)
		}
	}
}

// --- 8.2 styling survives a wrap ---

func TestWrappedStyledLineIsStyledOnEveryRow(t *testing.T) {

	// A header is styled end to end, so every row it occupies must carry the
	// styling. Relying on terminal state carrying across the newline breaks the
	// moment a viewport slices the rows apart.
	header := "## " + strings.Repeat("heading words ", 12)
	rows := strings.Split(renderMarkdown(header, 30), "\n")
	if len(rows) < 2 {
		t.Fatalf("expected the header to wrap, got %d row(s)", len(rows))
	}
	for i, row := range rows {
		if !strings.Contains(row, "\x1b[") {
			t.Errorf("row %d carries no styling: %q", i, row)
		}
	}
}

func TestWrappedBoldSpanIsStyledOnBothRows(t *testing.T) {

	line := "start **" + strings.Repeat("bold ", 12) + "** end"
	rows := strings.Split(renderMarkdown(line, 24), "\n")

	var styled int
	for _, row := range rows {
		if strings.Contains(row, "\x1b[1m") {
			styled++
		}
	}
	if styled < 2 {
		t.Errorf("the bold span covers several rows but only %d carry the opener:\n%s",
			styled, strings.Join(rows, "\n"))
	}
}

func TestReopenStylesAddsNoDisplayWidth(t *testing.T) {
	rows := []string{"\x1b[1mbold text", "continues here\x1b[0m"}
	out := reopenStyles(rows)
	for i := range rows {
		if ansi.StringWidth(out[i]) != ansi.StringWidth(rows[i]) {
			t.Errorf("row %d changed width: %d -> %d",
				i, ansi.StringWidth(rows[i]), ansi.StringWidth(out[i]))
		}
	}
}

// --- 8.3 every row is reachable ---

func TestEveryRowIsReachable(t *testing.T) {
	m := makeDocModel(t, numberedDoc(200), 30)

	height := m.docViewport.Height()
	if height < 5 {
		t.Fatalf("pane height %d is too small to be a meaningful test", height)
	}

	seen := map[string]bool{}
	for {
		for _, row := range visibleRows(m) {
			if trimmed := strings.TrimSpace(row); trimmed != "" {
				seen[trimmed] = true
			}
		}
		if m.docViewport.AtBottom() {
			break
		}
		m.docViewport.PageDown()
	}

	for i := 1; i <= 200; i++ {
		want := fmt.Sprintf("row-%03d", i)
		if !seen[want] {
			t.Fatalf("%s was never displayed: rows are being dropped", want)
		}
	}
}

func TestGotoBottomShowsTheLastRows(t *testing.T) {
	m := makeDocModel(t, numberedDoc(200), 30)
	m.docViewport.GotoBottom()

	rows := visibleRows(m)
	last := strings.TrimSpace(rows[len(rows)-1])
	if last != "row-200" {
		t.Errorf("last visible row is %q, want row-200", last)
	}

	// Nothing missing between the top of the pane and the end.
	height := m.docViewport.Height()
	for i := 0; i < height; i++ {
		want := fmt.Sprintf("row-%03d", 200-height+1+i)
		if strings.TrimSpace(rows[i]) != want {
			t.Errorf("row %d of the pane is %q, want %q", i, strings.TrimSpace(rows[i]), want)
		}
	}
}

// --- 8.4 key handling ---

func firstRow(m model) string {
	return strings.TrimSpace(visibleRows(m)[0])
}

func press(m model, k tea.KeyPressMsg) model {
	updated, _ := m.Update(k)
	return updated.(model)
}

func TestScrollKeyDistances(t *testing.T) {
	base := makeDocModel(t, numberedDoc(500), 30)
	height := base.docViewport.Height()

	tests := []struct {
		name string
		key  tea.KeyPressMsg
		want int
	}{
		{"down", tea.KeyPressMsg{Code: tea.KeyDown}, 1},
		{"j", tea.KeyPressMsg{Code: 'j', Text: "j"}, 1},
		{"pgdown", tea.KeyPressMsg{Code: tea.KeyPgDown}, height},
		{"ctrl+f", tea.KeyPressMsg{Code: 'f', Mod: tea.ModCtrl}, height},
		{"ctrl+d", tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl}, height / 2},
	}
	for _, tt := range tests {
		m := press(base, tt.key)
		if got := m.docViewport.YOffset(); got != tt.want {
			t.Errorf("%s moved to offset %d, want %d (pane height %d)",
				tt.name, got, tt.want, height)
		}
	}
}

func TestScrollUpKeysMirrorDown(t *testing.T) {
	base := makeDocModel(t, numberedDoc(500), 30)
	height := base.docViewport.Height()

	m := base
	m.docViewport.SetYOffset(200)

	for _, tt := range []struct {
		name string
		key  tea.KeyPressMsg
		want int
	}{
		{"up", tea.KeyPressMsg{Code: tea.KeyUp}, 199},
		{"k", tea.KeyPressMsg{Code: 'k', Text: "k"}, 199},
		{"pgup", tea.KeyPressMsg{Code: tea.KeyPgUp}, 200 - height},
		{"ctrl+b", tea.KeyPressMsg{Code: 'b', Mod: tea.ModCtrl}, 200 - height},
		{"ctrl+u", tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl}, 200 - height/2},
	} {
		got := press(m, tt.key)
		if got.docViewport.YOffset() != tt.want {
			t.Errorf("%s moved to offset %d, want %d", tt.name, got.docViewport.YOffset(), tt.want)
		}
	}
}

func TestJumpToEnds(t *testing.T) {
	m := makeDocModel(t, numberedDoc(500), 30)
	m = press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})
	if !m.docViewport.AtBottom() {
		t.Error("G should jump to the end")
	}

	m = press(m, tea.KeyPressMsg{Code: 'g', Text: "g"})
	m = press(m, tea.KeyPressMsg{Code: 'g', Text: "g"})
	if m.docViewport.YOffset() != 0 {
		t.Errorf("gg should jump to the top, offset is %d", m.docViewport.YOffset())
	}
}

func TestScrollStopsAtBounds(t *testing.T) {
	m := makeDocModel(t, numberedDoc(500), 30)

	// Already at the top.
	um := press(m, tea.KeyPressMsg{Code: tea.KeyUp})
	if got := um.docViewport.YOffset(); got != 0 {
		t.Errorf("scrolling up at the top moved to %d, want 0", got)
	}

	m.docViewport.GotoBottom()
	bottom := m.docViewport.YOffset()
	down := press(m, tea.KeyPressMsg{Code: tea.KeyDown})
	if got := down.docViewport.YOffset(); got != bottom {
		t.Errorf("scrolling down at the end moved to %d, want %d", got, bottom)
	}
}

func TestShortDocumentDoesNotScroll(t *testing.T) {
	m := makeDocModel(t, "one\ntwo\nthree", 30)
	for _, k := range []tea.KeyPressMsg{
		{Code: tea.KeyDown}, {Code: tea.KeyPgDown},
		{Code: 'd', Mod: tea.ModCtrl}, {Code: 'G', Text: "G"},
	} {
		um := press(m, k)
		if got := um.docViewport.YOffset(); got != 0 {
			t.Errorf("a document shorter than the pane moved to offset %d", got)
		}
	}
}

// --- 8.5 the lists are unaffected ---

func TestScrollKeysDoNotDisturbTheChangeList(t *testing.T) {
	// The change list has no document, so the scroll keys must find nothing to
	// move. Note that the half-page list paging the tasks mention belonged to
	// the old project list panel and its dead file listing, both removed by
	// project-picker; no list binds these keys any more.
	m := makeListModel()
	m.width, m.height = 80, 30
	m.recalcLayout()
	m.syncDocument()

	if m.docActive() {
		t.Fatal("the change list must not own the vertical axis")
	}
	for _, k := range []tea.KeyPressMsg{
		{Code: tea.KeyPgDown}, {Code: 'f', Mod: tea.ModCtrl}, {Code: 'd', Mod: tea.ModCtrl},
	} {
		got := press(m, k)
		if got.changeCursor != m.changeCursor {
			t.Errorf("a paging key moved the change cursor to %d", got.changeCursor)
		}
		if got.docViewport.YOffset() != 0 {
			t.Errorf("a paging key scrolled a document that is not displayed")
		}
	}
}

func TestJAndKStillMoveTheChangeCursor(t *testing.T) {
	m := makeListModel()
	m.width, m.height = 80, 30
	m.syncDocument()

	got := press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if got.changeCursor != 1 {
		t.Errorf("j moved the change cursor to %d, want 1", got.changeCursor)
	}
}

func TestPickerKeepsTheVerticalAxisWhileOpen(t *testing.T) {
	m := makeDocModel(t, numberedDoc(200), 30)
	m.pickerAll = testProjectRows()
	m.pickerOpen = true

	if m.docActive() {
		t.Error("an open picker must take the vertical axis from the document")
	}
	got := press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if got.docViewport.YOffset() != 0 {
		t.Error("j scrolled the document behind the picker")
	}
}

// --- 8.6 position follows the document ---

func TestDifferentDocumentStartsAtTheTop(t *testing.T) {
	m := makeDocModel(t, numberedDoc(200), 30)
	m.docViewport.SetYOffset(50)

	// Move to the other artifact sub-tab.
	m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
	if m.docViewport.YOffset() != 0 {
		t.Errorf("a different artifact opened at offset %d, want 0", m.docViewport.YOffset())
	}
}

func TestSameDocumentKeepsItsPosition(t *testing.T) {
	m := makeDocModel(t, numberedDoc(200), 30)
	m.docViewport.SetYOffset(50)

	// Away and back, without opening any other document.
	m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
	m = press(m, tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.docViewport.YOffset() != 0 {
		t.Errorf("returning via another artifact should start at the top, got %d",
			m.docViewport.YOffset())
	}

	// Scrolling then redisplaying the same document keeps the position.
	m.docViewport.SetYOffset(40)
	m.syncDocument()
	if m.docViewport.YOffset() != 40 {
		t.Errorf("the same document moved to offset %d, want 40", m.docViewport.YOffset())
	}
}

func TestRewrittenDocumentKeepsPositionClamped(t *testing.T) {
	m := makeDocModel(t, numberedDoc(200), 30)
	m.docViewport.SetYOffset(150)
	keyBefore := m.docKey

	// The watcher rewrites the file while it is being read. Same document, so
	// the position stays, but the content is now much shorter.
	st := m.projects["/p"]
	st.Info.Changes[0].ArtifactContents["proposal.md"] = numberedDoc(40)
	m.projects["/p"] = st
	m.syncDocument()

	if m.docKey != keyBefore {
		t.Fatal("rewriting the file must not change the document identity")
	}
	if m.docViewport.YOffset() > m.docViewport.TotalLineCount() {
		t.Errorf("offset %d points past the end of %d rows",
			m.docViewport.YOffset(), m.docViewport.TotalLineCount())
	}
	if !m.docViewport.AtBottom() {
		t.Errorf("a clamped offset should land at the end, offset %d of %d rows",
			m.docViewport.YOffset(), m.docViewport.TotalLineCount())
	}
}

func TestResizeRewrapsAndClamps(t *testing.T) {
	m := makeDocModel(t, strings.Repeat("some words to wrap ", 200), 30)
	m.docViewport.GotoBottom()

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 200, Height: 30})
	wide := updated.(model)

	// Asked for rather than restated: the panel's content width is derived in
	// one place, and a test that hardcodes the number has to be edited every
	// time the panel's chrome changes, which tells you nothing when it fails.
	if want := wide.contentBoxWidth(); wide.docViewport.Width() != want {
		t.Errorf("viewport width %d, want %d after the resize", wide.docViewport.Width(), want)
	}
	if wide.docViewport.YOffset() > wide.docViewport.TotalLineCount() {
		t.Errorf("offset %d points past the end of %d rows after re-wrapping",
			wide.docViewport.YOffset(), wide.docViewport.TotalLineCount())
	}
}

// --- 8.7 the position indicator ---

func TestTitleReportsPositionForATallDocument(t *testing.T) {
	m := makeDocModel(t, numberedDoc(200), 30)

	if pct := m.docScrollPercent(); pct != 0 {
		t.Errorf("at the top the position is %d%%, want 0%%", pct)
	}
	m.docViewport.GotoBottom()
	if pct := m.docScrollPercent(); pct != 100 {
		t.Errorf("at the end the position is %d%%, want 100%%", pct)
	}

	title := m.renderPanel(viewDetail, m.width-2, m.mainPanelHeight(), "x")
	if !strings.Contains(title, "100%") {
		t.Errorf("the panel title should carry the position:\n%s", strings.Split(title, "\n")[0])
	}
}

func TestTitleShowsNoPositionForAShortDocument(t *testing.T) {
	m := makeDocModel(t, "one\ntwo", 30)
	if pct := m.docScrollPercent(); pct != -1 {
		t.Errorf("a document that fits reported %d%%, want no indicator", pct)
	}
	title := strings.Split(m.renderPanel(viewDetail, m.width-2, m.mainPanelHeight(), "x"), "\n")[0]
	if strings.Contains(title, "%") {
		t.Errorf("a document that fits should show no indicator:\n%s", title)
	}
}

func TestTitleHasNoPositionWhenNoDocumentIsShown(t *testing.T) {
	m := makeListModel()
	m.width, m.height = 80, 30
	m.syncDocument()
	if pct := m.docScrollPercent(); pct != -1 {
		t.Errorf("the change list reported a position of %d%%", pct)
	}
}

// --- config tab ---

func TestConfigTabScrollsAndKeepsItsSourceLine(t *testing.T) {
	m := makeListModel()
	m.width, m.height = 80, 30
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		ConfigFile:    "project.md",
		ConfigContent: numberedDoc(200),
	}}}
	m.level = levelProject
	m.detailTab = tabConfig
	m.recalcLayout()
	m.syncDocument()

	if !m.docActive() {
		t.Fatal("the config tab should own the vertical axis")
	}

	m = press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
	if m.docViewport.YOffset() == 0 {
		t.Error("the config tab did not scroll")
	}

	// The source line names the content, so it is drawn above the content's
	// border rather than inside it. Asserted on the frame, because that is the
	// only place both the line and the border exist.
	var sawSource, sawBorder bool
	for _, l := range strings.Split(m.renderFrame(), "\n") {
		plain := ansi.Strip(l)
		if strings.Contains(plain, "openspec/project.md") {
			if sawBorder {
				t.Error("the file source line is inside the content border, not above it")
			}
			sawSource = true
		}
		if strings.Contains(plain, "╭") && !strings.HasPrefix(plain, "╭") {
			sawBorder = true
		}
	}
	if !sawSource {
		t.Error("the file source line is missing from the frame")
	}
	if !sawBorder {
		t.Error("the content border is missing from the frame")
	}
}

func TestConfigTabResetsWhenTheProjectChanges(t *testing.T) {
	m := makeListModel()
	m.width, m.height = 80, 30
	m.repoPaths = []string{"/p", "/q"}
	m.displayNames = []string{"p", "q"}
	info := scanner.ProjectInfo{ConfigFile: "project.md", ConfigContent: numberedDoc(200)}
	m.projects = scanner.ProjectMap{
		"/p": scanner.ProjectStatus{Info: info},
		"/q": scanner.ProjectStatus{Info: info},
	}
	m.level = levelProject
	m.detailTab = tabConfig
	m.recalcLayout()
	m.syncDocument()
	m.docViewport.SetYOffset(60)

	m.cursor = 1 // a different project, same file name and content
	m.syncDocument()

	if m.docViewport.YOffset() != 0 {
		t.Errorf("switching project left the offset at %d, want 0", m.docViewport.YOffset())
	}
}
