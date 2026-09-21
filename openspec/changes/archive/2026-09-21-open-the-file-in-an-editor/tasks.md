## 1. The root fix, first and on its own

- [x] 1.1 Write a failing test that toggles a task in a store-backed project and
      asserts the `tasks.md` under the store is written, and confirm it fails
      with the path error it fails with today
- [x] 1.2 Write a failing test that copies a change's path in a store-backed
      project and asserts the copied path is under the store and resolves to a
      directory, and confirm it fails
- [x] 1.3 Pass `m.currentRoot()` to `renderChangeArtifact` and to
      `changeDirPath`, and verify both tests now pass
- [x] 1.4 Verify the same tests pass unchanged for a project with no store, where
      the root and the starting directory are the same directory

## 2. Which file a pane is showing

- [x] 2.1 Set `document.path` for a change's artifact whatever the artifact is,
      not only for `tasks.md`, and verify a test asserts the path for proposal,
      design and tasks
- [x] 2.2 Set it for a spec on the specs tab, and verify a test asserts the path
      is `<root>/openspec/specs/<name>/spec.md`
- [x] 2.3 Set it for the project row of the properties tab from the configuration
      file that row reports, and verify a test asserts the path for both
      `config.yaml` and `project.md`
- [x] 2.4 Leave it empty for a change's specs sub-tab and for the schema and
      store rows, and verify a test asserts it is empty for all three
- [x] 2.5 Verify the task toggle still works, its guard being `docHasCursor()`
      which also requires `docLines`, so a path on every document does not give
      every document a cursor

## 3. Choosing the editor

- [x] 3.1 Add `src/ui/editor.go` resolving `$VISUAL` then `$EDITOR`, treating an
      empty value as unset, and verify a table test covers both set, only
      `EDITOR`, only `VISUAL`, and neither
- [x] 3.2 Split the chosen value on whitespace into a command and arguments with
      the file appended last, and verify a test covers a bare command, a command
      with one flag, and a command with two
- [x] 3.3 Return a reason rather than a command when neither variable is set, and
      verify a test asserts no process is described

## 4. The key

- [x] 4.1 Bind `E` to open `document.path` through `tea.ExecProcess`, and verify a
      test asserts the command built for a change artifact, a spec and the
      configuration row
- [x] 4.2 Do nothing when the document has no path, and verify a test presses `E`
      on the specs sub-tab of a change and on the store row and asserts no
      command
- [x] 4.3 Leave the key to the overlays, and verify a test presses `E` with the
      picker, a confirmation and the search prompt open and asserts no command
      and that the prompt received the character
- [x] 4.4 Report on the nav bar when no editor is configured, and verify a test
      asserts the message and that it clears on the next keystroke
- [x] 4.5 Report when the process could not be started and when it exits non-zero,
      and verify tests for both through the callback
- [x] 4.6 Rescan the open project in the callback, and verify a test asserts the
      command is issued

## 5. The nav bar

- [x] 5.1 List `E` exactly where the document has a path, and verify a test walks
      every pane and asserts the hint is present for the three that have a file
      and absent for the rest

## 6. Verification

- [x] 6.1 Revert the call-site change from 1.3 and confirm 1.1 and 1.2 fail again.
      Check the build succeeds first
- [x] 6.2 Assert no pane gained a cursor it did not have, by checking
      `docHasCursor()` for every pane before and after
- [x] 6.3 `nix flake check` passes, coverage floors included

## 7. Documentation

- [x] 7.1 Document `E` in README.md, including `$VISUAL` before `$EDITOR` and that
      a value with arguments works while a quoted path inside it does not
- [x] 7.2 Add the change to CHANGELOG.md, the fixes named apart from the feature
