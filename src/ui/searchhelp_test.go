package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// --- 1.2 to 1.5 the legend ---

func TestTheLegendNamesAllThreeMatchers(t *testing.T) {
	got := ansi.Strip(renderSearchPrompt("", true, 38, 38, 90))
	for _, want := range []string{"fuzzy name", "'exact", ":inside"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q:\n%s", want, got)
		}
	}
	if !strings.Contains(got, "38 of 38 shown") {
		t.Errorf("the count must survive beside it:\n%s", got)
	}
}

func TestTheLegendGoesOnTheFirstKeystroke(t *testing.T) {
	got := ansi.Strip(renderSearchPrompt("e", true, 3, 38, 90))
	if strings.Contains(got, ":inside") {
		t.Errorf("a query and its explanation must not share the line:\n%s", got)
	}
	if !strings.Contains(got, "3 of 38 shown") {
		t.Errorf("the count survives:\n%s", got)
	}
}

func TestTheLegendIsAbsentWhileUnfocused(t *testing.T) {
	got := ansi.Strip(renderSearchPrompt("", false, 38, 38, 90))
	if strings.Contains(got, ":inside") {
		t.Errorf("nothing is about to be typed:\n%s", got)
	}
}

func TestTheLegendIsDroppedRatherThanWrapped(t *testing.T) {
	// The change list's box is 52 columns at the 60-column minimum, and the
	// picker's is narrower still.
	for _, width := range []int{20, 30, 40} {
		got := ansi.Strip(renderSearchPrompt("", true, 38, 38, width))
		if strings.Contains(got, ":inside") {
			t.Errorf("width %d: the legend must be dropped, got:\n%s", width, got)
		}
		if !strings.Contains(got, "38 of 38 shown") {
			t.Errorf("width %d: the count is never dropped:\n%s", width, got)
		}
		if strings.Contains(got, "\n") {
			t.Errorf("width %d: the prompt must not wrap", width)
		}
	}
}

func TestTheLegendFitsAtTheMinimumWidth(t *testing.T) {
	m := groupedListModel(t)
	m.width, m.height = 60, 20
	m.recalcLayout()
	m.searchFocused = true

	line := ansi.Strip(m.renderChangesTab(m.contentBoxWidth(), 10))
	if len(strings.Split(line, "\n")) < 2 {
		t.Fatalf("expected a prompt line:\n%s", line)
	}
	for _, l := range strings.Split(line, "\n") {
		if len([]rune(l)) > m.contentBoxWidth() {
			t.Errorf("a line is %d columns in a %d box: %q", len([]rune(l)), m.contentBoxWidth(), l)
		}
	}
}

// --- 2.1 to 2.5 the failed search ---

func TestAFailedNameSearchSuggestsTheContentsSearch(t *testing.T) {
	for _, q := range []string{"inotify", "'inotify"} {
		got := ansi.Strip(noMatchMessage("changes", q))
		if !strings.Contains(got, "No changes match") {
			t.Errorf("%q: the message must still say so:\n%s", q, got)
		}
		if !strings.Contains(got, ":inotify") {
			t.Errorf("%q: the suggestion must name the term, not the sigil:\n%s", q, got)
		}
	}
}

func TestAFailedContentsSearchSuggestsNothing(t *testing.T) {
	got := ansi.Strip(noMatchMessage("changes", ":inotify"))
	if !strings.Contains(got, "No changes match") {
		t.Errorf("the message must still say so:\n%s", got)
	}
	if strings.Contains(got, "try ") {
		t.Errorf("there is no further matcher to try:\n%s", got)
	}
}

func TestTheSuggestionReachesBothSurfaces(t *testing.T) {
	if got := ansi.Strip(noMatchMessage("projects", "zipp")); !strings.Contains(got, ":zipp") {
		t.Errorf("the picker gets it too:\n%s", got)
	}

	m := groupedListModel(t)
	m.recalcLayout()
	m.searchInput.SetValue("zzzz")
	got := ansi.Strip(m.renderChangesTab(90, 12))
	if !strings.Contains(got, ":zzzz") {
		t.Errorf("the change list shows it:\n%s", got)
	}
}

// --- 3.x the match hint gets its gap ---

func TestTheMatchHintIsSeparatedFromAFullLastColumn(t *testing.T) {
	// The archive date fills its ten-wide column exactly, which is what made
	// the hint land flush against it.
	groups := buildGroups(groupedInfo(nil, map[string]string{"box-the-tab-content": "2026-09-18"}))
	rows := []filtered[changeRow]{{row: allRowsOf(groups)[0], matched: []string{"design"}}}

	out := ansi.Strip(renderGroupedTable(groups, rows, changeFieldDefs(defaultFields), 0, 100, 10))
	if strings.Contains(out, "2026-09-18design") {
		t.Errorf("the hint is flush against the date:\n%s", out)
	}
	if !strings.Contains(out, "2026-09-18"+columnSep+"design") {
		t.Errorf("want the column gap between them:\n%s", out)
	}
}

func TestTheMatchHintHeaderIsSeparatedToo(t *testing.T) {
	groups := buildGroups(groupedInfo(nil, map[string]string{"a": "2026-09-18"}))
	rows := []filtered[changeRow]{{row: allRowsOf(groups)[0], matched: []string{"design"}}}
	header := strings.Split(ansi.Strip(renderGroupedTable(groups, rows,
		changeFieldDefs(defaultFields), 0, 100, 10)), "\n")[0]

	if !strings.Contains(header, "archived"+columnSep) {
		t.Errorf("the header's gap must not depend on the last header's length:\n%q", header)
	}
	if !strings.Contains(header, "matched") {
		t.Errorf("the hint column is named:\n%q", header)
	}
}

func TestThePickerMatchHintIsSeparated(t *testing.T) {
	projects := scanner.ProjectMap{
		"/work/alpha": {Info: scanner.ProjectInfo{Root: "/work/alpha", Origin: "/work/alpha",
			SpecCount: 3, TasksTotal: 311, TasksDone: 117}},
	}
	rows := []filtered[projectRow]{{row: buildProjectRows(projects)[0], matched: []string{"a-spec/spec"}}}
	out := ansi.Strip(renderTable(rows, projectFields, 0, 100, 10))

	if !strings.Contains(out, "117/311"+columnSep) && !strings.Contains(out, " "+columnSep+"a-spec") {
		t.Errorf("both tables share the measurement:\n%s", out)
	}
	if strings.Contains(out, "117/311a-spec") {
		t.Errorf("the hint is flush:\n%s", out)
	}
}

// --- the legend reaches the picker ---

func TestThePickerPromptNamesTheMatchers(t *testing.T) {
	m := pickerPagingModel(t, 6)
	m.pickerFocused = true
	m.pickerInput.SetValue("")
	m.width, m.height = 110, 30
	m.recalcLayout()

	got := ansi.Strip(m.renderPicker())
	for _, want := range []string{"fuzzy name", "'exact", ":inside"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q from the picker prompt:\n%s", want, got)
		}
	}
}

func TestThePickerPromptSuggestsAfterAFailedSearch(t *testing.T) {
	m := pickerPagingModel(t, 6)
	m.width, m.height = 110, 30
	m.recalcLayout()
	m.pickerInput.SetValue("zzzz")
	m.pickerSync()

	got := ansi.Strip(m.renderPicker())
	if !strings.Contains(got, "No projects match") {
		t.Errorf("the message must say so:\n%s", got)
	}
	if !strings.Contains(got, ":zzzz") {
		t.Errorf("and suggest the contents search:\n%s", got)
	}
}

func TestTheLegendAppearsOnSlashAndGoesOnTheNextKey(t *testing.T) {
	m := newModel(&scanner.Config{}, true, "0.0.0")
	m.width, m.height = 100, 24
	m.repoPaths = []string{"/p"}
	m.displayNames = []string{"p"}
	m.detailTab = tabChanges
	m.focus = focusDetail
	m.fields = append([]string(nil), defaultFields...)
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{
		Info: groupedInfo([]string{"alpha", "beta"}, nil)}}
	m.recalcLayout()

	opened := press(m, tea.KeyPressMsg{Code: '/', Text: "/"})
	if !strings.Contains(ansi.Strip(opened.renderFrame()), ":inside") {
		t.Error("the legend appears on /")
	}
	typed := press(opened, tea.KeyPressMsg{Code: 'a', Text: "a"})
	if strings.Contains(ansi.Strip(typed.renderFrame()), ":inside") {
		t.Error("and goes on the next key")
	}
}
