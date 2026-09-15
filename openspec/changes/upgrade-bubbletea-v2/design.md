## Context

Explored twice, on 2026-09-15, first as `specgetty-p4e2` and again after
coverage work. Everything below was measured or read from the v2 sources rather
than inferred from the upgrade guide, which is summarised and incomplete on at
least one important point.

### The trap the upgrade guide does not mention

The v2 modules moved off GitHub import paths:

    github.com/charmbracelet/bubbletea/v2  ->  charm.land/bubbletea/v2
    github.com/charmbracelet/bubbles/v2    ->  charm.land/bubbles/v2
    github.com/charmbracelet/lipgloss/v2   ->  charm.land/lipgloss/v2

`go get` on the old path fails with "module declares its path as:
charm.land/...". Versions at the time of writing: bubbletea v2.0.9,
bubbles v2.2.1, lipgloss v2.0.6.

### The surface is smaller than the raw counts suggest

83 `tea.KeyMsg{}` literals sounds alarming. Every one of them is in a test.
Production code was checked for the things that actually break:

| Looked for                        | Found |
| --------------------------------- | ----- |
| `msg.Type`, `msg.Runes`, `msg.Alt` | none; only `msg.String()` is used |
| a binding for the space key        | none; v2's `String()` change for space is moot |
| mouse or focus reporting           | none; the removed program options do not apply |

So the struct-to-interface change for `KeyMsg` touches one type switch, not a
scattering of field accesses, and `key := msg.String()` with its `case "y":`
arms survives untouched.

### What the compiler will catch, and what it will not

| Change | Caught by |
| ------ | --------- |
| `YOffset` field becomes a method (31 uses) | compiler |
| `ViewDown`/`ViewUp` -> `PageDown`/`PageUp` | compiler |
| `LineDown`/`LineUp` -> `ScrollDown`/`ScrollUp` | compiler |
| `View() string` -> `View() tea.View` | compiler, for the signature |
| whether the returned `tea.View` is *correct* | **nothing** |
| whether the program still starts and draws | **nothing here** |

Most of this fails loudly. The part that does not is the shape of the value
`View` now returns, and the program lifecycle in `Run`.

## Goals / Non-Goals

**Goals**

- Behaviour identical before and after, including the alt screen
- The `withColor` helper deleted rather than ported
- No new dependency; three replaced

**Non-Goals**

- Light and dark colour adaptation. The upgrade makes it possible; choosing
  the colour pairs is a separate change.
- Mouse support. v2 moves mouse configuration into `View.MouseMode`, but any
  mouse reporting mode still takes native text selection away from the terminal.
  That is terminal behaviour, not a bubbletea choice, and the reason
  `scroll-markdown-documents` declined it is unchanged.
- `WindowTitle` and the native `ProgressBar` field. Both are a good fit for this
  tool and neither justifies the move.

## Decisions

### `tea.SetClipboard` is a trap, not a simplification

v2 ships a clipboard API. It is OSC 52, and it is fire and forget:

    case setClipboardMsg:
        p.execute(ansi.SetSystemClipboard(string(msg)))

No error return, no acknowledgement. Its own doc says OSC 52 is not supported in
all terminals.

`copy-change-name-and-path` requires that a copy which did not happen is
distinguishable from one that did. `atotto/clipboard` can do that: `WriteAll`
returns an error and `Unsupported` is a package flag. OSC 52 cannot, in
principle. Switching that feature to the built-in API would make one of its
requirements unimplementable and would report success while Ghostty silently
refused the write.

So after this upgrade lands, the clipboard feature keeps `atotto`. The temptation
to drop a dependency in favour of something built in is the reason this is
written down rather than left to judgement on the day.

### `View()` is the risk, and it is the reason for the ordering

`View` had no coverage at all until `specgetty-2pub`, and it is the function
whose signature changes. It is now at 100%, which is why this upgrade is
attemptable with a straight face.

The assertions there are structural: a modal appears when its state is set, the
frame is exactly the terminal height and never wider, overlays do not push the
frame out of shape. All of that holds across the `tea.View` change and will fail
if the new value is assembled wrongly.

It does not cover the alt screen, which is the one part of `View`'s new
responsibilities that has no test and cannot have one here.

### Colour profile handling inverts

`lipgloss.SetColorProfile` is removed, and `Style.Render` in v2 contains no
`isatty` and no writer check: it always emits, and downgrading happens at the
output writer. That is why the `withColor` helper, at 10 call sites across three
test files, is deleted rather than ported. Tests that force a colour profile to
avoid being vacuous stop needing to.

`lipgloss.Color("236")` still parses ANSI indices, and `Background` and
`Foreground` take the `color.Color` it returns, so the style definitions carry
over unchanged.

### bubbles is compatible where specgetty touches it

`textinput` keeps `New`, `Value`, `SetValue`, `Focus`, `Blur`, `Prompt` and
`Placeholder`. `spinner` keeps `New`, `Spinner`, `View` and `TickMsg`. Only
`viewport` renames things, and all of those are compile errors.

## Risks / Trade-offs

- **The alt screen is unverifiable from a test.** `WithAltScreen()` becomes a
  field on the value `View` returns. If it is missed, specgetty draws over the
  scrollback instead of taking the alternate screen, and no test can tell.
  Someone has to run it.
- **`Run` has no coverage and needs a TTY.** Program construction and teardown
  are 15 statements that only a human can exercise.
- **`charm.land` is a vanity import path.** A single point of failure that
  `github.com/` is not. Probably fine, worth one moment's thought.
- **Two vendor hashes.** `package.nix` and `flake.nix` both declare one, and a
  dependency change moves both. `scripts/release.sh` was fixed to update both;
  this change has to do the same by hand.

## Open Questions

- Whether to take `WindowTitle` in a follow-up, so the terminal tab shows the
  open project. Cheap once v2 is in, and out of scope here.
