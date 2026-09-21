## 1. What counts as a store

- [x] 1.1 Treat a directory as a store only when its `.openspec-store/store.yaml` id is one the registry resolves to that same directory. Proven by three fixtures: a registered store, a directory whose claimed id the registry points elsewhere, and one whose id the registry does not mention
- [x] 1.2 Report a directory with store metadata and no registry at all as an ordinary project. Proven against an empty data directory
- [x] 1.3 Leave the declaration path alone, which already required the registry key and the metadata id to agree. Proven by the existing `store_identity` tests still passing untouched

## 2. Discovery

- [x] 2.1 Report a repo whose `openspec/` holds a configuration and no content. Proven by a walk over a fixture with a pointer-only repo, which is the assertion `follow-the-store` wrote in the opposite direction
- [x] 2.2 Report both repos when two declare the same store id
- [x] 2.3 Skip a directory the registry resolves a store id to. Proven by a walk over a tree holding a registered store and the repos pointing at it, asserting the exact row set
- [x] 2.4 Skip nothing as a store when there is no registry. Proven by a walk with an empty data directory over that same tree
- [x] 2.5 Keep skipping an `openspec/` directory with no configuration and no content. Proven by the existing empty-shell test still passing
- [x] 2.6 Read the registry once per walk rather than once per candidate. Proven by deleting the registry file and passing the loaded map: the store is still excluded, which it could not be if the check went back to disk per candidate. Counting opens was the first plan and was dropped, since it needs production instrumentation to observe something a behavioural test already pins

## 3. What a discovered row carries

- [x] 3.1 Give a discovered project the statistics of the root its content resolves to, so a store-backed repo carries the store's counts and not zero. Proven by scanning a fixture and asserting the repo row's spec and change counts equal the store's
- [x] 3.2 List the store's `openspec/` tree as the repo's file listing, so a `:` search inside the project reaches the specs the project shows. Proven by asserting a spec filename from the store appears in the repo row's searchable paths
- [x] 3.3 Read a shared store once per scan rather than once per repo pointing at it. Proven by removing the store's specs between the two reads and asserting the second still reports them, which only a memo can do
- [x] 3.4 Keep reporting a repo whose declaration cannot be followed, carrying the problem. Proven by a fixture naming an unregistered store
- [x] 3.5 Leave the git state out of discovery, read only when a project is opened. Proven by asserting a scanned row carries no git state and an opened project does

## 4. The picker

- [x] 4.1 Name every row by its directory, never by a store it declares. Proven by a fixture whose repo directory and store id differ
- [x] 4.2 Replace the `kind` column's `store` marker with the name of the store a row reads from, blank for a project holding its own content. Proven by both shapes
- [x] 4.3 Show both repos, both naming the store and both carrying its statistics, when two declare the same one
- [x] 4.4 Let the cursor miss quietly when the open project has no row, which is what opening a registered store by `--path` gives. Proven by opening the picker from such a project and asserting it neither fails nor selects an unrelated row
- [x] 4.5 Keep basename disambiguation working. Proven by two repos sharing a basename

## 5. Verification

- [x] 5.1 Open the picker against the real `$HOME/gh.nivis-project` tree and confirm the five repos appear, the four store directories do not, `ospecs` appears under its own name, and no row is called `nivis` twice
- [x] 5.2 Confirm `spg` inside a store's own directory still opens it, with the header naming the store and carrying the mark, since the picker no longer lists it
- [x] 5.3 Confirm the config tab of a repo row opened from the picker has all three sub-tabs, which is the gap that prompted this change
- [x] 5.4 Revert task 2.3 and confirm the exclusion test fails. Check the build succeeds first
- [x] 5.5 Revert task 1.1 and confirm the unregistered-metadata test fails
- [x] 5.6 Confirm every plain-project behaviour is untouched: a tree with no stores gives the same rows, names and statistics as before
- [x] 5.7 `nix flake check` passes, coverage floors included

## 6. Notes

- [x] 6.1 Update the README if it describes what the picker lists
- [x] 6.2 `specgetty-a2zv`, the reverse index, is unaffected: it names a store's readers, and this change does not alter what a store knows. Confirm its description still reads correctly and amend it if not
