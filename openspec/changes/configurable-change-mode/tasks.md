## 1. Settle the vocabulary on "active"

The same three states are currently spelled two ways, and both reach the screen:
the nav bar says `f mode:open` while the empty state says "No active changes".
The specs contradict themselves inside a single requirement, and the README does
it across three lines.

"active" wins over "open" for three reasons. `openspec` itself says active, and
specgetty is a viewer for openspec. The scanner already says `ActiveChanges`;
only the UI layer invented "open". And "open" is overloaded here: an *open
change* is the one drilled into, which is why the detail header reads
`my-change (open)` about a change that is open on screen at that moment.

This is bundled with the configuration work because a config value is a
contract. Renaming a Go constant is free; renaming a value in someone's config
file is not.

- [ ] 1.1 Rename `modeOpen` to `modeActive`, 24 uses across `changerow.go`, `ui.go` and three test files
- [ ] 1.2 `listModeNames` becomes `active`, `archived`, `active+archived`
- [ ] 1.3 The `archived` column renders `active` rather than `open` (`fields.go`)
- [ ] 1.4 The open-change header renders `(active)` rather than `(open)` (`changelist.go`), which is the spelling that stops the header saying `(open)` about a change that is open on screen
- [ ] 1.5 Update the README: the `f` row currently says "Cycle open / archived / both" three lines from a sentence that says active, and the `archived` field row says "Whether the change is open or archived"

## 2. The configurable default

Mirrors `change_fields` exactly, which is the neighbour it sits beside.

- [ ] 2.1 Add `change_mode` to `scanner.Config` and to the commented example in `src/config.yml`
- [ ] 2.2 Add a `--change-mode` flag in `src/main.go`
- [ ] 2.3 Resolve as flag over config over default, with `active` as the default when nothing is set
- [ ] 2.4 Accept exactly the three names the nav bar shows, so what you see is what you write; report an unknown value with the valid ones rather than falling back silently
- [ ] 2.5 Apply it in `newModel`, replacing the hardcoded `listMode: modeOpen`
- [ ] 2.6 Apply it in `resetProjectState`, replacing the hardcoded `m.listMode = modeOpen`, so switching project returns to the configured default rather than to active

## 3. Tests

- [ ] 3.1 Mode resolution: flag beats config beats default
- [ ] 3.2 Each of the three names resolves to the right mode
- [ ] 3.3 An unknown name is an error naming the value and listing the valid ones
- [ ] 3.4 A model built with a configured default opens the change list in that mode
- [ ] 3.5 Switching project returns the mode to the configured default, not to active, with a default other than active so the test can tell the two apart
- [ ] 3.6 Pressing `f` still cycles through all three from whatever the default is
- [ ] 3.7 No occurrence of the word "open" as a mode name survives in what the user sees; assert on the nav bar hint, the state column and the open-change header

## 4. Verification

- [ ] 4.1 `go build ./...` and `go vet ./...` clean
- [ ] 4.2 Grep the test output for FAIL rather than counting PASS lines, which cannot see a failure
- [ ] 4.3 `grep -rn '"open"' src/` returns nothing that reaches the screen
- [ ] 4.4 `bash scripts/coverage-gate.sh` passes and the floors hold
- [ ] 4.5 `nix flake check` passes
- [ ] 4.6 Run `spg --change-mode=active+archived` and confirm the list opens with both, the nav bar reads `f mode:active+archived`, and `f` still cycles
- [ ] 4.7 Run `spg --change-mode=nonsense` and read the error

## 5. Notes

- [ ] 5.1 The `tinychange` schema comes from https://github.com/speclib/openspec-tinychange-schema and is vendored at `openspec/schemas/tinychange/`
