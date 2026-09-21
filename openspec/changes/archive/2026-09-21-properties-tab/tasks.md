## 1. Reading what schemas a project uses

- [x] 1.1 Read a change's `.openspec.yaml` for its `schema` key while parsing the change. Proven by a fixture change carrying one, and one carrying a file with no schema key
- [x] 1.2 Report a change with no `.openspec.yaml` as carrying no recorded schema, rather than assuming the project default. Proven by a fixture, and against the real `2026-04-20-discard-change-action`, which is the one change here that predates the file
- [x] 1.3 Compute the distinct schemas in use across active and archived changes, with a count each. Proven against a fixture with two schemas and an unrecorded change, asserting all three tallies
- [x] 1.4 Read the project default from the resolved root's configuration, falling back to `spec-driven` when there is no `schema` key, and mark which it is. Proven by a store-backed fixture whose repo and store name different schemas, asserting the store's wins
- [x] 1.5 Include the project default in the set even when no change uses it

## 2. Locating and reading a schema definition

- [x] 2.1 Run `openspec schema which <name> --json` with the resolved root as its working directory, and parse the path out of stdout. Proven by a fake `openspec` on PATH that records its arguments and working directory
- [x] 2.2 Read `<path>/schema.yaml` and parse description, ordered artifacts with `generates` and `requires`, and the `apply` block. Proven against both real schema files: the package's `spec-driven` and this project's `tinychange`
- [x] 2.3 Report a project schema as belonging to the project, and a built-in as not. Proven from the `source` field
- [x] 2.4 Report a project schema that shadows a built-in, naming what it overrides. Proven by a fixture where `shadows` is non-empty
- [x] 2.5 Bound the wait and abandon it when it elapses, leaving the interface responsive. Proven by a fake `openspec` that sleeps past the bound. The first attempt failed for a real reason worth recording: killing the process on a deadline does not close an output pipe a child inherited, so `Output()` waited for the child rather than for the bound. `cmd.WaitDelay` is what makes the bound real
- [x] 2.6 Report each failure distinctly: no binary on PATH, a schema that does not resolve (carrying the `available` list the CLI supplies), a timeout, and an unparseable `schema.yaml`. Proven by a test per case
- [x] 2.7 Never let a failure take down the interface or the other rows. Proven by asserting one failing schema leaves the others readable

## 3. When it runs

- [x] 3.1 Run nothing until the properties tab is opened. Proven by opening a project, never selecting the tab, and asserting the fake `openspec` was not invoked
- [x] 3.2 Run at most once per schema per project opened, serving later visits from what was read. Proven by opening the tab, leaving it, returning, and counting invocations
- [x] 3.3 Request several schemas concurrently. Proven by a fake `openspec` that sleeps, asserting two schemas complete in about the time of one rather than two
- [x] 3.4 Read again when a different project is opened, whether at startup or from the picker
- [x] 3.5 Do not reread when a schema file changes on disk while the project is open. Proven by editing a fixture schema mid-session and asserting the row is unchanged and no subprocess ran. This is the deliberate staleness, so it gets a test rather than a comment

## 4. The tab

- [x] 4.1 Rename the tab to `properties` in the tab bar and everywhere the UI names it. Proven by a rendered frame
- [x] 4.2 Draw the tab as a list and a content pane, each in its own border, reusing the specs tab's split. Proven by a line-count assertion at three widths including the 60-column minimum, which is the assertion a nested layout actually needs
- [x] 4.3 Size the list to its labels rather than to a share of the panel, and assert the content pane is wider than the specs rule would give it at 92 columns
- [x] 4.4 Rename `focusSpecsList` and `focusSpecsContent` to a generic pair, and check every comparison against them also establishes which tab is active. Proven by the existing specs tab focus tests passing unchanged
- [x] 4.5 Move focus between list, content and the log panel with the same key and the same ring the specs tab uses
- [x] 4.6 Light the border of whichever half holds the keyboard, and neither when the log does
- [x] 4.7 Move the selected section with the vertical keys while the list holds the keyboard, and scroll the document while the content does
- [x] 4.8 Start each section at the top when selected, and keep the active section's position across a tab switch. Proven by both moves

## 5. What each section shows

- [x] 5.1 Show the resolved root's configuration under `project`, naming the file's path. Proven by a store-backed fixture, asserting the store's configuration is shown and its path named. The path heads the document rather than sitting above the border: a split tab has no room for a naming line, and the list beside the content is what names the row
- [x] 5.2 Keep the existing markdown and YAML highlighting, and the message when there is no configuration at all
- [x] 5.3 Show `local` under `store` for a project holding its own content
- [x] 5.4 Show id, root, declaring file, and the registry's remote and branch under `store` for a store-backed project
- [x] 5.5 Keep the local git report: uncommitted changes, ahead and behind, and the line saying the comparison did not fetch
- [x] 5.6 Name the keys a declaring file carries in vain, with the reason. Proven by a fixture whose declaring file carries `schema`, `context`, `rules` and `operations`, and by one carrying only `store:` where nothing is warned about
- [x] 5.7 Keep reporting an unresolved declaration with its id, its file and the reason
- [x] 5.8 Give each used schema a row reporting source, path, artifact chain, apply rule and change count, and mark the project default
- [x] 5.9 Report the number of changes carrying no recorded schema, separately from any schema
- [x] 5.10 Say that a schema's details are being read while they are, and report the reason when they could not be

## 6. The open change

- [x] 6.1 Name an open change's schema beside its own name. Proven by a rendered frame for a change carrying one
- [x] 6.2 Name nothing for a change with no `.openspec.yaml`, rather than showing the project default as though it were recorded

## 7. Verification

- [x] 7.1 Open this project with the built binary and confirm the properties tab shows `spec-driven` with 24 changes, `tinychange` with 10, one change with no recorded schema, and `store: local`
- [x] 7.2 Open `~/gh.nivis-project/nivis-tunnel` and confirm `project` shows the store's configuration with its path, and `store` names the four keys the repo declares in vain
- [x] 7.3 Confirm with `openspec` removed from PATH that every row still renders and the schema rows report why their details are missing
- [x] 7.4 Revert task 3.1 and confirm the no-subprocess-until-opened test fails. Check the build succeeds first
- [x] 7.5 Revert task 5.6 and confirm the inert-declaration test fails
- [x] 7.6 Confirm every plain-project behaviour is untouched: a project with no store and one schema shows `project`, one schema row and `store: local`, with the same highlighting as before
- [x] 7.7 `nix flake check` passes, coverage floors included

## 8. Notes

- [x] 8.1 Update the README's properties section: the tab's name, its list, and that schema details need the `openspec` CLI
- [x] 8.2 Amend `specgetty-vru8` to record that its open question about resolving the schema is answered here, and that the change-list column can reuse this reader
- [x] 8.3 Note in `config-tab-display`'s Purpose that the tab is called `properties`, since the capability directory keeps its old name
