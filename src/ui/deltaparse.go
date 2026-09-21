package ui

import (
	"fmt"
	"regexp"
	"strings"
)

// A spec file in a change is a delta, and almost every rule that makes a main
// spec valid inverts in it. A delta header is the structure rather than a
// fault, `## Purpose` is written only for a capability the change introduces,
// and a requirement being removed carries a reason where a scenario would be.
//
// So this is an entry point of its own rather than a flag on parseSpec. The two
// share the mechanical half, the fence mask, the section reader and partsOf,
// and state their own validity rules, which is the part that genuinely differs.

var (
	// Any `## <word> Requirements` heading, so a delta carrying an operation
	// none of the four known ones matches is still read rather than dropped.
	reOpHeader = regexp.MustCompile(`(?i)^##\s+([A-Za-z]+)\s+Requirements\s*$`)
)

// knownOps are the four OpenSpec applies. An operation outside this set is
// carried through as written, so the reader sees what the file says.
var knownOps = map[string]bool{
	opAdded: true, opModified: true, opRemoved: true, opRenamed: true,
}

// needsScenarios reports whether a requirement under this operation must carry
// one. A removal names a reason instead, and a rename names only the two names.
func needsScenarios(op string) bool {
	switch op {
	case opRemoved, opRenamed:
		return false
	}
	return true
}

// parseDelta reads one spec file of a change, or reports why it cannot.
//
// The nodes carry the capability they came from, so a change's several deltas
// can be joined into one outline without losing which file a node opens.
func parseDelta(capability, content string) (specTree, []specProblem) {
	lines, heads := maskedLines(content)

	var problems []specProblem
	add := func(line int, format string, args ...any) {
		problems = append(problems, specProblem{line: line, text: fmt.Sprintf(format, args...)})
	}

	t := specTree{name: capability}

	// A `## Requirements` section is the main-spec form. OpenSpec reads only
	// delta headers in a change, so requirements under it are never applied.
	for i, l := range heads {
		if reRequirements.MatchString(l) {
			add(i+1, "`## Requirements` is the heading a main spec uses. In a change "+
				"openspec reads only `## ADDED`, `## MODIFIED`, `## REMOVED` and "+
				"`## RENAMED Requirements`, so nothing under this heading will be "+
				"applied when the change is archived.")
			break
		}
	}

	// Purpose, which a delta carries only for a capability the change
	// introduces. Absent is the ordinary case and not a fault.
	for i, l := range heads {
		if rePurpose.MatchString(l) {
			body := strings.TrimSpace(strings.Join(sectionBody(lines, heads, i, 2), "\n"))
			if body == "" {
				add(i+1, "The `## Purpose` section is empty. Either write one or "+
					"remove the heading; archive copies it into the new spec.")
				break
			}
			t.nodes = append(t.nodes, specNode{
				kind: nodePurpose, title: "Purpose", body: body,
				capability: capability, path: pathOf(capability, "purpose"),
			})
			break
		}
	}

	// Where each operation section starts, so a requirement can be given the
	// operation it sits under.
	type opAt struct {
		line int
		op   string
	}
	var ops []opAt
	for i, l := range heads {
		if m := reOpHeader.FindStringSubmatch(l); m != nil {
			ops = append(ops, opAt{line: i, op: strings.ToUpper(m[1])})
		}
	}
	opFor := func(line int) string {
		op := ""
		for _, o := range ops {
			if o.line < line {
				op = o.op
			}
		}
		return op
	}

	// Requirements, wherever they sit.
	type reqAt struct {
		line  int
		title string
	}
	var found []reqAt
	for i, l := range heads {
		if m := reRequirement.FindStringSubmatch(l); m != nil {
			found = append(found, reqAt{line: i, title: strings.TrimSpace(m[1])})
		}
	}
	if len(found) == 0 {
		add(0, "The file has no requirements. A delta needs at least one "+
			"`### Requirement: <name>` under a delta header.")
	}

	var orphans []int
	var noScenario []int
	for k, r := range found {
		op := opFor(r.line)
		if op == "" {
			orphans = append(orphans, r.line+1)
			continue
		}

		stop := len(lines)
		if k+1 < len(found) {
			stop = found[k+1].line
		}
		// An operation heading below this requirement ends it too, so the last
		// requirement of a section does not swallow the next section's header.
		for _, o := range ops {
			if o.line > r.line && o.line < stop {
				stop = o.line
			}
		}

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

		reqPath := pathOf(capability, "req/"+r.title)
		t.nodes = append(t.nodes, specNode{
			kind:       nodeRequirement,
			title:      r.title,
			body:       strings.TrimSpace(strings.Join(lines[r.line+1:bodyEnd], "\n")),
			op:         op,
			capability: capability,
			path:       reqPath,
		})

		kept := 0
		for n, at := range scenLines {
			end := stop
			if n+1 < len(scenLines) {
				end = scenLines[n+1]
			}
			body := lines[at+1 : end]
			if strings.TrimSpace(strings.Join(body, "\n")) == "" {
				continue
			}
			kept++
			title := scenarioName(heads[at])
			t.nodes = append(t.nodes, specNode{
				kind:       nodeScenario,
				title:      title,
				parts:      partsOf(body),
				capability: capability,
				path:       reqPath + "/scen/" + title,
			})
		}
		if kept == 0 && needsScenarios(op) {
			noScenario = append(noScenario, r.line+1)
		}
	}

	if len(orphans) > 0 {
		add(orphans[0], "%s under no delta header, at %s. openspec applies a "+
			"requirement only under `## ADDED`, `## MODIFIED`, `## REMOVED` or "+
			"`## RENAMED Requirements`.",
			plural(len(orphans), "requirement sits", "requirements sit"), lineList(orphans))
	}
	if len(noScenario) > 0 {
		add(noScenario[0], "%s no scenario, at %s. An added or modified requirement "+
			"needs at least one `#### ` heading with content under it; only a "+
			"removal or a rename may go without.",
			plural(len(noScenario), "requirement has", "requirements have"),
			lineList(noScenario))
	}

	if len(problems) > 0 {
		return specTree{}, problems
	}
	return t, nil
}

// pathOf prefixes a node path with the capability it came from, two deltas of
// one change being free to touch requirements of the same name.
func pathOf(capability, rest string) string { return "cap/" + capability + "/" + rest }

// maskedLines returns a file's lines and the same lines with every fenced one
// blanked, which is what every heading test reads.
//
// Lifted out of parseSpec so both grammars mask fences the same way. OpenSpec
// masks before any heading test, so a `#### Scenario:` shown as an example
// inside a fence is content rather than a scenario.
func maskedLines(content string) (lines, heads []string) {
	content = strings.TrimPrefix(content, "\ufeff")
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	raw := strings.Split(content, "\n")
	mask := fenceMask(raw)

	lines = make([]string, len(raw))
	heads = make([]string, len(raw))
	for i, l := range raw {
		lines[i] = strings.TrimRight(l, " \t")
		if !mask[i] {
			heads[i] = lines[i]
		}
	}
	return lines, heads
}
