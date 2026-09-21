package ui

import (
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
)

// The inline renderer draws one line of markdown, and every markdown surface
// goes through it: the specs tab, the detail card, a change's spec deltas, and
// the proposal and design documents.
//
// It is one pass rather than a mark at a time. Styles cannot be applied around
// one another, because the end of an inner style ends the outer one as well:
//
//	mdBoldStyle.Render("a " + mdCodeStyle.Render("b") + " c")
//	  -> "\x1b[1ma \x1b[36mb\x1b[m c\x1b[m"
//	                            ^^^^^^ the bold is over; " c" draws plain
//
// So every segment is emitted with exactly one style, and a segment is never
// handed back for another pass to look at.

// A spec is written in a small vocabulary that carries its meaning. The
// vocabulary and the roles are openspec.nvim's, whose `spec-highlighting`
// capability gives the same words capture groups for an editor: `@keyword` for
// the four that bind, and `@conditional`, `@property` and `@operator` for the
// three that shape a scenario. Reading the same spec in either tool should not
// mean learning two sets of colours.
//
// `GIVEN` is the one word that capability does not name. It opens a condition,
// so it shares `WHEN`'s role rather than earning a fourth.
var clauseKeywordDefs = []struct {
	word  string
	style func() lipgloss.Style
}{
	{"GIVEN", func() lipgloss.Style { return kwConditionStyle }},
	{"WHEN", func() lipgloss.Style { return kwConditionStyle }},
	{"THEN", func() lipgloss.Style { return kwAssertionStyle }},
	{"AND", func() lipgloss.Style { return kwContinuationStyle }},
}

// reBinding matches the four words that bind. Two-word forms lead, because Go's
// alternation prefers the first that matches at a position: `SHALL` first would
// draw a keyword beside a plain `NOT`.
//
// `SHOULD`, `MAY`, `REQUIRED`, `RECOMMENDED` and `OPTIONAL` are deliberately
// absent. OpenSpec's own schema instruction tells authors to use SHALL and MUST
// and to avoid the others, and drawing a word the format discourages would read
// as endorsement of it.
var reBinding = regexp.MustCompile(`\b(SHALL NOT|MUST NOT|SHALL|MUST)\b`)

// renderInlineMarkdown draws one line: code spans, bold, italics and keywords.
func renderInlineMarkdown(line string) string {
	var b strings.Builder

	rest := line
	if styled, tail, ok := leadingClauseKeyword(line); ok {
		b.WriteString(styled)
		rest = tail
	}

	var prose strings.Builder
	flush := func() {
		if prose.Len() > 0 {
			b.WriteString(withBinding(prose.String(), lipgloss.Style{}))
			prose.Reset()
		}
	}

	for i := 0; i < len(rest); {
		// A code span binds tighter than everything else, and its contents are
		// never looked at again. That is what markdown says, and it is also
		// what keeps a keyword named rather than used from being drawn as one.
		if rest[i] == '`' {
			if end := strings.IndexByte(rest[i+1:], '`'); end >= 0 {
				flush()
				b.WriteString(mdCodeStyle.Render(rest[i+1 : i+1+end]))
				i += end + 2
				continue
			}
		}
		if strings.HasPrefix(rest[i:], "**") {
			if end := strings.Index(rest[i+2:], "**"); end >= 0 {
				flush()
				b.WriteString(withBinding(rest[i+2:i+2+end], mdBoldStyle))
				i += end + 4
				continue
			}
		}
		if rest[i] == '_' {
			if end := strings.IndexByte(rest[i+1:], '_'); end >= 0 {
				flush()
				b.WriteString(withBinding(rest[i+1:i+1+end], mdItalicStyle))
				i += end + 2
				continue
			}
		}
		// An unclosed mark is a character like any other. The rest of the line
		// is prose and deserves to survive.
		prose.WriteByte(rest[i])
		i++
	}
	flush()

	return b.String()
}

// withBinding draws text in base, with the words that bind drawn as keywords.
//
// A keyword keeps its own style rather than inheriting the run it sits in, so a
// `SHALL` inside a bold span reads as a keyword and not as bold. Each emitted
// run still carries exactly one style.
func withBinding(text string, base lipgloss.Style) string {
	locs := reBinding.FindAllStringIndex(text, -1)
	if len(locs) == 0 {
		return render(base, text)
	}

	var b strings.Builder
	at := 0
	for _, l := range locs {
		if l[0] > at {
			b.WriteString(render(base, text[at:l[0]]))
		}
		b.WriteString(kwBindingStyle.Render(text[l[0]:l[1]]))
		at = l[1]
	}
	if at < len(text) {
		b.WriteString(render(base, text[at:]))
	}
	return b.String()
}

// render draws text in a style, leaving it alone when the style is the empty
// one, so that a line with nothing to draw carries no escape sequences at all.
func render(style lipgloss.Style, text string) string {
	if text == "" {
		return ""
	}
	if style.String() == "" {
		return text
	}
	return style.Render(text)
}

// leadingClauseKeyword draws the keyword that opens a line, if one does.
//
// Only at the start of a line, with or without a list marker and with or
// without bold marks around it. That is openspec.nvim's rule, and it is what
// keeps the forty mid-sentence uppercase `AND`s in the local corpus out of the
// highlighting without a case of their own: in the middle of a sentence these
// are ordinary words.
func leadingClauseKeyword(line string) (styled, rest string, ok bool) {
	i := 0
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	indent := line[:i]
	rest = line[i:]

	marker := ""
	if len(rest) > 1 && (rest[0] == '-' || rest[0] == '*' || rest[0] == '+') && rest[1] == ' ' {
		marker, rest = rest[:2], rest[2:]
	}
	bold := false
	if strings.HasPrefix(rest, "**") {
		bold, rest = true, rest[2:]
	}

	for _, kw := range clauseKeywordDefs {
		if !strings.HasPrefix(rest, kw.word) {
			continue
		}
		after := rest[len(kw.word):]
		if bold {
			if !strings.HasPrefix(after, "**") {
				continue
			}
			after = after[2:]
		}
		// A word, not a prefix of one: SHALLOW and ANDROID are not keywords.
		if after != "" && after[0] != ' ' && after[0] != ':' && after[0] != ',' {
			continue
		}
		return indent + marker + kw.style().Render(kw.word), after, true
	}
	return "", line, false
}

// opensWithClauseKeyword reports whether a line begins with one, without
// drawing it.
func opensWithClauseKeyword(line string) bool {
	_, _, ok := leadingClauseKeyword(line)
	return ok
}

// clauseStyleFor is the style a clause keyword carries on the detail card,
// where the keyword sits on a row of its own rather than in a line of prose.
func clauseStyleFor(keyword string) lipgloss.Style {
	for _, kw := range clauseKeywordDefs {
		if kw.word == keyword {
			return kw.style()
		}
	}
	return kwConditionStyle
}
