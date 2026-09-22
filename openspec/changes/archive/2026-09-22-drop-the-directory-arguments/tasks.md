# Tasks

Schema source: https://github.com/speclib/openspec-tinychange-schema

## 1. Implementation

- [x] 1.1 Remove the positional-argument branch from `runApp` in `src/main.go`:
      the `c.Args().Len() > 0` case, the `Arguments given, skipping config`
      message, and the assignment of `config.ScanDirs.Include` from the
      arguments. `expandScanDirs` then runs on every path rather than only when
      no argument was given
- [x] 1.2 Refuse an argument instead of ignoring it. urfave/cli passes
      positional arguments through without complaint, so without this `spg
      ~/work` would silently open the working directory and look like it had
      honoured the argument. The message names `--path` for opening one project
      and `scandirs.include` for what the picker searches
- [x] 1.3 Drop the nil-config guard that existed only for this branch. It was
      added because arguments suppressed a configuration parse error, which made
      `spg --config broken.yml somedir` dereference a nil config. With the branch
      gone a broken configuration always errors before anything reads it
- [x] 1.4 Update `README.md`: the usage line `spg [ <directories...> ]` and the
      paragraph saying arguments override `scandirs.include`

## 2. Verification

- [x] 2.1 Replace `TestRunAppTakesArgumentsOverTheConfig` in `src/main_test.go`
      with one asserting an argument is refused, that the error names `--path`,
      and that nothing was scanned
- [x] 2.2 Verify a broken configuration now errors with an argument present,
      which is the case that used to segfault
- [x] 2.3 Verify `--path`, `--view` and `--change-fields` are unaffected, by the
      tests that already cover them passing unchanged
- [x] 2.4 `nix flake check` passes, coverage floors included
