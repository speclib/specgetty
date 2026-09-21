package ui

import (
	"fmt"
	"strings"
)

// A spec is read as one markdown document today, which for `change-list-view`
// is 425 rows in the 49 columns the specs tab gives it. The structure is in the
// file and the renderer throws it away.
//
// Across the live specs in this project there are no headings outside four
// shapes, no requirement without its scenarios, and no scenario without its
// keyword bullets. That is reliable enough to navigate by, and unreliable
// enough that a file which does not fit it must be refused rather than
// half-parsed.

// What a node in the outline is.
const (
	nodePurpose = iota
	nodeRequirement
	nodeScenario
)

// specClause is one keyword bullet of a scenario: `GIVEN`, `WHEN`, `THEN` or
// `AND`, and the text that follows it.
type specClause struct {
	keyword string
	text    string
}

// specNode is one entry in the outline, and one card.
type specNode struct {
	kind  int
	title string
	// body is the prose under a requirement heading, or the Purpose text. A
	// scenario carries clauses instead.
	body    string
	clauses []specClause
	// path identifies the node across a re-parse, so a rescan that rewrites the
	// file keeps the cursor where the eye is. A title is enough: a requirement
	// heading is unique in a spec, and a scenario's is unique under it.
	path string
}

// specTree is a parsed spec.
type specTree struct {
	name  string
	nodes []specNode
}

// requirements and scenarios counts what a tree holds, which is what the
// corpus test asserts over every live spec.
func (t specTree) counts() (requirements, scenarios int) {
	for _, n := range t.nodes {
		switch n.kind {
		case nodeRequirement:
			requirements++
		case nodeScenario:
			scenarios++
		}
	}
	return
}

// indexOfPath finds a node by the path it carried before a re-parse, or -1.
func (t specTree) indexOfPath(path string) int {
	for i, n := range t.nodes {
		if n.path == path {
			return i
		}
	}
	return -1
}

// clauseKeywords are the four a scenario bullet may open with.
var clauseKeywords = []string{"GIVEN", "WHEN", "THEN", "AND"}

// parseSpec reads a spec into an outline, or refuses it.
//
// A file that yields no requirements, or a requirement with no scenarios, is not
// a spec this view can navigate, and showing an outline of nothing would be
// worse than saying so.
func parseSpec(name, content string) (specTree, error) {
	t := specTree{name: name}

	var cur *specNode
	flush := func() {
		if cur != nil {
			cur.body = strings.TrimSpace(cur.body)
			t.nodes = append(t.nodes, *cur)
			cur = nil
		}
	}

	var pending *specClause
	flushClause := func() {
		if pending != nil && cur != nil {
			pending.text = strings.TrimSpace(pending.text)
			cur.clauses = append(cur.clauses, *pending)
		}
		pending = nil
	}

	var lastRequirement string
	for _, raw := range strings.Split(content, "\n") {
		line := strings.TrimRight(raw, " \t")
		switch {
		case strings.HasPrefix(line, "## Purpose"):
			flushClause()
			flush()
			cur = &specNode{kind: nodePurpose, title: "Purpose", path: "purpose"}

		case strings.HasPrefix(line, "### Requirement: "):
			flushClause()
			flush()
			title := strings.TrimSpace(strings.TrimPrefix(line, "### Requirement: "))
			lastRequirement = title
			cur = &specNode{kind: nodeRequirement, title: title, path: "req/" + title}

		case strings.HasPrefix(line, "#### Scenario: "):
			flushClause()
			flush()
			title := strings.TrimSpace(strings.TrimPrefix(line, "#### Scenario: "))
			cur = &specNode{kind: nodeScenario, title: title,
				path: "req/" + lastRequirement + "/scen/" + title}

		case strings.HasPrefix(line, "# "), strings.HasPrefix(line, "## "):
			// The spec's title line, and the `## Requirements` heading that
			// separates the Purpose from the rest. Neither is a node, and both
			// end whatever was being read: without this the Purpose card
			// absorbs the heading below it.
			flushClause()
			flush()

		default:
			if cur == nil {
				continue
			}
			if cur.kind == nodeScenario {
				if kw, rest, ok := clauseOf(line); ok {
					flushClause()
					pending = &specClause{keyword: kw, text: rest}
					continue
				}
				// A clause may span several source lines; a blank line ends it.
				if pending != nil {
					if strings.TrimSpace(line) == "" {
						flushClause()
					} else {
						pending.text += " " + strings.TrimSpace(line)
					}
					continue
				}
				continue
			}
			cur.body += line + "\n"
		}
	}
	flushClause()
	flush()

	if req, scen := t.counts(); req == 0 || scen == 0 {
		return specTree{}, fmt.Errorf("%s has %d requirements and %d scenarios; not a navigable spec", name, req, scen)
	}
	for i, n := range t.nodes {
		if n.kind != nodeRequirement {
			continue
		}
		// A requirement with nothing under it before the next requirement.
		next := i + 1
		if next >= len(t.nodes) || t.nodes[next].kind != nodeScenario {
			return specTree{}, fmt.Errorf("%s: requirement %q has no scenarios", name, n.title)
		}
	}
	return t, nil
}

// clauseOf reads a keyword bullet, or reports that the line is not one.
func clauseOf(line string) (keyword, text string, ok bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "- ") {
		return "", "", false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
	for _, kw := range clauseKeywords {
		for _, form := range []string{"**" + kw + "**", kw} {
			if strings.HasPrefix(rest, form) {
				return kw, strings.TrimSpace(strings.TrimPrefix(rest, form)), true
			}
		}
	}
	return "", "", false
}
