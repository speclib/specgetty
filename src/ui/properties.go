package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// The properties tab answers "what is this project and where do its parts come
// from". None of its sections are the project's content: changes and specs are
// what the project says, while its configuration, its workflow schemas and the
// store its content lives in are how it is wired.
//
// The sections are a vertical list because one of them repeats: a project runs
// as many schemas as its changes record, and a row each keeps every one of them
// a document of its own.
const (
	sectionConfig = iota // the one configuration that applies
	sectionSchema        // one per schema the project uses
	sectionStore         // where the content comes from
)

// propSection is one row of the properties list.
type propSection struct {
	label  string
	kind   int
	schema string // the schema this row reports, for sectionSchema
	source string // the file it came from, for the document sections
	md     bool   // render as markdown rather than YAML
}

// propSections returns the rows the properties tab offers for a project.
//
// The shape is the same for every project: the configuration, a row per schema,
// and the store. A project that holds its own content still gets a store row,
// reading `local`, because a row that says so teaches the concept where an
// absent one teaches nothing.
func propSections(info scanner.ProjectInfo) []propSection {
	sections := []propSection{{
		label:  "project",
		kind:   sectionConfig,
		source: configSourceLabel(info),
		md:     strings.HasSuffix(info.ConfigFile, ".md"),
	}}
	for _, u := range info.SchemaUsage {
		sections = append(sections, propSection{
			label: u.Name, kind: sectionSchema, schema: u.Name,
		})
	}
	sections = append(sections, propSection{label: "store", kind: sectionStore})
	return sections
}

// configSourceLabel names the file the configuration was read from.
//
// For a store-backed project this is the store's file, which is the point of
// naming it: the content did not come from the directory the user is standing
// in, and nothing else on screen would say so.
func configSourceLabel(info scanner.ProjectInfo) string {
	if info.ConfigFile == "" {
		return ""
	}
	if info.ResolvedElsewhere() {
		return info.Root + "/openspec/" + info.ConfigFile
	}
	return "openspec/" + info.ConfigFile
}

// renderSchemaSection reports one schema: what it is, where it came from, and
// how much of the project runs on it.
func renderSchemaSection(info scanner.ProjectInfo, name string, state schemaState) string {
	var b strings.Builder

	for _, u := range info.SchemaUsage {
		if u.Name != name {
			continue
		}
		b.WriteString("schema: " + u.Name + "\n")
		if u.IsDefault {
			b.WriteString("default: yes\n")
		}
		b.WriteString(fmt.Sprintf("changes: %d\n", u.Changes))
	}
	if info.UnrecordedChanges > 0 {
		b.WriteString(fmt.Sprintf("# %d change(s) in this project record no schema at all\n", info.UnrecordedChanges))
	}
	b.WriteString("\n")

	switch {
	case state.problem != nil:
		p := state.problem
		b.WriteString("problem: its definition could not be read\n")
		b.WriteString("reason: " + p.Detail + "\n")
		if len(p.Available) > 0 {
			b.WriteString("available: " + strings.Join(p.Available, ", ") + "\n")
		}
		return b.String()
	case state.details == nil:
		b.WriteString("# reading its definition...\n")
		return b.String()
	}

	d := state.details
	b.WriteString("source: " + d.Source + "\n")
	b.WriteString("path: " + d.Path + "\n")
	if d.Description != "" {
		b.WriteString("description: " + firstLine(d.Description) + "\n")
	}
	for _, sh := range d.Shadows {
		b.WriteString("overrides: " + sh.Source + " " + sh.Path + "\n")
	}

	b.WriteString("\nartifacts:\n")
	for _, a := range d.Artifacts {
		line := "  " + a.ID
		if a.Generates != "" {
			line += " -> " + a.Generates
		}
		if len(a.Requires) > 0 {
			line += " (after " + strings.Join(a.Requires, ", ") + ")"
		}
		b.WriteString(line + "\n")
	}

	b.WriteString("\napply:\n")
	if len(d.Apply.Requires) > 0 {
		b.WriteString("  requires: " + strings.Join(d.Apply.Requires, ", ") + "\n")
	}
	if d.Apply.Tracks != "" {
		b.WriteString("  tracks: " + d.Apply.Tracks + "\n")
	}
	return b.String()
}

// firstLine keeps a description to one line. A schema's prose runs to
// paragraphs, and this tab reports what a schema is rather than instructing
// whoever writes its artifacts.
func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

// renderStoreSection reports where the project's content comes from.
func renderStoreSection(info scanner.ProjectInfo) string {
	var b strings.Builder

	if p := info.StoreProblem; p != nil {
		b.WriteString("problem: the declared store could not be followed\n")
		if p.ID != "" {
			b.WriteString("declared: " + p.ID + "\n")
		}
		if p.File != "" {
			b.WriteString("declared_in: " + p.File + "\n")
		}
		b.WriteString("reason: " + p.Detail + "\n")
		return b.String()
	}

	s := info.Store
	if s == nil {
		b.WriteString("store: local\n")
		b.WriteString("# this project holds its own specs and changes\n")
		return b.String()
	}

	b.WriteString("store: " + s.ID + "\n")
	b.WriteString("root: " + s.Root + "\n")
	if s.Origin != "" && s.Origin != s.Root {
		b.WriteString("declared_in: " + s.Origin + "/openspec/\n")
	}
	if s.Remote != "" {
		b.WriteString("registered_remote: " + s.Remote + "\n")
	}
	if s.Branch != "" {
		b.WriteString("registered_branch: " + s.Branch + "\n")
	}
	if s.Canonical != "" {
		b.WriteString("canonical_remote: " + s.Canonical + "\n")
	}

	if len(info.InertKeys) > 0 {
		b.WriteString("\n# the declaring file also sets " + strings.Join(info.InertKeys, ", ") + "\n")
		b.WriteString("# none of it has any effect: OpenSpec reads that file for\n")
		b.WriteString("# `store:` alone and takes the rest from the store above\n")
	}

	b.WriteString("\n")
	g := s.Git
	if g == nil || !g.IsRepo {
		b.WriteString("git: the store root is not a git working copy\n")
		return b.String()
	}

	b.WriteString("git:\n")
	if g.OriginURL != "" {
		b.WriteString("  origin_url: " + g.OriginURL + "\n")
	}
	if g.DirtyKnown {
		if g.Dirty {
			b.WriteString("  uncommitted_changes: yes\n")
		} else {
			b.WriteString("  uncommitted_changes: no\n")
		}
	}
	if g.TrackingKnown {
		b.WriteString(fmt.Sprintf("  ahead: %d\n", g.Ahead))
		b.WriteString(fmt.Sprintf("  behind: %d\n", g.Behind))
		// Said plainly, because a number that looks fetched and was not is
		// worse than no number at all.
		b.WriteString("  # compared against the last known upstream ref, without fetching\n")
	} else {
		b.WriteString("  # no upstream to compare against\n")
	}
	return b.String()
}

// storeProblemLine is the one-line report shown where a tab would otherwise say
// the project is empty.
func storeProblemLine(info scanner.ProjectInfo) string {
	p := info.StoreProblem
	if p == nil {
		return ""
	}
	if p.ID != "" {
		return fmt.Sprintf("Store %q could not be followed: %s", p.ID, p.Detail)
	}
	return "The store declaration could not be followed: " + p.Detail
}

// storeMarkStyleRef keeps the header mark's one job in one place: it says the
// content came from a store and nothing else. The id, the path, the remote and
// the git state are answers to a question asked once per project, not once per
// glance, and they live on the properties tab's store row.
var _ = lipgloss.NewStyle

// storeMark is the header's one sign that the content came from a store.
func storeMark(info scanner.ProjectInfo) string {
	if !info.FromStore() {
		return ""
	}
	return " " + storeMarkStyle.Render("store")
}

// headerName is what the header calls the open project.
//
// The user is standing in the repo, so that is what the header names. A store
// opened on its own account has no repo to name and is called by its id, which
// is what every other OpenSpec surface calls it.
func headerName(info scanner.ProjectInfo, key string) string {
	if info.ResolvedElsewhere() {
		return info.Origin
	}
	if info.StoreID != "" {
		return info.StoreID
	}
	return key
}
