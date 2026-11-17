package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

const projectsLive = `# projects Specification

## Purpose
Hold the projects of one tenancy.

## Requirements

### Requirement: Projects belong to a tenancy
The system SHALL keep every project inside a tenancy.

#### Scenario: A tenancy is opened
- **WHEN** a tenancy is opened
- **THEN** its projects SHALL be listed
`

const projectsDelta = `# projects Delta

## MODIFIED Requirements

### Requirement: Projects belong to a tenancy
The system SHALL keep every project inside a tenancy, and SHALL name it.

#### Scenario: A tenancy is opened
- **WHEN** a tenancy is opened
- **THEN** its projects SHALL be listed
`

func nestedProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"openspec/config.yaml":                                          "schema: spec-driven\n",
		"openspec/specs/tenancy/projects/spec.md":                       projectsLive,
		"openspec/specs/tenancy/identity/spec.md":                       "# identity Specification\n",
		"openspec/specs/content/pages/spec.md":                          "# pages Specification\n",
		"openspec/specs/archive-from-ui/spec.md":                        "# archive-from-ui Specification\n",
		"openspec/changes/nested-change/proposal.md":                    "# Proposal\n",
		"openspec/changes/nested-change/specs/tenancy/projects/spec.md": projectsDelta,
	}
	for path, body := range files {
		full := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestNestedProjectFromDiskToScreen(t *testing.T) {
	dir := nestedProject(t)

	var m model
	msg := m.doScanSingle(dir)()
	sm, ok := msg.(scanMsg)
	if !ok || sm.err != nil {
		t.Fatalf("scan: %+v", msg)
	}
	st, found := sm.projects[dir]
	if !found {
		t.Fatalf("not filed under %q: %v", dir, keysOfMap(sm.projects))
	}
	info := st.Info

	if info.SpecCount != 4 {
		t.Fatalf("SpecCount = %d, want 4: %v", info.SpecCount, info.SpecNames)
	}
	want := []string{
		"archive-from-ui", "content/pages", "tenancy/identity", "tenancy/projects",
	}
	if strings.Join(info.SpecNames, " ") != strings.Join(want, " ") {
		t.Fatalf("spec names = %v, want %v", info.SpecNames, want)
	}

	m.projects = sm.projects
	m.repoPaths = []string{dir}
	m.level = levelProject
	m.detailTab = tabSpecs
	m.width, m.height = 100, 30
	m.focus = m.defaultFocus()
	m.recalcLayout()
	m.syncDocument()

	out := ansi.Strip(m.renderSpecsTab(m.width-2, 20))
	for _, shown := range []string{"archive-from-ui", "content (1)", "tenancy (2)", "projects", "identity"} {
		if !strings.Contains(out, shown) {
			t.Errorf("the list does not show %q:\n%s", shown, out)
		}
	}
	if strings.Contains(out, "tenancy/projects") {
		t.Errorf("a row still carries its group:\n%s", out)
	}
	if strings.TrimSpace(m.docViewport.View()) == "" {
		t.Error("the selected capability shows no content")
	}

	var change *scanner.ChangeInfo
	for i := range info.Changes {
		if info.Changes[i].Name == "nested-change" {
			change = &info.Changes[i]
		}
	}
	if change == nil {
		t.Fatalf("change not scanned: %v", info.ActiveChanges)
	}
	if len(change.SpecNames) != 1 || change.SpecNames[0] != "tenancy/projects" {
		t.Fatalf("delta names = %v, want tenancy/projects", change.SpecNames)
	}

	live := map[string]string{}
	for _, n := range info.SpecNames {
		live[n] = info.SpecContents[n]
	}
	tree, problems := buildChangeSpecTree(change.Name, change.SpecNames, change.SpecContents, live)
	if len(problems) > 0 {
		t.Fatalf("outline problems: %v", problems)
	}

	capability, marked := "", 0
	for _, n := range tree.nodes {
		switch n.kind {
		case nodeCapability:
			capability = n.title
		case nodeRequirement:
			if n.op != opModified {
				continue
			}
			if n.old == nil {
				t.Errorf("MODIFIED %q has no original, so it is drawn unmarked", n.title)
			}
			if sigil, _ := outlineMark(n); sigil != "~" {
				t.Errorf("MODIFIED %q is drawn %q, want ~", n.title, sigil)
			}
			marked++
		}
	}
	if capability != "tenancy/projects" {
		t.Errorf("delta capability = %q, want tenancy/projects", capability)
	}
	if marked != 1 {
		t.Errorf("marked requirements = %d, want 1", marked)
	}
}
