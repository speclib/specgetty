package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The demo is recorded against the fixture harbour under demo/. A fixture that
// stopped parsing, or a delta that stopped exercising a mark, would produce a
// recording that quietly shows less than it claims to, and the only thing that
// would notice is a person looking at a GIF.

func demoPath(parts ...string) string {
	return filepath.Join(append([]string{"..", "..", "demo"}, parts...)...)
}

func readDemo(t *testing.T, parts ...string) string {
	t.Helper()
	b, err := os.ReadFile(demoPath(parts...))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// TestEveryFixtureSpecParses is what stops a recording opening on a report
// instead of an outline.
func TestEveryFixtureSpecParses(t *testing.T) {
	var seen int
	err := filepath.Walk(demoPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Base(path) != "spec.md" {
			return nil
		}
		if strings.Contains(filepath.ToSlash(path), "/changes/") {
			return nil // a delta, read by the other grammar below
		}
		b, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		seen++
		name := filepath.Base(filepath.Dir(path))
		if _, problems := parseSpec(name, string(b)); len(problems) > 0 {
			t.Errorf("%s does not parse: %v", path, problems)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if seen < 4 {
		t.Errorf("found %d fixture specs, want the harbour's four or more", seen)
	}
}

func TestEveryFixtureDeltaParses(t *testing.T) {
	var seen int
	err := filepath.Walk(demoPath("stores"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Base(path) != "spec.md" {
			return nil
		}
		if !strings.Contains(filepath.ToSlash(path), "/changes/") {
			return nil
		}
		b, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		seen++
		cap := filepath.Base(filepath.Dir(path))
		if _, problems := parseDelta(cap, string(b)); len(problems) > 0 {
			t.Errorf("%s does not parse as a delta: %v", path, problems)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if seen < 4 {
		t.Errorf("found %d fixture deltas, want the active change and three archived", seen)
	}
}

// TestTheActiveChangeExercisesEveryMark is the fixture's whole reason for
// existing. The recording of a change's deltas is worth making only if the
// delta carries all three operations and a modified requirement whose
// scenarios are unchanged, edited and added.
func TestTheActiveChangeExercisesEveryMark(t *testing.T) {
	const change = "light-the-dial-at-night"
	delta := readDemo(t, "stores", "tideclock", "openspec", "changes", change,
		"specs", "tide-display", "spec.md")
	live := readDemo(t, "stores", "tideclock", "openspec", "specs", "tide-display", "spec.md")

	tree, problems := buildChangeSpecTree(change, []string{"tide-display"},
		map[string]string{"tide-display": delta},
		map[string]string{"tide-display": live})
	if len(problems) > 0 {
		t.Fatalf("the fixture delta must build: %v", problems)
	}

	ops := map[string]bool{}
	marks := map[int]bool{}
	for _, n := range tree.nodes {
		switch n.kind {
		case nodeRequirement:
			ops[n.op] = true
		case nodeScenario:
			marks[n.mark] = true
		}
	}
	for _, op := range []string{opAdded, opModified, opRemoved} {
		if !ops[op] {
			t.Errorf("the delta carries no %s requirement, so the outline cannot show that mark", op)
		}
	}
	for name, mark := range map[string]int{
		"unchanged": markUnchanged, "edited": markEdited, "added": markAdded,
	} {
		if !marks[mark] {
			t.Errorf("no scenario is marked %s, so the recording cannot show it", name)
		}
	}
}

// TestTheActiveChangeStaysActive. Archived, its comparison disappears and the
// most interesting shot in the demo shows nothing.
func TestTheActiveChangeStaysActive(t *testing.T) {
	active := demoPath("stores", "tideclock", "openspec", "changes", "light-the-dial-at-night")
	if _, err := os.Stat(active); err != nil {
		t.Fatalf("the active change must stay outside archive/: %v", err)
	}
	archived := demoPath("stores", "tideclock", "openspec", "changes", "archive")
	entries, err := os.ReadDir(archived)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Errorf("got %d archived changes, want enough for both groups to be worth showing",
			len(entries))
	}
}

// TestTheRecordingCannotPublishTheRecordersMachine. The project header prints
// the project's absolute path, so recording in the checkout published the
// recorder's home directory in every frame. The fixture is staged somewhere
// neutral instead, which is also what makes two machines record the same GIF.
func TestTheRecordingCannotPublishTheRecordersMachine(t *testing.T) {
	script := readDemo(t, "..", "scripts", "record-demo.sh")

	if !strings.Contains(script, "SPECGETTY_DEMO_STAGE:-/tmp/specgetty-demo") {
		t.Error("the fixture should be staged at a fixed, neutral path")
	}
	if !strings.Contains(script, `export SPECGETTY_DEMO="$stage"`) {
		t.Error("the config should resolve against the staging directory, not the checkout")
	}
	if !strings.Contains(script, `cd "$stage" && vhs`) {
		t.Error("the tapes should run from the staging directory")
	}

	// The recordings are committed, so a GIF that exists only after `make demo`
	// would be a broken image in the README for everyone else.
	rec, err := os.ReadDir(demoPath("recordings"))
	if err != nil {
		t.Fatalf("the recordings must be tracked, not generated: %v", err)
	}
	var gifs int
	for _, e := range rec {
		if strings.HasSuffix(e.Name(), ".gif") {
			gifs++
		}
	}
	if gifs < 7 {
		t.Errorf("got %d recordings, want one per tape", gifs)
	}

	cfg := readDemo(t, "config.yml")
	if strings.Contains(cfg, "$HOME") || strings.Contains(cfg, "~") {
		t.Error("the demo config must not reach into a home directory")
	}
	if !strings.Contains(cfg, "$SPECGETTY_DEMO/projects") {
		t.Error("the scan should find the fixture projects and nothing else")
	}

	for _, name := range mustReadDir(t, demoPath("tapes")) {
		tape := readDemo(t, "tapes", name)
		if strings.Contains(tape, "/home/") {
			t.Errorf("%s names a home directory", name)
		}
		if !strings.Contains(tape, "$SPECGETTY_DEMO/config.yml") {
			t.Errorf("%s does not pass the demo config, so it would scan the machine", name)
		}
		// Every tape pins where it starts. A tape run from elsewhere reaches
		// the "No OpenSpec project here" prompt, which then swallows its keys.
		if !strings.Contains(tape, "cd $SPECGETTY_DEMO/projects/harbour-app") {
			t.Errorf("%s does not pin its starting directory", name)
		}
	}
}

func mustReadDir(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if len(names) < 7 {
		t.Errorf("got %d tapes, want one per feature plus the hero", len(names))
	}
	return names
}
