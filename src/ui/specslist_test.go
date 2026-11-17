package ui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func nestedSpecs(groups ...string) map[string]string {
	docs := map[string]string{}
	for _, g := range groups {
		for i := 0; i < 3; i++ {
			docs[fmt.Sprintf("%s/cap-%d", g, i)] = "# Spec"
		}
	}
	return docs
}

func TestSpecLinesGroupByFirstSegment(t *testing.T) {
	names := []string{
		"archive-from-ui",
		"content/editing",
		"content/pages",
		"media/images",
	}
	lines := specLines(names)

	want := []struct {
		kind  int
		label string
		index int
	}{
		{specLineCapability, "archive-from-ui", 0},
		{specLineSpacer, "", -1},
		{specLineHeader, "content (2)", -1},
		{specLineCapability, "editing", 1},
		{specLineCapability, "pages", 2},
		{specLineSpacer, "", -1},
		{specLineHeader, "media (1)", -1},
		{specLineCapability, "images", 3},
	}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %+v", len(lines), len(want), lines)
	}
	for i, w := range want {
		got := lines[i]
		if got.kind != w.kind || got.label != w.label || got.specIndex != w.index {
			t.Errorf("line %d = {%d %q %d}, want {%d %q %d}",
				i, got.kind, got.label, got.specIndex, w.kind, w.label, w.index)
		}
	}
}

func TestSpecLinesLeaveAFlatProjectAlone(t *testing.T) {
	lines := specLines([]string{"alpha", "beta", "gamma"})
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want one per spec", len(lines))
	}
	for i, l := range lines {
		if l.kind != specLineCapability {
			t.Errorf("line %d is chrome (%d) in a project with no groups", i, l.kind)
		}
	}
}

func TestSpecsTabDrawsGroupsBesideBareSpecs(t *testing.T) {
	docs := nestedSpecs("content", "media")
	docs["archive-from-ui"] = "# Archive"
	m := makeSpecsModel(t, docs)

	out := ansi.Strip(m.renderSpecsTab(m.width-2, 24))
	for _, want := range []string{"archive-from-ui", "content (3)", "media (3)", "cap-0", "cap-2"} {
		if !strings.Contains(out, want) {
			t.Errorf("the list does not show %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "content/cap-0") {
		t.Errorf("the header's segment is repeated in its rows:\n%s", out)
	}

	rows := strings.Split(out, "\n")
	blankBefore := func(label string) bool {
		for i, row := range rows {
			if i == 0 || !strings.Contains(row, label) {
				continue
			}
			left := strings.Split(rows[i-1], "│")
			return len(left) > 1 && strings.TrimSpace(left[1]) == ""
		}
		return false
	}
	for _, group := range []string{"content (3)", "media (3)"} {
		if !blankBefore(group) {
			t.Errorf("no blank line above %q:\n%s", group, out)
		}
	}
}

func TestSpecsListTruncatesAWideRow(t *testing.T) {
	long := "acceptance/reference-instance-of-a-capability"
	m := makeSpecsModel(t, map[string]string{long: "# Spec"})

	out := m.renderSpecsTab(m.width-2, 24)
	for _, row := range strings.Split(ansi.Strip(out), "\n") {
		if got := ansi.StringWidth(row); got > m.width-2 {
			t.Errorf("a row is %d cells wide in a %d cell panel: %q", got, m.width-2, row)
		}
	}
	if strings.Contains(ansi.Strip(out), "of-a-capability") {
		t.Errorf("the row was not truncated:\n%s", ansi.Strip(out))
	}

	narrow := makeSpecsModel(t, map[string]string{long: "# Spec"})
	narrow.width = 40
	narrow.recalcLayout()
	narrow.syncDocument()
	narrowed := narrow.renderSpecsTab(narrow.width-2, 20)
	for _, row := range strings.Split(ansi.Strip(narrowed), "\n") {
		if got := ansi.StringWidth(row); got > narrow.width-2 {
			t.Errorf("a row is %d cells wide in a %d cell panel: %q",
				got, narrow.width-2, row)
		}
	}
	if strings.Contains(ansi.Strip(narrowed), "reference-instance") {
		t.Errorf("the row widened a narrow pane:\n%s", ansi.Strip(narrowed))
	}
}

func TestSpecPageCountsTheHeadersItDraws(t *testing.T) {
	m := makeSpecsModel(t, nestedSpecs("alpha", "beta", "gamma", "delta"))
	rows := m.mainPanelHeight() - 5 - boxRows
	lines := specLines(m.currentSpecNames())

	page := m.listPage()
	if want := ownersOfSpecs(lines).fit(0, rows); page != want {
		t.Errorf("a page moves %d specs, want the %d that fit in %d rows", page, want, rows)
	}
	if page >= rows {
		t.Errorf("a page moves %d specs in %d rows, so the headers cost nothing", page, rows)
	}
}

func TestTheSpecCursorNeverLandsOnAHeader(t *testing.T) {
	m := makeSpecsModel(t, nestedSpecs("content", "media"))
	lines := specLines(m.currentSpecNames())
	listOuter, _ := specsSplit(m.panelContentWidth())
	width := listOuter - boxChrome

	for i, name := range m.currentSpecNames() {
		out := renderSpecList(lines, i, width, 40, false)
		selected := ""
		for _, row := range strings.Split(out, "\n") {
			if strings.Contains(row, "\x1b[30;42m") {
				selected = strings.TrimSpace(ansi.Strip(row))
			}
		}
		if want := name[strings.IndexByte(name, '/')+1:]; selected != want {
			t.Errorf("cursor %d highlights %q, want the capability %q", i, selected, want)
		}
	}
}

func TestAGroupFurtherDownScrollsIntoView(t *testing.T) {
	m := makeSpecsModel(t, nestedSpecs("alpha", "beta", "gamma", "delta"))
	lines := specLines(m.currentSpecNames())
	listOuter, _ := specsSplit(m.panelContentWidth())
	width := listOuter - boxChrome
	last := len(m.currentSpecNames()) - 1

	var selected string
	for _, row := range strings.Split(renderSpecList(lines, last, width, 5, false), "\n") {
		if strings.Contains(row, "\x1b[30;42m") {
			selected = strings.TrimSpace(ansi.Strip(row))
		}
	}
	if selected != "cap-2" {
		t.Errorf("the selected capability is %q, want cap-2 on screen", selected)
	}

	out := ansi.Strip(renderSpecList(lines, last, width, 5, false))
	for _, want := range []string{"gamma (3)", "cap-0", "cap-2"} {
		if !strings.Contains(out, want) {
			t.Errorf("jumping to the last capability does not show %q:\n%s", want, out)
		}
	}
}
