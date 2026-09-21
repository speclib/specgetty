package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// groupedInfo builds a project with the given active names and archived
// name/date pairs.
func groupedInfo(active []string, archived map[string]string) scanner.ProjectInfo {
	info := scanner.ProjectInfo{Root: "/p", Origin: "/p"}
	for _, n := range active {
		info.Changes = append(info.Changes, scanner.ChangeInfo{Name: n, DirName: n})
	}
	for n, d := range archived {
		ci := scanner.ChangeInfo{Name: n, DirName: d + "-" + n}
		if t, err := time.Parse("2006-01-02", d); err == nil {
			ci.ArchiveDate = t
		}
		info.ArchivedChanges = append(info.ArchivedChanges, ci)
	}
	return info
}

// --- 2.1 to 2.6 the groups and their orders ---

func TestActiveGroupLeads(t *testing.T) {
	groups := buildGroups(groupedInfo([]string{"one"}, map[string]string{"old": "2026-01-01"}))
	if groups[0].label != groupActive {
		t.Errorf("got %q first, want the active group", groups[0].label)
	}
}

func TestActiveGroupIsOrderedByName(t *testing.T) {
	info := scanner.ProjectInfo{Changes: []scanner.ChangeInfo{
		{Name: "zebra"}, {Name: "alpha"}, {Name: "middle"},
	}}
	got := plainNames(buildGroups(info)[0].rows)
	want := []string{"alpha", "middle", "zebra"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestArchivedGroupIsNewestFirst(t *testing.T) {
	// The order that made the change archived most recently sit at the bottom
	// of thirty-five rows, reversed.
	info := groupedInfo(nil, map[string]string{
		"oldest": "2026-03-31", "middle": "2026-06-15", "newest": "2026-09-21",
	})
	got := plainNames(buildGroups(info)[1].rows)
	want := []string{"newest", "middle", "oldest"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestArchivedChangesSharingADateKeepAStableOrder(t *testing.T) {
	info := groupedInfo(nil, nil)
	same, _ := time.Parse("2006-01-02", "2026-09-21")
	for _, n := range []string{"first", "second", "third"} {
		info.ArchivedChanges = append(info.ArchivedChanges,
			scanner.ChangeInfo{Name: n, ArchiveDate: same})
	}
	one := plainNames(buildGroups(info)[1].rows)
	two := plainNames(buildGroups(info)[1].rows)
	for i := range one {
		if one[i] != two[i] {
			t.Fatalf("two renders disagree: %v then %v", one, two)
		}
	}
}

func TestArchivedChangeWithNoDateSortsAfterThoseWithOne(t *testing.T) {
	info := groupedInfo(nil, map[string]string{"dated": "2026-01-01"})
	info.ArchivedChanges = append(info.ArchivedChanges,
		scanner.ChangeInfo{Name: "undated", DirName: "undated"})

	got := plainNames(buildGroups(info)[1].rows)
	if len(got) != 2 {
		t.Fatalf("got %v, want both listed", got)
	}
	if got[len(got)-1] != "undated" {
		t.Errorf("got %v, want the undated change last rather than dropped or floating", got)
	}
}

func TestAGroupKeepsItsHeaderWhenEmpty(t *testing.T) {
	for _, tc := range []struct {
		name string
		info scanner.ProjectInfo
	}{
		{"no active", groupedInfo(nil, map[string]string{"old": "2026-01-01"})},
		{"no archived", groupedInfo([]string{"one"}, nil)},
		{"neither", groupedInfo(nil, nil)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			groups := buildGroups(tc.info)
			if len(groups) != 2 {
				t.Fatalf("got %d groups, want both always", len(groups))
			}
			lines := groupLines(groups, wrap(allRowsOf(groups)))
			headers := 0
			for _, l := range lines {
				if l.kind == lineHeader {
					headers++
				}
			}
			if headers != 2 {
				t.Errorf("got %d headers, want one per group", headers)
			}
		})
	}
}

// --- 3.3 the counts ---

func TestGroupHeadersCountTheirRows(t *testing.T) {
	groups := buildGroups(groupedInfo([]string{"a", "b"},
		map[string]string{"x": "2026-01-01", "y": "2026-02-01", "z": "2026-03-01"}))
	lines := groupLines(groups, wrap(allRowsOf(groups)))

	var headers []string
	for _, l := range lines {
		if l.kind == lineHeader {
			headers = append(headers, l.header)
		}
	}
	if headers[0] != "ACTIVE (2)" || headers[1] != "ARCHIVED (3)" {
		t.Errorf("got %v, want the counts", headers)
	}
}

func TestGroupHeadersCountWhatAFilterLeft(t *testing.T) {
	groups := buildGroups(groupedInfo([]string{"keep-me", "drop-me"},
		map[string]string{"keep-me-too": "2026-01-01", "drop-me-too": "2026-02-01"}))
	all := allRowsOf(groups)

	var kept []filtered[changeRow]
	for _, r := range all {
		if strings.HasPrefix(r.ci.Name, "keep") {
			kept = append(kept, filtered[changeRow]{row: r})
		}
	}
	lines := groupLines(groups, kept)

	var headers []string
	for _, l := range lines {
		if l.kind == lineHeader {
			headers = append(headers, l.header)
		}
	}
	if headers[0] != "ACTIVE (1)" || headers[1] != "ARCHIVED (1)" {
		t.Errorf("got %v, want what the filter left in each", headers)
	}
}

// --- 3.4 and 3.5 the index-to-line mapping, where every off-by-one lives ---

func TestLineOfRowCountsHeaders(t *testing.T) {
	groups := buildGroups(groupedInfo([]string{"a", "b"},
		map[string]string{"x": "2026-02-01", "y": "2026-01-01"}))
	lines := groupLines(groups, wrap(allRowsOf(groups)))

	// index  line   what is drawn
	//   -     0     ACTIVE (2)
	//   0     1     a
	//   1     2     b
	//   -     3     (blank)
	//   -     4     ARCHIVED (2)
	//   2     5     x
	//   3     6     y
	for row, wantLine := range map[int]int{0: 1, 1: 2, 2: 5, 3: 6} {
		if got := lineOfRow(lines, row); got != wantLine {
			t.Errorf("change %d is drawn on line %d, want %d", row, got, wantLine)
		}
	}
	if got := lineOfRow(lines, 99); got != -1 {
		t.Errorf("a change the filter removed has no line, got %d", got)
	}
}

func TestScrollingKeepsTheSelectedChangeVisibleAcrossAGroupBoundary(t *testing.T) {
	// Only shows up when the list is longer than the panel and the cursor is
	// near a boundary, which is why it is constructed rather than stumbled on.
	archived := map[string]string{}
	for i := 0; i < 12; i++ {
		archived[fmt.Sprintf("arch-%02d", i)] = fmt.Sprintf("2026-01-%02d", i+1)
	}
	groups := buildGroups(groupedInfo([]string{"act-a", "act-b"}, archived))
	rows := wrap(allRowsOf(groups))
	lines := groupLines(groups, rows)

	const height = 8 // one header row plus seven body lines
	for cursor := range rows {
		out := renderGroupedTable(groups, rows, changeFieldDefs(defaultFields),
			cursor, 80, height)
		body := strings.Split(ansi.Strip(out), "\n")[1:]

		want := rows[cursor].row.ci.Name
		var seen bool
		for _, l := range body {
			if strings.Contains(l, want) {
				seen = true
			}
		}
		if !seen {
			t.Errorf("cursor %d (%s, line %d): the selected change is off screen:\n%s",
				cursor, want, lineOfRow(lines, cursor), strings.Join(body, "\n"))
		}
		if len(body) > height-1 {
			t.Fatalf("cursor %d: %d body rows, want at most %d", cursor, len(body), height-1)
		}
	}
}

// --- 3.1 and 3.2 drawing ---

func TestGroupedTableKeepsOneColumnHeader(t *testing.T) {
	groups := buildGroups(groupedInfo([]string{"a"}, map[string]string{"x": "2026-01-01"}))
	out := ansi.Strip(renderGroupedTable(groups, wrap(allRowsOf(groups)),
		changeFieldDefs(defaultFields), 0, 80, 10))
	lines := strings.Split(out, "\n")

	if !strings.Contains(lines[0], "name") || !strings.Contains(lines[0], "tasks") {
		t.Errorf("the first line is the column header: %q", lines[0])
	}
	if strings.Count(out, "name") != 1 {
		t.Error("one column header serves both groups")
	}
}

func TestGroupHeadersUseTheExistingSectionStyle(t *testing.T) {
	groups := buildGroups(groupedInfo([]string{"a"}, nil))
	out := renderGroupedTable(groups, wrap(allRowsOf(groups)),
		changeFieldDefs(defaultFields), 0, 80, 10)
	want := sectionHeaderStyle.Render("x")
	prefix := want[:strings.Index(want, "x")]
	if !strings.Contains(out, prefix+"ACTIVE") {
		t.Error("group headers use the style spec names already use")
	}
}

// --- 5.x the columns ---

func TestTheDateIsADefaultColumnAndBlankOnActiveRows(t *testing.T) {
	groups := buildGroups(groupedInfo([]string{"an-active-change"},
		map[string]string{"an-archived-change": "2026-09-21"}))
	out := ansi.Strip(renderGroupedTable(groups, wrap(allRowsOf(groups)),
		changeFieldDefs(defaultFields), 0, 100, 10))

	if !strings.Contains(out, "2026-09-21") {
		t.Errorf("the archived date must be shown by default:\n%s", out)
	}
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, "an-active-change") && strings.Contains(l, "2026") {
			t.Errorf("an active row carries no date: %q", l)
		}
	}
}

func TestTheStateColumnIsNoLongerADefaultButRemainsAvailable(t *testing.T) {
	for _, id := range defaultFields {
		if id == "archived" {
			t.Error("the group header says it; the column is redundant")
		}
	}
	if _, ok := knownFields["archived"]; !ok {
		t.Error("it must remain configurable")
	}
	got := changeFieldDefs([]string{"name", "archived"})
	if len(got) != 2 || got[1].id != "archived" {
		t.Errorf("got %+v, want it usable when named", got)
	}
}

func TestTheNameColumnSurvivesFourDefaultsAtTheMinimumWidth(t *testing.T) {
	kept, widths := layoutFields(changeFieldDefs(defaultFields), 60-8)
	if len(kept) != len(defaultFields) {
		t.Fatalf("got %d columns at the 60-column minimum, want all %d", len(kept), len(defaultFields))
	}
	if widths[0] < minFlexWidth {
		t.Errorf("the name column is %d wide, want at least %d", widths[0], minFlexWidth)
	}
}

// --- 3.1 the picker is untouched ---

func TestThePickerStillUsesThePlainTable(t *testing.T) {
	projects := scanner.ProjectMap{
		"/work/alpha": {Info: scanner.ProjectInfo{Root: "/work/alpha", Origin: "/work/alpha", SpecCount: 3}},
		"/work/beta":  {Info: scanner.ProjectInfo{Root: "/work/beta", Origin: "/work/beta", SpecCount: 5}},
	}
	out := ansi.Strip(renderTable(wrap(buildProjectRows(projects)), projectFields, 0, 80, 10))
	for _, unwanted := range []string{"ACTIVE", "ARCHIVED", "("} {
		if strings.Contains(out, unwanted) {
			t.Errorf("the picker has no groups, found %q:\n%s", unwanted, out)
		}
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Errorf("both projects must be listed:\n%s", out)
	}
}

// --- 4.x moving through it ---

func groupedListModel(t *testing.T) model {
	t.Helper()
	archived := map[string]string{"arch-old": "2026-01-01", "arch-new": "2026-06-01"}
	m := model{width: 100, height: 30,
		repoPaths:    []string{"/p"},
		displayNames: []string{"p"},
		detailTab:    tabChanges,
		focus:        focusDetail,
		fields:       append([]string(nil), defaultFields...),
		projects: scanner.ProjectMap{"/p": scanner.ProjectStatus{
			Info: groupedInfo([]string{"act-a", "act-b"}, archived)}},
	}
	m.recalcLayout()
	return m
}

func TestTheCursorStepsBetweenChangesAcrossTheBoundary(t *testing.T) {
	// The cursor indexes changes, never lines. Stepping from the last active
	// change lands on the first archived one, with no stop on the header.
	m := groupedListModel(t)
	names := plainNames(m.allRows())
	if len(names) != 4 {
		t.Fatalf("got %v, want four changes", names)
	}

	for i := 1; i < len(names); i++ {
		updated, _ := m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
		m = updated.(model)
		if m.changeCursor != i {
			t.Fatalf("after %d steps the cursor is %d, want %d", i, m.changeCursor, i)
		}
		if r, ok := m.selectedRow(); !ok || r.ci.Name != names[i] {
			t.Fatalf("step %d selected %v, want %q", i, r.ci.Name, names[i])
		}
	}
	// And stops at the end rather than walking onto a header.
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if got := updated.(model).changeCursor; got != len(names)-1 {
		t.Errorf("cursor ran past the last change to %d", got)
	}
}

func TestEveryDrawnRowIsExactlyOneOfTheThreeKinds(t *testing.T) {
	m := groupedListModel(t)
	groups := m.currentGroups()
	for i, l := range groupLines(groups, wrap(allRowsOf(groups))) {
		switch l.kind {
		case lineChange:
			if l.rowIndex < 0 || l.header != "" {
				t.Errorf("line %d is a change and something else: %+v", i, l)
			}
		case lineHeader:
			if l.header == "" || l.rowIndex >= 0 {
				t.Errorf("line %d is a header and something else: %+v", i, l)
			}
		case lineSpacer:
			if l.header != "" || l.rowIndex >= 0 {
				t.Errorf("line %d is a spacer and something else: %+v", i, l)
			}
		default:
			t.Errorf("line %d has no kind: %+v", i, l)
		}
	}
}

func TestTheCursorStopsAtBothEndsRatherThanOnAHeader(t *testing.T) {
	// `g` and `G` are document keys and have never been bound on the change
	// list, so there is no jump for a header to catch. What must hold is that
	// stepping to either end stops on a change.
	m := groupedListModel(t)
	last := len(m.allRows()) - 1

	for i := 0; i < last+3; i++ {
		updated, _ := m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
		m = updated.(model)
	}
	if m.changeCursor != last {
		t.Errorf("stepping down ended at %d, want the last change at %d", m.changeCursor, last)
	}
	if r, ok := m.selectedRow(); !ok || !r.archived {
		t.Error("the last change is in the archived group")
	}

	for i := 0; i < last+3; i++ {
		updated, _ := m.Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
		m = updated.(model)
	}
	if m.changeCursor != 0 {
		t.Errorf("stepping up ended at %d, want the first change", m.changeCursor)
	}
	if r, ok := m.selectedRow(); !ok || r.archived {
		t.Error("the first change is in the active group")
	}
}

func TestEnterOpensTheChangeUnderTheCursorInEitherGroup(t *testing.T) {
	for _, cursor := range []int{0, 3} {
		m := groupedListModel(t)
		m.changeCursor = cursor
		m.rememberSelection()
		want, _ := m.selectedRow()

		updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		um := updated.(model)
		if um.level != levelChange {
			t.Fatalf("cursor %d: did not descend", cursor)
		}
		if got, _ := um.selectedRow(); got.ci.Name != want.ci.Name {
			t.Errorf("cursor %d: opened %q, want %q", cursor, got.ci.Name, want.ci.Name)
		}
	}
}

// --- 6.3 archiving moves a change between groups ---

func TestArchivingMovesAChangeToTheTopOfTheArchivedGroup(t *testing.T) {
	before := groupedInfo([]string{"act-a", "moving"},
		map[string]string{"arch-old": "2026-01-01"})
	if got := plainNames(buildGroups(before)[0].rows); len(got) != 2 {
		t.Fatalf("got %v, want two active", got)
	}

	// What the next read sees once the change has been archived today.
	after := groupedInfo([]string{"act-a"},
		map[string]string{"moving": "2026-09-21", "arch-old": "2026-01-01"})
	groups := buildGroups(after)

	if got := plainNames(groups[0].rows); len(got) != 1 || got[0] != "act-a" {
		t.Errorf("active = %v, want the change gone from it", got)
	}
	archived := plainNames(groups[1].rows)
	if len(archived) != 2 || archived[0] != "moving" {
		t.Errorf("archived = %v, want the change at the top, newest first", archived)
	}
}

func TestGroupedChangeListFrameGeometry(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {92, 30}, {120, 50}} {
		m := groupedListModel(t)
		m.width, m.height = size.w, size.h
		m.recalcLayout()
		m.syncDocument()

		lines := strings.Split(m.renderFrame(), "\n")
		if len(lines) != size.h {
			t.Errorf("%dx%d: got %d rows, want %d", size.w, size.h, len(lines), size.h)
		}
		for _, l := range lines {
			if w := ansi.StringWidth(l); w > size.w {
				t.Errorf("%dx%d: a row is %d columns wide", size.w, size.h, w)
				break
			}
		}
	}
}

func TestNoMatchStillShowsBothGroupsAndSaysWhy(t *testing.T) {
	m := groupedListModel(t)
	m.searchInput.SetValue("zzzz")
	out := ansi.Strip(m.renderChangesTab(90, 12))

	for _, want := range []string{"ACTIVE (0)", "ARCHIVED (0)", `No changes match "zzzz"`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
}

// --- the blank line between groups ---

func TestOneBlankLineSitsBetweenTheGroups(t *testing.T) {
	m := groupedListModel(t)
	m.recalcLayout()
	lines := strings.Split(ansi.Strip(m.renderChangesTab(90, 20)), "\n")

	activeAt, archivedAt := -1, -1
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "ACTIVE") {
			activeAt = i
		}
		if strings.HasPrefix(strings.TrimSpace(l), "ARCHIVED") {
			archivedAt = i
		}
	}
	if activeAt < 0 || archivedAt < 0 {
		t.Fatalf("both headers must be on screen:\n%s", strings.Join(lines, "\n"))
	}
	if strings.TrimSpace(lines[archivedAt-1]) != "" {
		t.Errorf("the line above ARCHIVED is %q, want a blank one", lines[archivedAt-1])
	}
	if strings.TrimSpace(lines[activeAt-1]) == "" {
		t.Errorf("the first group follows the column header directly, got a blank line above ACTIVE")
	}
}

func TestTheBlankLineIsDrawnWhenAGroupIsEmpty(t *testing.T) {
	spacers := func(info scanner.ProjectInfo, rows []filtered[changeRow]) int {
		groups := buildGroups(info)
		if rows == nil {
			rows = wrap(allRowsOf(groups))
		}
		n := 0
		for _, l := range groupLines(groups, rows) {
			if l.kind == lineSpacer {
				n++
			}
		}
		return n
	}

	t.Run("no active changes", func(t *testing.T) {
		if got := spacers(groupedInfo(nil, map[string]string{"x": "2026-01-01"}), nil); got != 1 {
			t.Errorf("got %d spacers, want one", got)
		}
	})
	t.Run("no archived changes", func(t *testing.T) {
		if got := spacers(groupedInfo([]string{"a"}, nil), nil); got != 1 {
			t.Errorf("got %d spacers, want one", got)
		}
	})
	t.Run("neither", func(t *testing.T) {
		if got := spacers(groupedInfo(nil, nil), nil); got != 1 {
			t.Errorf("got %d spacers, want one", got)
		}
	})
	t.Run("emptied by a filter", func(t *testing.T) {
		info := groupedInfo([]string{"a"}, map[string]string{"x": "2026-01-01"})
		if got := spacers(info, nil); got != 1 {
			t.Errorf("got %d spacers, want one", got)
		}
		if got := spacers(info, []filtered[changeRow]{}); got != 1 {
			t.Errorf("with everything filtered away, got %d spacers, want one", got)
		}
	})
}

func TestTheBlankLineIsCountedByTheScrollArithmetic(t *testing.T) {
	// Stepping the cursor across the boundary in a pane too short to hold both
	// groups is where a line nobody budgeted for would lose the selection.
	archived := map[string]string{}
	for i := 0; i < 12; i++ {
		archived[fmt.Sprintf("arch-%02d", i)] = fmt.Sprintf("2026-01-%02d", i+1)
	}
	groups := buildGroups(groupedInfo([]string{"act-a", "act-b"}, archived))
	rows := wrap(allRowsOf(groups))

	const height = 8
	for cursor := range rows {
		out := renderGroupedTable(groups, rows, changeFieldDefs(defaultFields), cursor, 80, height)
		body := strings.Split(ansi.Strip(out), "\n")[1:]

		want := rows[cursor].row.ci.Name
		var seen bool
		for _, l := range body {
			if strings.Contains(l, want) {
				seen = true
			}
		}
		if !seen {
			t.Errorf("cursor %d (%s): the selected change is off screen:\n%s",
				cursor, want, strings.Join(body, "\n"))
		}
	}
}

func TestThePaneNeverOpensOnABlankLine(t *testing.T) {
	archived := map[string]string{}
	for i := 0; i < 12; i++ {
		archived[fmt.Sprintf("arch-%02d", i)] = fmt.Sprintf("2026-01-%02d", i+1)
	}
	groups := buildGroups(groupedInfo([]string{"act-a", "act-b"}, archived))
	rows := wrap(allRowsOf(groups))

	// Heights where the boundary can reach the top of the pane.
	for _, height := range []int{5, 6, 7, 8, 9} {
		for cursor := range rows {
			out := renderGroupedTable(groups, rows, changeFieldDefs(defaultFields), cursor, 80, height)
			body := strings.Split(ansi.Strip(out), "\n")[1:]
			if len(body) == 0 {
				continue
			}
			if strings.TrimSpace(body[0]) == "" {
				t.Fatalf("height %d, cursor %d: the pane opens on a blank line:\n%s",
					height, cursor, strings.Join(body, "\n"))
			}
		}
	}
}
