package ui

import (
	"fmt"
	"strings"

	"github.com/mipmip/specgetty/src/scanner"
)

// A project that reads its content from a store has more than one
// configuration to show. The repo keeps its own context, rules and operations;
// the store keeps what everything using it shares. Showing only the root's
// would drop the repo's half off the screen, which is the regression these
// panes exist to prevent.
const (
	paneDocument = iota // a configuration file, rendered as its own format
	paneDetails         // the store report, generated rather than read
)

// configPane is one sub-tab of the config tab.
type configPane struct {
	label   string // what the sub-tab says
	source  string // the file it came from, for the single-pane label
	kind    int
	content string
	md      bool // render as markdown rather than YAML
}

// configPanes returns what the config tab has to show for a project, in order.
//
// A plain project has one pane and draws no sub-tab row, which is exactly the
// tab as it was before stores existed.
func configPanes(info scanner.ProjectInfo) []configPane {
	var panes []configPane

	if info.ResolvedElsewhere() && info.OriginConfigFile != "" {
		panes = append(panes, configPane{
			label:   "repo",
			source:  "openspec/" + info.OriginConfigFile,
			kind:    paneDocument,
			content: info.OriginConfigContent,
			md:      strings.HasSuffix(info.OriginConfigFile, ".md"),
		})
	}

	if info.ConfigFile != "" {
		label := "config"
		if info.FromStore() {
			label = "store"
		}
		panes = append(panes, configPane{
			label:   label,
			source:  "openspec/" + info.ConfigFile,
			kind:    paneDocument,
			content: info.ConfigContent,
			md:      strings.HasSuffix(info.ConfigFile, ".md"),
		})
	}

	if info.FromStore() {
		panes = append(panes, configPane{
			label:   "store details",
			source:  "store",
			kind:    paneDetails,
			content: storeDetails(info),
		})
	}

	return panes
}

// storeDetails renders the report behind the details sub-tab.
//
// This is where the header's single mark is explained, which is what makes
// that mark affordable. Everything here is read from local files: the registry,
// the store's identity file, and git's own refs. Nothing fetches.
func storeDetails(info scanner.ProjectInfo) string {
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
		if info.Origin != "" {
			b.WriteString("origin: " + info.Origin + "\n")
		}
		return b.String()
	}

	s := info.Store
	if s == nil {
		return "problem: no store information\n"
	}

	b.WriteString("store: " + s.ID + "\n")
	b.WriteString("root: " + s.Root + "\n")
	if s.Origin != "" && s.Origin != s.Root {
		b.WriteString("origin: " + s.Origin + "\n")
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

// storeProblemLine is the one-line report shown where a tab would otherwise
// say the project is empty.
//
// An unfollowed declaration leaves zero specs and zero changes, which reads as
// a project with nothing in it. That is the exact failure this change exists to
// fix, so the emptiness is never shown without its reason.
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
