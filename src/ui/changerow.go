package ui

import (
	"fmt"
	"sort"

	"github.com/mipmip/specgetty/src/scanner"
)

// List modes for the merged change list. Active and archived changes live in
// one list; this selects which of them it contains.
const (
	modeOpen = iota
	modeArchived
	modeBoth
)

var listModeNames = []string{"open", "archived", "open+archived"}

// changeRow is one line of the change list. It wraps a ChangeInfo with the two
// things the list needs that the scanner does not record: whether the change
// came from the archive, and which of its files a body search matched.
type changeRow struct {
	ci           scanner.ChangeInfo
	archived     bool
	matchedFiles []string
}

// key identifies a row across a re-filter or a rescan. The archived flag is
// part of it because an archived change keeps the name it had while active, so
// a name on its own is not unique in modeBoth.
func (r changeRow) key() string {
	if r.archived {
		return "archived/" + r.ci.Name
	}
	return "open/" + r.ci.Name
}

// buildRows merges the two slices the scanner produces into one list, in the
// order the requested mode implies: active changes first, then archived.
func buildRows(info scanner.ProjectInfo, mode int) []changeRow {
	var rows []changeRow
	if mode == modeOpen || mode == modeBoth {
		for _, ci := range info.Changes {
			rows = append(rows, changeRow{ci: ci})
		}
	}
	if mode == modeArchived || mode == modeBoth {
		for _, ci := range info.ArchivedChanges {
			rows = append(rows, changeRow{ci: ci, archived: true})
		}
	}
	return rows
}

// emptyListMessage explains why the list is empty, which differs by mode. An
// empty list with no explanation reads as a broken scan.
func emptyListMessage(mode int) string {
	switch mode {
	case modeArchived:
		return "No archived changes"
	case modeBoth:
		return "No changes, active or archived"
	default:
		return "No active changes"
	}
}

// indexOfKey finds the row carrying the given key, or -1.
func indexOfKey(rows []changeRow, key string) int {
	for i, r := range rows {
		if r.key() == key {
			return i
		}
	}
	return -1
}

// artifactTabNames lists the sub-tabs for an open change: one per .md file,
// plus a specs tab when the change carries spec deltas.
func (r changeRow) artifactTabNames() []string {
	var names []string
	for _, f := range r.ci.ArtifactFiles {
		names = append(names, trimMarkdownSuffix(f))
	}
	if len(r.ci.SpecNames) > 0 {
		names = append(names, "specs")
	}
	return names
}

// searchableFiles returns every named body of text in the change, keyed by the
// label a search hint should show for it.
func (r changeRow) searchableFiles() []string {
	var names []string
	for _, f := range r.ci.ArtifactFiles {
		names = append(names, f)
	}
	for _, s := range r.ci.SpecNames {
		names = append(names, s)
	}
	sort.Strings(names)
	return names
}

// taskLabel renders task progress, or an empty string when the change has no
// tasks.md to count.
func (r changeRow) taskLabel() string {
	if r.ci.TasksTotal == 0 {
		return ""
	}
	return fmt.Sprintf("%d/%d", r.ci.TasksDone, r.ci.TasksTotal)
}
