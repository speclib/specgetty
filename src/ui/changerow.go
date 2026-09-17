package ui

import (
	"fmt"

	"github.com/mipmip/specgetty/src/scanner"
)

// List modes for the merged change list. Active and archived changes live in
// one list; this selects which of them it contains.
const (
	modeActive = iota
	modeArchived
	modeBoth
)

// The names of the three states, used on screen and accepted in configuration.
// One spelling everywhere: what the nav bar shows is what you write in the
// config file.
var listModeNames = []string{"active", "archived", "active+archived"}

// changeRow is one line of the change list. It wraps a ChangeInfo with the one
// thing the list needs that the scanner does not record: whether the change came
// from the archive.
type changeRow struct {
	ci       scanner.ChangeInfo
	archived bool
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
	if mode == modeActive || mode == modeBoth {
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

// searchName is the fuzzy and literal target for a change: its name.
func (r changeRow) searchName() string { return r.ci.Name }

// searchBodies exposes every artifact and spec in the change to a ':' query,
// labelled the way the match hint should show it.
func (r changeRow) searchBodies() map[string]string {
	bodies := make(map[string]string, len(r.ci.ArtifactContents)+len(r.ci.SpecContents))
	for name, content := range r.ci.ArtifactContents {
		bodies[trimMarkdownSuffix(name)] = content
	}
	for name, content := range r.ci.SpecContents {
		bodies[name] = content
	}
	return bodies
}

// taskLabel renders task progress, or an empty string when the change has no
// tasks.md to count.
func (r changeRow) taskLabel() string {
	if r.ci.TasksTotal == 0 {
		return ""
	}
	return fmt.Sprintf("%d/%d", r.ci.TasksDone, r.ci.TasksTotal)
}
