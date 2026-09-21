package ui

import (
	"fmt"
	"sort"

	"github.com/mipmip/specgetty/src/scanner"
)

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

// changeGroup is a run of rows under one heading.
//
// Active and archived changes are always both shown. The group a row sits in is
// what says which it is, by position and permanently, where a filter mode said
// it with a label somewhere else and could be entered without noticing.
type changeGroup struct {
	label string
	rows  []changeRow
}

// Group labels. One spelling, used on screen and in the tests.
const (
	groupActive   = "ACTIVE"
	groupArchived = "ARCHIVED"
)

// buildGroups turns what the scanner produced into the two groups, in the order
// they are shown.
//
// Active changes lead, because they are what the developer is working on.
// Within the groups the orders are fixed and different: active alphabetically,
// archived most recent first, so that the work most recently finished is at the
// top of its group rather than thirty rows down.
func buildGroups(info scanner.ProjectInfo) []changeGroup {
	active := make([]changeRow, 0, len(info.Changes))
	for _, ci := range info.Changes {
		active = append(active, changeRow{ci: ci})
	}
	sort.SliceStable(active, func(i, j int) bool {
		return active[i].ci.Name < active[j].ci.Name
	})

	archived := make([]changeRow, 0, len(info.ArchivedChanges))
	for _, ci := range info.ArchivedChanges {
		archived = append(archived, changeRow{ci: ci, archived: true})
	}
	// Stable, so two changes archived on the same day keep the order the
	// scanner read them in and two renders agree. A change whose directory
	// carries no parsable date has the zero time, which sorts last here rather
	// than first, so it lands after those that do.
	sort.SliceStable(archived, func(i, j int) bool {
		return archived[i].ci.ArchiveDate.After(archived[j].ci.ArchiveDate)
	})

	return []changeGroup{
		{label: groupActive, rows: active},
		{label: groupArchived, rows: archived},
	}
}

// allRowsOf flattens the groups into the selection order, which is what the
// cursor indexes and what the filter narrows.
func allRowsOf(groups []changeGroup) []changeRow {
	var rows []changeRow
	for _, g := range groups {
		rows = append(rows, g.rows...)
	}
	return rows
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
