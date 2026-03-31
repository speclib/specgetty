## Context

The walker in `find.go` currently emits a project whenever it finds a directory named `openspec`. This is too loose — any nested directory with that name triggers a match, including inside `.git/`, test fixtures, etc.

## Goals / Non-Goals

**Goals:**
- Only report directories that contain a genuine OpenSpec project
- A valid openspec directory must contain: (`config.yaml` OR `project.md`) AND (`specs/` OR `archive/`)

**Non-Goals:**
- Deep validation of OpenSpec file contents (parsing YAML, checking spec structure)
- Changing walker traversal order or performance characteristics

## Decisions

### Validation function in find.go
Add `isValidOpenSpecDir(path string) bool` that checks for required marker files/dirs using `os.Stat`. This keeps the validation close to the detection logic and avoids coupling to the scanner package's data model.

The check: at least one of `config.yaml` or `project.md` must exist, AND at least one of `specs/` or `archive/` must exist as a directory.

**Alternative considered**: Checking only for `config.yaml`. Rejected because some OpenSpec projects use `project.md` as the primary marker, and requiring both would be too strict.

## Risks / Trade-offs

- **Risk**: Additional `os.Stat` calls per candidate directory → Mitigation: Only 2-4 stat calls per `openspec/` directory found, negligible overhead.
