package scanner

import (
	"context"
	"path/filepath"
	"sort"
	"testing"
)

// walkCollecting is collectWalk without the fatal on error: these tests are
// about what Walk does when the configuration points somewhere unusable, so
// the error is the thing under test rather than a reason to stop.
//
// It also proves the results channel gets closed on every path. The drain
// goroutine below only finishes when it does, so a Walk that returns early
// without closing would hang here rather than fail.
func walkCollecting(t *testing.T, config *Config, ignoreDirErrors bool) ([]string, error) {
	t.Helper()
	results := make(chan string, 64)
	var found []string
	done := make(chan struct{})
	go func() {
		for p := range results {
			found = append(found, p)
		}
		close(done)
	}()
	err := Walk(context.Background(), config, results, ignoreDirErrors)
	<-done
	sort.Strings(found)
	return found, err
}

func TestWalkCarriesOnPastAGlobIntoAMissingDirectory(t *testing.T) {
	// Expanding a glob used to call log.Fatal, which exits the process. Nothing
	// a test asserts afterwards would run, so this passing at all is part of
	// the point.
	root := t.TempDir()
	makeProjectAt(t, root, "alpha")
	missing := filepath.Join(t.TempDir(), "gone", "c*")

	got, err := walkCollecting(t, configFor([]string{root, missing}, nil, false), true)
	if err != nil {
		t.Fatalf("Walk: %v, want the unreadable include skipped", err)
	}
	if len(got) != 1 {
		t.Errorf("found %v, want the one project under the readable include", got)
	}
}

func TestWalkReturnsTheGlobErrorWhenDirectoryErrorsAreNotIgnored(t *testing.T) {
	root := t.TempDir()
	makeProjectAt(t, root, "alpha")
	missing := filepath.Join(t.TempDir(), "gone", "c*")

	_, err := walkCollecting(t, configFor([]string{root, missing}, nil, false), false)
	if err == nil {
		t.Error("want the error reported to the caller when errors are not ignored")
	}
}

func TestWalkIgnoresAnEmptyInclude(t *testing.T) {
	// A YAML dash with nothing after it parses to an empty string, and the
	// glob test used to slice the last byte off whatever it was given.
	root := t.TempDir()
	makeProjectAt(t, root, "alpha")

	got, err := walkCollecting(t, configFor([]string{"", root}, nil, false), true)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("found %v, want the one project under the non-empty include", got)
	}
}

func TestSkipIgnoresEmptyExcludeEntries(t *testing.T) {
	if skip("/a/b/c", []string{""}) {
		t.Error("an empty exclude entry should exclude nothing")
	}
	if !skip("/a/b/vendor", []string{"", "vendor"}) {
		t.Error("an empty entry should not stop the entries after it from matching")
	}
}
