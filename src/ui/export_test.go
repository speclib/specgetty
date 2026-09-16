package ui

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// writeFile creates a file and the directories above it.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// makeChangeDir builds a change directory with a couple of nested files, so a
// zip test can assert the internal structure survived.
func makeChangeDir(t *testing.T, dir string) {
	t.Helper()
	writeFile(t, filepath.Join(dir, "proposal.md"), "# why\n")
	writeFile(t, filepath.Join(dir, "tasks.md"), "- [x] done\n")
	writeFile(t, filepath.Join(dir, "specs", "thing", "spec.md"), "# thing\n")
}

// zipEntries reads back the names inside a zip, sorted.
func zipEntries(t *testing.T, path string) []string {
	t.Helper()
	r, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	defer r.Close()

	names := make([]string, 0, len(r.File))
	for _, f := range r.File {
		names = append(names, filepath.ToSlash(f.Name))
	}
	sort.Strings(names)
	return names
}

func runExport(t *testing.T, project, dirName, semanticName string, archived bool) exportMsg {
	t.Helper()
	msg := doExportChange(project, dirName, semanticName, archived)()
	got, ok := msg.(exportMsg)
	if !ok {
		t.Fatalf("got %T, want exportMsg", msg)
	}
	return got
}

// TestExportArchivedChangeFindsTheDatedDirectory is the regression test for a
// defect present since the feature shipped in f2f6105.
//
// The scanner stores an archived change under its display name, with the date
// prefix stripped (scan.go: displayName = dirName[11:]). The exporter was
// joining that stripped name onto the archive directory, so it looked for a
// path that never exists.
//
// openspec/specs/export-change/spec.md already required this behaviour.
func TestExportArchivedChangeFindsTheDatedDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	project := t.TempDir()
	makeChangeDir(t, filepath.Join(project, "openspec", "changes", "archive", "2026-04-02-my-feature"))

	// "my-feature" is what the scanner reports and therefore what the UI hands
	// to the exporter; the directory on disk still carries its date prefix.
	got := runExport(t, project, "2026-04-02-my-feature", "my-feature", true)
	if !got.ok {
		t.Fatalf("export failed: %s", got.output)
	}

	dest := filepath.Join(home, "my-feature-"+time.Now().Format("2006-01-02")+".zip")
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("no zip at %s: %v", dest, err)
	}

	want := []string{
		"my-feature/proposal.md",
		"my-feature/specs/thing/spec.md",
		"my-feature/tasks.md",
	}
	got2 := zipEntries(t, dest)
	if strings.Join(got2, ",") != strings.Join(want, ",") {
		t.Errorf("zip contains %v, want %v", got2, want)
	}
}

func TestExportActiveChange(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	project := t.TempDir()
	makeChangeDir(t, filepath.Join(project, "openspec", "changes", "my-feature"))

	got := runExport(t, project, "my-feature", "my-feature", false)
	if !got.ok {
		t.Fatalf("export failed: %s", got.output)
	}

	dest := filepath.Join(home, "my-feature-"+time.Now().Format("2006-01-02")+".zip")
	want := []string{
		"my-feature/proposal.md",
		"my-feature/specs/thing/spec.md",
		"my-feature/tasks.md",
	}
	if got2 := zipEntries(t, dest); strings.Join(got2, ",") != strings.Join(want, ",") {
		t.Errorf("zip contains %v, want %v", got2, want)
	}
}

func TestExportMissingChangeReportsFailure(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	project := t.TempDir()

	got := runExport(t, project, "nope", "nope", false)
	if got.ok {
		t.Error("exporting a change that does not exist should fail")
	}
	if !strings.Contains(got.output, "Source not found") {
		t.Errorf("output = %q, want it to name the missing source", got.output)
	}

	// A failed export must not leave a zip behind.
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("home contains %d file(s) after a failed export, want none", len(entries))
	}
}

func TestExportOverwritesAnExistingZip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	project := t.TempDir()
	makeChangeDir(t, filepath.Join(project, "openspec", "changes", "my-feature"))

	dest := filepath.Join(home, "my-feature-"+time.Now().Format("2006-01-02")+".zip")
	if err := os.WriteFile(dest, []byte("stale content, not a zip"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := runExport(t, project, "my-feature", "my-feature", false)
	if !got.ok {
		t.Fatalf("export failed: %s", got.output)
	}
	// The spec requires an overwrite, so the stale file must be a real zip now.
	if entries := zipEntries(t, dest); len(entries) != 3 {
		t.Errorf("zip has %d entries, want 3: the existing file was not overwritten", len(entries))
	}
}

func TestExportDestPath(t *testing.T) {
	t.Setenv("HOME", "/home/someone")
	got := exportDestPath("my-feature")
	want := filepath.Join("/home/someone", "my-feature-"+time.Now().Format("2006-01-02")+".zip")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestExportKeyPassesTheDirectoryName is the integration half of the archived
// export fix: doExportChange can only find the directory if the key handler
// hands it DirName rather than the display name.
func TestExportKeyPassesTheDirectoryName(t *testing.T) {
	m := makeListModel()
	m.listMode = modeArchived
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		ArchivedChanges: []scanner.ChangeInfo{
			{Name: "my-feature", DirName: "2026-04-02-my-feature"},
		},
	}}}
	m.syncDocument()

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	um := updated.(model)

	if um.exportState != exportConfirming {
		t.Fatalf("exportState = %d, want exportConfirming", um.exportState)
	}
	if um.exportChangeName != "my-feature" {
		t.Errorf("exportChangeName = %q, want the display name", um.exportChangeName)
	}
	if um.exportDirName != "2026-04-02-my-feature" {
		t.Errorf("exportDirName = %q, want the on-disk directory name", um.exportDirName)
	}
	if !um.exportIsArchived {
		t.Error("exportIsArchived = false, want true")
	}
}

func TestDiscardChangeMovesItUnderDiscarded(t *testing.T) {
	project := t.TempDir()
	changes := filepath.Join(project, "openspec", "changes")
	makeChangeDir(t, filepath.Join(changes, "abandoned"))

	msg := doDiscardChange(project, "abandoned")()
	got, ok := msg.(discardMsg)
	if !ok {
		t.Fatalf("got %T, want discardMsg", msg)
	}
	if !got.ok {
		t.Fatalf("discard failed: %s", got.output)
	}

	target := filepath.Join(changes, "discarded", time.Now().Format("2006-01-02")+"-abandoned")
	if _, err := os.Stat(filepath.Join(target, "proposal.md")); err != nil {
		t.Errorf("change did not land at %s: %v", target, err)
	}
	if _, err := os.Stat(filepath.Join(changes, "abandoned")); !os.IsNotExist(err) {
		t.Error("the original directory should be gone after a discard")
	}
}

func TestDiscardRefusesWhenTheTargetExists(t *testing.T) {
	project := t.TempDir()
	changes := filepath.Join(project, "openspec", "changes")
	makeChangeDir(t, filepath.Join(changes, "abandoned"))

	// Something already discarded under the same name today.
	occupied := filepath.Join(changes, "discarded", time.Now().Format("2006-01-02")+"-abandoned")
	writeFile(t, filepath.Join(occupied, "keep.md"), "do not clobber\n")

	msg := doDiscardChange(project, "abandoned")().(discardMsg)
	if msg.ok {
		t.Error("discard should refuse rather than overwrite an existing target")
	}
	if _, err := os.Stat(filepath.Join(changes, "abandoned")); err != nil {
		t.Error("the change should still be in place after a refused discard")
	}
	if b, _ := os.ReadFile(filepath.Join(occupied, "keep.md")); string(b) != "do not clobber\n" {
		t.Error("the existing discarded directory was modified")
	}
}

func TestArchiveWithoutTheOpenspecBinary(t *testing.T) {
	// The success path shells out to openspec and is left to the CLI's own
	// tests; the missing-binary branch is the one this code owns.
	t.Setenv("PATH", "")

	msg := doArchiveChange(t.TempDir(), "whatever")().(archiveMsg)
	if msg.ok {
		t.Error("archiving without the openspec binary should fail")
	}
	if !strings.Contains(msg.output, "openspec CLI not found") {
		t.Errorf("output = %q, want it to name the missing binary", msg.output)
	}
}
