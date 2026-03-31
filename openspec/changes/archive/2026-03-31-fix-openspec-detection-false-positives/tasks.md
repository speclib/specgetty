## 1. Validation logic

- [x] 1.1 Add `isValidOpenSpecDir(path string) bool` to `find.go` that checks for (`config.yaml` OR `project.md`) AND (`specs/` OR `archive/`)
- [x] 1.2 Update the walker callback to call `isValidOpenSpecDir` before emitting a result — only emit and skip if valid, otherwise continue walking

## 2. Tests

- [x] 2.1 Add test for `isValidOpenSpecDir` with valid directory (config.yaml + specs/)
- [x] 2.2 Add test for `isValidOpenSpecDir` with valid directory (project.md + archive/)
- [x] 2.3 Add test for `isValidOpenSpecDir` with invalid directory (empty or missing markers)
- [x] 2.4 Verify build and existing tests pass
