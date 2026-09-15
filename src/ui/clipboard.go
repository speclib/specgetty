package ui

import (
	"errors"
	"path/filepath"

	"github.com/atotto/clipboard"
)

// writeClipboard is the seam the tests replace.
//
// It exists so the suite never touches the developer's real clipboard: a test
// that clobbers what you had copied is a bad neighbour, and one that depends on
// wl-copy being installed fails on machines that are fine.
var writeClipboard = func(s string) error {
	// Unsupported is set at init time when no backend was found at all, which
	// is a different problem from a write that failed, and deserves a message
	// that says which.
	if clipboard.Unsupported {
		return errors.New("no clipboard tool found (install wl-copy, xclip or xsel)")
	}
	return clipboard.WriteAll(s)
}

// changeDirPath returns the absolute path of a change's directory.
//
// The path is built from DirName, never from Name. For an archived change the
// directory on disk keeps its `YYYY-MM-DD-` prefix while the displayed name has
// it stripped, so building from Name yields a path that does not exist. That is
// the defect doExportChange shipped with, and on the clipboard it would fail
// later still, in a shell, far from the tool that produced it.
func changeDirPath(projectPath string, r changeRow) string {
	if r.archived {
		return filepath.Join(projectPath, "openspec", "changes", "archive", r.ci.DirName)
	}
	return filepath.Join(projectPath, "openspec", "changes", r.ci.DirName)
}

// copyToClipboard puts text on the clipboard and returns the status message to
// show for it, so that a copy which did not happen never looks like one that
// did.
func copyToClipboard(what, text string) string {
	if err := writeClipboard(text); err != nil {
		return "could not copy: " + err.Error()
	}
	return "copied " + what + ": " + text
}
