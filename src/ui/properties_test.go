package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// twoSchemaModel is a project running two workflow schemas, with a change that
// records none.
func twoSchemaModel() model {
	root := "/work/specgetty"
	return model{
		width: 100, height: 40,
		repoPaths: []string{root},
		detailTab: tabProperties,
		focus:     focusListPane,
		projects: scanner.ProjectMap{
			root: scanner.ProjectStatus{Info: scanner.ProjectInfo{
				Root: root, Origin: root,
				ConfigFile:    "config.yaml",
				ConfigContent: "schema: spec-driven\n",
				DefaultSchema: "spec-driven",
				SchemaUsage: []scanner.SchemaUsage{
					{Name: "spec-driven", Changes: 24, IsDefault: true},
					{Name: "tinychange", Changes: 10},
				},
				UnrecordedChanges: 1,
			}},
		},
		displayNames: []string{"specgetty"},
	}
}

// --- 5.8 to 5.10 what a schema row reports ---

func TestSchemaRowReportsTheDefinition(t *testing.T) {
	m := twoSchemaModel()
	info := m.projects[m.currentKey()].Info
	state := schemaState{details: &scanner.SchemaDetails{
		Name: "spec-driven", Source: "package", Path: "/pkg/schemas/spec-driven",
		Description: "Default OpenSpec workflow\nsecond line ignored",
		Artifacts: []scanner.SchemaArtifact{
			{ID: "proposal", Generates: "proposal.md"},
			{ID: "tasks", Generates: "tasks.md", Requires: []string{"proposal"}},
		},
		Apply: scanner.SchemaApply{Requires: []string{"tasks"}, Tracks: "tasks.md"},
	}}

	out := renderSchemaSection(info, "spec-driven", state)
	for _, want := range []string{
		"schema: spec-driven",
		"default: yes",
		"changes: 24",
		"source: package",
		"path: /pkg/schemas/spec-driven",
		"proposal -> proposal.md",
		"tasks -> tasks.md (after proposal)",
		"requires: tasks",
		"tracks: tasks.md",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "second line ignored") {
		t.Error("a schema's prose is guidance for whoever writes the artifact, not a property of the project")
	}
}

func TestSchemaRowReportsChangesCarryingNoSchema(t *testing.T) {
	m := twoSchemaModel()
	out := renderSchemaSection(m.projects[m.currentKey()].Info, "tinychange", schemaState{})
	if !strings.Contains(out, "1 change(s) in this project record no schema") {
		t.Errorf("changes with no metadata are counted apart:\n%s", out)
	}
	if strings.Contains(out, "default: yes") {
		t.Error("tinychange is not the project default here")
	}
	if !strings.Contains(out, "changes: 10") {
		t.Errorf("its own count must be reported:\n%s", out)
	}
}

func TestSchemaRowWhileItsDefinitionIsBeingRead(t *testing.T) {
	m := twoSchemaModel()
	out := renderSchemaSection(m.projects[m.currentKey()].Info, "spec-driven", schemaState{})
	if !strings.Contains(out, "reading its definition") {
		t.Errorf("an unfinished read says so:\n%s", out)
	}
	// What is already known does not wait on the subprocess.
	if !strings.Contains(out, "changes: 24") {
		t.Errorf("the count comes from the scan, not the CLI:\n%s", out)
	}
}

func TestSchemaRowWhenTheDefinitionCouldNotBeRead(t *testing.T) {
	m := twoSchemaModel()
	state := schemaState{problem: &scanner.SchemaProblem{
		Code: scanner.SchemaNotFound, Name: "spec-driven",
		Detail:    "Schema 'spec-driven' not found",
		Available: []string{"tinychange"},
	}}
	out := renderSchemaSection(m.projects[m.currentKey()].Info, "spec-driven", state)
	if !strings.Contains(out, "not found") {
		t.Errorf("the reason must be given:\n%s", out)
	}
	if !strings.Contains(out, "available: tinychange") {
		t.Errorf("the alternatives the CLI supplies must be shown:\n%s", out)
	}
}

func TestSchemaRowReportsAnOverride(t *testing.T) {
	m := twoSchemaModel()
	state := schemaState{details: &scanner.SchemaDetails{
		Name: "spec-driven", Source: "project", Path: "/p/openspec/schemas/spec-driven",
		Shadows: []scanner.SchemaShadow{{Source: "package", Path: "/pkg/schemas/spec-driven"}},
	}}
	out := renderSchemaSection(m.projects[m.currentKey()].Info, "spec-driven", state)
	if !strings.Contains(out, "overrides: package /pkg/schemas/spec-driven") {
		t.Errorf("a project schema hiding a built-in must say so:\n%s", out)
	}
}

// --- 3.1 to 3.4 when the reads happen ---

func TestNothingIsReadUntilThePropertiesTabIsOpened(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabChanges
	m.focus = focusDetail
	m.recalcLayout()

	if cmds := m.enterTab(); cmds != nil {
		t.Errorf("got %d commands on the changes tab, want none", len(cmds))
	}
	if m.schemas != nil {
		t.Error("nothing is marked as being read")
	}
}

func TestOpeningThePropertiesTabStartsOneReadPerSchema(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabProperties
	cmds := m.enterTab()
	if len(cmds) != 2 {
		t.Fatalf("got %d commands, want one per schema in use", len(cmds))
	}
	if len(m.schemas) != 2 {
		t.Fatalf("got %d marked as reading, want 2", len(m.schemas))
	}
	for name, st := range m.schemas {
		if st.loaded() {
			t.Errorf("%s is marked finished before its command ran", name)
		}
	}
}

func TestReturningToTheTabStartsNothingFurther(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabProperties
	if n := len(m.enterTab()); n != 2 {
		t.Fatalf("first visit started %d reads, want 2", n)
	}
	if n := len(m.enterTab()); n != 0 {
		t.Errorf("second visit started %d reads, want none", n)
	}
}

func TestOpeningADifferentProjectReadsAgain(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabProperties
	m.enterTab()

	// A genuinely different project: same schema name, different key, so the
	// cache has to be dropped rather than reused.
	other := storeBackedModel()
	other.detailTab = tabProperties
	other.schemas = m.schemas
	other.schemaFor = m.schemaFor

	// Asked of the schema loader rather than of enterTab, which also refreshes
	// a store's git state now and would make this a count of two unrelated
	// things.
	cmds := other.ensureSchemasLoaded()
	if len(cmds) != 1 {
		t.Fatalf("got %d reads for the new project, want one per schema it uses", len(cmds))
	}
	if other.schemaFor != other.currentKey() {
		t.Errorf("the cache belongs to %q, want %q", other.schemaFor, other.currentKey())
	}
	if _, stale := other.schemas["tinychange"]; stale {
		t.Error("the previous project's schemas must not be carried over")
	}
}

func TestAnAnswerForAClosedProjectIsDropped(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabProperties
	m.enterTab()

	updated, _ := m.Update(schemaMsg{
		project: "/some/other/project", name: "spec-driven",
		details: &scanner.SchemaDetails{Name: "spec-driven", Source: "package"},
	})
	if updated.(model).schemas["spec-driven"].loaded() {
		t.Error("an answer that arrived after a project switch must not be shown against the new one")
	}
}

func TestAnAnswerForTheOpenProjectIsKept(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabProperties
	m.enterTab()

	updated, _ := m.Update(schemaMsg{
		project: m.currentKey(), name: "spec-driven",
		details: &scanner.SchemaDetails{Name: "spec-driven", Source: "package", Path: "/pkg"},
	})
	um := updated.(model)
	st := um.schemaStateOf("spec-driven")
	if !st.loaded() || st.details == nil || st.details.Source != "package" {
		t.Errorf("got %+v, want the details kept", st)
	}
}

func TestASchemaFileEditedWhileOpenIsNotReread(t *testing.T) {
	// Deliberate staleness. A schema is a workflow definition rather than
	// content: it changes once in a project's life, where changes and specs
	// move all day. The watcher covers the tree and will refresh the other
	// tabs; these rows keep what they read until the project is opened again.
	m := twoSchemaModel()
	m.detailTab = tabProperties
	m.enterTab()
	updated, _ := m.Update(schemaMsg{project: m.currentKey(), name: "spec-driven",
		details: &scanner.SchemaDetails{Name: "spec-driven", Source: "package"}})
	m = updated.(model)

	before := m.schemaStateOf("spec-driven")
	if cmds := m.enterTab(); len(cmds) != 0 {
		t.Errorf("a filesystem change started %d rereads, want none", len(cmds))
	}
	if after := m.schemaStateOf("spec-driven"); after.details != before.details {
		t.Error("what was read must be kept as it was")
	}
}

// --- the tab entry hooks ---

func TestNumberKeyOpeningThePropertiesTabStartsTheReads(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabChanges
	m.focus = focusDetail
	m.recalcLayout()

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '3', Text: "3"})
	um := updated.(model)
	if um.detailTab != tabProperties {
		t.Fatalf("got tab %d, want properties", um.detailTab)
	}
	if len(um.schemas) != 2 {
		t.Errorf("got %d schemas being read, want 2", len(um.schemas))
	}
	if cmd == nil {
		t.Error("the reads must actually be dispatched")
	}
}

func TestArrowIntoThePropertiesTabStartsTheReads(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabSpecs
	m.focus = focusDetail
	m.recalcLayout()

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if len(updated.(model).schemas) != 2 {
		t.Errorf("got %d schemas being read, want 2", len(updated.(model).schemas))
	}
}

// --- 4.1 the tab is named properties ---

func TestTheThirdTabIsNamedProperties(t *testing.T) {
	if tabNames[tabProperties] != "properties" {
		t.Errorf("got %q, want properties", tabNames[tabProperties])
	}
	m := twoSchemaModel()
	m.recalcLayout()
	m.syncDocument()
	out := ansi.Strip(m.renderFrame())
	if !strings.Contains(out, "properties") {
		t.Error("the tab bar must name it")
	}
	// Asked of the tab bar rather than of the whole frame: `config` is a row
	// label in this tab's own list now, and the question is what the tab is
	// called.
	if bar := ansi.Strip(m.renderTabHeader(m.width)); strings.Contains(bar, "config") {
		t.Errorf("the old name must be gone from the tab bar: %q", bar)
	}
}

// --- 6.1 and 6.2 an open change names its schema ---

func TestOpenChangeNamesItsSchema(t *testing.T) {
	r := changeRow{ci: scanner.ChangeInfo{Name: "some-change", Schema: "tinychange"}}
	out := ansi.Strip(changeDetailHeader(t, r))
	if !strings.Contains(out, "schema tinychange") {
		t.Errorf("the schema belongs beside the name:\n%s", out)
	}
}

func TestOpenChangeWithNoRecordedSchemaNamesNone(t *testing.T) {
	r := changeRow{ci: scanner.ChangeInfo{Name: "some-change"}}
	out := ansi.Strip(changeDetailHeader(t, r))
	if strings.Contains(out, "schema") {
		t.Errorf("a change that recorded none must not borrow the project default:\n%s", out)
	}
}

func changeDetailHeader(t *testing.T, r changeRow) string {
	t.Helper()
	m := twoSchemaModel()
	m.level = levelChange
	m.recalcLayout()
	return strings.SplitN(m.renderChangeDetail(r, 0, 80, 20), "\n", 2)[0]
}

// --- 4.2 and 4.6 the frame and the borders ---

func TestPropertiesFrameGeometry(t *testing.T) {
	// A nested split is exactly the shape of change a width assertion misses:
	// lipgloss pads a wrapped remainder back out to full width, so only the
	// line count catches an overflow.
	for _, size := range []struct{ w, h int }{{60, 20}, {92, 30}, {120, 50}} {
		m := twoSchemaModel()
		m.width, m.height = size.w, size.h
		m.detailTab = tabProperties
		m.recalcLayout()
		m.syncDocument()

		lines := strings.Split(m.renderFrame(), "\n")
		if len(lines) != size.h {
			t.Errorf("%dx%d: got %d rows, want %d", size.w, size.h, len(lines), size.h)
		}
		for _, l := range lines {
			if w := ansi.StringWidth(l); w > size.w {
				t.Errorf("%dx%d: a row is %d columns wide", size.w, size.h, w)
				break
			}
		}
	}
}

func TestPropertiesBordersSayWhereTheKeyboardIs(t *testing.T) {
	const lit = "\x1b[32m"
	count := func(m model) int {
		m.recalcLayout()
		m.syncDocument()
		n := 0
		for _, l := range strings.Split(m.renderFrame(), "\n") {
			n += strings.Count(l, lit+"╭") + strings.Count(l, lit+"╰")
		}
		return n
	}

	m := twoSchemaModel()
	m.detailTab = tabProperties

	m.focus = focusListPane
	onList := count(m)
	m.focus = focusContentPane
	onContent := count(m)

	if onList == 0 || onContent == 0 {
		t.Errorf("one half is always lit: list %d, content %d", onList, onContent)
	}
	if onList != onContent {
		t.Errorf("exactly one half is lit either way: list %d, content %d", onList, onContent)
	}
}

// --- 4.8 each row is a document of its own ---

func TestSelectingAnotherRowStartsAtTheTop(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabProperties
	m.recalcLayout()
	m.syncDocument()
	first := m.docKey

	m.propSection = 1
	m.syncDocument()
	if m.docKey == first {
		t.Fatal("a different row is a different document")
	}
	if m.docViewport.YOffset() != 0 {
		t.Errorf("a different document starts at the top, got offset %d", m.docViewport.YOffset())
	}
}

func TestPropertiesTabKeepsItsRowAcrossATabSwitch(t *testing.T) {
	m := twoSchemaModel()
	m.detailTab = tabProperties
	m.propSection = 2
	m.recalcLayout()
	m.syncDocument()
	before := m.docKey

	m.detailTab = tabChanges
	m.syncDocument()
	m.detailTab = tabProperties
	m.syncDocument()

	if m.docKey != before {
		t.Errorf("the same row is the same document: got %q, want %q", m.docKey, before)
	}
	if m.propSection != 2 {
		t.Errorf("the active row survives a tab switch: got %d", m.propSection)
	}
}
