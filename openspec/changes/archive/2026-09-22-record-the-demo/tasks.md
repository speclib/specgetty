## 1. The fixture projects

- [x] 1.1 Add `demo/stores/tideclock/` with its `.openspec-store/store.yaml` and
      an `openspec/` holding two specs, and verify `spg --path` opens it and
      lists both
- [x] 1.2 Write the two specs so they fill an outline: at least two requirements
      each with two scenarios, in prose that reads on screen rather than
      `the system SHALL do the thing`, and verify the spec detail view opens both
      without a report
- [x] 1.3 Add `demo/projects/harbour-app/openspec/config.yaml` declaring
      `store: tideclock`, and verify the header carries the store mark
- [x] 1.4 Add `demo/projects/ferry-times/` and `demo/projects/quay-signage/` as
      plain projects, and verify the picker lists three rows
- [x] 1.5 Add archived changes to the store with dated directory names, and
      verify the change list groups them under ARCHIVED with their dates

## 2. The change that stays active

- [x] 2.1 Add one active change to the store with a proposal, a design and a
      tasks file, some tasks checked, and verify the artifact sub-tabs show all
      three and the task counts read
- [x] 2.2 Give it a spec delta with an ADDED and a REMOVED requirement, and
      verify the outline marks them `+` and `-`
- [x] 2.3 Give it a MODIFIED requirement whose original is in the store's live
      specs, restating it with one scenario untouched, one edited and one added,
      and verify the outline marks all three and the card offers diff, old and
      new
- [x] 2.4 Verify the change stays active, by asserting `demo/stores/tideclock/
      openspec/changes/` holds it outside `archive/` and that nothing in the
      repository archives it

## 3. Recording

- [x] 3.1 Add `demo/config.yml` with `scandirs.include` naming only
      `demo/projects` and `export_dir` naming `demo/out`, and verify a picker
      scan finds three projects and nothing else
- [x] 3.2 Add a `demo` target that generates `demo/data/openspec/stores/
      registry.yaml` with an absolute `local_path` from `$(PWD)`, and verify the
      store resolves from a fresh clone at a different path
- [x] 3.3 Make the target create `demo/recordings/`, stamp the version recorded
      at into it, and report what is missing when `vhs` is absent rather than
      failing obscurely
- [x] 3.4 Track `demo/recordings/` rather than ignoring it: the README shows the
      recordings, so a GIF nobody committed is a broken image. Ignore the
      staging scratch instead, and verify a recording changes only the GIFs

## 4. The tapes

- [x] 4.1 Give every tape a hidden opening block that enters the fixture project
      and exports `XDG_DATA_HOME`, and verify no tape reaches the
      `No OpenSpec project here` prompt
- [x] 4.2 `hero.tape`: open a spec in detail, move through the outline, about
      fifteen seconds, and verify its first frame shows the outline beside the
      card rather than a shell prompt
- [x] 4.3 `changes.tape`: the grouped list, an active change opened, its artifact
      sub-tabs, then an archived one
- [x] 4.4 `change-specs.tape`: the delta outline, the marks, and the card moved
      through diff, old and new on the modified requirement
- [x] 4.5 `specs.tape`: the specs tab, then `enter` into the detail view
- [x] 4.6 `properties.tape`: the grouped rows and the content beside them
- [x] 4.7 `picker.tape`: `p`, the three fixture projects, the store column
- [x] 4.8 `export.tape`: `e`, the directory prompt, the result, and verify the
      zip lands in `demo/out/`

## 5. Verification

- [x] 5.1 Record every tape and view each result, asserting by eye that the
      feature named in the tape is what the recording shows
- [x] 5.2 Verify no recording contains a path outside the repository, by
      grepping the extracted frames' text for the home directory
- [x] 5.3 Establish what determinism is available. VHS output is not
      reproducible: two recordings of one tape differ in frame count and in every
      pixel, because the GIF encoder quantises a palette per file. The content is
      fixed at the source instead, asserted by tests over the fixture, and that
      is what this task verifies
- [x] 5.4 Verify a clone at a different path records the same content, by
      recording from one and reading the frame. This is what staging the fixture
      at a fixed path is for, and it is checked by eye because bytes cannot be
      compared across encodes
- [x] 5.5 `nix flake check` passes, coverage floors included, this change adding
      no Go code

## 6. Documentation

- [x] 6.1 Put the hero at the top of README.md and the six others beside the
      sections they illustrate
- [x] 6.2 Document `make demo` and what it needs, in README.md or CONTRIBUTING
- [x] 6.3 Add the change to CHANGELOG.md under Added
