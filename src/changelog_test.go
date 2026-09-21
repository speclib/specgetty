package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// scripts/changelog-entry.sh is what the release script and the release
// workflow both read a changelog section with. It is shell rather than Go
// because both callers are shell, and it is tested here because this is the
// suite the gate runs.

func entry(t *testing.T, label, file string) (string, int) {
	t.Helper()
	cmd := exec.Command("bash",
		filepath.Join("..", "scripts", "changelog-entry.sh"),
		label, filepath.Join("testdata", file))
	out, err := cmd.Output()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return string(out), code
}

func TestAnEntryIsTheSectionUnderItsHeading(t *testing.T) {
	got, code := entry(t, "0.8.0", "changelog-sample.md")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	for _, want := range []string{"The middle entry", "A second category"} {
		if !strings.Contains(got, want) {
			t.Errorf("the entry should carry %q:\n%s", want, got)
		}
	}
}

// The heading above and the heading below are both excluded: an entry is what
// sits between them.
func TestAnEntryStopsAtTheNeighbouringHeadings(t *testing.T) {
	got, _ := entry(t, "0.8.0", "changelog-sample.md")
	for _, unwanted := range []string{
		"## [0.8.0]",      // its own heading
		"## [Unreleased]", // the one above
		"## [0.7.20]",     // the one below
		"Something not yet released",
		"A version whose number begins",
	} {
		if strings.Contains(got, unwanted) {
			t.Errorf("the entry should not carry %q:\n%s", unwanted, got)
		}
	}
}

// A version named in prose cannot start a section, and neither can a line that
// looks like a heading inside one: only a line beginning `## [` is a heading,
// and an entry containing one is already malformed.
func TestProseNamingAVersionDoesNotStartASection(t *testing.T) {
	got, _ := entry(t, "0.8.0", "changelog-sample.md")
	if !strings.Contains(got, "It mentions 0.7.1 in its prose") {
		t.Errorf("the sentence naming another version belongs to this entry:\n%s", got)
	}
}

// 0.7.2 must not match the heading of 0.7.20, which is what the closing bracket
// in the match is for.
func TestAVersionDoesNotMatchALongerOne(t *testing.T) {
	got, code := entry(t, "0.7.2", "changelog-sample.md")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(got, "The last entry in the file") {
		t.Errorf("0.7.2 should be its own entry:\n%s", got)
	}
	if strings.Contains(got, "begins with another version") {
		t.Errorf("0.7.2 matched the heading of 0.7.20:\n%s", got)
	}
}

func TestTheLastEntryRunsToTheEndOfTheFile(t *testing.T) {
	got, _ := entry(t, "0.7.2", "changelog-sample.md")
	if !strings.Contains(got, "The last entry in the file") {
		t.Errorf("an entry with no heading below it still reads:\n%s", got)
	}
}

func TestUnreleasedIsReadTheSameWay(t *testing.T) {
	got, code := entry(t, "Unreleased", "changelog-sample.md")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(got, "Something not yet released") {
		t.Errorf("Unreleased is a section like any other:\n%s", got)
	}
	if strings.Contains(got, "The middle entry") {
		t.Errorf("it stops at the version below it:\n%s", got)
	}
}

// What the release script's refusal rests on: an Unreleased section holding
// nothing is how "the changelog was not written" is observable.
func TestAnEmptySectionReadsAsEmpty(t *testing.T) {
	got, code := entry(t, "Unreleased", "changelog-empty.md")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if strings.TrimSpace(got) != "" {
		t.Errorf("got %q, want nothing", got)
	}
}

func TestAMissingHeadingIsAnError(t *testing.T) {
	_, code := entry(t, "9.9.9", "changelog-sample.md")
	if code == 0 {
		t.Error("a version with no section should fail rather than print nothing")
	}
}

// The real file, so the script is asserted against the shape it actually reads
// rather than only against a fixture written to suit it.
func TestTheProjectsOwnChangelogIsReadable(t *testing.T) {
	cmd := exec.Command("bash",
		filepath.Join("..", "scripts", "changelog-entry.sh"),
		"0.7.1", filepath.Join("..", "CHANGELOG.md"))
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("0.7.1 should be readable from the real changelog: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "$VISUAL") || !strings.Contains(got, "could not save") {
		t.Errorf("0.7.1's entry should carry what it shipped:\n%s", got)
	}
	if strings.Contains(got, "## [0.7.0]") {
		t.Error("it should stop at the version below it")
	}
}

// --- the release plumbing that reads the entry ---

// TestTheWorkflowPublishesTheChangelogEntry asserts the wiring rather than the
// text: that both goreleaser runs are handed the notes file, that the file is
// built from the tag being released, and that neither run is still assembling
// notes from commit subjects. Nothing else in this repository would notice if
// one of the four moved.
func TestTheWorkflowPublishesTheChangelogEntry(t *testing.T) {
	wf := readFile(t, filepath.Join("..", ".github", "workflows", "release.yml"))

	if n := strings.Count(wf, "--release-notes=release-notes.md"); n != 2 {
		t.Errorf("got %d goreleaser runs given the notes, want both", n)
	}
	if n := strings.Count(wf, "changelog-entry.sh"); n != 2 {
		t.Errorf("got %d jobs extracting the entry, want both", n)
	}
	if !strings.Contains(wf, `version="${GITHUB_REF_NAME#v}"`) {
		t.Error("the notes should be extracted for the tag being released")
	}
	// The release is created once. The second run adds its artifacts to it,
	// which is what stops two runs publishing two accounts.
	if !strings.Contains(wf, "needs: release-linux") {
		t.Error("the darwin job must wait on the linux job")
	}

	for _, f := range []string{".goreleaser-linux.yaml", ".goreleaser-darwin.yaml"} {
		cfg := readFile(t, filepath.Join("..", f))
		if !strings.Contains(cfg, "disable: true") {
			t.Errorf("%s still assembles notes from commits", f)
		}
		for _, gone := range []string{"sort: asc", "^docs:", "Merge pull request"} {
			if strings.Contains(cfg, gone) {
				t.Errorf("%s still carries the commit-derived changelog (%q)", f, gone)
			}
		}
	}
}

// TestTheReleaseScriptRefusesAnEmptyEntry asserts the check exists and runs
// before anything is edited. An empty entry caught after the tag is pushed is
// caught too late: the workflow has started and the release exists.
func TestTheReleaseScriptRefusesAnEmptyEntry(t *testing.T) {
	sh := readFile(t, filepath.Join("..", "scripts", "release.sh"))

	check := strings.Index(sh, "[Unreleased] section is empty")
	if check < 0 {
		t.Fatal("the release script should refuse an empty entry")
	}
	for _, later := range []string{
		"insert_version_heading \"$NEW_VERSION\"", // the first file edit
		"> \"$VERSION_FILE\"",                     // the other one
		"git tag -a",                              // and the two that publish
		"jj tag set",
	} {
		at := strings.Index(sh, later)
		if at >= 0 && at < check {
			t.Errorf("the check runs after %q; it must come first", later)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
