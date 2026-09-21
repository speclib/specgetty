package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// mkChangeWithSchema creates a change directory recording a schema. An empty
// schema means no metadata file at all, which is a different thing.
func mkChangeWithSchema(t *testing.T, root, where, name, schema string) {
	t.Helper()
	dir := filepath.Join(root, "openspec", "changes", where, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "tasks.md"), []byte("- [x] done\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if schema == "" {
		return
	}
	body := "schema: " + schema + "\ncreated: 2026-09-21\n"
	if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// --- 1.1 and 1.2 a change's recorded schema ---

func TestChangeRecordsItsSchema(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkChangeWithSchema(t, root, ".", "with-schema", "tinychange")
	mkChangeWithSchema(t, root, ".", "without-schema", "")

	info := ParseProjectInfo(root)
	got := map[string]string{}
	for _, ci := range info.Changes {
		got[ci.Name] = ci.Schema
	}
	if got["with-schema"] != "tinychange" {
		t.Errorf("got %q, want tinychange", got["with-schema"])
	}
	if got["without-schema"] != "" {
		t.Errorf("got %q, want empty for a change with no metadata file", got["without-schema"])
	}
}

func TestChangeWithUnreadableMetadata(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkChangeWithSchema(t, root, ".", "broken", "")
	bad := filepath.Join(root, "openspec", "changes", "broken", ".openspec.yaml")
	if err := os.WriteFile(bad, []byte("schema: [unclosed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	info := ParseProjectInfo(root)
	if len(info.Changes) != 1 || info.Changes[0].Schema != "" {
		t.Errorf("unreadable metadata reads as unrecorded: %+v", info.Changes)
	}
}

// --- 1.3 to 1.5 the usage tally ---

func TestSchemaUsageCountsEveryChange(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkChangeWithSchema(t, root, ".", "a", "spec-driven")
	mkChangeWithSchema(t, root, ".", "b", "tinychange")
	mkChangeWithSchema(t, root, ".", "c", "")
	mkChangeWithSchema(t, root, "archive", "2026-01-01-d", "tinychange")

	info := ParseProjectInfo(root)
	got := map[string]int{}
	for _, u := range info.SchemaUsage {
		got[u.Name] = u.Changes
	}
	if got["spec-driven"] != 1 {
		t.Errorf("spec-driven: got %d, want 1", got["spec-driven"])
	}
	if got["tinychange"] != 2 {
		t.Errorf("tinychange: got %d, want 2 across active and archived", got["tinychange"])
	}
	if info.UnrecordedChanges != 1 {
		t.Errorf("unrecorded: got %d, want 1", info.UnrecordedChanges)
	}
}

func TestSchemaUsageLeadsWithTheDefault(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: zzz-last-alphabetically\n")
	mkChangeWithSchema(t, root, ".", "a", "aaa-first-alphabetically")

	info := ParseProjectInfo(root)
	if len(info.SchemaUsage) != 2 {
		t.Fatalf("got %+v, want both", info.SchemaUsage)
	}
	if !info.SchemaUsage[0].IsDefault || info.SchemaUsage[0].Name != "zzz-last-alphabetically" {
		t.Errorf("the default leads: got %+v", info.SchemaUsage)
	}
}

func TestSchemaUsageIncludesAnUnusedDefault(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkChangeWithSchema(t, root, ".", "a", "tinychange")

	info := ParseProjectInfo(root)
	var found bool
	for _, u := range info.SchemaUsage {
		if u.Name == "spec-driven" {
			found = true
			if u.Changes != 0 || !u.IsDefault {
				t.Errorf("got %+v, want the default with no changes", u)
			}
		}
	}
	if !found {
		t.Errorf("the project default must be reported even when unused: %+v", info.SchemaUsage)
	}
}

func TestDefaultSchemaFallsBackToSpecDriven(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "# no schema key here\n")
	if got := ParseProjectInfo(root).DefaultSchema; got != DefaultSchemaName {
		t.Errorf("got %q, want %q", got, DefaultSchemaName)
	}
}

func TestDefaultSchemaComesFromTheResolvedRoot(t *testing.T) {
	// OpenSpec reads the schema from the root and nowhere else, so a pointing
	// repo's own `schema:` key never decides anything.
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: store-schema\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	repo := mkPointer(t, t.TempDir(), "schema: repo-schema\nstore: alpha\n")

	_, st, err := ScanResolved(repo)
	if err != nil {
		t.Fatal(err)
	}
	if st.Info.DefaultSchema != "store-schema" {
		t.Errorf("got %q, want the store's schema", st.Info.DefaultSchema)
	}
}

// --- the inert keys ---

func TestInertKeysOfADeclaringRepo(t *testing.T) {
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	repo := mkPointer(t, t.TempDir(),
		"schema: spec-driven\nstore: alpha\ncontext: |\n  x\nrules:\n  proposal:\n    - y\n")

	_, st, err := ScanResolved(repo)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"schema": true, "context": true, "rules": true}
	if len(st.Info.InertKeys) != len(want) {
		t.Fatalf("got %v, want %v", st.Info.InertKeys, want)
	}
	for _, k := range st.Info.InertKeys {
		if !want[k] {
			t.Errorf("unexpected inert key %q", k)
		}
	}
}

func TestNoInertKeysWhenOnlyTheStoreIsDeclared(t *testing.T) {
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	repo := mkPointer(t, t.TempDir(), "store: alpha\n")

	_, st, _ := ScanResolved(repo)
	if len(st.Info.InertKeys) != 0 {
		t.Errorf("got %v, want none", st.Info.InertKeys)
	}
}

func TestNoInertKeysForAPlainProject(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: spec-driven\ncontext: |\n  used\n")
	_, st, _ := ScanResolved(root)
	if len(st.Info.InertKeys) != 0 {
		t.Errorf("a project that is its own root uses its own config: %v", st.Info.InertKeys)
	}
}

// --- 2.x resolving a definition ---

// fakeCLI puts a scripted `openspec` at the front of PATH.
func fakeCLI(t *testing.T, script string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "openspec")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return dir
}

// mkSchemaDir writes a schema definition and returns its directory.
func mkSchemaDir(t *testing.T, name, body string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "schema.yaml"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

const sampleSchema = `name: spec-driven
version: 1
description: Default OpenSpec workflow
artifacts:
  - id: proposal
    generates: proposal.md
    template: proposal.md
    description: Why this change is needed
    requires: []
  - id: tasks
    generates: tasks.md
    template: tasks.md
    requires:
      - proposal
apply:
  requires:
    - tasks
  tracks: tasks.md
`

func TestResolveSchemaReadsTheLocatedDefinition(t *testing.T) {
	schemaDir := mkSchemaDir(t, "spec-driven", sampleSchema)
	fakeCLI(t, `printf '{"name":"spec-driven","source":"package","path":"`+schemaDir+`","shadows":[]}'`)

	d, p := ResolveSchema(t.TempDir(), "spec-driven")
	if p != nil {
		t.Fatalf("unexpected problem: %+v", p)
	}
	if d.Source != "package" || d.Path != schemaDir {
		t.Errorf("location not carried: %+v", d)
	}
	if d.Description != "Default OpenSpec workflow" {
		t.Errorf("description: got %q", d.Description)
	}
	if len(d.Artifacts) != 2 {
		t.Fatalf("got %d artifacts, want 2", len(d.Artifacts))
	}
	if d.Artifacts[1].ID != "tasks" || d.Artifacts[1].Generates != "tasks.md" {
		t.Errorf("artifact: %+v", d.Artifacts[1])
	}
	if len(d.Artifacts[1].Requires) != 1 || d.Artifacts[1].Requires[0] != "proposal" {
		t.Errorf("requires: %+v", d.Artifacts[1].Requires)
	}
	if d.Apply.Tracks != "tasks.md" || len(d.Apply.Requires) != 1 {
		t.Errorf("apply: %+v", d.Apply)
	}
	if d.Overrides() {
		t.Error("nothing is shadowed here")
	}
}

func TestResolveSchemaRunsWhereTheContentLives(t *testing.T) {
	root := t.TempDir()
	schemaDir := mkSchemaDir(t, "s", sampleSchema)
	record := filepath.Join(t.TempDir(), "args")
	fakeCLI(t, `{ pwd; echo "$@"; } > `+record+`
printf '{"name":"s","source":"project","path":"`+schemaDir+`","shadows":[]}'`)

	if _, p := ResolveSchema(root, "s"); p != nil {
		t.Fatalf("unexpected problem: %+v", p)
	}
	b, err := os.ReadFile(record)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if got, _ := filepath.EvalSymlinks(lines[0]); got != mustEval(t, root) {
		t.Errorf("working directory: got %q, want the resolved root %q", lines[0], root)
	}
	if lines[1] != "schema which s --json" {
		t.Errorf("arguments: got %q", lines[1])
	}
}

func mustEval(t *testing.T, p string) string {
	t.Helper()
	r, err := filepath.EvalSymlinks(p)
	if err != nil {
		return p
	}
	return r
}

func TestResolveSchemaReportsAnOverride(t *testing.T) {
	schemaDir := mkSchemaDir(t, "spec-driven", sampleSchema)
	fakeCLI(t, `printf '{"name":"spec-driven","source":"project","path":"`+schemaDir+
		`","shadows":[{"source":"package","path":"/pkg/schemas/spec-driven"}]}'`)

	d, p := ResolveSchema(t.TempDir(), "spec-driven")
	if p != nil {
		t.Fatalf("unexpected problem: %+v", p)
	}
	if !d.Overrides() || d.Shadows[0].Source != "package" {
		t.Errorf("shadows not carried: %+v", d.Shadows)
	}
}

func TestResolveSchemaProblems(t *testing.T) {
	t.Run("no binary on PATH", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		_, p := ResolveSchema(t.TempDir(), "spec-driven")
		if p == nil || p.Code != SchemaNoCLI {
			t.Fatalf("got %+v, want a no-CLI problem", p)
		}
		if p.Name != "spec-driven" {
			t.Errorf("the schema must be named: %+v", p)
		}
	})

	t.Run("a schema that does not resolve", func(t *testing.T) {
		fakeCLI(t, `printf '{"error":"Schema '"'"'nope'"'"' not found","available":["spec-driven","tinychange"]}'
exit 1`)
		_, p := ResolveSchema(t.TempDir(), "nope")
		if p == nil || p.Code != SchemaNotFound {
			t.Fatalf("got %+v, want a not-found problem", p)
		}
		if len(p.Available) != 2 {
			t.Errorf("the alternatives the CLI supplies must be kept: %+v", p.Available)
		}
		if !strings.Contains(p.Detail, "not found") {
			t.Errorf("detail: %q", p.Detail)
		}
	})

	t.Run("the command does not answer", func(t *testing.T) {
		fakeCLI(t, "sleep 5")
		start := time.Now()
		_, p := resolveSchemaWith(t.TempDir(), "slow", 200*time.Millisecond)
		if p == nil || p.Code != SchemaTimeout {
			t.Fatalf("got %+v, want a timeout", p)
		}
		if elapsed := time.Since(start); elapsed > 3*time.Second {
			t.Errorf("waited %v, want the bound to hold", elapsed)
		}
	})

	t.Run("output that is not JSON", func(t *testing.T) {
		fakeCLI(t, `echo "something went wrong"; exit 1`)
		_, p := ResolveSchema(t.TempDir(), "s")
		if p == nil || p.Code != SchemaNotFound {
			t.Fatalf("got %+v, want a not-found problem", p)
		}
	})

	t.Run("a definition that cannot be parsed", func(t *testing.T) {
		dir := mkSchemaDir(t, "broken", "artifacts: [unclosed\n")
		fakeCLI(t, `printf '{"name":"broken","source":"project","path":"`+dir+`","shadows":[]}'`)
		_, p := ResolveSchema(t.TempDir(), "broken")
		if p == nil || p.Code != SchemaUnreadable {
			t.Fatalf("got %+v, want an unreadable problem", p)
		}
		if !strings.Contains(p.Detail, dir) {
			t.Errorf("the path must be named: %q", p.Detail)
		}
	})

	t.Run("a definition that is not there", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "gone")
		fakeCLI(t, `printf '{"name":"s","source":"project","path":"`+missing+`","shadows":[]}'`)
		_, p := ResolveSchema(t.TempDir(), "s")
		if p == nil || p.Code != SchemaUnreadable {
			t.Fatalf("got %+v, want an unreadable problem", p)
		}
	})
}

func TestResolveSchemaAgainstTheRealCLI(t *testing.T) {
	// The definitions this project actually uses, read through the real binary.
	// Skipped where it is not installed, so it never gates a build.
	if _, err := os.Stat("../../openspec/schemas/tinychange/schema.yaml"); err != nil {
		t.Skip("not this project")
	}
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"spec-driven", "tinychange"} {
		d, p := ResolveSchema(root, name)
		if p != nil {
			if p.Code == SchemaNoCLI {
				t.Skip("openspec CLI not installed")
			}
			t.Fatalf("%s: %+v", name, p)
		}
		if d.Name != name || len(d.Artifacts) == 0 || d.Apply.Tracks == "" {
			t.Errorf("%s came back thin: %+v", name, d)
		}
	}
}

func TestSchemaProblemErrorIsItsDetail(t *testing.T) {
	p := &SchemaProblem{Code: SchemaNotFound, Name: "x", Detail: "it is not there"}
	if p.Error() != "it is not there" {
		t.Errorf("got %q", p.Error())
	}
	var nilProblem *SchemaProblem
	if nilProblem.Error() != "" {
		t.Error("a nil problem says nothing")
	}
}

func TestReadConfigSchemaEdgeCases(t *testing.T) {
	t.Run("no configuration at all", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0o755); err != nil {
			t.Fatal(err)
		}
		if got := readConfigSchema(dir); got != DefaultSchemaName {
			t.Errorf("got %q, want the fallback", got)
		}
	})

	t.Run("a configuration that is not YAML", func(t *testing.T) {
		dir := t.TempDir()
		writeConfig(t, dir, "config.yaml", "schema: [unclosed\n")
		if got := readConfigSchema(dir); got != DefaultSchemaName {
			t.Errorf("got %q, want the fallback", got)
		}
	})

	t.Run("an empty schema key", func(t *testing.T) {
		dir := t.TempDir()
		writeConfig(t, dir, "config.yaml", "schema: \"\"\n")
		if got := readConfigSchema(dir); got != DefaultSchemaName {
			t.Errorf("got %q, want the fallback", got)
		}
	})
}

func TestInertConfigKeysEdgeCases(t *testing.T) {
	t.Run("no configuration", func(t *testing.T) {
		dir := t.TempDir()
		if got := inertConfigKeys(dir); got != nil {
			t.Errorf("got %v, want none", got)
		}
	})

	t.Run("a configuration that is not YAML", func(t *testing.T) {
		dir := t.TempDir()
		writeConfig(t, dir, "config.yaml", "store: [unclosed\n")
		if got := inertConfigKeys(dir); got != nil {
			t.Errorf("got %v, want none", got)
		}
	})

	t.Run("references counts as inert too", func(t *testing.T) {
		// The one key OpenSpec's own doctor warns about, reported here beside
		// the four it says nothing about.
		dir := t.TempDir()
		writeConfig(t, dir, "config.yaml", "store: alpha\nreferences:\n  - other\n")
		got := inertConfigKeys(dir)
		if len(got) != 1 || got[0] != "references" {
			t.Errorf("got %v, want references", got)
		}
	})
}

func TestCountSchemaUsageWithNoChangesAtAll(t *testing.T) {
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	info := ParseProjectInfo(root)
	if len(info.SchemaUsage) != 1 || info.SchemaUsage[0].Changes != 0 {
		t.Errorf("got %+v, want the default alone with no changes", info.SchemaUsage)
	}
	if info.UnrecordedChanges != 0 {
		t.Errorf("got %d unrecorded, want none", info.UnrecordedChanges)
	}
}
