---
# specgetty-9xiu
title: Nested capability folders under specs/ are not shown
status: in-progress
type: bug
priority: normal
created_at: 2026-11-16T18:45:05Z
updated_at: 2026-11-16T22:27:46Z
---

Capabilities nested in subdirectories are invisible. OpenSpec allows a
capability path with more than one segment (`identity/user-auth`); its own
artifact instructions use that example, and `openspec list --specs` and
`openspec validate --strict` accept it.

The scanner reads only one level under `specs/`:

- `src/scanner/scan.go` `parseChangeDir`: `os.ReadDir(<change>/specs)`, then
  `<dir>/spec.md`.
- `src/scanner/scan.go` `ParseProjectInfo`: the same for `openspec/specs/`.
- `src/ui/docview.go` builds `specs/<name>/spec.md` from that one-segment name.

So `specs/onboarding/import/spec.md` shows up as a capability named
`onboarding` with no content, and two deltas under `specs/tenancy/...`
collapse into one name.

Seen in chatcms-app, change `hosted-until-connected`: three deltas
(`onboarding/import`, `tenancy/projects`, `publish/act`), none visible.
The main specs there are all nested too (`platform/release`,
`content/editing`, and 49 more).

Expected: walk `specs/` recursively and treat every directory holding a
`spec.md` as a capability, named by its path relative to `specs/`.


## Implementation checklist

- [x] 1.1 capability discovery helper
- [x] 1.2 both scan call sites use it
- [x] 1.3 delta and live names agree
- [x] 2.1 grouped specs list
- [x] 2.2 paging over drawn lines
- [x] 2.3 changelog
- [x] 3.1 tests and build
- [x] 3.2 verify against chatcms-app


- [x] 3.3 End-to-end fixture test: scan a nested project from disk, render the grouped specs tab, compare the delta
