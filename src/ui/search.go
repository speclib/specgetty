package ui

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sahilm/fuzzy"
)

// Matcher kinds, selected by a prefix sigil on the query.
const (
	matchFuzzyName = iota
	matchLiteralName
	matchBody
)

// query is a parsed search term. Names are matched by fuzzy subsequence because
// they are short and half-remembered; bodies are matched literally because a
// fuzzy subsequence matches nearly any document of real length.
type query struct {
	kind          int
	term          string
	caseSensitive bool
}

// parseQuery reads the sigil grammar:
//
//	export    fuzzy on the change name
//	'export   literal substring on the change name
//	:export   literal substring in artifact and spec text
func parseQuery(raw string) query {
	q := query{kind: matchFuzzyName, term: raw}
	switch {
	case strings.HasPrefix(raw, "'"):
		q.kind = matchLiteralName
		q.term = raw[1:]
	case strings.HasPrefix(raw, ":"):
		q.kind = matchBody
		q.term = raw[1:]
	}
	q.caseSensitive = hasUpper(q.term)
	return q
}

func (q query) empty() bool { return q.term == "" }

func hasUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// contains applies the smart-case rule: an all-lowercase term matches
// case-insensitively, a term carrying any uppercase matches exactly.
func (q query) contains(haystack string) bool {
	if q.caseSensitive {
		return strings.Contains(haystack, q.term)
	}
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(q.term))
}

// filterRows narrows the list to the rows matching the query, reporting for
// each survivor which of its bodies matched. Fuzzy results come back ranked
// strongest first; literal results keep the order they came in with, since
// there is no score to rank them by.
func filterRows[T tableRow](rows []T, q query) []filtered[T] {
	if q.empty() {
		return wrap(rows)
	}

	switch q.kind {
	case matchLiteralName:
		out := make([]filtered[T], 0, len(rows))
		for _, r := range rows {
			if q.contains(r.searchName()) {
				out = append(out, filtered[T]{row: r})
			}
		}
		return out

	case matchBody:
		out := make([]filtered[T], 0, len(rows))
		for _, r := range rows {
			labels, ok := q.bodyMatch(r)
			if !ok {
				continue
			}
			out = append(out, filtered[T]{row: r, matched: labels})
		}
		return out

	default:
		names := make([]string, len(rows))
		for i, r := range rows {
			names[i] = r.searchName()
		}
		matches := fuzzy.Find(q.term, names)
		out := make([]filtered[T], 0, len(matches))
		for _, m := range matches {
			// fuzzy.Find compares with equalFold and has no case-sensitive
			// mode, so smart case is enforced here by checking that every
			// character it matched also matches in case.
			if q.caseSensitive && !exactCaseMatch(m, q.term) {
				continue
			}
			out = append(out, filtered[T]{row: rows[m.Index]})
		}
		return out
	}
}

// exactCaseMatch reports whether the characters fuzzy matched agree with the
// pattern in case, not merely when folded.
func exactCaseMatch(m fuzzy.Match, pattern string) bool {
	pr := []rune(pattern)
	if len(m.MatchedIndexes) != len(pr) {
		return false
	}
	for i, idx := range m.MatchedIndexes {
		if idx < 0 || idx >= len(m.Str) {
			return false
		}
		r, _ := utf8.DecodeRuneInString(m.Str[idx:])
		if r != pr[i] {
			return false
		}
	}
	return true
}

// bodyMatch searches the row's name and every named body of text it exposes.
// It reports which labels hit so the row can explain itself: a row that matched
// only on the text of proposal.md shows a name the user will not recognise
// otherwise.
func (q query) bodyMatch(r tableRow) ([]string, bool) {
	matched := q.contains(r.searchName())

	bodies := r.searchBodies()
	var labels []string
	for _, label := range sortedLabels(bodies) {
		if q.contains(bodies[label]) {
			labels = append(labels, label)
			matched = true
		}
	}
	return labels, matched
}

func trimMarkdownSuffix(name string) string {
	return strings.TrimSuffix(name, ".md")
}
