package ui

import (
	"fmt"
	"regexp"
	"strings"
)

// A spec is read by the rules OpenSpec's own parser uses, transcribed from
// openspec 1.10.0, `dist/core/parsers/`. Reading by the same rules keeps two
// questions apart that were being answered as one: whether a file is a spec,
// which OpenSpec decides, and how to draw it, which this view decides.
//
// The rules, and where each comes from:
//
//	## Requirements          /^##\s+Requirements\s*$/i     spec-structure.js
//	## Purpose               a section titled Purpose      markdown-parser.js
//	### Requirement: <name>  /^###\s+Requirement:\s*(.+)/i spec-structure.js
//	a scenario               /^####\s+/ with content       requirement-text.js
//	a delta header           an error in a main spec       spec-structure.js
//	code fences              masked everywhere             code-fence.js
//
// What OpenSpec does not define is how a scenario's content is written. Its
// model is `{ rawText }`: there is no GIVEN/WHEN/THEN in it at all, and the
// `- **WHEN**` form lives in a template and an error message rather than in a
// rule. So the content is carried in full, whatever shape it is in, and laying
// a clause out under its keyword is a rendering of that text and not a
// condition of showing it.

// What a node in the outline is.
const (
	nodePurpose = iota
	nodeRequirement
	nodeScenario
	// nodeCapability roots one delta file inside a change, a change carrying
	// one per capability it touches. A main spec never produces one.
	nodeCapability
)

// What a change does to a requirement, and what a delta header names.
//
// A string rather than an enum: a delta may carry a heading none of these
// match, and the requirements under it are listed marked with what the heading
// actually said rather than hidden for being unrecognised.
const (
	opAdded    = "ADDED"
	opModified = "MODIFIED"
	opRemoved  = "REMOVED"
	opRenamed  = "RENAMED"
)

// What a scenario of a MODIFIED requirement does to the one it restates.
const (
	markNone = iota // nothing to compare against
	markUnchanged
	markEdited
	markAdded
)

// What a piece of a scenario is.
const (
	partClause = iota
	partProse
)

// specPart is one piece of a scenario's content, in the order the file gives it.
//
// Ordered rather than clauses beside a body: the two interleave in real files,
// which carry horizontal rules and `**Rationale**:` paragraphs between clauses.
// Grouping them would reorder a behaviour contract, which is worse than showing
// it plainly.
type specPart struct {
	kind    int
	keyword string // for partClause
	text    string
}

// specNode is one entry in the outline, and one card.
type specNode struct {
	kind  int
	title string
	// body is the prose under a requirement heading, or the Purpose text. A
	// scenario carries parts instead.
	body  string
	parts []specPart
	// path identifies the node across a re-parse, so a rescan that rewrites the
	// file keeps the cursor where the eye is. A title is enough: a requirement
	// heading is unique in a spec, and a scenario's is unique under it. In a
	// change the capability is prefixed, two deltas being free to modify
	// requirements of the same name in different capabilities.
	path string

	// The rest is set only for a node read from a change's delta.

	// op is the delta operation a requirement carries, empty elsewhere.
	op string
	// capability is the delta file this node came from, which is what `E`
	// opens and what the comparison is looked up against.
	capability string
	// mark says what a scenario of a MODIFIED requirement does to the one it
	// restates. markNone when there was nothing to compare against.
	mark int
	// old is the node this one modifies, when the original could be found.
	// Its presence is what gives a card its old, new and difference views.
	old *specNode
}

// specTree is a parsed spec.
type specTree struct {
	name  string
	nodes []specNode
}

// specProblem is one reason a file is not a spec, and the line it was found on.
//
// Several rather than one: a file that does not fit usually does not fit in
// more than one way, and naming only the first sends the reader back for the
// next after each repair.
type specProblem struct {
	line int // 1-based; 0 when the fault is that something is absent
	text string
}

// requirements and scenarios counts what a tree holds.
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

// The grammar, as OpenSpec writes it.
var (
	reHeading      = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	reRequirements = regexp.MustCompile(`(?i)^##\s+Requirements\s*$`)
	rePurpose      = regexp.MustCompile(`(?i)^##\s+Purpose\s*$`)
	reRequirement  = regexp.MustCompile(`(?i)^###\s+Requirement:\s*(.+?)\s*$`)
	reScenario     = regexp.MustCompile(`^####\s+`)
	reDelta        = regexp.MustCompile(`(?i)^##\s+(ADDED|MODIFIED|REMOVED|RENAMED)\s+Requirements\s*$`)
	reFence        = regexp.MustCompile("^\\s*(```|~~~)")
	reATXClose     = regexp.MustCompile(`[ \t]#+\s*$`)
	reScenPrefix   = regexp.MustCompile(`(?i)^Scenario:\s*`)
)

// fenceMask marks every line inside a fenced code block, and the fences.
//
// OpenSpec masks fences before any heading test, so a `#### Scenario:` shown as
// an example inside a fence is content rather than a scenario. Without this a
// spec that documents its own format parses as something it is not.
func fenceMask(lines []string) []bool {
	mask := make([]bool, len(lines))
	inFence := false
	for i, l := range lines {
		if reFence.MatchString(l) {
			mask[i] = true
			inFence = !inFence
			continue
		}
		mask[i] = inFence
	}
	return mask
}

// headingLevel returns the level of a heading line, or 0.
func headingLevel(line string) int {
	if m := reHeading.FindStringSubmatch(line); m != nil {
		return len(m[1])
	}
	return 0
}

// scenarioName is the label the author reads.
//
// The heading text, less an optional CommonMark closing `#` run and an optional
// `Scenario:` prefix. OpenSpec counts any level-four heading as a scenario, and
// its own comment is explicit that a header not literally reading `Scenario:`
// is still one, so the prefix is stripped where present rather than required.
func scenarioName(line string) string {
	t := reScenario.ReplaceAllString(line, "")
	t = reATXClose.ReplaceAllString(t, "")
	t = reScenPrefix.ReplaceAllString(strings.TrimSpace(t), "")
	return strings.TrimSpace(t)
}

// parseSpec reads a spec into an outline, or says why it is not one.
func parseSpec(name, content string) (specTree, []specProblem) {
	// Two views of the same file. lines is the content, which a fence is part
	// of; heads is the content with every fenced line blanked, and is what
	// every heading test reads, so none of them has to remember the mask. A
	// spec that shows its own format in a fenced block reaches this. The delta
	// grammar reads the same two views, which is why this is shared.
	lines, heads := maskedLines(content)

	var problems []specProblem
	add := func(line int, format string, args ...any) {
		problems = append(problems, specProblem{line: line, text: fmt.Sprintf(format, args...)})
	}

	// The Purpose section, and its content.
	purposeAt := -1
	for i, l := range heads {
		if rePurpose.MatchString(l) {
			purposeAt = i
			break
		}
	}
	purpose := ""
	if purposeAt >= 0 {
		purpose = strings.TrimSpace(strings.Join(sectionBody(lines, heads, purposeAt, 2), "\n"))
	}
	switch {
	case purposeAt < 0:
		add(0, "There is no `## Purpose` section. Every spec needs one: a sentence "+
			"or two on what the capability is for.")
	case purpose == "":
		add(purposeAt+1, "The `## Purpose` section is empty. OpenSpec reads an empty "+
			"Purpose as no Purpose at all.")
	}

	// A delta header is an error in a main spec, and the reason the section
	// below it is never parsed.
	for i, l := range heads {
		if m := reDelta.FindStringSubmatch(l); m != nil {
			add(i+1, "`%s` is a delta header. It belongs in a change, under "+
				"`openspec/changes/<name>/specs/`. A main spec keeps its requirements "+
				"under `## Requirements`, and openspec parses only that section, so "+
				"everything below this heading is invisible to validate, list and "+
				"archive.", strings.TrimSpace(l))
		}
	}

	// The requirements section.
	reqAt := -1
	for i, l := range heads {
		if reRequirements.MatchString(l) {
			reqAt = i
			break
		}
	}
	reqEnd := len(lines)
	if reqAt >= 0 {
		for i := reqAt + 1; i < len(lines); i++ {
			if l := headingLevel(heads[i]); l > 0 && l <= 2 {
				reqEnd = i
				break
			}
		}
	} else {
		add(0, "There is no `## Requirements` section. A main spec keeps every "+
			"requirement under that one heading.")
	}

	// Requirements, wherever they are, so one outside the section can be named.
	type reqAtLine struct {
		line  int
		title string
	}
	var found []reqAtLine
	for i, l := range heads {
		if m := reRequirement.FindStringSubmatch(l); m != nil {
			found = append(found, reqAtLine{line: i, title: strings.TrimSpace(m[1])})
		}
	}
	if len(found) == 0 {
		add(0, "The file has no requirements. A requirement is a "+
			"`### Requirement: <name>` heading.")
	}
	var outside []int
	for _, r := range found {
		if !(reqAt >= 0 && r.line > reqAt && r.line < reqEnd) {
			outside = append(outside, r.line+1)
		}
	}
	if len(outside) > 0 && reqAt >= 0 {
		add(outside[0], "%s outside the `## Requirements` section, at %s. Main specs "+
			"parse requirements only inside that section, so these are invisible to "+
			"validate, list and archive.",
			plural(len(outside), "requirement sits", "requirements sit"), lineList(outside))
	}

	// The tree, built only from what is inside the section.
	t := specTree{name: name}
	if purpose != "" {
		t.nodes = append(t.nodes, specNode{
			kind: nodePurpose, title: "Purpose", body: purpose, path: "purpose"})
	}

	var noScenario []int
	for k, r := range found {
		if !(reqAt >= 0 && r.line > reqAt && r.line < reqEnd) {
			continue
		}
		stop := reqEnd
		if k+1 < len(found) && found[k+1].line < stop {
			stop = found[k+1].line
		}

		// The requirement's own prose: up to its first scenario.
		bodyEnd := stop
		var scenLines []int
		for i := r.line + 1; i < stop; i++ {
			if reScenario.MatchString(heads[i]) {
				if len(scenLines) == 0 {
					bodyEnd = i
				}
				scenLines = append(scenLines, i)
			}
		}
		t.nodes = append(t.nodes, specNode{
			kind:  nodeRequirement,
			title: r.title,
			body:  strings.TrimSpace(strings.Join(lines[r.line+1:bodyEnd], "\n")),
			path:  "req/" + r.title,
		})

		kept := 0
		for n, at := range scenLines {
			end := stop
			if n+1 < len(scenLines) {
				end = scenLines[n+1]
			}
			content := lines[at+1 : end]
			if strings.TrimSpace(strings.Join(content, "\n")) == "" {
				// OpenSpec does not count a scenario heading with nothing under
				// it, so neither does this.
				continue
			}
			kept++
			title := scenarioName(heads[at])
			t.nodes = append(t.nodes, specNode{
				kind:  nodeScenario,
				title: title,
				parts: partsOf(content),
				path:  "req/" + r.title + "/scen/" + title,
			})
		}
		if kept == 0 {
			noScenario = append(noScenario, r.line+1)
		}
	}
	if len(noScenario) > 0 {
		add(noScenario[0], "%s no scenario, at %s. Every requirement needs at least "+
			"one `#### ` heading with content under it.",
			plural(len(noScenario), "requirement has", "requirements have"),
			lineList(noScenario))
	}

	if len(problems) > 0 {
		return specTree{}, problems
	}
	return t, nil
}

// sectionBody returns the lines under a heading, up to the next heading at or
// above its level, which is how OpenSpec delimits a section.
func sectionBody(lines, heads []string, at, level int) []string {
	for i := at + 1; i < len(lines); i++ {
		if l := headingLevel(heads[i]); l > 0 && l <= level {
			return lines[at+1 : i]
		}
	}
	return lines[at+1:]
}

// isThematicBreak reports a CommonMark horizontal rule: three or more of `-`,
// `_` or `*`, all the same character, with nothing else but spaces.
//
// Written out rather than matched, because the rule needs a backreference and
// Go's regexp has none.
func isThematicBreak(s string) bool {
	var mark rune
	count := 0
	for _, r := range s {
		switch {
		case r == ' ' || r == '\t':
		case r == '-' || r == '_' || r == '*':
			if mark == 0 {
				mark = r
			}
			if r != mark {
				return false
			}
			count++
		default:
			return false
		}
	}
	return count >= 3
}

// partsOf reads a scenario's content into ordered parts.
//
// Nothing is discarded for being in an unrecognised shape. A line that opens
// with a keyword bullet starts a clause, a line after one continues it, a blank
// line ends it, and anything else is prose in the place it was written.
func partsOf(content []string) []specPart {
	var parts []specPart
	var pending *specPart

	flush := func() {
		if pending != nil {
			pending.text = strings.TrimSpace(pending.text)
			if pending.text != "" || pending.keyword != "" {
				parts = append(parts, *pending)
			}
			pending = nil
		}
	}

	for _, line := range content {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			flush()
		case isThematicBreak(trimmed):
			// A separator between requirements in the file, not part of the
			// behaviour described.
			flush()
		default:
			if kw, rest, ok := clauseOf(line); ok {
				flush()
				pending = &specPart{kind: partClause, keyword: kw, text: rest}
				continue
			}
			if pending == nil {
				pending = &specPart{kind: partProse}
			}
			if pending.text != "" {
				pending.text += " "
			}
			pending.text += trimmed
		}
	}
	flush()
	return parts
}

// clauseKeywords are the four a scenario bullet may open with.
var clauseKeywords = []string{"GIVEN", "WHEN", "THEN", "AND"}

// clauseOf reads a keyword bullet, or reports that the line is not one.
//
// A bulleted, upper-case keyword, which is what the OpenSpec template shows and
// what 5080 of the 5262 clause lines in the local corpus are written as. It is
// deliberately not widened: since nothing is dropped any more, which shapes get
// the keyword laid out on its own row is presentation, and the other shapes are
// shown as the prose they are.
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

// plural picks the singular or plural form, with the count in front.
func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}

// lineList names up to four lines, and counts the rest.
func lineList(lines []int) string {
	const show = 4
	var parts []string
	for i, l := range lines {
		if i == show {
			return strings.Join(parts, ", ") +
				fmt.Sprintf(" and %d more", len(lines)-show)
		}
		parts = append(parts, fmt.Sprintf("line %d", l))
	}
	return strings.Join(parts, ", ")
}
