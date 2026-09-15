## 1. Fix: exporting an archived change has never worked

The scanner stores the archived change name with its date prefix already
stripped (`src/scanner/scan.go`, `displayName = dirName[11:]`), but
`doExportChange` builds `openspec/changes/archive/<that name>`, which is not a
directory that exists. Present since the feature shipped in `f2f6105`.
`exportSemanticName` strips a prefix the name never carries, which is the fossil
of the assumption that caused it.

`openspec/specs/export-change/spec.md` already requires the right behaviour
("all files and subdirectories from `openspec/changes/archive/2026-04-02-my-feature/`"),
so the spec does not change; the code is brought into line with it.

- [ ] 1.1 Give `ChangeInfo` the on-disk directory name for archived changes, so the exporter can find the directory without reconstructing it from a date and a name
- [ ] 1.2 Make `doExportChange` resolve the archived source from that name rather than from the display name
- [ ] 1.3 Keep the zip root folder and the zip filename at the stripped name, which is what the spec requires and what `exportSemanticName` already produces
- [ ] 1.4 Decide whether `exportSemanticName` still earns its place once the source path is correct, and remove it if it does not

## 2. Cover the export and discard logic (specgetty-17c3)

- [ ] 2.1 `exportSemanticName` and `exportDestPath` as pure functions: prefixed and unprefixed names, and the `~/<name>-<date>.zip` shape
- [ ] 2.2 `doExportChange` on an active change in a `t.TempDir()`: read the zip back and assert the root folder and the full file list
- [ ] 2.3 `doExportChange` on an archived change: assert it finds the dated directory and that the zip root is the stripped name. This is the test that fails before task 1
- [ ] 2.4 `doExportChange` on a change that does not exist: reports a failure rather than writing an empty zip
- [ ] 2.5 Overwriting an existing zip at the destination, which the spec requires
- [ ] 2.6 `doDiscardChange` moves the directory under `discarded/` with today's date prefix, and refuses when the target already exists
- [ ] 2.7 `doArchiveChange` when the `openspec` binary is absent: reports the clean error it already has. The success path shells out and is left alone

## 3. Cover the scanner walk and scan path (specgetty-sreo)

- [ ] 3.1 Lift the temp-project helper from `src/ui/loadprojects_test.go`; remember that `isValidOpenSpecDir` needs a `config.yaml` or `project.md` marker beside `specs/` or `changes/`, or nothing is detected
- [ ] 3.2 `Walk` finds projects under an include directory and ignores directories that are not projects
- [ ] 3.3 `exclude` in both its forms: a leading `/` compares the whole path, anything else compares the basename only (`src/scanner/find.go`, `skip`)
- [ ] 3.4 `followsymlinks` both ways
- [ ] 3.5 `ListOpenSpecContents` returns relative paths with the directory flag set correctly
- [ ] 3.6 `ScanPaths` parses a known list without walking, and skips a path that has since been deleted
- [ ] 3.7 `Scan` end to end: walk plus parse, with the task and spec counts that come out of it
- [ ] 3.8 `DumpConfig` round-trips a config

## 4. Remove the unused unwrap helper (specgetty-4wjm)

- [ ] 4.1 Delete `unwrap` from `src/ui/table.go`; it has no callers and was added speculatively during the generics extraction. Its counterpart `wrap` is used and stays

## 5. Verification

- [ ] 5.1 Write the archived-export test first and watch it fail, so the fix in task 1 is demonstrated rather than assumed
- [ ] 5.2 Export a real archived change from this repository and open the resulting zip, since the whole defect is that a path did not exist on disk
- [ ] 5.3 `go build ./...` and `go vet ./...` clean
- [ ] 5.4 `bash scripts/coverage-gate.sh` passes, and the floors are raised to the new numbers
- [ ] 5.5 Report how close scanner and ui land to the 80% the gate names as the target, since `specgetty-p4e2` waits on it
- [ ] 5.6 `nix flake check` passes

## 6. Notes

- [ ] 6.1 The `tinychange` schema comes from https://github.com/speclib/openspec-tinychange-schema and is vendored at `openspec/schemas/tinychange/`
