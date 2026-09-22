# specgetty

[![Check](https://github.com/speclib/specgetty/actions/workflows/check.yml/badge.svg)](https://github.com/speclib/specgetty/actions/workflows/check.yml)
[![coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/speclib/specgetty/gh-pages/badges/coverage.json)](https://github.com/speclib/specgetty/actions/workflows/check.yml)
![specs](https://raw.githubusercontent.com/speclib/specgetty/gh-pages/badges/number_of_specs.svg)
![requirements](https://raw.githubusercontent.com/speclib/specgetty/gh-pages/badges/number_of_requirements.svg)
![tasks](https://raw.githubusercontent.com/speclib/specgetty/gh-pages/badges/tasks_status.svg)
[![Go](https://img.shields.io/github/go-mod/go-version/speclib/specgetty)](go.mod)
[![MIT](https://img.shields.io/badge/licence-MIT-blue)](LICENSE)

`Specgetty`, spg for short, is a text-mode UI for reviewing
[OpenSpec](https://github.com/Fission-AI/OpenSpec) changes and specifications.
It is built for speed and focus.

![A spec opened as an outline beside a card](demo/recordings/hero.gif)

## Why

Writing software has moved from typing code to reviewing what a machine
proposes. The bottleneck moved with it. There are more specs and more changes to
read than there used to be, and they are read in editors built for writing code
rather than for reading a behaviour contract.

specgetty is built for that reading, so that you can get through more of it.
[OpenSpec](https://github.com/Fission-AI/OpenSpec) is the format it reads: a
project keeps its requirements under `openspec/specs/` and its proposed work
under `openspec/changes/`, as markdown.

## Features

- Fast navigation through changes/specs/archive
- Review specs in focussed card view
- Change deltas marked added, modified, removed
- Task checkboxes ticked off in place
- Every openspec project on your machine, searchable
- Archive, discard and export a change
- Shared spec stores followed
- Unreadable specs reported with the reason and the line

## Install

With nix, which needs nothing installed first:

```bash
nix run github:speclib/specgetty
```

A released binary, from the
[latest release](https://github.com/speclib/specgetty/releases/latest). The
archives are named per version and platform, `spg_<version>_linux_amd64.tar.gz`
and the same for `darwin` and `arm64`:

```bash
gh release download --repo speclib/specgetty --pattern 'spg_*_linux_amd64.tar.gz'
tar xzf spg_*_linux_amd64.tar.gz
```

With Go. The command is `.../src` rather than the repository root, because that
is where the program lives, and it installs a binary named `src`:

```bash
go install github.com/mipmip/specgetty/src@latest
mv "$(go env GOPATH)/bin/src" "$(go env GOPATH)/bin/spg"
```

## Configuration

Copy [config.yml](src/config.yml) to `~/.config/specgetty/config.yml` and edit
to your needs. The path follows the XDG Base Directory Specification: with
`$XDG_CONFIG_HOME` set, the config is read from
`$XDG_CONFIG_HOME/specgetty/config.yml`.

## Running

```bash
spg
```

`spg` opens the OpenSpec project you are standing in. It resolves it from the
working directory and starts immediately, without scanning your disk. If there
is no project there, it offers the project picker.

| Flag                    | Effect                                              |
| ----------------------- | --------------------------------------------------- |
| `--view=single`         | Open the project at the working directory (default) |
| `--view=all`            | Open at the project picker                          |
| `--path <dir>`          | Open that project; implies `--view=single`          |
| `--change-fields=<a,b>` | Choose the change list columns                      |

The command line takes no positional arguments. To choose where the picker
looks, edit `scandirs.include` in the configuration file.

## Documentation

| Page                                    | What is in it                                      |
| --------------------------------------- | -------------------------------------------------- |
| [Keys](docs/keys.md)                    | every key, by where the keyboard is                 |
| [Reading specs](docs/reading-specs.md)  | the three levels, the vocabulary, the report view   |
| [Projects](docs/projects.md)            | the picker, stores, the cache, the properties tab   |
| [The change list](docs/change-list.md)  | searching, columns, grouping, exporting             |
| [The editor](docs/editor.md)            | `E`, and the environment it reads                   |

## Contributing

Issues and pull requests are welcome at
[github.com/speclib/specgetty](https://github.com/speclib/specgetty/issues). Every push runs the same
gate a change is shipped through, `nix flake check`, which builds the package,
vets it, runs the suite and enforces a coverage ratchet.

specgetty is itself written with OpenSpec: its behaviour lives in
`openspec/specs/`, and a change starts as a proposal under `openspec/changes/`.
The badges above count them.

## Related

- [OpenSpec](https://github.com/Fission-AI/OpenSpec), the format specgetty reads
- [openspec.nvim](https://github.com/speclib/openspec.nvim), the same specs in
  neovim, whose keyword vocabulary specgetty follows
- [awesome-openspec](https://github.com/speclib/awesome-openspec), what else
  exists around the format

## Development

```bash
make lint
```

## Licence

MIT. See [LICENSE](LICENSE).
