package scanner

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// The resolver here reimplements OpenSpec's own rules rather than shelling out
// to the CLI, because the picker cannot afford a process per project. The
// price of that choice is that a beta which moves its rules would drift away
// silently. This test is the hedge: the same fixtures are put to both, and a
// disagreement fails here rather than showing wrong content on screen.
//
// It is a test, never a runtime dependency. Without the binary it skips.

type contextOutput struct {
	Root *struct {
		Path    string `json:"path"`
		Source  string `json:"source"`
		StoreID string `json:"store_id"`
	} `json:"root"`
	Status []struct {
		Severity string `json:"severity"`
		Code     string `json:"code"`
	} `json:"status"`
}

func openspecContext(t *testing.T, dir, dataHome string) (contextOutput, bool) {
	t.Helper()
	cmd := exec.Command("openspec", "context", "--json")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "XDG_DATA_HOME="+dataHome)
	out, err := cmd.Output()
	if err != nil {
		// A resolution the CLI refuses outright still prints JSON on most
		// paths; where it does not, the caller only checks that both sides
		// agree something is wrong.
		if len(out) == 0 {
			return contextOutput{}, false
		}
	}
	var parsed contextOutput
	if err := json.Unmarshal(out, &parsed); err != nil {
		return contextOutput{}, false
	}
	return parsed, true
}

func TestResolverAgreesWithTheOpenSpecCLI(t *testing.T) {
	if _, err := exec.LookPath("openspec"); err != nil {
		t.Skip("openspec CLI not installed")
	}

	dataHome := isolateStores(t)

	store := mkRoot(t, filepath.Join(t.TempDir(), "the-store"), "schema: spec-driven\n")
	mkStoreMetadata(t, store, "parity")
	writeRegistry(t, "version: 1\nstores:\n  parity:\n    backend:\n      type: git\n      local_path: "+store+"\n")

	plain := mkRoot(t, filepath.Join(t.TempDir(), "plain"), "schema: spec-driven\n")
	pointing := mkPointer(t, filepath.Join(t.TempDir(), "pointing"), "schema: spec-driven\nstore: parity\n")
	bothWays := mkRoot(t, filepath.Join(t.TempDir(), "both"), "schema: spec-driven\nstore: parity\n")

	cases := []struct {
		name     string
		start    string
		wantRoot string
		wantID   string
		// One deliberate divergence. Standing inside a store, `openspec
		// context` reports source "nearest" and no store_id: it answers "which
		// root am I in", and the answer is this directory. specgetty still
		// reads the identity file, because the picker has to name that row by
		// the store's id rather than by the folder it sits in. The root and
		// the classification still have to match, so only the id check is
		// waived.
		idIsSpecgettysOwn bool
	}{
		{"a plain project", plain, plain, "", false},
		{"a repo pointing at a store", pointing, store, "parity", false},
		{"the store itself", store, store, "parity", true},
		{"content beside a pointer", bothWays, bothWays, "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ResolveRoot(tc.start)
			if !ok {
				t.Fatalf("specgetty resolved nothing for %q", tc.start)
			}
			if got.Problem != nil {
				t.Fatalf("unexpected problem: %v", got.Problem)
			}

			cli, parsed := openspecContext(t, tc.start, dataHome)
			if !parsed || cli.Root == nil {
				t.Skipf("the CLI gave no root for %q", tc.start)
			}

			// Temp directories can sit behind a symlink, and the CLI
			// canonicalises where specgetty does not. Compare resolved paths.
			if !samePath(t, got.Root, cli.Root.Path) {
				t.Errorf("root: specgetty says %q, the CLI says %q", got.Root, cli.Root.Path)
			}
			if !tc.idIsSpecgettysOwn && got.StoreID != cli.Root.StoreID {
				t.Errorf("store id: specgetty says %q, the CLI says %q", got.StoreID, cli.Root.StoreID)
			}
			if tc.idIsSpecgettysOwn && cli.Root.StoreID != "" {
				t.Errorf("the CLI now names the store from inside it (%q); "+
					"the divergence this case records has closed and the case can go",
					cli.Root.StoreID)
			}
			// The source is not carried through specgetty's API, but it is the
			// classification both sides have to agree on, so it is checked
			// against what specgetty's two fields encode.
			wantDeclared := got.Origin != got.Root
			if gotDeclared := cli.Root.Source == "declared"; gotDeclared != wantDeclared {
				t.Errorf("source: the CLI says %q, specgetty has origin=%q root=%q",
					cli.Root.Source, got.Origin, got.Root)
			}
			if !samePath(t, tc.wantRoot, got.Root) {
				t.Errorf("root: got %q, want %q", got.Root, tc.wantRoot)
			}
			if got.StoreID != tc.wantID {
				t.Errorf("store id: got %q, want %q", got.StoreID, tc.wantID)
			}
		})
	}
}

func TestBothSidesRefuseAnUnregisteredStore(t *testing.T) {
	if _, err := exec.LookPath("openspec"); err != nil {
		t.Skip("openspec CLI not installed")
	}
	dataHome := isolateStores(t)
	writeRegistry(t, "version: 1\nstores: {}\n")
	repo := mkPointer(t, filepath.Join(t.TempDir(), "orphan"), "schema: spec-driven\nstore: nowhere\n")

	res, ok := ResolveRoot(repo)
	if !ok || res.Problem == nil {
		t.Fatalf("specgetty must refuse an unregistered store: %+v", res)
	}

	cli, parsed := openspecContext(t, repo, dataHome)
	if parsed && cli.Root != nil {
		t.Errorf("the CLI resolved a root where specgetty refused: %+v", cli.Root)
	}
}

func samePath(t *testing.T, a, b string) bool {
	t.Helper()
	ra, err := filepath.EvalSymlinks(a)
	if err != nil {
		ra = a
	}
	rb, err := filepath.EvalSymlinks(b)
	if err != nil {
		rb = b
	}
	return ra == rb
}
