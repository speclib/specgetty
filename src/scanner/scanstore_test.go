package scanner

import (
	"context"
	"os"
	"path/filepath"
	"sort"
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

func TestWalkListsPointerOnlyReposAndSkipsTheStore(t *testing.T) {
	// The rule that decides what the picker is a list of. A repo declaring a
	// store is where a person works, and is the only place that repo's own
	// context and rules can be read from; the store is where content lives and
	// already has a row for every repo reading it.
	isolateStores(t)
	base := t.TempDir()

	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	one := mkPointer(t, filepath.Join(base, "repo-one"), "store: alpha\n")
	two := mkPointer(t, filepath.Join(base, "repo-two"), "store: alpha\n")

	found := walkAll(t, base)

	want := map[string]bool{one: true, two: true}
	for _, d := range found {
		if d == store {
			t.Errorf("the registered store must not be listed: %q", d)
		}
		if !want[d] {
			t.Errorf("unexpected row %q", d)
		}
		delete(want, d)
	}
	for d := range want {
		t.Errorf("missing row %q", d)
	}
}

func TestWalkListsAStoreWhenNoRegistryConfirmsIt(t *testing.T) {
	// A directory carrying an identity file the registry does not point at is
	// an ordinary project, not a store, so nothing excludes it.
	isolateStores(t)
	base := t.TempDir()
	leftover := mkRoot(t, filepath.Join(base, "leftover"), "schema: spec-driven\n")
	mkStoreMetadata(t, leftover, "alpha")
	writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: /somewhere/else\n")

	found := walkAll(t, base)
	if len(found) != 1 || found[0] != leftover {
		t.Errorf("got %v, want the leftover listed as a project", found)
	}
}

func TestWalkExcludesNothingWithNoRegistry(t *testing.T) {
	isolateStores(t)
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")

	found := walkAll(t, base)
	if len(found) != 1 || found[0] != store {
		t.Errorf("got %v, want discovery unaffected when no stores are registered", found)
	}
}

func TestValidityUsesThePassedRegistryNotTheFile(t *testing.T) {
	// Every candidate directory asks the same registry question, so the answer
	// is fetched once for the walk and handed down. Deleting the file and
	// passing the map proves the check does not go back to disk per candidate:
	// if it did, the store would stop being recognised here.
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)

	stores, _, err := LoadRegistry(RegistryPath())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(RegistryPath()); err != nil {
		t.Fatal(err)
	}

	openspecDir := filepath.Join(store, "openspec")
	if isValidOpenSpecDir(openspecDir, stores) {
		t.Error("the store must be excluded from the registry handed in, with no file on disk")
	}
	if !isValidOpenSpecDir(openspecDir, nil) {
		t.Error("with no registry at all nothing is excluded")
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
	registerStore(t, "alpha", store)
	mkSpec(t, store, "shared-capability")
	repo := mkPointer(t, filepath.Join(base, "repo"), "store: alpha\n")

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	paths := walkAll(t, base)
	if len(paths) != 1 || paths[0] != repo {
		t.Fatalf("got %v, want the repo alone", paths)
	}

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

	// A cached path is a repo, and parsing it still follows the declaration,
	// so the row carries the store's statistics rather than the nothing the
	// repo holds itself.
	projects := ScanPaths(loaded)
	st, found := projects[repo]
	if !found {
		t.Fatalf("the repo is missing from %v", loaded)
	}
	if st.Info.StoreID != "alpha" {
		t.Errorf("store: got %q, want alpha", st.Info.StoreID)
	}
	if st.Info.SpecCount != 1 {
		t.Errorf("specs: got %d, want the store's 1", st.Info.SpecCount)
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

	key, st, err := ScanResolved(repo)
	if err != nil {
		t.Fatal(err)
	}
	// Filed under where the reading started, so that two repos sharing a store
	// are two projects. Where it ended travels inside the info as Root.
	if key != repo {
		t.Fatalf("key: got %q, want the repo %q", key, repo)
	}
	if st.Info.Root != store {
		t.Fatalf("root: got %q, want the store %q", st.Info.Root, store)
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

	key, st, err := ScanResolved(repo)
	if err != nil {
		t.Fatalf("an unfollowed declaration is a report, not an error: %v", err)
	}
	if key != repo {
		t.Errorf("key: got %q, want the repo %q", key, repo)
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

	key, st, err := ScanResolved(dir)
	if err != nil {
		t.Fatal(err)
	}
	if key != dir {
		t.Errorf("got %q, want %q", key, dir)
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
	key, _, err := ScanResolved(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if key != "" {
		t.Errorf("got %q, want nothing outside a project", key)
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
	_, st, _ := ScanResolved(repo)
	if st.Info.Root != alpha || st.Info.SpecCount != 1 {
		t.Fatalf("got root=%q specs=%d, want alpha with 1", st.Info.Root, st.Info.SpecCount)
	}

	writeConfig(t, repo, "config.yaml", "store: beta\n")
	_, st, _ = ScanResolved(repo)
	if st.Info.Root != beta {
		t.Errorf("root: got %q, want beta %q after repointing", st.Info.Root, beta)
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

// --- what a discovered row carries ---

func TestScanGivesAStoreBackedRepoTheStoresStatistics(t *testing.T) {
	// A row reporting zero for a repo whose specs are one directory away would
	// be worse than no row at all.
	isolateStores(t)
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	mkSpec(t, store, "one")
	mkSpec(t, store, "two")
	mkChange(t, store, "a-change", "- [x] 1.1 done\n- [ ] 1.2 todo\n")
	repo := mkPointer(t, filepath.Join(base, "repo"), "store: alpha\n")

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	projects, err := Scan(cfg, true)
	if err != nil {
		t.Fatal(err)
	}

	st, ok := projects[repo]
	if !ok {
		t.Fatalf("the repo is missing from %v", keysOf(projects))
	}
	if st.Info.SpecCount != 2 {
		t.Errorf("specs: got %d, want the store's 2", st.Info.SpecCount)
	}
	if len(st.Info.ActiveChanges) != 1 {
		t.Errorf("changes: got %d, want the store's 1", len(st.Info.ActiveChanges))
	}
	if st.Info.TasksTotal != 2 || st.Info.TasksDone != 1 {
		t.Errorf("tasks: got %d/%d, want 1/2", st.Info.TasksDone, st.Info.TasksTotal)
	}
	if st.Info.StoreID != "alpha" {
		t.Errorf("store: got %q, want alpha", st.Info.StoreID)
	}
}

func TestScanListsTheStoresFilesForAStoreBackedRepo(t *testing.T) {
	// The `:` query searches a project's file paths and contents. Listing the
	// repo's own openspec/ would make a search inside the project miss the
	// specs the project shows.
	isolateStores(t)
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	mkSpec(t, store, "findable-capability")
	repo := mkPointer(t, filepath.Join(base, "repo"), "store: alpha\n")

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	projects, _ := Scan(cfg, true)

	var joined string
	for _, f := range projects[repo].Files {
		joined += f.Path + "\n"
	}
	if !strings.Contains(joined, "findable-capability") {
		t.Errorf("the listing must be the store's tree, got:\n%s", joined)
	}
}

func TestScanReadsASharedStoreOncePerScan(t *testing.T) {
	// Two repos on one store is the ordinary case now that repos are what get
	// listed, so the scan must not cost twice as much for it. Removing the
	// store's content between the two reads shows whether the second was
	// served from the memo or went back to disk.
	isolateStores(t)
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	mkSpec(t, store, "only-spec")
	one := mkPointer(t, filepath.Join(base, "repo-one"), "store: alpha\n")
	two := mkPointer(t, filepath.Join(base, "repo-two"), "store: alpha\n")

	cache := newRootCache()
	_, first, err := scanResolved(one, cache, false)
	if err != nil {
		t.Fatal(err)
	}
	if first.Info.SpecCount != 1 {
		t.Fatalf("specs: got %d, want 1", first.Info.SpecCount)
	}

	if err := os.RemoveAll(filepath.Join(store, "openspec", "specs")); err != nil {
		t.Fatal(err)
	}

	_, second, err := scanResolved(two, cache, false)
	if err != nil {
		t.Fatal(err)
	}
	if second.Info.SpecCount != 1 {
		t.Errorf("specs: got %d, want the memoised 1: the store was read twice", second.Info.SpecCount)
	}
	if second.Info.Origin != two {
		t.Errorf("origin: got %q, want %q: the memo must not carry the first repo's identity", second.Info.Origin, two)
	}
}

func TestTwoReposSharingAStoreKeepTheirOwnIdentity(t *testing.T) {
	isolateStores(t)
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	one := mkPointer(t, filepath.Join(base, "repo-one"), "store: alpha\ncontext: one\n")
	two := mkPointer(t, filepath.Join(base, "repo-two"), "store: alpha\ncontext: two\n")

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	projects, _ := Scan(cfg, true)

	if len(projects) != 2 {
		t.Fatalf("got %v, want a row per repo", keysOf(projects))
	}
	for path, want := range map[string]string{one: "one", two: "two"} {
		st, ok := projects[path]
		if !ok {
			t.Fatalf("%q missing", path)
		}
		if st.Info.Origin != path {
			t.Errorf("origin: got %q, want %q", st.Info.Origin, path)
		}
		if st.Info.Root != store {
			t.Errorf("root: got %q, want the shared store", st.Info.Root)
		}
		if !strings.Contains(st.Info.OriginConfigContent, want) {
			t.Errorf("%q carries the wrong repo configuration: %q", path, st.Info.OriginConfigContent)
		}
		if st.Info.Store == nil || st.Info.Store.Origin != path {
			t.Errorf("%q shares a store record with its sibling", path)
		}
	}
}

func TestScanKeepsARepoWhoseDeclarationCannotBeFollowed(t *testing.T) {
	isolateStores(t)
	base := t.TempDir()
	repo := mkPointer(t, filepath.Join(base, "orphan"), "store: nowhere\n")

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	projects, _ := Scan(cfg, true)

	st, ok := projects[repo]
	if !ok {
		t.Fatalf("a repo with an unfollowable declaration is still a project: %v", keysOf(projects))
	}
	if st.Info.StoreProblem == nil {
		t.Error("it must carry the problem rather than look empty")
	}
}

func TestDiscoveryLeavesTheGitStateAlone(t *testing.T) {
	// Several subprocesses per discovered store is what the picker cannot
	// afford. The git state belongs to opening a project, not to finding one.
	isolateStores(t)
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	repo := mkPointer(t, filepath.Join(base, "repo"), "store: alpha\n")

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{base}
	projects, _ := Scan(cfg, true)
	if s := projects[repo].Info.Store; s == nil || s.Git != nil {
		t.Errorf("discovery must not read git state: %+v", s)
	}

	_, opened, err := ScanResolved(repo)
	if err != nil {
		t.Fatal(err)
	}
	if s := opened.Info.Store; s == nil || s.Git == nil {
		t.Errorf("opening a project does read it: %+v", s)
	}
}

func keysOf(m ProjectMap) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func TestWalkExcludesNothingWhenTheRegistryCannotBeRead(t *testing.T) {
	// Failing to list projects because the registry is broken would be worse
	// than listing a store alongside them, so a bad registry excludes nothing.
	isolateStores(t)
	writeRegistry(t, "stores: [unclosed\n")
	base := t.TempDir()
	store := mkRoot(t, filepath.Join(base, "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")

	found := walkAll(t, base)
	if len(found) != 1 || found[0] != store {
		t.Errorf("got %v, want discovery to carry on", found)
	}
}

func TestScanPathsSkipsAPathThatResolvesToNothing(t *testing.T) {
	isolateStores(t)
	good := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	gone := filepath.Join(t.TempDir(), "never-existed")

	projects := ScanPaths([]string{good, gone})
	if len(projects) != 1 {
		t.Fatalf("got %v, want the readable path alone", keysOf(projects))
	}
	if _, ok := projects[good]; !ok {
		t.Errorf("got %v, want %q", keysOf(projects), good)
	}
}

func TestScanResolvedReportsAnUnreadableRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads anything")
	}
	isolateStores(t)
	// The specs directory, not openspec/ itself: locking the latter would make
	// the root stop qualifying and the walk would simply carry on upward,
	// which is a different path from the one under test.
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	specsDir := filepath.Join(dir, "openspec", "specs")
	if err := os.Chmod(specsDir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(specsDir, 0o755) })

	if _, _, err := ScanResolved(dir); err == nil {
		t.Error("a root that cannot be listed is an error, not an empty project")
	}
}

// --- settings a released version accepted and this one does not ---

func TestRetiredKeysIn(t *testing.T) {
	write := func(body string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "config.yml")
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}

	t.Run("a configuration still carrying one", func(t *testing.T) {
		got := RetiredKeysIn(write("scandirs:\n  include: [x]\nchange_mode: active\n"))
		if len(got) != 1 || got[0] != "change_mode" {
			t.Errorf("got %v, want change_mode", got)
		}
		if RetiredConfigKeys["change_mode"] == "" {
			t.Error("a retired key must carry a reason to report")
		}
	})

	t.Run("a configuration without it", func(t *testing.T) {
		if got := RetiredKeysIn(write("scandirs:\n  include: [x]\n")); len(got) != 0 {
			t.Errorf("got %v, want nothing", got)
		}
	})

	t.Run("a file that is not there", func(t *testing.T) {
		if got := RetiredKeysIn(filepath.Join(t.TempDir(), "gone.yml")); got != nil {
			t.Errorf("got %v, want nothing", got)
		}
	})

	t.Run("a file that is not YAML", func(t *testing.T) {
		if got := RetiredKeysIn(write("scandirs: [unclosed\n")); got != nil {
			t.Errorf("got %v, want nothing", got)
		}
	})

	t.Run("the key is ignored by the parser either way", func(t *testing.T) {
		// Which is the whole reason it has to be reported: YAML discards keys
		// a program does not know, with no error.
		cfg, err := ParseConfigFile(write("scandirs:\n  include: [x]\nchange_mode: active\n"), "")
		if err != nil {
			t.Fatalf("a retired key must not break parsing: %v", err)
		}
		if len(cfg.ScanDirs.Include) != 1 {
			t.Errorf("the rest of the configuration must still load: %+v", cfg.ScanDirs)
		}
	})
}

// --- a glob include contributes only what it expands to ---

// walkStrict walks base without ignoring directory errors, which is what the
// glob pattern used to fail.
func walkStrict(t *testing.T, includes []string) ([]string, error) {
	t.Helper()
	cfg := &Config{}
	cfg.ScanDirs.Include = includes
	results := make(chan string, 64)
	errCh := make(chan error, 1)
	go func() { errCh <- Walk(context.Background(), cfg, results, false) }()
	var found []string
	for d := range results {
		found = append(found, d)
	}
	sort.Strings(found)
	return found, <-errCh
}

func TestGlobIncludeIsNotWalkedAsAPath(t *testing.T) {
	// The pattern is not a directory. Walking it failed on every scan, and the
	// failure was swallowed by a log panel nobody opened.
	isolateStores(t)
	base := t.TempDir()
	one := mkRoot(t, filepath.Join(base, "gh.one"), "schema: spec-driven\n")
	two := mkRoot(t, filepath.Join(base, "gh.two"), "schema: spec-driven\n")
	mkRoot(t, filepath.Join(base, "other"), "schema: spec-driven\n")

	found, err := walkStrict(t, []string{filepath.Join(base, "gh.*")})
	if err != nil {
		t.Fatalf("a glob include must not fail the walk: %v", err)
	}
	want := []string{one, two}
	if len(found) != 2 || found[0] != want[0] || found[1] != want[1] {
		t.Errorf("got %v, want %v: the matches are walked and the pattern is not", found, want)
	}
}

func TestGlobIncludeThatMatchesNothing(t *testing.T) {
	isolateStores(t)
	base := t.TempDir()
	mkRoot(t, filepath.Join(base, "other"), "schema: spec-driven\n")

	found, err := walkStrict(t, []string{filepath.Join(base, "nomatch.*")})
	if err != nil {
		t.Fatalf("a glob that matches nothing must not fail the walk: %v", err)
	}
	if len(found) != 0 {
		t.Errorf("got %v, want nothing walked for that include", found)
	}
}

func TestPlainIncludeIsStillWalkedAsThePathItIs(t *testing.T) {
	isolateStores(t)
	base := t.TempDir()
	only := mkRoot(t, filepath.Join(base, "plain"), "schema: spec-driven\n")

	found, err := walkStrict(t, []string{only})
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 || found[0] != only {
		t.Errorf("got %v, want %q", found, only)
	}
}

func TestAGlobCompletesAScanWithDirectoryErrorsNotIgnored(t *testing.T) {
	// This is the case that failed outright: the pattern was walked, could not
	// be stated, and the error was returned rather than swallowed.
	isolateStores(t)
	base := t.TempDir()
	mkRoot(t, filepath.Join(base, "gh.one"), "schema: spec-driven\n")

	if _, err := walkStrict(t, []string{filepath.Join(base, "gh.*")}); err != nil {
		t.Errorf("got %v, want a completed scan", err)
	}
}

func TestAnUnreadableGlobParentStillBehaves(t *testing.T) {
	// The existing rule is about the parent, not the pattern, so it survives.
	isolateStores(t)
	missing := filepath.Join(t.TempDir(), "gone", "gh.*")

	if _, err := walkStrict(t, []string{missing}); err == nil {
		t.Error("an unreadable glob parent must still be returned when errors are not ignored")
	}

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{missing}
	results := make(chan string, 8)
	go func() {
		if err := Walk(context.Background(), cfg, results, true); err != nil {
			t.Error(err)
		}
	}()
	for range results {
	}
}

func TestAPlainIncludeThatDoesNotExist(t *testing.T) {
	// Until the glob fix, the pattern itself was the thing that always failed
	// here, which is what exercised this path. A plain include that is not
	// there is now the only way to reach it, and it is the case the flag exists
	// for.
	isolateStores(t)
	base := t.TempDir()
	real := mkRoot(t, filepath.Join(base, "real"), "schema: spec-driven\n")
	missing := filepath.Join(base, "not-here")

	if _, err := walkStrict(t, []string{real, missing}); err == nil {
		t.Error("with directory errors not ignored, an unwalkable include is returned")
	}

	cfg := &Config{}
	cfg.ScanDirs.Include = []string{real, missing}
	results := make(chan string, 8)
	errCh := make(chan error, 1)
	go func() { errCh <- Walk(context.Background(), cfg, results, true) }()
	var found []string
	for d := range results {
		found = append(found, d)
	}
	if err := <-errCh; err != nil {
		t.Errorf("with them ignored the scan carries on: %v", err)
	}
	if len(found) != 1 || found[0] != real {
		t.Errorf("got %v, want the readable include alone", found)
	}
}
