package ui

import (
	"sort"
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

// filterRows narrows the list to the rows matching the query. Fuzzy results
// come back ranked strongest first; literal results keep the order they came
// in with, since there is no score to rank them by.
func filterRows(rows []changeRow, q query) []changeRow {
	if q.empty() {
		return rows
	}

	switch q.kind {
	case matchLiteralName:
		out := make([]changeRow, 0, len(rows))
		for _, r := range rows {
			if q.contains(r.ci.Name) {
				r.matchedFiles = nil
				out = append(out, r)
			}
		}
		return out

	case matchBody:
		out := make([]changeRow, 0, len(rows))
		for _, r := range rows {
			files, ok := q.bodyMatch(r)
			if !ok {
				continue
			}
			r.matchedFiles = files
			out = append(out, r)
		}
		return out

	default:
		names := make([]string, len(rows))
		for i, r := range rows {
			names[i] = r.ci.Name
		}
		matches := fuzzy.Find(q.term, names)
		out := make([]changeRow, 0, len(matches))
		for _, m := range matches {
			// fuzzy.Find compares with equalFold and has no case-sensitive
			// mode, so smart case is enforced here by checking that every
			// character it matched also matches in case.
			if q.caseSensitive && !exactCaseMatch(m, q.term) {
				continue
			}
			r := rows[m.Index]
			r.matchedFiles = nil
			out = append(out, r)
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

// bodyMatch searches the change name and the text of every artifact and spec
// file. It reports which files hit so the row can explain itself: a change that
// matched only on the text of proposal.md shows a name the user will not
// recognise otherwise.
func (q query) bodyMatch(r changeRow) ([]string, bool) {
	matched := false
	var files []string

	if q.contains(r.ci.Name) {
		matched = true
	}
	for name, content := range r.ci.ArtifactContents {
		if q.contains(content) {
			files = append(files, trimMarkdownSuffix(name))
			matched = true
		}
	}
	for name, content := range r.ci.SpecContents {
		if q.contains(content) {
			files = append(files, name)
			matched = true
		}
	}
	sort.Strings(files)
	return files, matched
}

func trimMarkdownSuffix(name string) string {
	return strings.TrimSuffix(name, ".md")
}
