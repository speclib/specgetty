## 1. The registry is watched

- [ ] 1.1 Add the registry's directory to `watchDirs` for a project whose
      configuration declares a store, resolved or not, and verify a test asserts
      the watch set for a store-backed project, an unresolved-store project and a
      project with no store
- [ ] 1.2 Verify a test registering a store for an open project whose declaration
      could not be resolved ends with the project's specs and changes on screen
- [ ] 1.3 Verify a test repointing a registered store moves both the content and
      the watch set, which is the existing `sameDirs` path

## 2. A signal during a scan is not lost

- [ ] 2.1 Write a failing test that delivers a filesystem signal while a scan is
      in flight and asserts a second scan follows, and confirm it fails today
- [ ] 2.2 Add the pending flag beside `m.scanning`, set when a signal arrives
      during a scan and consumed when that scan lands, and verify 2.1 passes
- [ ] 2.3 Verify a test asserts exactly one further scan however many signals
      arrived during the one in flight, a scan reading everything
- [ ] 2.4 Verify a test asserts no further scan when nothing arrived during it
- [ ] 2.5 Verify the watcher's send stays non-blocking, by asserting it keeps
      reading events while the consumer is busy

## 3. Git state on entering the properties tab

- [ ] 3.1 Read the store's git state in `enterTab` for a store-backed project,
      alongside the schema loading already there, and verify a test asserts it is
      read on entry
- [ ] 3.2 Verify a test asserts nothing is read for a project with no store, and
      nothing at all for a session that never enters the tab
- [ ] 3.3 Verify a test asserts the read is local only, with no fetch, which is
      the existing `config-tab-display` requirement and must not move

## 4. The key goes

- [ ] 4.1 Remove the `s` case from the key handler, and verify a test presses `s`
      on every surface and asserts no scan is issued and no state changes
- [ ] 4.2 Remove all four `{"s", "scan"}` nav bar hints, and verify a test walks
      every surface and asserts the hint is absent
- [ ] 4.3 Verify `rescanCurrent` still has its three remaining callers, the
      watcher and the archive and discard results, and that a test covers each

## 5. Verification

- [ ] 5.1 Assert the picker still re-resolves and re-reads the project it opens,
      including a store's git state, this being the manual path that remains
- [ ] 5.2 Revert 2.2 and confirm 2.1 fails again. Check the build succeeds first
- [ ] 5.3 `nix flake check` passes, coverage floors included

## 6. Documentation

- [ ] 6.1 Remove the `s` row from the README's key table and point the reader at
      the picker for a manual re-read
- [ ] 6.2 Fix the sentence in the README's editor section that tells the reader
      their edit appears "without pressing `s`", which names a key that is gone
- [ ] 6.3 Add the change to CHANGELOG.md, the removal named apart from the three
      fixes that make it safe
