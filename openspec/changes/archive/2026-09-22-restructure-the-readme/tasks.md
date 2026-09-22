## 1. The landing page

- [x] 1.1 Write the opening: what the tool is, and why in the terms this change
      settles, that reviewing specs is now the work and there is more of it to
      read than there used to be. Explain OpenSpec in a sentence and link it, it
      being the noun everything else rests on and the one the page never defines
- [x] 1.2 Add all five badges: licence, Go version from `go.mod`, build status,
      coverage and the OpenSpec metrics. Every URL was fetched rather than
      trusted: all six return 200, and the Go badge reads v1.25.0, which is what
      `go.mod` says
- [x] 1.3 Replace the installation section with the three paths that work, each
      run before being written down. The section this replaced did not work at
      all: `go install github.com/mipmip/specgetty@master` fails twice over,
      there being no `master` revision and no package at the module root. What
      works is `.../specgetty/src@latest`, which installs a binary named `src`,
      so the README says so and renames it. The release archives are named per
      version, so `releases/latest/download/<name>` 404s and the command uses a
      pattern instead. `nix run github:speclib/specgetty` works as written
- [x] 1.4 Keep the features list and the hero recording, add a quick start of
      three lines, and add a table linking each documentation page
- [x] 1.5 Add contributing, where to get help, the licence and related projects,
      three or four lines each. `LICENSE` is MIT and is never mentioned today
- [x] 1.6 Verify the page is under about 140 lines and loads one recording rather
      than seven. 130 lines and the hero alone, 334 KB against 2.5 MB

## 2. The documentation pages

- [x] 2.1 Move every key table into `docs/keys.md`, from all five levels, and
      give it a table of contents
- [x] 2.2 Move the three navigation levels, the keyword vocabulary and the report
      view into `docs/reading-specs.md`
- [x] 2.3 Move the picker, stores, the cache and the properties tab into
      `docs/projects.md`
- [x] 2.4 Move search, columns, grouping and export into `docs/change-list.md`
- [x] 2.5 Move `E` and the editor environment into `docs/editor.md`
- [x] 2.6 Move each recording to the page that documents it, leaving the hero on
      the README, and verify every image path resolves from its new directory.
      Each of the seven now appears exactly once. `changes.gif` came along with
      the levels section but shows the change list, so it went to that page
      rather than staying where it was carried

## 3. Corrections carried out while moving

- [x] 3.1 Remove the `f` key row from the change list keys. The key, the
      `change_mode` setting and the `--change-mode` flag were all removed by
      `group-the-change-list`, and the section below already explains that
- [x] 3.2 Correct the specs tab passage saying a spec that does not fit reports
      on the status line. `follow-the-spec-grammar` replaced that with the report
      view, which the same document already describes
- [x] 3.3 Read the whole of both documents against the running binary and note
      anything else that has drifted. Nothing else had: the two known passages
      were the whole of it, and the heading comparison confirms no section was
      lost in the move

## 4. Verification

- [x] 4.1 Verify every internal link in the README and in `docs/` resolves, by a
      link checker rather than by eye. Anchors too, since the split created
      cross-page ones: 0 broken of both kinds
- [x] 4.2 Verify no section was lost in the move, by comparing the set of
      headings before and after. Three differ and all three are deliberate
      renames: `UI` to `The levels`, `Source-mode installation` to `Install`,
      and `Opening a file in your editor` becoming its page's title
- [x] 4.3 Read the landing page as a stranger would, in order, and check that
      each of these is answered before the documentation table: what it is, why
      it exists, what OpenSpec is, what it does, how to install it, how to run
      it. All six, within the first 92 lines
- [x] 4.4 `nix flake check` passes, nothing here touching the software

## 5. Notes

- [x] 5.1 This change ships after `add-continuous-integration`, which publishes
      the build status, the coverage number and the OpenSpec metrics. Shipping
      them the other way round would mean editing the README twice for badges
      and adding three that point at nothing in between
- [x] 5.2 The module path is `github.com/mipmip/specgetty` while the repository
      is `speclib/specgetty`. The install path resolves through a GitHub
      redirect and works, so it is left alone here, but the README will name one
      of the two and the reader may notice the other
