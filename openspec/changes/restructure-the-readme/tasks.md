## 1. The landing page

- [ ] 1.1 Write the opening: what the tool is, and why in the terms this change
      settles, that reviewing specs is now the work and there is more of it to
      read than there used to be. Explain OpenSpec in a sentence and link it, it
      being the noun everything else rests on and the one the page never defines
- [ ] 1.2 Add all five badges: licence, Go version from `go.mod`, build status,
      coverage and the four OpenSpec metrics. The last three are published by
      `add-continuous-integration`, which ships first. Verify each renders by
      opening the pushed page rather than by trusting the markdown
- [ ] 1.3 Replace the installation section with the three paths that work: the
      released archive from `goreleaser`, `go install` at a tagged version rather
      than `@master`, and the nix flake. Verify each by running it on a clean
      machine or a container, because an installation line that does not work is
      the one error a reader cannot recover from
- [ ] 1.4 Keep the features list and the hero recording, add a quick start of
      three lines, and add a table linking each documentation page
- [ ] 1.5 Add contributing, where to get help, the licence and related projects,
      three or four lines each. `LICENSE` is MIT and is never mentioned today
- [ ] 1.6 Verify the page is under about 140 lines and loads one recording rather
      than seven

## 2. The documentation pages

- [ ] 2.1 Move every key table into `docs/keys.md`, from all five levels, and
      give it a table of contents
- [ ] 2.2 Move the three navigation levels, the keyword vocabulary and the report
      view into `docs/reading-specs.md`
- [ ] 2.3 Move the picker, stores, the cache and the properties tab into
      `docs/projects.md`
- [ ] 2.4 Move search, columns, grouping and export into `docs/change-list.md`
- [ ] 2.5 Move `E` and the editor environment into `docs/editor.md`
- [ ] 2.6 Move each recording to the page that documents it, leaving the hero on
      the README, and verify every image path resolves from its new directory:
      `demo/recordings/...` becomes `../demo/recordings/...` from inside `docs/`

## 3. Corrections carried out while moving

- [ ] 3.1 Remove the `f` key row from the change list keys. The key, the
      `change_mode` setting and the `--change-mode` flag were all removed by
      `group-the-change-list`, and the section below already explains that
- [ ] 3.2 Correct the specs tab passage saying a spec that does not fit reports
      on the status line. `follow-the-spec-grammar` replaced that with the report
      view, which the same document already describes
- [ ] 3.3 Read the whole of both documents against the running binary and note
      anything else that has drifted, rather than assuming these two are all of
      it

## 4. Verification

- [ ] 4.1 Verify every internal link in the README and in `docs/` resolves, by a
      link checker rather than by eye
- [ ] 4.2 Verify no section was lost in the move, by comparing the set of
      headings before and after
- [ ] 4.3 Read the landing page as a stranger would, in order, and check that
      each of these is answered before the documentation table: what it is, why
      it exists, what OpenSpec is, what it does, how to install it, how to run it
- [ ] 4.4 `nix flake check` passes, nothing here touching the software

## 5. Notes

- [ ] 5.1 This change ships after `add-continuous-integration`, which publishes
      the build status, the coverage number and the OpenSpec metrics. Shipping
      them the other way round would mean editing the README twice for badges
      and adding three that point at nothing in between
- [ ] 5.2 The module path is `github.com/mipmip/specgetty` while the repository
      is `speclib/specgetty`. The install path resolves through a GitHub
      redirect and works, so it is left alone here, but the README will name one
      of the two and the reader may notice the other
