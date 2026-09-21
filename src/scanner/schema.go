package scanner

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// A workflow schema decides which artifacts a change needs. Project schemas sit
// under the resolved root's `openspec/schemas/`, but the built-in ones ship
// inside the OpenSpec installation, at a path that is a package directory here
// and a store hash there and moves on every upgrade.
//
// So the location is asked for and the content is read: `openspec schema which`
// answers where, and the parser below answers what. That is the same division
// the store work settled on, with a subprocess standing in for the registry
// because nothing about the package's location is documented or stable.

// Schema problem codes. Named because the display and the tests both key on
// them, and because a person needs to be told which of these happened.
const (
	SchemaNoCLI      = "no_cli"
	SchemaNotFound   = "not_found"
	SchemaTimeout    = "timeout"
	SchemaUnreadable = "unreadable"
)

// SchemaTimeout bounds the wait for the CLI. A hung subprocess must not take
// the interface with it, and a second is already long for a terminal.
const schemaWait = 10 * time.Second

// SchemaArtifact is one step of a workflow.
type SchemaArtifact struct {
	ID          string   `yaml:"id"`
	Generates   string   `yaml:"generates"`
	Template    string   `yaml:"template"`
	Description string   `yaml:"description"`
	Requires    []string `yaml:"requires"`
}

// SchemaApply is what a schema requires before implementation and what it
// tracks progress in.
type SchemaApply struct {
	Requires []string `yaml:"requires"`
	Tracks   string   `yaml:"tracks"`
}

// SchemaShadow is a copy of a schema that another one overrides.
type SchemaShadow struct {
	Source string `json:"source"`
	Path   string `json:"path"`
}

// SchemaDetails is a workflow schema as read from its definition.
type SchemaDetails struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	// Source is OpenSpec's own word: `project` for one under the root,
	// `package` for one shipped with the CLI.
	Source    string
	Path      string
	Artifacts []SchemaArtifact `yaml:"artifacts"`
	Apply     SchemaApply      `yaml:"apply"`
	Shadows   []SchemaShadow
}

// Overrides reports whether this schema hides another of the same name.
func (d SchemaDetails) Overrides() bool { return len(d.Shadows) > 0 }

// SchemaProblem is a schema whose definition could not be read.
type SchemaProblem struct {
	Code   string
	Name   string
	Detail string
	// Available is the list the CLI supplies when it refuses a name, which is
	// the most useful thing to show beside the one that failed.
	Available []string
}

func (p *SchemaProblem) Error() string {
	if p == nil {
		return ""
	}
	return p.Detail
}

// whichOutput is what `openspec schema which --json` prints, on success and on
// failure. The CLI puts its experimental notice on stderr, so stdout is clean
// either way.
type whichOutput struct {
	Name      string         `json:"name"`
	Source    string         `json:"source"`
	Path      string         `json:"path"`
	Shadows   []SchemaShadow `json:"shadows"`
	Error     string         `json:"error"`
	Available []string       `json:"available"`
}

// ResolveSchema reads a workflow schema's definition.
//
// root is the resolved root, which is where the command runs so that a
// store-backed project resolves the store's schemas. Exactly one of the two
// return values is non-nil.
func ResolveSchema(root, name string) (*SchemaDetails, *SchemaProblem) {
	return resolveSchemaWith(root, name, schemaWait)
}

func resolveSchemaWith(root, name string, wait time.Duration) (*SchemaDetails, *SchemaProblem) {
	if _, err := exec.LookPath("openspec"); err != nil {
		return nil, &SchemaProblem{Code: SchemaNoCLI, Name: name,
			Detail: "the openspec CLI is not installed, so a schema's definition cannot be located"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	cmd := exec.CommandContext(ctx, "openspec", "schema", "which", name, "--json")
	cmd.Dir = root
	// Nothing here should ever wait on a person or a network.
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "NO_COLOR=1")
	// Killing the process on a timeout is not enough on its own: a child that
	// inherited the output pipe keeps it open, and Output() would wait for that
	// child rather than for the deadline. WaitDelay forces the pipes shut
	// shortly after the context is done, which is what makes the bound real.
	cmd.WaitDelay = 500 * time.Millisecond
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return nil, &SchemaProblem{Code: SchemaTimeout, Name: name,
			Detail: "the openspec CLI did not answer within " + wait.String()}
	}

	var parsed whichOutput
	// A refused name exits non-zero and still prints JSON, so the output is
	// parsed before the exit status is judged.
	if jsonErr := json.Unmarshal(out, &parsed); jsonErr != nil {
		detail := "the openspec CLI gave no readable answer"
		if err != nil {
			detail += ": " + err.Error()
		}
		return nil, &SchemaProblem{Code: SchemaNotFound, Name: name, Detail: detail}
	}
	if parsed.Error != "" || parsed.Path == "" {
		detail := parsed.Error
		if detail == "" {
			detail = "the openspec CLI named no path for it"
		}
		return nil, &SchemaProblem{Code: SchemaNotFound, Name: name,
			Detail: detail, Available: parsed.Available}
	}

	details, readErr := readSchemaDefinition(parsed.Path)
	if readErr != nil {
		return nil, &SchemaProblem{Code: SchemaUnreadable, Name: name,
			Detail: "its definition at " + parsed.Path + " could not be read: " + readErr.Error()}
	}
	details.Source = parsed.Source
	details.Path = parsed.Path
	details.Shadows = parsed.Shadows
	if details.Name == "" {
		details.Name = name
	}
	return details, nil
}

// readSchemaDefinition parses a schema.yaml from a located schema directory.
func readSchemaDefinition(dir string) (*SchemaDetails, error) {
	b, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "schema.yaml"))
	if err != nil {
		return nil, err
	}
	var d SchemaDetails
	if err := yaml.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	d.Description = strings.TrimSpace(d.Description)
	for i := range d.Artifacts {
		d.Artifacts[i].Description = strings.TrimSpace(d.Artifacts[i].Description)
	}
	return &d, nil
}
