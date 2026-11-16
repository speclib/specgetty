---
# specgetty-9xiu
title: Nested capability folders under specs/ are not shown
status: todo
type: bug
created_at: 2026-11-16T18:45:05Z
updated_at: 2026-11-16T18:45:05Z
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
