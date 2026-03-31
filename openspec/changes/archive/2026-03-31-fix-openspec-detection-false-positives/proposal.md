## Why

The scanner matches any directory named `openspec` anywhere in the tree, causing false positives. For example, scanning `/home/pim/cForks/OpenSpec` produces four hits — the real project plus matches inside `.git/refs/`, `.git/logs/`, and `test/fixtures/`. Only the root project is a real OpenSpec source.

## What Changes

- Validate that a detected `openspec/` directory is a real OpenSpec project by checking for required contents: it must contain `config.yaml` or `project.md`, and at least one of `specs/` or `archive/`
- Stop descending into the parent directory after a valid openspec project is found (already done via `SkipThis`, but currently fires on invalid matches too)

## Capabilities

### New Capabilities

### Modified Capabilities
- `openspec-scanning`: Detection now validates openspec directory contents before reporting a project

## Impact

- **scanner/find.go**: Add validation function, update callback to use it
- **scanner/scan.go**: No changes needed
- **Tests**: Add test cases for false positive scenarios
