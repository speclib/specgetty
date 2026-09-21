package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// isolateStores points the registry at a temporary directory, so no test ever
// reads the machine's real store registry.
func isolateStores(t *testing.T) string {
	t.Helper()
	data := t.TempDir()
	t.Setenv("XDG_DATA_HOME", data)
	return data
}

// writeRegistry writes a registry mapping store ids to local paths.
func writeRegistry(t *testing.T, body string) {
	t.Helper()
	path := RegistryPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// mkRoot creates an OpenSpec root with content of its own.
func mkRoot(t *testing.T, dir, config string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "openspec", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "openspec", "changes", "archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, dir, "config.yaml", config)
	return dir
}

// mkPointer creates a repo that keeps no content and only declares a store.
func mkPointer(t *testing.T, dir, body string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, dir, "config.yaml", body)
	return dir
}

func writeConfig(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "openspec", name), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// registerStore records a store in the registry, which is what makes its
// identity file mean anything.
func registerStore(t *testing.T, id, root string) {
	t.Helper()
	writeRegistry(t, "version: 1\nstores:\n  "+id+":\n    backend:\n      type: git\n      local_path: "+root+"\n")
}

// mkStoreMetadata makes a root a store by giving it an identity file.
func mkStoreMetadata(t *testing.T, root, id string) {
	t.Helper()
	dir := filepath.Join(root, storeMetadataDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "version: 1\nid: " + id + "\n"
	if err := os.WriteFile(filepath.Join(dir, storeMetadataFile), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// --- 1.1 the registry is read from the OpenSpec data directory ---

func TestDataDirHonoursXDGDataHome(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/somewhere/share")
	want := filepath.Join("/somewhere/share", "openspec")
	if got := DataDir(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	wantReg := filepath.Join(want, "stores", "registry.yaml")
	if got := RegistryPath(); got != wantReg {
		t.Errorf("registry: got %q, want %q", got, wantReg)
	}
}

func TestDataDirFallsBackToLocalShare(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home directory")
	}
	want := filepath.Join(home, ".local", "share", "openspec")
	if got := DataDir(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLoadRegistryReportsAbsenceApartFromEmptiness(t *testing.T) {
	isolateStores(t)

	stores, present, err := LoadRegistry(RegistryPath())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if present {
		t.Error("a missing registry must not report as present")
	}
	if len(stores) != 0 {
		t.Errorf("got %d stores, want none", len(stores))
	}

	writeRegistry(t, "version: 1\nstores: {}\n")
	_, present, err = LoadRegistry(RegistryPath())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !present {
		t.Error("an empty registry is still present")
	}
}

func TestLoadRegistryReadsTheBackend(t *testing.T) {
	isolateStores(t)
	writeRegistry(t, `version: 1
stores:
  alpha:
    backend:
      type: git
      local_path: /srv/alpha
      remote: git@example.com:team/alpha.git
      branch: main
`)
	stores, present, err := LoadRegistry(RegistryPath())
	if err != nil || !present {
		t.Fatalf("present=%v err=%v", present, err)
	}
	be, ok := stores["alpha"]
	if !ok {
		t.Fatal("alpha missing from the registry")
	}
	if be.LocalPath != "/srv/alpha" {
		t.Errorf("local_path: got %q", be.LocalPath)
	}
	if be.Remote != "git@example.com:team/alpha.git" {
		t.Errorf("remote: got %q", be.Remote)
	}
	if be.Branch != "main" {
		t.Errorf("branch: got %q", be.Branch)
	}
}

// --- 1.2 the store declaration ---

func TestReadStorePointer(t *testing.T) {
	cases := []struct {
		name  string
		body  string
		state PointerState
		id    string
	}{
		{"a declared id", "schema: spec-driven\nstore: alpha\n", PointerDeclared, "alpha"},
		{"no store key", "schema: spec-driven\n", PointerAbsent, ""},
		{"a list", "store:\n  - alpha\n  - beta\n", PointerNotAString, ""},
		{"a mapping", "store:\n  id: alpha\n", PointerNotAString, ""},
		{"an empty value", "store: \"\"\n", PointerNotAString, ""},
		{"broken yaml", "store: [unclosed\n", PointerUnreadable, ""},
		{"a document that is not a mapping", "- one\n- two\n", PointerAbsent, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeConfig(t, dir, "config.yaml", tc.body)
			p := ReadStorePointer(dir)
			if p.State != tc.state {
				t.Errorf("state: got %v, want %v", p.State, tc.state)
			}
			if p.ID != tc.id {
				t.Errorf("id: got %q, want %q", p.ID, tc.id)
			}
			if p.State != PointerAbsent && p.File == "" {
				t.Error("a problem must name the file it was read from")
			}
		})
	}
}

func TestReadStorePointerAcceptsTheYmlSpelling(t *testing.T) {
	// OpenSpec probes config.yaml then config.yml. The declaration lives in
	// that file, so knowing only the first spelling would lose the pointer.
	dir := t.TempDir()
	writeConfig(t, dir, "config.yml", "store: alpha\n")
	p := ReadStorePointer(dir)
	if p.State != PointerDeclared || p.ID != "alpha" {
		t.Fatalf("got state=%v id=%q, want a declared alpha", p.State, p.ID)
	}
}

func TestConfigFilePathPrefersYaml(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, "config.yml", "store: fromYml\n")
	writeConfig(t, dir, "config.yaml", "store: fromYaml\n")
	if got := ReadStorePointer(dir).ID; got != "fromYaml" {
		t.Errorf("got %q, want the .yaml spelling to win", got)
	}
}

func TestReadStorePointerWithNoConfigAtAll(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}
	if p := ReadStorePointer(dir); p.State != PointerAbsent || p.File != "" {
		t.Errorf("got state=%v file=%q, want an absent pointer with no file", p.State, p.File)
	}
}

// --- 1.3 store identity ---

func TestReadStoreMetadata(t *testing.T) {
	t.Run("a store", func(t *testing.T) {
		root := t.TempDir()
		mkStoreMetadata(t, root, "alpha")
		md, ok := ReadStoreMetadata(root)
		if !ok || md.ID != "alpha" {
			t.Errorf("got %+v ok=%v, want alpha", md, ok)
		}
	})

	t.Run("not a store", func(t *testing.T) {
		if _, ok := ReadStoreMetadata(t.TempDir()); ok {
			t.Error("a directory with no metadata file is not a store")
		}
	})

	t.Run("unreadable metadata", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, storeMetadataDir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, storeMetadataFile), []byte("id: [broken\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, ok := ReadStoreMetadata(root); ok {
			t.Error("corrupt metadata must not read as a store")
		}
	})

	t.Run("metadata with no id", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, storeMetadataDir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, storeMetadataFile), []byte("version: 1\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, ok := ReadStoreMetadata(root); ok {
			t.Error("metadata naming no id must not read as a store")
		}
	})
}

// --- 1.4 and 1.5 resolution ---

func TestResolveRootContentOfItsOwn(t *testing.T) {
	isolateStores(t)
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")

	res, ok := ResolveRoot(dir)
	if !ok {
		t.Fatal("expected a resolution")
	}
	if res.Root != dir || res.Origin != dir {
		t.Errorf("got root=%q origin=%q, want both %q", res.Root, res.Origin, dir)
	}
	if res.StoreID != "" || res.Problem != nil {
		t.Errorf("a plain project carries no store: %+v", res)
	}
}

func TestResolveRootFromBeneath(t *testing.T) {
	isolateStores(t)
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	nested := filepath.Join(dir, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	res, ok := ResolveRoot(nested)
	if !ok || res.Root != dir {
		t.Errorf("got root=%q ok=%v, want %q", res.Root, ok, dir)
	}
}

func TestResolveRootIgnoresAPointerBesideContent(t *testing.T) {
	// Content present locally outranks a pointer. OpenSpec warns about this
	// case rather than following it, and a project that read the store instead
	// would show specs the user cannot find on disk.
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: "+store+"\n")

	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\nstore: alpha\n")

	res, ok := ResolveRoot(dir)
	if !ok {
		t.Fatal("expected a resolution")
	}
	if res.Root != dir {
		t.Errorf("got root=%q, want the local root %q: a pointer beside content is dead", res.Root, dir)
	}
	if res.StoreID != "" {
		t.Errorf("got storeID %q, want none", res.StoreID)
	}
}

func TestResolveRootFollowsAPointer(t *testing.T) {
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	writeRegistry(t, `version: 1
stores:
  alpha:
    backend:
      type: git
      local_path: `+store+`
      remote: git@example.com:team/alpha.git
      branch: trunk
`)
	repo := mkPointer(t, t.TempDir(), "schema: spec-driven\nstore: alpha\n")

	res, ok := ResolveRoot(repo)
	if !ok {
		t.Fatal("expected a resolution")
	}
	if res.Root != store {
		t.Errorf("root: got %q, want the store %q", res.Root, store)
	}
	if res.Origin != repo {
		t.Errorf("origin: got %q, want the repo %q", res.Origin, repo)
	}
	if res.StoreID != "alpha" {
		t.Errorf("storeID: got %q, want alpha", res.StoreID)
	}
	if res.Problem != nil {
		t.Errorf("unexpected problem: %v", res.Problem)
	}
	if res.Store == nil {
		t.Fatal("expected store details")
	}
	if res.Store.Remote != "git@example.com:team/alpha.git" || res.Store.Branch != "trunk" {
		t.Errorf("registry details not carried: %+v", res.Store)
	}
}

func TestResolveRootTwoReposOneStore(t *testing.T) {
	// The mapping is many to one. Both resolve to the same content and each
	// keeps its own origin, which is what the header names.
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "shared")
	writeRegistry(t, "version: 1\nstores:\n  shared:\n    backend:\n      type: git\n      local_path: "+store+"\n")

	one := mkPointer(t, t.TempDir(), "store: shared\n")
	two := mkPointer(t, t.TempDir(), "store: shared\n")

	a, _ := ResolveRoot(one)
	b, _ := ResolveRoot(two)

	if a.Root != store || b.Root != store {
		t.Errorf("both must resolve to %q: got %q and %q", store, a.Root, b.Root)
	}
	if a.Origin != one || b.Origin != two {
		t.Errorf("each keeps its own origin: got %q and %q", a.Origin, b.Origin)
	}
}

func TestResolveRootFreshlyInitialisedProject(t *testing.T) {
	// A configuration, no content and no pointer: still a root of its own.
	isolateStores(t)
	dir := mkPointer(t, t.TempDir(), "schema: spec-driven\n")
	res, ok := ResolveRoot(dir)
	if !ok || res.Root != dir || res.Problem != nil {
		t.Errorf("got %+v ok=%v, want the directory itself", res, ok)
	}
}

func TestResolveRootSkipsAnEmptyShell(t *testing.T) {
	isolateStores(t)
	outer := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outer, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(outer, "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, ok := ResolveRoot(nested); ok {
		t.Error("an openspec/ directory with nothing in it is not a root")
	}
}

func TestResolveRootNamesAStoreOpenedDirectly(t *testing.T) {
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)

	res, ok := ResolveRoot(store)
	if !ok {
		t.Fatal("expected a resolution")
	}
	if res.StoreID != "alpha" {
		t.Errorf("got %q, want alpha: a store knows its own name", res.StoreID)
	}
	if res.Origin != res.Root {
		t.Error("a store opened directly has no separate origin")
	}
}

// --- 1.6 every unresolvable declaration is reported distinctly ---

func TestResolveRootProblems(t *testing.T) {
	t.Run("no registry on the machine", func(t *testing.T) {
		isolateStores(t)
		repo := mkPointer(t, t.TempDir(), "store: alpha\n")
		res, _ := ResolveRoot(repo)
		assertProblem(t, res, ProblemNoRegistry, "alpha")
	})

	t.Run("the declared store is not registered", func(t *testing.T) {
		isolateStores(t)
		writeRegistry(t, "version: 1\nstores:\n  beta:\n    backend:\n      type: git\n      local_path: /srv/beta\n")
		repo := mkPointer(t, t.TempDir(), "store: alpha\n")
		res, _ := ResolveRoot(repo)
		assertProblem(t, res, ProblemUnknownStore, "alpha")
	})

	t.Run("the registered path is gone", func(t *testing.T) {
		isolateStores(t)
		writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: /no/such/place\n")
		repo := mkPointer(t, t.TempDir(), "store: alpha\n")
		res, _ := ResolveRoot(repo)
		assertProblem(t, res, ProblemMissingPath, "alpha")
	})

	t.Run("the registered path holds no identity", func(t *testing.T) {
		isolateStores(t)
		store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
		writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: "+store+"\n")
		repo := mkPointer(t, t.TempDir(), "store: alpha\n")
		res, _ := ResolveRoot(repo)
		assertProblem(t, res, ProblemStoreIdentity, "alpha")
	})

	t.Run("the identity names a different store", func(t *testing.T) {
		isolateStores(t)
		store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
		mkStoreMetadata(t, store, "beta")
		writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: "+store+"\n")
		repo := mkPointer(t, t.TempDir(), "store: alpha\n")
		res, _ := ResolveRoot(repo)
		assertProblem(t, res, ProblemStoreIdentity, "alpha")
	})

	t.Run("the declaration is not a single id", func(t *testing.T) {
		isolateStores(t)
		repo := mkPointer(t, t.TempDir(), "store:\n  - alpha\n")
		res, _ := ResolveRoot(repo)
		assertProblem(t, res, ProblemMalformedPointer, "")
	})

	t.Run("the configuration is not yaml", func(t *testing.T) {
		isolateStores(t)
		repo := mkPointer(t, t.TempDir(), "store: [unclosed\n")
		res, _ := ResolveRoot(repo)
		assertProblem(t, res, ProblemMalformedPointer, "")
	})
}

func assertProblem(t *testing.T, res Resolution, code, id string) {
	t.Helper()
	if res.Problem == nil {
		t.Fatalf("expected a problem, got %+v", res)
	}
	if res.Problem.Code != code {
		t.Errorf("code: got %q, want %q", res.Problem.Code, code)
	}
	if res.Problem.ID != id {
		t.Errorf("id: got %q, want %q", res.Problem.ID, id)
	}
	if res.Problem.File == "" {
		t.Error("a problem must name the file that declared the store")
	}
	if res.Problem.Detail == "" {
		t.Error("a problem must say why")
	}
	if !res.FromStore() {
		t.Error("an unfollowed declaration is still a project reading from a store")
	}
}

// --- git state is local only ---

func TestReadStoreGitOnSomethingThatIsNotARepository(t *testing.T) {
	dir := t.TempDir()
	g := ReadStoreGit(dir)
	if g.IsRepo {
		t.Error("a plain directory is not a git working copy")
	}
	if g.TrackingKnown || g.DirtyKnown {
		t.Error("nothing is known about a directory that is not a repository")
	}
}

// --- a store is what the registry says it is ---

func TestStoreIdentityNeedsTheRegistryToAgree(t *testing.T) {
	// An identity file is kept by a clone, and by a directory the registry has
	// since moved on from. Trusting it alone let an unregistered leftover wear
	// the registered store's name: gh.nivis-project/ospecs claims id `nivis`
	// while the registered `nivis` is somewhere else entirely.
	isolateStores(t)
	registered := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	leftover := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, registered, "alpha")
	mkStoreMetadata(t, leftover, "alpha")
	registerStore(t, "alpha", registered)

	if info := storeInfoAt(registered); info == nil || info.ID != "alpha" {
		t.Errorf("the registered directory is the store: %+v", info)
	}
	if info := storeInfoAt(leftover); info != nil {
		t.Errorf("a directory the registry points away from is not a store: %+v", info)
	}
}

func TestStoreIdentityWithNoRegistryEntry(t *testing.T) {
	isolateStores(t)
	writeRegistry(t, "version: 1\nstores:\n  other:\n    backend:\n      type: git\n      local_path: /srv/other\n")
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, dir, "alpha")

	if info := storeInfoAt(dir); info != nil {
		t.Errorf("an id the registry never mentions is not a store: %+v", info)
	}
}

func TestStoreIdentityWithNoRegistryAtAll(t *testing.T) {
	isolateStores(t)
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, dir, "alpha")

	if info := storeInfoAt(dir); info != nil {
		t.Errorf("with no registry nothing is a store: %+v", info)
	}
}

func TestStoreIdentityMatchesThroughASymlink(t *testing.T) {
	// The registry records a canonicalised path while a walk reports whatever
	// it descended through, so a plain string comparison would fail to match a
	// store reached behind a link.
	isolateStores(t)
	real := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, real, "alpha")
	registerStore(t, "alpha", real)

	link := filepath.Join(t.TempDir(), "linked")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if info := storeInfoAt(link); info == nil {
		t.Error("the same directory reached by a link is the same store")
	}
}

func TestResolveRootStillFollowsADeclarationAfterTheIdentityChange(t *testing.T) {
	// The declaration path already made the registry and the metadata agree.
	// Tightening what counts as a store met head on must not disturb it.
	isolateStores(t)
	store := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "alpha")
	registerStore(t, "alpha", store)
	repo := mkPointer(t, t.TempDir(), "store: alpha\n")

	res, ok := ResolveRoot(repo)
	if !ok || res.Root != store || res.StoreID != "alpha" || res.Problem != nil {
		t.Errorf("got %+v, want the store followed", res)
	}
}

func TestSamePath(t *testing.T) {
	dir := t.TempDir()
	cases := []struct {
		name string
		a, b string
		want bool
	}{
		{"identical strings", "/x/y", "/x/y", true},
		{"different", "/x/y", "/x/z", false},
		{"one empty", "", "/x", false},
		{"both empty", "", "", true},
		{"unclean but equal", dir, dir + "/.", true},
		{"neither exists", "/no/such/a", "/no/such/b", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := samePath(tc.a, tc.b); got != tc.want {
				t.Errorf("samePath(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestStoreInfoAtWithAnUnreadableRegistry(t *testing.T) {
	isolateStores(t)
	writeRegistry(t, "stores: [unclosed\n")
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, dir, "alpha")

	if info := storeInfoAt(dir); info != nil {
		t.Errorf("a registry that cannot be read confirms nothing: %+v", info)
	}
}

func TestStoreInfoWithNilRegistry(t *testing.T) {
	dir := mkRoot(t, t.TempDir(), "schema: spec-driven\n")
	mkStoreMetadata(t, dir, "alpha")
	if info := StoreInfoWith(nil, dir); info != nil {
		t.Errorf("no registry confirms nothing: %+v", info)
	}
	if info := StoreInfoWith(map[string]StoreBackend{}, t.TempDir()); info != nil {
		t.Errorf("a directory with no metadata is not a store: %+v", info)
	}
}

func TestRegistryPathWithNoHomeAndNoXDG(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("HOME", "")
	if got := DataDir(); got != "" {
		t.Skipf("this platform still resolves a home directory: %q", got)
	}
	if got := RegistryPath(); got != "" {
		t.Errorf("got %q, want empty when there is nowhere to look", got)
	}
}
