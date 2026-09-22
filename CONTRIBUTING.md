# Contributing

Issues and pull requests are welcome.

## The gate

Every push and every pull request runs `nix flake check`. It builds the package,
runs `go vet`, runs the suite, and enforces a coverage ratchet whose floors may
only ever rise. It is the same command `scripts/ship-change.sh` runs before a
change is allowed to ship, so what CI asks of you and what shipping asks of you
are the same thing.

Run it yourself with:

```bash
nix flake check
```

## specgetty is written with OpenSpec

Its behaviour lives in `openspec/specs/`, one file per capability, and a change
starts as a proposal under `openspec/changes/` before any code is written. A
change that alters behaviour carries a spec delta saying which requirements it
adds, modifies or removes; a change that does not, such as a documentation move,
says so with `skip_specs: true` rather than inventing a requirement.

```bash
openspec list                  # what is in flight
openspec validate --specs      # are the specs well formed
```

If that sounds like a lot for a small fix, it is not required of you: open a pull
request and it can be written up on the way in.

## Getting help

Ask in an [issue](https://github.com/speclib/specgetty/issues). A bug report that
names the terminal size and what was on screen is worth three that do not.
