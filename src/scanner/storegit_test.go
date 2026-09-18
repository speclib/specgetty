package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// git state is read from local files and refs only. These tests build real
// repositories in temp directories, including a bare one standing in for a
// remote, so the ahead and behind counts are exercised without any network.

func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
}

// newRepo makes a working copy with one commit and returns it.
func newRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	gitRun(t, dir, "init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README"), []byte("one\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-qm", "first")
	return dir
}

func TestReadStoreGitOnACleanRepository(t *testing.T) {
	requireGit(t)
	dir := newRepo(t)

	g := ReadStoreGit(dir)
	if !g.IsRepo {
		t.Fatal("want a git working copy")
	}
	if !g.DirtyKnown || g.Dirty {
		t.Errorf("a committed tree is clean: %+v", g)
	}
	if g.OriginURL != "" {
		t.Errorf("no remote was added, got %q", g.OriginURL)
	}
	if g.TrackingKnown {
		t.Error("there is no upstream to compare against")
	}
}

func TestReadStoreGitSeesUncommittedChanges(t *testing.T) {
	requireGit(t)
	dir := newRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "README"), []byte("two\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	g := ReadStoreGit(dir)
	if !g.DirtyKnown || !g.Dirty {
		t.Errorf("an edited tree is dirty: %+v", g)
	}
}

func TestReadStoreGitScopesDirtinessToTheStoreSubtree(t *testing.T) {
	// Several stores commonly share one working copy, so a sibling store's
	// edits are not this store's state.
	requireGit(t)
	repo := newRepo(t)
	mine := filepath.Join(repo, "mine")
	theirs := filepath.Join(repo, "theirs")
	for _, d := range []string{mine, theirs} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "f"), []byte("x\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, repo, "add", "-A")
	gitRun(t, repo, "commit", "-qm", "two stores")

	if err := os.WriteFile(filepath.Join(theirs, "f"), []byte("changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if g := ReadStoreGit(mine); g.Dirty {
		t.Error("a sibling's edit must not make this store dirty")
	}
	if g := ReadStoreGit(theirs); !g.Dirty {
		t.Error("the edited store is dirty")
	}
}

func TestReadStoreGitReportsRemoteAndTracking(t *testing.T) {
	requireGit(t)
	remote := t.TempDir()
	gitRun(t, remote, "init", "-q", "--bare", "-b", "main")

	dir := newRepo(t)
	gitRun(t, dir, "remote", "add", "origin", remote)
	gitRun(t, dir, "push", "-q", "-u", "origin", "main")

	g := ReadStoreGit(dir)
	if g.OriginURL != remote {
		t.Errorf("origin url: got %q, want %q", g.OriginURL, remote)
	}
	if !g.TrackingKnown {
		t.Fatal("an upstream exists, so the comparison is available")
	}
	if g.Ahead != 0 || g.Behind != 0 {
		t.Errorf("just pushed, got ahead=%d behind=%d", g.Ahead, g.Behind)
	}

	if err := os.WriteFile(filepath.Join(dir, "README"), []byte("two\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-qm", "second")

	g = ReadStoreGit(dir)
	if g.Ahead != 1 || g.Behind != 0 {
		t.Errorf("one unpushed commit: got ahead=%d behind=%d", g.Ahead, g.Behind)
	}
	// The count came from local refs. Nothing was fetched, and the bare
	// repository still holds only the first commit.
	if head := gitRun(t, remote, "rev-list", "--count", "main"); head != "1" {
		t.Errorf("the remote must be untouched, it has %s commits", head)
	}
}

func TestStoreProblemErrorIsItsDetail(t *testing.T) {
	p := &StoreProblem{Code: ProblemUnknownStore, ID: "alpha", Detail: "it is not registered"}
	if p.Error() != "it is not registered" {
		t.Errorf("got %q", p.Error())
	}
	var nilProblem *StoreProblem
	if nilProblem.Error() != "" {
		t.Error("a nil problem says nothing")
	}
}

func TestLoadRegistryReportsUnreadableYAML(t *testing.T) {
	isolateStores(t)
	writeRegistry(t, "stores: [unclosed\n")
	_, present, err := LoadRegistry(RegistryPath())
	if err == nil {
		t.Error("want an error for a registry that is not YAML")
	}
	if !present {
		t.Error("the file exists, so it is present even when unreadable")
	}
}

func TestLoadRegistryWithNoPath(t *testing.T) {
	stores, present, err := LoadRegistry("")
	if err != nil || present || stores != nil {
		t.Errorf("got %v %v %v, want an empty answer", stores, present, err)
	}
}

func TestResolveRootReportsAnUnreadableRegistry(t *testing.T) {
	isolateStores(t)
	writeRegistry(t, "stores: [unclosed\n")
	repo := mkPointer(t, t.TempDir(), "store: alpha\n")

	res, _ := ResolveRoot(repo)
	assertProblem(t, res, ProblemNoRegistry, "alpha")
}

func TestResolveRootReportsAnEmptyLocalPath(t *testing.T) {
	isolateStores(t)
	writeRegistry(t, "version: 1\nstores:\n  alpha:\n    backend:\n      type: git\n      local_path: \"\"\n")
	repo := mkPointer(t, t.TempDir(), "store: alpha\n")

	res, _ := ResolveRoot(repo)
	if res.Problem == nil {
		t.Fatal("a registry entry with no path cannot be followed")
	}
}

func TestGitAtOnANonexistentDirectory(t *testing.T) {
	requireGit(t)
	if _, err := gitAt(filepath.Join(t.TempDir(), "nope"), "status"); err == nil {
		t.Error("want an error for a directory that is not there")
	}
}
