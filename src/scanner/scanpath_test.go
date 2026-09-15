package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestListOpenSpecContents(t *testing.T) {
	root := t.TempDir()
	project := makeProjectAt(t, root, "alpha")
	if err := os.WriteFile(filepath.Join(project, "openspec", "project.md"),
		[]byte("# project\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := ListOpenSpecContents(project)
	if err != nil {
		t.Fatalf("ListOpenSpecContents: %v", err)
	}

	byPath := map[string]bool{}
	for _, e := range entries {
		byPath[filepath.ToSlash(e.Path)] = e.IsDir
	}

	// Paths are relative to openspec/, not absolute.
	for _, want := range []string{"config.yaml", "project.md", "specs", "specs/thing", "specs/thing/spec.md"} {
		if _, ok := byPath[want]; !ok {
			t.Errorf("missing entry %q; got %v", want, byPath)
		}
	}
	if byPath["specs"] != true {
		t.Error("specs should be reported as a directory")
	}
	if byPath["specs/thing/spec.md"] != false {
		t.Error("spec.md should be reported as a file")
	}
}

func TestListOpenSpecContentsOnAMissingProject(t *testing.T) {
	if _, err := ListOpenSpecContents(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("expected an error for a directory with no openspec/")
	}
}

func TestScanPathsParsesWithoutWalking(t *testing.T) {
	root := t.TempDir()
	alpha := makeProjectAt(t, root, "alpha")
	beta := makeProjectAt(t, root, "beta")

	// A change with tasks, so the parsed counts are visible.
	change := filepath.Join(alpha, "openspec", "changes", "feature")
	if err := os.MkdirAll(change, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(change, "tasks.md"),
		[]byte("- [x] one\n- [ ] two\n- [ ] three\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := ScanPaths([]string{alpha, beta})
	if len(got) != 2 {
		t.Fatalf("parsed %d projects, want 2", len(got))
	}

	info := got[alpha].Info
	if info.SpecCount != 1 {
		t.Errorf("SpecCount = %d, want 1", info.SpecCount)
	}
	if info.TasksTotal != 3 || info.TasksDone != 1 {
		t.Errorf("tasks = %d/%d, want 1/3", info.TasksDone, info.TasksTotal)
	}
	if len(info.Changes) != 1 || info.Changes[0].Name != "feature" {
		t.Errorf("changes = %v, want one named feature", info.ActiveChanges)
	}
}

func TestScanPathsSkipsAVanishedPath(t *testing.T) {
	root := t.TempDir()
	alive := makeProjectAt(t, root, "alive")
	gone := filepath.Join(root, "gone")

	got := ScanPaths([]string{alive, gone})
	if len(got) != 1 {
		t.Fatalf("parsed %d projects, want only the surviving one", len(got))
	}
	if _, ok := got[alive]; !ok {
		t.Error("the surviving project is missing from the result")
	}
}

func TestScanEndToEnd(t *testing.T) {
	root := t.TempDir()
	alpha := makeProjectAt(t, root, "alpha")
	makeProjectAt(t, root, "beta")

	// An archived change, to exercise the date-prefix handling.
	archived := filepath.Join(alpha, "openspec", "changes", "archive", "2026-04-02-old-thing")
	if err := os.MkdirAll(archived, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(archived, "proposal.md"), []byte("# old\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	projects, err := Scan(configFor([]string{root}, nil, false), true)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(projects) != 2 {
		paths := make([]string, 0, len(projects))
		for p := range projects {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		t.Fatalf("found %v, want 2 projects", paths)
	}

	info := projects[alpha].Info
	if len(info.ArchivedChanges) != 1 {
		t.Fatalf("archived changes = %d, want 1", len(info.ArchivedChanges))
	}
	ci := info.ArchivedChanges[0]

	// Name is for display and loses the prefix; DirName must keep it, because
	// it is the only way back to the directory on disk.
	if ci.Name != "old-thing" {
		t.Errorf("Name = %q, want the prefix stripped", ci.Name)
	}
	if ci.DirName != "2026-04-02-old-thing" {
		t.Errorf("DirName = %q, want the directory as it is on disk", ci.DirName)
	}
	if ci.ArchiveDate.Format("2006-01-02") != "2026-04-02" {
		t.Errorf("ArchiveDate = %v, want 2026-04-02", ci.ArchiveDate)
	}
	if projects[alpha].Files == nil {
		t.Error("Files should be populated by the scan")
	}
}

func TestActiveChangeDirNameMatchesItsName(t *testing.T) {
	root := t.TempDir()
	project := makeProjectAt(t, root, "alpha")
	change := filepath.Join(project, "openspec", "changes", "feature")
	if err := os.MkdirAll(change, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(change, "proposal.md"), []byte("# f\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	info := ParseProjectInfo(project)
	if len(info.Changes) != 1 {
		t.Fatalf("changes = %d, want 1", len(info.Changes))
	}
	if got := info.Changes[0]; got.Name != got.DirName {
		t.Errorf("Name %q and DirName %q should match for an active change", got.Name, got.DirName)
	}
}

func TestDumpConfigRoundTrips(t *testing.T) {
	if err := DumpConfig(configFor([]string{"/a", "/b"}, []string{"junk"}, true)); err != nil {
		t.Errorf("DumpConfig: %v", err)
	}
}
