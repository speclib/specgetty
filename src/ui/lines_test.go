package ui

import "testing"

// --- 1.1 to 1.3 the shared arithmetic ---

// oneLineEach is the change list's shape: every item occupies one drawn line.
func oneLineEach(n int) itemLines { return ownedLines(n) }

// withChrome is the properties list's shape: a header before each group.
func withChrome() itemLines {
	return itemLines{lineOwner, 0, 1, lineOwner, lineOwner, 2, 3}
}

// tall is the spec outline's shape: item 1 wraps onto three lines.
func tall() itemLines { return itemLines{0, 1, 1, 1, 2} }

func TestSpanFindsWhereAnItemIsDrawn(t *testing.T) {
	for _, c := range []struct {
		name              string
		lines             itemLines
		item, first, last int
	}{
		{"one line each", oneLineEach(4), 2, 2, 2},
		{"after chrome", withChrome(), 0, 1, 1},
		{"after two chrome lines", withChrome(), 2, 5, 5},
		{"an item spanning three lines", tall(), 1, 1, 3},
		{"an item that is not drawn", oneLineEach(2), 9, -1, -1},
	} {
		first, last := c.lines.span(c.item)
		if first != c.first || last != c.last {
			t.Errorf("%s: got %d..%d, want %d..%d", c.name, first, last, c.first, c.last)
		}
	}
}

// TestFitIsTheRowCountForOneLineItems is the property the change list's page
// key relies on, so it is asserted rather than assumed.
func TestFitIsTheRowCountForOneLineItems(t *testing.T) {
	l := oneLineEach(50)
	for _, rows := range []int{1, 5, 13, 40} {
		if got := l.fit(0, rows); got != rows {
			t.Errorf("%d rows: got %d items, want %d", rows, got, rows)
		}
	}
	// Past the end it stops at what is there rather than running on.
	if got := l.fit(48, 13); got != 2 {
		t.Errorf("near the end: got %d, want the 2 items that remain", got)
	}
}

func TestFitCountsTallItemsAsTheLinesTheyTake(t *testing.T) {
	l := tall() // items of 1, 3 and 1 lines
	if got := l.fit(0, 4); got != 2 {
		t.Errorf("got %d items in 4 lines, want 2: one line plus three", got)
	}
	if got := l.fit(0, 5); got != 3 {
		t.Errorf("got %d items in 5 lines, want all 3", got)
	}
	if got := l.fit(1, 2); got != 1 {
		t.Errorf("got %d, want 1: an item taller than the pane still moves the cursor", got)
	}
}

func TestFitAlwaysMovesAtLeastOne(t *testing.T) {
	// Otherwise a page key would be inert wherever the next item is taller
	// than the pane, which is exactly where the reader needs it most.
	if got := tall().fit(1, 1); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestOffsetForBringsTheWholeItemIntoView(t *testing.T) {
	for _, c := range []struct {
		name         string
		lines        itemLines
		item, height int
		want         int
	}{
		{"already visible", oneLineEach(10), 2, 5, 0},
		{"below the fold", oneLineEach(10), 7, 5, 3},
		// Item 1 spans lines 1..3. In a 2-line pane the bottom-aligned offset
		// would be 2 and would cut off its first line, so it clamps to 1.
		{"an item taller than the pane shows its top", tall(), 1, 2, 1},
		{"a tall item that does fit", tall(), 1, 4, 0},
		{"an item much taller than the pane", itemLines{0, 1, 1, 1, 1, 1}, 1, 3, 1},
		{"an item that is not drawn", oneLineEach(3), 9, 5, 0},
	} {
		if got := c.lines.offsetFor(c.item, c.height); got != c.want {
			t.Errorf("%s: got offset %d, want %d", c.name, got, c.want)
		}
	}
}

// TestOffsetForReducesToTheOldExpression is task 1.2: for one-line items it is
// max(0, last-height+1), which is the arithmetic the change list shipped with.
func TestOffsetForReducesToTheOldExpression(t *testing.T) {
	l := oneLineEach(40)
	for item := 0; item < 40; item++ {
		for _, height := range []int{1, 7, 13} {
			want := 0
			if item >= height {
				want = item - height + 1
			}
			if got := l.offsetFor(item, height); got != want {
				t.Errorf("item %d in %d rows: got %d, want %d", item, height, got, want)
			}
		}
	}
}

func TestOffsetForClampsDownToTheItemsFirstLine(t *testing.T) {
	// An item taller than the pane is read from its top. Scrolled to its last
	// row it would read as a fragment with no beginning.
	l := itemLines{0, 1, 1, 1, 1, 2}
	if got := l.offsetFor(1, 2); got != 1 {
		t.Errorf("got %d, want the item's first line", got)
	}
}
