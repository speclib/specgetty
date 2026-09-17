package ui

import (
	"fmt"
	"sort"
	"strings"
)

// knownFields is the whole vocabulary accepted by config and --change-fields.
//
// Fields needing scanner capability that does not exist yet (schema,
// completeness against the schema, spec delta counts, created and updated
// dates) are deliberately absent. They are tracked in bean specgetty-vru8 and
// will be added here behind the same mechanism.
var knownFields = map[string]fieldDef[changeRow]{
	"name": {
		id: "name", header: "name", width: 0,
		value: func(r changeRow) string { return r.ci.Name },
	},
	"tasks": {
		id: "tasks", header: "tasks", width: 7,
		value: func(r changeRow) string { return r.taskLabel() },
	},
	"specs": {
		id: "specs", header: "specs", width: 5,
		value: func(r changeRow) string {
			if len(r.ci.SpecNames) == 0 {
				return ""
			}
			return fmt.Sprintf("%d", len(r.ci.SpecNames))
		},
	},
	"archived": {
		id: "archived", header: "state", width: 8,
		value: func(r changeRow) string {
			if r.archived {
				return "archived"
			}
			return "active"
		},
	},
	"date": {
		id: "date", header: "archived", width: 10,
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

// changeFieldDefs turns resolved field names into column definitions.
func changeFieldDefs(names []string) []fieldDef[changeRow] {
	defs := make([]fieldDef[changeRow], 0, len(names))
	for _, n := range names {
		if d, ok := knownFields[n]; ok {
			defs = append(defs, d)
		}
	}
	if len(defs) == 0 {
		defs = append(defs, knownFields["name"])
	}
	return defs
}

// ResolveListMode picks the change list's starting filter: the command-line
// flag wins over the config file, which wins over active only.
//
// The names accepted here are exactly the ones the nav bar shows, so what you
// see is what you write. An unknown value is reported with the valid ones
// rather than falling back silently, because a typo that quietly does nothing
// looks identical to a setting that is being ignored.
func ResolveListMode(flagValue, configValue string) (int, error) {
	source := "--change-mode"
	raw := strings.TrimSpace(flagValue)
	if raw == "" {
		source = "change_mode in the config file"
		raw = strings.TrimSpace(configValue)
	}
	if raw == "" {
		return modeActive, nil
	}

	for i, name := range listModeNames {
		if raw == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("unknown mode %q in %s; valid modes are: %s",
		raw, source, strings.Join(listModeNames, ", "))
}
