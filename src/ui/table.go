package ui

import (
	"fmt"
	"sort"
	"strings"
)

// tableRow is what both lists have in common: something a cursor can hold onto
// across a re-filter, a short name to match fuzzily, and named bodies of text to
// search literally.
//
// The change list was the first consumer and the project picker is the second,
// which is what justifies the abstraction. Generalising from one example would
// have been guesswork.
type tableRow interface {
	// key identifies the row across a re-filter or a rescan.
	key() string
	// searchName is the short target for fuzzy and literal name matching.
	searchName() string
	// searchBodies maps a display label to a body of text that a ':' query
	// searches. The label is what the match hint shows.
	searchBodies() map[string]string
}

// filtered pairs a row with the reason it survived a filter. Returning the
// reason is cleaner than writing it back onto the row: a row that matched on
// its name must not keep hints from an earlier body search.
type filtered[T tableRow] struct {
	row     T
	matched []string
}

// wrap lifts plain rows into unmatched entries.
func wrap[T tableRow](rows []T) []filtered[T] {
	out := make([]filtered[T], len(rows))
	for i, r := range rows {
		out[i] = filtered[T]{row: r}
	}
	return out
}

// indexOfKey finds the row carrying the given key, or -1.
func indexOfKey[T tableRow](rows []filtered[T], key string) int {
	for i, r := range rows {
		if r.row.key() == key {
			return i
		}
	}
	return -1
}

// fieldDef describes one column.
//
// width 0 means the column is flexible: it absorbs whatever space the
// fixed-width columns leave. Exactly one field per table is flexible.
type fieldDef[T tableRow] struct {
	id     string
	header string
	width  int
	value  func(T) string
}

// minFlexWidth is the narrowest the flexible column may become before columns
// are dropped instead of squeezed further.
const minFlexWidth = 12

// columnGap is how many blank columns separate one column from the next. One
// was not enough: a value that filled its column sat a single space from the
// value beside it, which read as one run-on field rather than two.
const columnGap = 2

// columnSep is what cells are joined with. It must be columnGap columns wide;
// layoutFields charges for it by the same constant.
var columnSep = strings.Repeat(" ", columnGap)

// layoutFields decides which columns fit and how wide each is. Fixed columns
// are served first and the flexible column takes the remainder. When that
// remainder would fall below minFlexWidth, columns are dropped from the right.
func layoutFields[T tableRow](defs []fieldDef[T], width int) (kept []fieldDef[T], widths []int) {
	kept = append([]fieldDef[T](nil), defs...)

	for {
		widths = widths[:0]
		fixed := 0
		flexCount := 0
		for _, d := range kept {
			if d.width == 0 {
				flexCount++
			} else {
				fixed += d.width
			}
		}
		// columnGap columns per gap. This and the separator the cells are
		// joined with below are the same fact written twice: if they disagree
		// the row comes out one or two columns wrong per gap, and nothing
		// downstream notices because every row is wrong by the same amount.
		gaps := 0
		if len(kept) > 1 {
			gaps = columnGap * (len(kept) - 1)
		}
		remaining := width - fixed - gaps

		if flexCount == 0 {
			if remaining >= 0 || len(kept) <= 1 {
				for _, d := range kept {
					widths = append(widths, d.width)
				}
				return kept, widths
			}
			kept = kept[:len(kept)-1]
			continue
		}

		flexEach := remaining / flexCount
		if flexEach >= minFlexWidth || len(kept) <= 1 {
			if flexEach < 1 {
				flexEach = 1
			}
			for _, d := range kept {
				if d.width == 0 {
					widths = append(widths, flexEach)
				} else {
					widths = append(widths, d.width)
				}
			}
			return kept, widths
		}
		kept = kept[:len(kept)-1]
	}
}

// fitCell pads or truncates a value to exactly w columns, marking truncation
// with an ellipsis so a clipped value is visibly clipped.
func fitCell(s string, w int) string {
	if w <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) > w {
		if w == 1 {
			return "…"
		}
		return string(runes[:w-1]) + "…"
	}
	return s + strings.Repeat(" ", w-len(runes))
}

// renderTable draws a header row and one row per entry, across the full width.
func renderTable[T tableRow](rows []filtered[T], defs []fieldDef[T], cursor, width, height int) string {
	// A body search needs room to say which files matched, otherwise the
	// flexible column absorbs the whole width and the hint never shows. The
	// column only appears while such a search is active.
	hintWidth := 0
	for _, r := range rows {
		if len(r.matched) > 0 {
			hintWidth = min(28, width/3)
			break
		}
	}
	// The hint is a column like any other, so it is separated like one. It used
	// to be appended flush, and looked right only because the last default
	// column held values shorter than its width; the archive date fills its
	// column exactly.
	if hintWidth > 0 {
		hintWidth += columnGap
	}
	tableWidth := width - hintWidth
	if tableWidth < 1 {
		tableWidth = width
		hintWidth = 0
	}

	kept, widths := layoutFields(defs, tableWidth)

	var b strings.Builder

	headerCells := make([]string, len(kept))
	for i, d := range kept {
		headerCells[i] = fitCell(d.header, widths[i])
	}
	header := fitCell(strings.Join(headerCells, columnSep), tableWidth)
	if hintWidth > 0 {
		header += columnSep + fitCell("matched", hintWidth-columnGap)
	}
	b.WriteString(dimStyle.Render(header))

	bodyHeight := height - 1
	if bodyHeight < 1 {
		return b.String()
	}

	offset := 0
	if cursor >= bodyHeight {
		offset = cursor - bodyHeight + 1
	}
	end := offset + bodyHeight
	if end > len(rows) {
		end = len(rows)
	}

	for i := offset; i < end; i++ {
		b.WriteString("\n")

		cells := make([]string, len(kept))
		for j, d := range kept {
			cells[j] = fitCell(d.value(rows[i].row), widths[j])
		}
		base := fitCell(strings.Join(cells, columnSep), tableWidth)

		hint := ""
		if hintWidth > 0 {
			hint = columnSep + fitCell(strings.Join(rows[i].matched, ", "), hintWidth-columnGap)
		}

		if i == cursor {
			b.WriteString(selectedStyle.Width(width).Render(fitCell(base+hint, width)))
		} else {
			b.WriteString(normalStyle.Render(base))
			if hint != "" {
				b.WriteString(dimStyle.Render(hint))
			}
		}
	}

	return b.String()
}

// searchLegend names the three matchers the query grammar accepts.
//
// Two of them are prefix sigils and neither can be guessed at. One wording
// serves both surfaces: the contents sigil reaches artifact and spec text in a
// change and file paths and contents in a project, and `inside` is the shortest
// word true of both.
const searchLegend = "fuzzy name  'exact  :inside"

// renderSearchPrompt draws the filter line under a table. It stays visible for
// as long as a query is applied, so a narrowed list always shows why.
//
// While the prompt is focused and empty it also names the matchers, which is
// the moment between asking to search and knowing what to type. The first
// keystroke replaces the naming with the query, so it costs a fast typist
// nothing and is dropped entirely when the line is too narrow to hold it.
func renderSearchPrompt(input string, focused bool, shown, total int, width int) string {
	var b strings.Builder
	b.WriteString(navBarKeyStyle.Render("/"))
	b.WriteString(navBarStyle.Render(input))
	if focused {
		b.WriteString(navBarStyle.Render("_"))
	}

	count := fmt.Sprintf("   %d of %d shown", shown, total)
	if focused && input == "" {
		legend := "  " + searchLegend
		// 1 for the sigil, 1 for the cursor, and the count, which is never
		// dropped: it is what the prompt is for.
		if width >= 2+len(legend)+len(count) {
			b.WriteString(dimStyle.Render(legend))
		}
	}
	b.WriteString(dimStyle.Render(count))
	return b.String()
}

// noMatchMessage explains an empty result, and suggests the contents search
// when the query that failed was a name search.
//
// A name search that found nothing is the moment a person is already looking
// for another way, which is a better time to be told about one than any other.
// A contents query that failed suggests nothing: there is no further matcher.
func noMatchMessage(noun, raw string) string {
	msg := fmt.Sprintf("No %s match %q", noun, raw)
	q := parseQuery(raw)
	if q.kind == matchBody || q.term == "" {
		return msg
	}
	return msg + "\n" + dimStyle.Render("try :"+q.term+" to search inside them")
}

// sortedLabels returns map keys in a stable order, so match hints do not
// reshuffle between renders.
func sortedLabels(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
