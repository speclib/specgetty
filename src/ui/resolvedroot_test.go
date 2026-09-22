package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// storeBackedModel builds the case every operation here has to get right: a
// repository that declares a store, whose `openspec/` lives in the store and
// not under the repository at all.
//
// The project map is keyed by the directory the user started in, which is the
// repository. Anything that joins that key with `openspec/...` produces a path
// that does not exist, which is the defect these tests are about.
func onDiskStoreModel(t *testing.T, tasks string) (m model, origin, root, tasksPath string) {
	t.Helper()
	base := t.TempDir()
	origin = filepath.Join(base, "tunnel-repo")
	root = filepath.Join(base, "stores", "nivis-tunnel")

	active := filepath.Join(root, "openspec", "changes", "my-change")
	archived := filepath.Join(root, "openspec", "changes", "archive", "2026-04-02-old-thing")
	spec := filepath.Join(root, "openspec", "specs", "some-capability")
	for _, d := range []string{origin, active, archived, spec} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	tasksPath = filepath.Join(active, "tasks.md")
	if err := os.WriteFile(tasksPath, []byte(tasks), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{
		filepath.Join(active, "proposal.md"),
		filepath.Join(spec, "spec.md"),
		filepath.Join(root, "openspec", "project.md"),
	} {
		if err := os.WriteFile(f, []byte("# content\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	done, total := scanner.ParseTaskStats(tasks)
	m = makeListModel()
	m.width, m.height = 90, 24
	m.repoPaths = []string{origin}
	m.displayNames = []string{"tunnel-repo"}
	m.projects = scanner.ProjectMap{origin: scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Root: root, Origin: origin, StoreID: "nivis-tunnel",
		ConfigFile:    "project.md",
		ConfigContent: "# content\n",
		SpecCount:     1,
		SpecNames:     []string{"some-capability"},
		SpecContents:  map[string]string{"some-capability": "# content\n"},
		Changes: []scanner.ChangeInfo{{
			Name: "my-change", DirName: "my-change",
			ArtifactFiles: []string{"proposal.md", "tasks.md"},
			ArtifactContents: map[string]string{
				"proposal.md": "# content\n",
				"tasks.md":    tasks,
			},
			TasksDone: done, TasksTotal: total,
		}},
		ArchivedChanges: []scanner.ChangeInfo{{
			Name: "old-thing", DirName: "2026-04-02-old-thing",
		}},
	}}}
	m.recalcLayout()
	m.syncDocument()
	return m, origin, root, tasksPath
}

// moveToTask puts the document cursor on the first task.
//
// The cursor selects tasks and nothing else, so it is already there. The helper
// stays because it says what the tests below depend on, and fails loudly on a
// document that has no task to select.
func moveToTask(t *testing.T, m model) model {
	t.Helper()
	if _, ok := m.selectedTask(); !ok {
		t.Fatal("no task found in the document")
	}
	return m
}

// TestToggleWritesTheFileUnderTheStore is task 1.1.
//
// It fails before the fix with "could not save: no such file or directory",
// because the path was built from the repository the user started in.
func TestToggleWritesTheFileUnderTheStore(t *testing.T) {
	m, origin, root, tasksPath := onDiskStoreModel(t, threeTasks)
	m.level = levelChange
	m.changeArtifactTab = 1 // tasks.md
	m.recalcLayout()
	m.syncDocument()

	if m.docPath != tasksPath {
		t.Errorf("the document's path is %q, want the file under the store %q",
			m.docPath, tasksPath)
	}
	if strings.HasPrefix(m.docPath, origin+string(filepath.Separator)) {
		t.Errorf("the path is under the repository the user started in: %q", m.docPath)
	}

	if !m.docHasCursor() {
		t.Fatal("the tasks artifact must have a cursor for space to act on")
	}
	um := press(moveToTask(t, m), tea.KeyPressMsg{Code: ' ', Text: " "})
	if strings.Contains(um.statusMsg, "could not save") {
		t.Fatalf("the toggle failed: %s", um.statusMsg)
	}

	after, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), "- [x] 1.1 first") {
		t.Errorf("the file under the store was not written:\n%s", after)
	}
	// And nothing was created under the repository, which has no openspec/ at
	// all and must not grow one.
	if _, err := os.Stat(filepath.Join(origin, "openspec")); !os.IsNotExist(err) {
		t.Error("a directory was created under the pointing repository")
	}
	_ = root
}

// TestCopiedPathIsUnderTheStore is task 1.2.
func TestCopiedPathIsUnderTheStore(t *testing.T) {
	got := fakeClipboard(t, nil)
	m, origin, root, _ := onDiskStoreModel(t, threeTasks)

	press(m, tea.KeyPressMsg{Code: 'Y', Text: "Y"})

	want := filepath.Join(root, "openspec", "changes", "my-change")
	if *got != want {
		t.Errorf("copied %q, want %q", *got, want)
	}
	if strings.HasPrefix(*got, origin+string(filepath.Separator)) {
		t.Errorf("the copied path is under the repository the user started in: %q", *got)
	}
	if _, err := os.Stat(*got); err != nil {
		t.Errorf("the copied path does not resolve: %v", err)
	}
}

func TestCopiedPathOfAnArchivedChangeIsUnderTheStore(t *testing.T) {
	got := fakeClipboard(t, nil)
	m, _, root, _ := onDiskStoreModel(t, threeTasks)
	m.changeCursor = len(m.allRows()) - 1
	m.rememberSelection()

	press(m, tea.KeyPressMsg{Code: 'Y', Text: "Y"})

	want := filepath.Join(root, "openspec", "changes", "archive", "2026-04-02-old-thing")
	if *got != want {
		t.Errorf("copied %q, want %q", *got, want)
	}
	if _, err := os.Stat(*got); err != nil {
		t.Errorf("the copied path does not resolve: %v", err)
	}
}

// TestTheRootRuleHoldsWhereThereIsNoStore is task 1.4: the same operations on a
// project whose root and starting directory are one directory.
func TestTheRootRuleHoldsWhereThereIsNoStore(t *testing.T) {
	clip := fakeClipboard(t, nil)
	m, path := taskModel(t, threeTasks)
	project := m.repoPaths[0]

	if m.docPath != path {
		t.Errorf("the document's path is %q, want %q", m.docPath, path)
	}
	um := press(moveToTask(t, m), tea.KeyPressMsg{Code: ' ', Text: " "})
	if strings.Contains(um.statusMsg, "could not save") {
		t.Fatalf("the toggle failed: %s", um.statusMsg)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), "- [x] 1.1 first") {
		t.Errorf("the file was not written:\n%s", after)
	}

	list := um
	list.level = levelProject
	list.syncDocument()
	press(list, tea.KeyPressMsg{Code: 'Y', Text: "Y"})
	if want := filepath.Join(project, "openspec", "changes", "my-change"); *clip != want {
		t.Errorf("copied %q, want %q", *clip, want)
	}
}

// --- 2.x which file a pane is showing ---

func TestEveryArtifactOfAChangeNamesItsFile(t *testing.T) {
	m, _, root, _ := onDiskStoreModel(t, threeTasks)
	m.level = levelChange
	m.recalcLayout()

	want := map[int]string{
		0: filepath.Join(root, "openspec", "changes", "my-change", "proposal.md"),
		1: filepath.Join(root, "openspec", "changes", "my-change", "tasks.md"),
	}
	for tab, path := range want {
		m.changeArtifactTab = tab
		m.syncDocument()
		if m.docPath != path {
			t.Errorf("artifact tab %d: path %q, want %q", tab, m.docPath, path)
		}
		if _, err := os.Stat(m.docPath); err != nil {
			t.Errorf("artifact tab %d: the path does not resolve: %v", tab, err)
		}
	}
}

func TestASpecNamesItsFile(t *testing.T) {
	m, _, root, _ := onDiskStoreModel(t, threeTasks)
	m.detailTab = tabSpecs
	m.recalcLayout()
	m.syncDocument()

	want := filepath.Join(root, "openspec", "specs", "some-capability", "spec.md")
	if m.docPath != want {
		t.Errorf("path %q, want %q", m.docPath, want)
	}
	if _, err := os.Stat(m.docPath); err != nil {
		t.Errorf("the path does not resolve: %v", err)
	}
}

func TestTheProjectRowNamesItsConfigurationFile(t *testing.T) {
	for _, name := range []string{"project.md", "config.yaml"} {
		m, _, root, _ := onDiskStoreModel(t, threeTasks)
		info := m.projects[m.repoPaths[0]].Info
		info.ConfigFile = name
		m.projects[m.repoPaths[0]] = scanner.ProjectStatus{Info: info}
		m.detailTab = tabProperties
		m.recalcLayout()
		m.syncDocument()

		if want := filepath.Join(root, "openspec", name); m.docPath != want {
			t.Errorf("%s: path %q, want %q", name, m.docPath, want)
		}
	}
}

func TestAPaneWithNoSingleFileNamesNone(t *testing.T) {
	t.Run("a change's specs sub-tab", func(t *testing.T) {
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		info := m.projects[m.repoPaths[0]].Info
		info.Changes[0].SpecNames = []string{"one", "two"}
		info.Changes[0].SpecContents = map[string]string{"one": "# a\n", "two": "# b\n"}
		m.projects[m.repoPaths[0]] = scanner.ProjectStatus{Info: info}
		m.level = levelChange
		r, ok := m.selectedRow()
		if !ok {
			t.Fatal("no change selected")
		}
		m.changeArtifactTab = len(r.artifactTabNames()) - 1
		m.recalcLayout()
		m.syncDocument()

		if m.docPath != "" {
			t.Errorf("path %q, want none: the pane shows every delta at once", m.docPath)
		}
	})

	t.Run("the schema and store rows", func(t *testing.T) {
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		info := m.projects[m.repoPaths[0]].Info
		info.SchemaUsage = []scanner.SchemaUsage{{Name: "spec-driven", IsDefault: true, Changes: 1}}
		m.projects[m.repoPaths[0]] = scanner.ProjectStatus{Info: info}
		m.detailTab = tabProperties
		m.recalcLayout()

		sections := m.currentSections()
		var checked int
		for i, sec := range sections {
			if sec.kind == sectionConfig {
				continue
			}
			m.propSection = i
			m.syncDocument()
			if m.docPath != "" {
				t.Errorf("row %q: path %q, want none: it is a report, not a file",
					sec.label, m.docPath)
			}
			checked++
		}
		if checked < 2 {
			t.Fatalf("checked %d rows, want the schema row and the store row", checked)
		}
	})
}

// TestNoPaneGainedACursor is task 2.5 and task 6.2.
//
// `document.path` now means "the file this came from" rather than "the file a
// toggle writes to". The toggle's guard is `docHasCursor()`, which also wants
// mapped source lines, and only the tasks artifact has those.
func TestNoPaneGainedACursor(t *testing.T) {
	type pane struct {
		name string
		set  func(m *model)
		want bool
	}
	panes := []pane{
		{"a change's proposal", func(m *model) {
			m.level = levelChange
			m.changeArtifactTab = 0
		}, false},
		{"a change's tasks", func(m *model) {
			m.level = levelChange
			m.changeArtifactTab = 1
		}, true},
		{"the specs tab", func(m *model) { m.detailTab = tabSpecs }, false},
		{"the properties tab", func(m *model) { m.detailTab = tabProperties }, false},
		{"the change list", func(m *model) { m.detailTab = tabChanges }, false},
	}

	for _, p := range panes {
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		p.set(&m)
		m.focus = focusContentPane
		m.recalcLayout()
		m.syncDocument()

		if got := m.docHasCursor(); got != p.want {
			t.Errorf("%s: docHasCursor() = %v, want %v (path %q, %d tasks)",
				p.name, got, p.want, m.docPath, m.docTasks.count())
		}
	}
}

func TestSpaceStillWritesNothingOnAPaneWithAPathButNoCursor(t *testing.T) {
	m, _, root, _ := onDiskStoreModel(t, threeTasks)
	m.detailTab = tabSpecs
	m.recalcLayout()
	m.syncDocument()

	specPath := filepath.Join(root, "openspec", "specs", "some-capability", "spec.md")
	before, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	if m.docPath != specPath {
		t.Fatalf("the test needs the spec's path set, got %q", m.docPath)
	}

	press(m, tea.KeyPressMsg{Code: ' ', Text: " "})

	after, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("space wrote to a pane that has no task cursor")
	}
}
