## 1. Resolution

- [x] 1.1 Read the store registry from `$XDG_DATA_HOME/openspec/stores/registry.yaml`, falling back to `~/.local/share`, returning the id to `backend.local_path` mapping plus the recorded `remote` and `branch`. Proven by a test setting `XDG_DATA_HOME` at a fixture and one leaving it unset
- [x] 1.2 Read a store declaration from `openspec/config.yaml`, and from `openspec/config.yml` when the first is absent. Distinguish absent, a string, a non-string and unparseable YAML. Proven by one test per case
- [x] 1.3 Read `.openspec-store/store.yaml` and return its id. Proven by a test on a fixture store and one on a directory without the file
- [x] 1.4 Resolve a root from a starting directory by the rules in `store-resolution`: walk up to the nearest qualifying `openspec/`, prefer local content over a pointer, follow a pointer through the registry otherwise. Proven by a table test covering each branch of the walk, including a repo carrying both `changes/` and `store:` where the pointer must be ignored
- [x] 1.5 Return the resolution's origin alongside its root, equal to the root when no pointer was followed. Proven by asserting origin on both a plain project and a store-backed one
- [x] 1.6 Report every unresolvable declaration distinctly: unregistered id, no registry file, malformed value, `local_path` gone, id mismatch between the registry key and `store.yaml`. Proven by a test per case asserting the declared id and the declaring file appear in the message

## 2. Scanning

- [x] 2.1 Carry the store id on a scanned root, read from `.openspec-store/store.yaml`, empty for a plain project. Proven by a scan over a fixture tree holding both
- [x] 2.2 Treat unreadable or idless store metadata as a plain project rather than dropping the root. Proven by a fixture with a corrupt `store.yaml`
- [x] 2.3 Pin the existing exclusion as deliberate: a directory whose `openspec/` holds a config file and neither `specs/` nor `changes/` is not reported. Proven by a walk test over a store-backed repo fixture, which fails if `isValidOpenSpecDir` is loosened
- [x] 2.4 Confirm the cache round-trips unchanged, since its paths are roots before and after. Proven by an existing cache test still passing against a tree containing a store

## 3. Reading and acting on the root

- [x] 3.1 Resolve at startup and read specs, changes and archive from the resolved root. Proven by opening a store-backed repo fixture and asserting the spec and change counts are the store's, not zero
- [x] 3.2 Point `doArchiveChange`'s `cmd.Dir` at the root. Proven by asserting the working directory handed to the command
- [x] 3.3 Point `doDiscardChange` at the root, and assert no directory is created under the origin. The second half is the defect being fixed, so it is its own assertion
- [x] 3.4 Point `doExportChange` at the root. Proven by exporting a change from a store-backed fixture
- [x] 3.5 Re-resolve on rescan, so an edited `store:` key is picked up by `s`. Proven by editing the pointer in a fixture and rescanning

## 4. What the user sees

- [x] 4.1 Mark the header when the content came from a store, showing the origin path when there is one and the store id when there is not. Proven by a rendered-frame assertion for all three header states: plain, resolved-through, store-itself
- [x] 4.2 Assert the mark carries nothing else: no id, path, remote or git state on the header. This is the scenario that keeps the header from growing a row
- [x] 4.3 Name a store row in the picker by its declared id and mark it. Proven by a picker render over a fixture where a store's id and folder differ
- [x] 4.4 Disambiguate a store id colliding with another row the way duplicate basenames are already disambiguated
- [x] 4.5 Rest the picker cursor on the row whose root holds the open content, which for a store-backed project is the store's row. Proven by opening the picker from a store-backed project
- [x] 4.6 `esc` from that picker leaves the view underneath unchanged, origin included

## 5. The config tab

- [x] 5.1 Give the config tab sub-tabs when there is more than one configuration, drawn where the filename is drawn today, above the content border. Proven by a line-count assertion on the tab at three widths, since a nested row is exactly the shape of change that width assertions miss
- [x] 5.2 Keep the single dimmed filename, and no sub-tab row, for a project with one configuration. Proven by asserting the plain-project frame is exactly its terminal's rows and columns at three widths including the 60-column minimum, that the filename is shown, and that the content box opens below it. Byte-identity against the previous build was the first plan and was dropped: it needs a golden captured from code that no longer exists, where a line-count assertion catches the one failure mode a nested row actually has
- [x] 5.3 Show the origin repo's configuration behind the repo sub-tab, with the YAML highlighting an ordinary configuration gets. Proven against a fixture whose repo config carries context and rules the store's does not, asserting those lines are reachable
- [x] 5.4 Show the store's configuration behind the store sub-tab
- [x] 5.5 Offer no repo sub-tab for a store opened directly from the picker
- [x] 5.6 Report id, root path, origin, and the registry's remote and branch behind the details sub-tab
- [x] 5.7 Report the store's local git state there: uncommitted changes, and ahead or behind its upstream ref. Read from local refs only, with the report saying the comparison is against the last known upstream
- [x] 5.8 Report a store that is not a git working copy without git state rather than as an error
- [x] 5.9 Report an unresolved declaration behind the details sub-tab: the declared id, the declaring file, and the reason. Proven by a fixture per failure mode from 1.6
- [x] 5.10 Make each sub-tab its own document, so selecting another starts at the top while leaving the config tab and returning keeps the active sub-tab's position. Proven by all three moves. This replaces the original task, which asked for a remembered position per sub-tab: the codebase settles every reset with one rule, that a different document starts at the top, and the spec list and artifact sub-tabs already follow it. A position map for this tab alone would have contradicted three archived specs

## 6. Watching

- [x] 6.1 Watch the store's `openspec/` tree and the origin's `openspec/` tree when the two differ, and exactly one tree when they do not. Proven by counting watched trees for both shapes
- [x] 6.2 Re-resolve and re-watch when the origin's configuration changes. Proven by repointing `store:` in a fixture and asserting the displayed content and the watched trees both follow
- [x] 6.3 Stop every tree of the previous project before starting the next. Proven by asserting no watcher survives a project switch

## 7. Verification

- [x] 7.1 Compare the resolver against `openspec context --json` over the fixture set, asserting the same root, the same store id and the same source classification. The comparison is a test, never a runtime dependency, and it is skipped when no `openspec` binary is present. One deliberate divergence is recorded rather than papered over: standing inside a store, the CLI reports source `nearest` and no store id, while specgetty still reads the identity file because the picker has to name that row. The test fails if the CLI ever starts reporting it, so the divergence cannot outlive its reason
- [x] 7.2 Open `~/gh.nivis-project/nivis-tunnel` with the built binary and confirm the store's real specs and archived changes appear, the header carries the mark, and all three config sub-tabs render
- [x] 7.3 Revert task 3.3 and confirm the discard test fails, checking the build succeeds first
- [x] 7.4 Revert task 2.3 and confirm the exclusion test fails
- [x] 7.5 Confirm every plain-project behaviour is untouched: the picker over a tree with no stores, the header, the config tab, the three actions and the watcher
- [x] 7.6 `nix flake check` passes, coverage floors included

## 8. Notes

- [x] 8.1 Created as specgetty-a2zv: a follow-up bean for the reverse index: which repos point at a given store, shown on the details sub-tab. The walk already passes every pointing repo's `openspec/` directory, so collecting the pointers costs one file read per skipped candidate
- [x] 8.2 Created as specgetty-lexl: a follow-up bean for `--store <id>` and the machine-wide `defaultStore`, the two resolution sources this change leaves out
