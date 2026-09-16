package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// Note for anyone wondering where the colour-forcing helper went: lipgloss v2
// has no global colour profile. Style.Render always emits, and downgrading
// happens at the output writer. Under v1 these tests needed a helper to stop
// being silently vacuous, because Render checked for a terminal and stripped
// styling when there was none. That trap is gone.

func lineWidths(s string) []int {
	lines := strings.Split(s, "\n")
	out := make([]int, len(lines))
	for i, l := range lines {
		out[i] = lipgloss.Width(l)
	}
	return out
}

func TestPickerLinesAllHaveTheSameWidth(t *testing.T) {
	// Note: this alone does NOT catch an overflowing row. lipgloss pads the
	// wrapped remainder, so every line still measures the same width. The
	// detector for that is TestPickerDoesNotOverflowItsBox, which counts lines.
	// This test guards ragged output from other causes.

	for _, cursor := range []int{0, 1} {
		m := makePickerModel()
		m.pickerOpen = true
		m.width, m.height = 100, 30
		m.pickerCursor = cursor

		widths := lineWidths(m.renderPicker())
		for i, w := range widths {
			if w != widths[0] {
				t.Errorf("cursor %d: line %d is %d columns, line 0 is %d; the highlighted row must line up with the rest",
					cursor, i, w, widths[0])
			}
		}
	}
}

func TestPickerDoesNotOverflowItsBox(t *testing.T) {

	m := makePickerModel()
	m.pickerOpen = true
	m.width, m.height = 100, 30

	boxWidth, boxHeight := m.pickerBox()
	lines := strings.Split(m.renderPicker(), "\n")

	if got := lipgloss.Width(lines[0]); got != boxWidth {
		t.Errorf("box rendered %d columns wide, want %d", got, boxWidth)
	}
	// An overflowing row is wrapped rather than rejected, so extra lines are
	// the symptom. lipgloss Style.Width covers content plus padding and
	// excludes the border, which is the trap this guards.
	if len(lines) != boxHeight {
		t.Errorf("box rendered %d lines, want %d: extra lines mean a row overflowed and wrapped",
			len(lines), boxHeight)
	}
}

func TestPickerWidthsHoldWithAMatchHintColumn(t *testing.T) {

	m := makePickerModel()
	m.pickerOpen = true
	m.width, m.height = 100, 30
	m.pickerInput.SetValue(":flake")
	m.pickerSync()

	widths := lineWidths(m.renderPicker())
	for i, w := range widths {
		if w != widths[0] {
			t.Errorf("line %d is %d columns, line 0 is %d, with a hint column present", i, w, widths[0])
		}
	}
}

func TestPickerWidthsHoldAtAwkwardTerminalSizes(t *testing.T) {

	for _, size := range []struct{ w, h int }{
		{60, 20}, {61, 21}, {80, 24}, {100, 30}, {201, 51},
	} {
		m := makePickerModel()
		m.pickerOpen = true
		m.width, m.height = size.w, size.h
		m.pickerCursor = 1

		widths := lineWidths(m.renderPicker())
		for i, w := range widths {
			if w != widths[0] {
				t.Errorf("at %dx%d: line %d is %d columns, line 0 is %d",
					size.w, size.h, i, w, widths[0])
				break
			}
		}
	}
}

func TestPanelTopBorderMatchesTheBoxWidth(t *testing.T) {

	// The top border is assembled by hand rather than by lipgloss, so it can
	// drift from the box it caps. It did: it was one column short.
	for _, width := range []int{40, 58, 80, 98, 120} {
		m := makeListModel()
		m.width, m.height = width+2, 30

		widths := lineWidths(m.renderPanel(viewDetail, width, 5, "content"))
		for i, w := range widths {
			if w != widths[0] {
				t.Errorf("width %d: line %d is %d columns, top border is %d",
					width, i, w, widths[0])
				break
			}
		}
	}
}

func TestPanelTopBorderHandlesALongTitle(t *testing.T) {

	// A title longer than the box must not make the border negative or ragged.
	m := makeListModel()
	m.width, m.height = 30, 20
	m.displayNames = []string{strings.Repeat("very-long-project-name", 3)}

	widths := lineWidths(m.renderPanel(viewDetail, 28, 4, "content"))
	for i, w := range widths {
		if w < widths[len(widths)-1] {
			t.Errorf("line %d is %d columns, narrower than the bottom border at %d",
				i, w, widths[len(widths)-1])
		}
	}
}
