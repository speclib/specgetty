## 1. The fixture projects

- [ ] 1.1 Add `demo/stores/tideclock/` with its `.openspec-store/store.yaml` and
      an `openspec/` holding two specs, and verify `spg --path` opens it and
      lists both
- [ ] 1.2 Write the two specs so they fill an outline: at least two requirements
      each with two scenarios, in prose that reads on screen rather than
      `the system SHALL do the thing`, and verify the spec detail view opens both
      without a report
- [ ] 1.3 Add `demo/projects/harbour-app/openspec/config.yaml` declaring
      `store: tideclock`, and verify the header carries the store mark
- [ ] 1.4 Add `demo/projects/ferry-times/` and `demo/projects/quay-signage/` as
      plain projects, and verify the picker lists three rows
- [ ] 1.5 Add archived changes to the store with dated directory names, and
      verify the change list groups them under ARCHIVED with their dates

## 2. The change that stays active

- [ ] 2.1 Add one active change to the store with a proposal, a design and a
      tasks file, some tasks checked, and verify the artifact sub-tabs show all
      three and the task counts read
- [ ] 2.2 Give it a spec delta with an ADDED and a REMOVED requirement, and
      verify the outline marks them `+` and `-`
- [ ] 2.3 Give it a MODIFIED requirement whose original is in the store's live
      specs, restating it with one scenario untouched, one edited and one added,
      and verify the outline marks all three and the card offers diff, old and
      new
- [ ] 2.4 Verify the change stays active, by asserting `demo/stores/tideclock/
      openspec/changes/` holds it outside `archive/` and that nothing in the
      repository archives it

## 3. Recording

- [ ] 3.1 Add `demo/config.yml` with `scandirs.include` naming only
      `demo/projects` and `export_dir` naming `demo/out`, and verify a picker
      scan finds three projects and nothing else
- [ ] 3.2 Add a `demo` target that generates `demo/data/openspec/stores/
      registry.yaml` with an absolute `local_path` from `$(PWD)`, and verify the
      store resolves from a fresh clone at a different path
- [ ] 3.3 Make the target create `demo/out/`, stamp the version recorded at into
      it, and report what is missing when `vhs` is absent rather than failing
      obscurely
- [ ] 3.4 Add the generated registry, `demo/data/`, `demo/out/` and the recorded
      files to `.gitignore`, and verify `git status` is clean after a recording

## 4. The tapes

- [ ] 4.1 Give every tape a hidden opening block that enters the fixture project
      and exports `XDG_DATA_HOME`, and verify no tape reaches the
      `No OpenSpec project here` prompt
- [ ] 4.2 `hero.tape`: open a spec in detail, move through the outline, about
      fifteen seconds, and verify its first frame shows the outline beside the
      card rather than a shell prompt
- [ ] 4.3 `changes.tape`: the grouped list, an active change opened, its artifact
      sub-tabs, then an archived one
- [ ] 4.4 `change-specs.tape`: the delta outline, the marks, and the card moved
      through diff, old and new on the modified requirement
- [ ] 4.5 `specs.tape`: the specs tab, then `enter` into the detail view
- [ ] 4.6 `properties.tape`: the grouped rows and the content beside them
- [ ] 4.7 `picker.tape`: `p`, the three fixture projects, the store column
- [ ] 4.8 `export.tape`: `e`, the directory prompt, the result, and verify the
      zip lands in `demo/out/`

## 5. Verification

- [ ] 5.1 Record every tape and view each result, asserting by eye that the
      feature named in the tape is what the recording shows
- [ ] 5.2 Verify no recording contains a path outside the repository, by
      grepping the extracted frames' text for the home directory
- [ ] 5.3 Record twice and verify the two runs produce the same frames, which is
      what says the fixtures are deterministic
- [ ] 5.4 Verify a fresh clone at a different path records identically, which is
      what the generated registry is for
- [ ] 5.5 `nix flake check` passes, coverage floors included, this change adding
      no Go code

## 6. Documentation

- [ ] 6.1 Put the hero at the top of README.md and the six others beside the
      sections they illustrate
- [ ] 6.2 Document `make demo` and what it needs, in README.md or CONTRIBUTING
- [ ] 6.3 Add the change to CHANGELOG.md under Added
