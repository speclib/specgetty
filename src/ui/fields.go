package ui

import (
	"fmt"
	"sort"
	"strings"
)

// fieldDef describes one column of the change table.
//
// width 0 means the column is flexible: it absorbs whatever space the
// fixed-width columns leave. Exactly one field is flexible (name).
type fieldDef struct {
	header string
	width  int
	value  func(r changeRow) string
}

// knownFields is the whole vocabulary accepted by config and --change-fields.
//
// Fields needing scanner capability that does not exist yet (schema,
// completeness against the schema, spec delta counts, created and updated
// dates) are deliberately absent. They are tracked in bean specgetty-vru8 and
// will be added here behind the same mechanism.
var knownFields = map[string]fieldDef{
	"name": {
		header: "name",
		width:  0,
		value:  func(r changeRow) string { return r.ci.Name },
	},
	"tasks": {
		header: "tasks",
		width:  7,
		value:  func(r changeRow) string { return r.taskLabel() },
	},
	"specs": {
		header: "specs",
		width:  5,
		value: func(r changeRow) string {
			if len(r.ci.SpecNames) == 0 {
				return ""
			}
			return fmt.Sprintf("%d", len(r.ci.SpecNames))
		},
	},
	"archived": {
		header: "state",
		width:  8,
		value: func(r changeRow) string {
			if r.archived {
				return "archived"
			}
			return "open"
		},
	},
	"date": {
		header: "archived",
		width:  10,
		value: func(r changeRow) string {
			if r.ci.ArchiveDate.IsZero() {
				return ""
			}
			return r.ci.ArchiveDate.Format("2006-01-02")
		},
	},
}

var defaultFields = []string{"name", "tasks", "specs"}

func validFieldNames() []string {
	names := make([]string, 0, len(knownFields))
	for n := range knownFields {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// ResolveFields picks the column set: the command-line flag wins over the
// config file, which wins over the default. An unknown name is reported with
// the valid ones rather than silently dropped, because a typo would otherwise
// look like a missing column.
func ResolveFields(flagValue, configValue string) ([]string, error) {
	source := "--change-fields"
	raw := strings.TrimSpace(flagValue)
	if raw == "" {
		source = "change_fields in the config file"
		raw = strings.TrimSpace(configValue)
	}
	if raw == "" {
		return append([]string(nil), defaultFields...), nil
	}

	var fields []string
	for _, part := range strings.Split(raw, ",") {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		if _, ok := knownFields[name]; !ok {
			return nil, fmt.Errorf("unknown field %q in %s; valid fields are: %s",
				name, source, strings.Join(validFieldNames(), ", "))
		}
		fields = append(fields, name)
	}
	if len(fields) == 0 {
		return append([]string(nil), defaultFields...), nil
	}
	return fields, nil
}

// minFlexWidth is the narrowest the name column may become before columns are
// dropped instead of squeezed further.
const minFlexWidth = 12

// layoutFields decides which of the requested columns fit in the available
// width and how wide each is. Fixed columns are served first; the flexible
// column takes the remainder. When the remainder would fall below
// minFlexWidth, columns are dropped from the right until it fits.
func layoutFields(fields []string, width int) (kept []string, widths []int) {
	kept = append([]string(nil), fields...)

	for {
		widths = widths[:0]
		fixed := 0
		flexCount := 0
		for _, f := range kept {
			def := knownFields[f]
			if def.width == 0 {
				flexCount++
			} else {
				fixed += def.width
			}
		}
		gaps := 0
		if len(kept) > 1 {
			gaps = len(kept) - 1
		}
		remaining := width - fixed - gaps

		if flexCount == 0 {
			// No flexible column: fixed columns alone must fit.
			if remaining >= 0 || len(kept) <= 1 {
				for _, f := range kept {
					widths = append(widths, knownFields[f].width)
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
			for _, f := range kept {
				def := knownFields[f]
				if def.width == 0 {
					widths = append(widths, flexEach)
				} else {
					widths = append(widths, def.width)
				}
			}
			return kept, widths
		}
		kept = kept[:len(kept)-1]
	}
}

// fitCell pads or truncates a value to exactly w columns, marking truncation
// with an ellipsis so a clipped name is visibly clipped.
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
