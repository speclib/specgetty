---
# specgetty-nu5a
title: move the program to cmd/spg so go install names it spg
status: todo
type: task
priority: normal
created_at: 2026-09-22T16:03:11Z
updated_at: 2026-09-22T16:03:11Z
parent: specgetty-ab9u
---

`main` lives at `src/`, so the package path ends in `src` and every tool that
builds specgetty has to rename the binary afterwards. Three already do:

| build path  | what it does about the name              |
| ----------- | ----------------------------------------- |
| `Makefile`  | `go build -o spg ./src`                   |
| goreleaser  | `main: ./src` with `binary: spg`          |
| `package.nix` | `mv $out/bin/src $out/bin/spg`          |
| `go install` | installs `src`, and cannot be told otherwise |

`go install` is the only one that cannot, because it takes the name from the
last element of the package path. So the README has to say:

```bash
go install github.com/mipmip/specgetty/src@latest
mv "$(go env GOPATH)/bin/src" "$(go env GOPATH)/bin/spg"
```

Moving the program to `cmd/spg` makes that

```bash
go install github.com/mipmip/specgetty/cmd/spg@latest
```

and lets `package.nix` drop its `postInstall` rename. Three workarounds for one
directory name.

Found while shipping `restructure-the-readme`, where every install path was run
before being written down.

## What moves

`//go:embed` is the part that makes this more than a `git mv`: `src/main.go`
embeds `config.yml` and `VERSION`, and an embed can only reach files beside the
source or below it. So those two travel with it.

```
  src/main.go          ->  cmd/spg/main.go
  src/main_test.go     ->  cmd/spg/main_test.go
  src/changelog_test.go->  cmd/spg/changelog_test.go
  src/testdata/        ->  cmd/spg/testdata/        (changelog_test.go reads it)
  src/config.yml       ->  cmd/spg/config.yml       (embedded)
  src/VERSION          ->  cmd/spg/VERSION          (embedded)

  src/ui/  src/scanner/  src/watcher/               stay where they are
```

## What has to be updated with it

| file                        | what it says today                       |
| --------------------------- | ----------------------------------------- |
| `Makefile`                  | `go build -o spg ./src`                   |
| `package.nix`               | `subPackages = [ "src" ]`, and the rename |
| `.goreleaser-linux.yaml`    | `main: ./src`                             |
| `.goreleaser-darwin.yaml`   | `main: ./src`                             |
| `flake.nix`                 | `readFile ./src/VERSION`                  |
| `scripts/release.sh`        | `VERSION_FILE="src/VERSION"`, read, written and `git add`ed |
| `scripts/record-demo.sh`    | `$root/src/VERSION` and `go build ./src`  |
| `scripts/coverage-gate.sh`  | the floor keyed on `.../src`              |
| `README.md`                 | the install command, and `[config.yml](src/config.yml)` |

## The sharp edge

`scripts/coverage-gate.sh` keys its floors on the full package path, and a
package it has no floor for is not failed:

```sh
if [ -z "$floor" ]; then
  printf '  %-45s %6s%%   (no floor set, new package)\n' "$label" "$actual"
  return
fi
```

So renaming the package without renaming its floor key does not fail the gate.
It silently stops enforcing coverage for that package, and the run still says
`Coverage gate passed`. The floor rename belongs in the same commit, and the
verification is to read the gate's output and see a floor beside `cmd/spg`
rather than the words `new package`.

## No spec impact

No spec names a path under `src/`. `release-process` describes what a release
does, not where the version file lives, and the rest describe the software's
behaviour. This is a pure refactor: `skip_specs: true` rather than a requirement
invented to justify it.

## Two questions worth settling at the same time

**The module path.** It is `github.com/mipmip/specgetty` while the repository is
`speclib/specgetty`. It resolves through a GitHub redirect, so it works, but the
install command names one org and the page the reader is on names another. If it
is going to change, the moment the install command changes anyway is the cheapest
one.

**Whether `src/ui`, `src/scanner` and `src/watcher` become `internal/`.** They
are not meant to be imported by anyone, and `internal/` is how Go says so. It is
a larger change than this one and touches every import in the tree, so it is
worth deciding separately, but the coverage floors would be renamed twice if the
two are done apart.

## Verification when it is done

- `go install github.com/mipmip/specgetty/cmd/spg@latest` installs a binary
  called `spg`, run from a clean `GOBIN`
- `nix run github:speclib/specgetty` still prints the version
- `make build` still produces `./spg`
- the coverage gate names a floor for `cmd/spg`, not `new package`
- `scripts/release.sh --dry-run`, or its tests, still find the version file
