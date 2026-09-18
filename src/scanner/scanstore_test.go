package scanner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// --- 2.1 and 2.2 a discovered root reports whether it is a store ---

func TestParseProjectInfoLabelsAStore(t *testing.T) {
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: "+store+"\n      remote: git@example.com:t/a.git\n")

	info := ParseProjectInfo(store)
	if info.StoreID != "alpha" {
		t.Errorf("storeID: got %q, want alpha", info.StoreID)
	}
	if info.Store == nil || info.Store.Remote != "git@example.com:t/a.git" {
		t.Errorf("registry details not carried: %+v", info.Store)
	}
	if !info.FromStore() {
		t.Error("a store reads from a store")
	}
	if info.Root != store || info.Origin != store {
		t.Errorf("root and origin: got %q and %q, want %q", info.Root, info.Origin, store)
	}
}

func TestParseProjectInfoLeavesAnOrdinaryProjectAlone(t *testing.T) {
	isolateStores(t)
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	info := ParseProjectInfo(dir)
	if info.StoreID != "" || info.Store != nil {
		t.Errorf("a project with no metadata file is not a store: %+v", info.Store)
	}
	if info.FromStore() {
		t.Error("an ordinary project does not read from a store")
	}
}

func TestParseProjectInfoKeepsARootWithCorruptStoreMetadata(t *testing.T) {
	// Unreadable metadata must hide nothing: the root is still a root, it just
	// is not reported as a store.
	isolateStores(t)
	root := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	dir := filepath.Join(root, storeMetadataDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, storeMetadataFile), []byte("id: [broken\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	info := ParseProjectInfo(root)
	if info.StoreID != "" {
		t.Errorf("got storeID %q, want none", info.StoreID)
	}
	if info.ConfigFile == "" {
		t.Error("the root must still be read as a project")
	}
}

func TestParseProjectInfoReadsTheYmlSpelling(t *testing.T) {
	isolateStores(t)
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "openspec", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, dir, "config.yml", "schema: spec-driven\n")

	info := ParseProjectInfo(dir)
	if info.ConfigFile != "config.yml" {
		t.Errorf("got %q, want config.yml", info.ConfigFile)
	}
}

// --- 2.3 a repo that only points at a store is not a row of its own ---

func TestWalkSkipsAPointerOnlyRepo(t *testing.T) {
	// This is the property that makes the picker show stores rather than every
	// repo that reads from one. It holds today by accident of the validity
	// rule; this test makes it deliberate, and fails if that rule is loosened.
	isolateStores(t)
	base := t.TempDir()

	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	mkPointer(t, filepath.Join(base, "repo-one"), "store: alpha\n")
	mkPointer(t, filepath.Join(base, "repo-two"), "store: alpha\n")

	found := walkAll(t, base)

	if len(found) != 1 {
		t.Fatalf("got %v, want the store alone", found)
	}
	if found[0] != store {
		t.Errorf("got %q, want the store %q", found[0], store)
	}
}

func TestWalkFindsAProjectUsingTheYmlSpelling(t *testing.T) {
	isolateStores(t)
	base := t.TempDir()
	dir := filepath.Join(base, "ymlproject")
	if err := os.MkdirAll(filepath.Join(dir, "openspec", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, dir, "config.yml", "schema: spec-driven\n")

	found := walkAll(t, base)
	if len(found) != 1 || found[0] != dir {
		t.Errorf("got %v, want %q", found, dir)
	}
}

func walkAll(t *testing.T, base string) []string {
	t.Helper()
	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	results := make(chan string, 16)
	go func() {
		if err := Walk(context.Background(), cfg, results, true); err != nil {
			t.Error(err)
		}
	}()
	var found []string
	for d := range results {
		found = append(found, d)
	}
	return found
}

// --- 2.4 the cache still round-trips, because its paths are roots ---

func TestCacheRoundTripsOverATreeHoldingAStore(t *testing.T) {
	isolateStores(t)
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	mkPointer(t, filepath.Join(base, "repo"), "store: alpha\n")

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	paths := walkAll(t, base)

	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := SaveCache(cachePath, cfg, paths, time.Now()); err != nil {
		t.Fatal(err)
	}
	loaded, ok := LoadCache(cachePath, cfg)
	if !ok {
		t.Fatal("cache did not load")
	}
	if len(loaded) != len(paths) {
		t.Fatalf("got %v, want %v", loaded, paths)
	}

	// A cached path is a root, so parsing it still labels the store.
	projects := ScanPaths(loaded)
	st, found := projects[store]
	if !found {
		t.Fatalf("the store is missing from %v", loaded)
	}
	if st.Info.StoreID != "alpha" {
		t.Errorf("got %q, want alpha", st.Info.StoreID)
	}
}

// --- 3.1 reading a store-backed project ---

func TestScanResolvedReadsTheStoresContent(t *testing.T) {
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n# shared\n")
	mkStoreMetadata(t, store, "alpha")
	mkChange(t, store, "do-a-thing", "- [x] 1.1 one\n- [ ] 1.2 two\n")
	mkSpec(t, store, "some-capability")
	writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: "+store+"\n")

	repo := mkPointer(t, t.TempDir(), "schema: spec-driven\nstore: alpha\ncontext: |\n  repo only\n")

	root, st, err := ScanResolved(repo)
	if err != nil {
		t.Fatal(err)
	}
	if root != store {
		t.Fatalf("root: got %q, want the store %q", root, store)
	}

	info := st.Info
	if info.SpecCount != 1 {
		t.Errorf("specs: got %d, want the store's 1, not zero", info.SpecCount)
	}
	if len(info.ActiveChanges) != 1 {
		t.Errorf("changes: got %d, want the store's 1", len(info.ActiveChanges))
	}
	if info.TasksTotal != 2 || info.TasksDone != 1 {
		t.Errorf("tasks: got %d/%d, want 1/2", info.TasksDone, info.TasksTotal)
	}
	if info.Origin != repo {
		t.Errorf("origin: got %q, want %q", info.Origin, repo)
	}
	if !info.ResolvedElsewhere() {
		t.Error("the content came from somewhere other than the origin")
	}
	if info.OriginConfigContent == "" || info.OriginConfigFile != "config.yaml" {
		t.Error("the repo's own configuration must travel with the project")
	}
	if !contains(info.OriginConfigContent, "repo only") {
		t.Error("the repo's configuration must be the repo's, not the store's")
	}
	if !contains(info.ConfigContent, "shared") {
		t.Error("the root's configuration must be the store's")
	}
	if len(st.Files) == 0 {
		t.Error("the file listing must come from the store")
	}
}

func TestScanResolvedOnAProblemStillReportsTheProject(t *testing.T) {
	isolateStores(t)
	repo := mkPointer(t, t.TempDir(), "schema: spec-driven\nstore: missing\n")

	root, st, err := ScanResolved(repo)
	if err != nil {
		t.Fatalf("an unfollowed declaration is a report, not an error: %v", err)
	}
	if root != repo {
		t.Errorf("root: got %q, want the repo %q", root, repo)
	}
	if st.Info.StoreProblem == nil {
		t.Fatal("expected the problem to travel with the project")
	}
	if st.Info.StoreProblem.ID != "missing" {
		t.Errorf("got %q, want missing", st.Info.StoreProblem.ID)
	}
}

func TestScanResolvedOnAPlainProject(t *testing.T) {
	isolateStores(t)
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkSpec(t, dir, "thing")

	root, st, err := ScanResolved(dir)
	if err != nil {
		t.Fatal(err)
	}
	if root != dir {
		t.Errorf("got %q, want %q", root, dir)
	}
	if st.Info.FromStore() || st.Info.ResolvedElsewhere() {
		t.Error("a plain project reads from itself")
	}
	if st.Info.OriginConfigFile != "" {
		t.Error("a plain project has no second configuration")
	}
	if st.Info.SpecCount != 1 {
		t.Errorf("specs: got %d, want 1", st.Info.SpecCount)
	}
}

func TestScanResolvedOutsideAnyProject(t *testing.T) {
	isolateStores(t)
	root, _, err := ScanResolved(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if root != "" {
		t.Errorf("got %q, want nothing outside a project", root)
	}
}

// --- 3.5 a repointed declaration is picked up ---

func TestScanResolvedFollowsARepointedDeclaration(t *testing.T) {
	isolateStores(t)
	alpha := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, alpha, "alpha")
	mkSpec(t, alpha, "from-alpha")

	beta := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, beta, "beta")
	mkSpec(t, beta, "from-beta")
	mkSpec(t, beta, "also-from-beta")

	writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: "+alpha+
		"\n  beta:\n    backend:\n      type: git\n      local_path: "+beta+"\n")

	repo := mkPointer(t, t.TempDir(), "store: alpha\n")
	root, st, _ := ScanResolved(repo)
	if root != alpha || st.Info.SpecCount != 1 {
		t.Fatalf("got root=%q specs=%d, want alpha with 1", root, st.Info.SpecCount)
	}

	writeConfig(t, repo, "config.yaml", "store: beta\n")
	root, st, _ = ScanResolved(repo)
	if root != beta {
		t.Errorf("root: got %q, want beta %q after repointing", root, beta)
	}
	if st.Info.SpecCount != 2 {
		t.Errorf("specs: got %d, want beta's 2", st.Info.SpecCount)
	}
}

// helpers

func mkChange(t *testing.T, root, name, tasks string) {
	t.Helper()
	dir := filepath.Join(root, "openspec", "changes", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "tasks.md"), []byte(tasks), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mkSpec(t *testing.T, root, name string) {
	t.Helper()
	dir := filepath.Join(root, "openspec", "specs", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# "+name+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
