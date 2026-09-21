## 1. The glob

- [x] 1.1 Stop adding a glob include to the walk list as a path of its own; contribute only what it expands to. Proven by a fixture whose include ends in `*`, asserting the matched directories are walked and the pattern is not
- [x] 1.2 Walk nothing, and report nothing, for a glob that matches no directory
- [x] 1.3 Leave a plain include walked as the path it is
- [x] 1.4 Complete a scan with `--ignore_dir_errors=false` and a glob include, which fails today. Proven against a fixture, and against this author's own configuration, which has three globs and fails outright
- [x] 1.5 Confirm the existing scenarios about an unreadable glob parent still hold, since they are about the parent rather than the pattern

## 2. Log output

- [x] 2.1 Discard log output while the terminal interface is running. Deleting the redirection instead would send it to stderr, which writes over the frame, so this is a replacement rather than a removal
- [x] 2.2 Leave `--debug` logging exactly as it is. Proven by asserting its output still carries the per-project lines and the two timings
- [x] 2.3 Confirm nothing is written to stderr while the interface runs, which is what a corrupted frame would look like

## 3. The panel

- [x] 3.1 Remove the viewport, the content, the message type, the writer and the shown-once flag from the model
- [x] 3.2 Remove the `l` key and its nav bar hint
- [x] 3.3 Remove `focusLog` from the focus ring, and with it every guard that asks whether the log holds the keyboard
- [x] 3.4 Remove the two-view panel switch, the log panel height, and the log branch of the half-page helper
- [x] 3.5 Simplify the `tab` ring to the two halves of a split tab, with nothing beyond them
- [x] 3.6 Simplify the border rules, which no longer have a third region to dim for
- [x] 3.7 Remove the panel from the frame's vertical layout, so the nav bar sits directly below the main panel

## 4. Verification

- [x] 4.1 Assert the frame is exactly its terminal's rows and columns at three widths including the 60-column minimum, which is where a removed row would show as an off-by-one
- [x] 4.2 Assert `tab` cycles only the two halves on both split tabs, and does nothing on a single-pane tab
- [x] 4.3 Assert `l` does nothing
- [x] 4.4 Assert the paging keys reach a list or a document and have nowhere else to go
- [x] 4.5 Assert exactly one border is lit on a split tab and two on a single-pane tab, with no state in which none is
- [x] 4.6 Open this project with the built binary, confirm no stray output appears over the frame during a scan, and confirm `--debug` still prints what it printed before
- [x] 4.7 Revert task 1.1 and confirm the glob test fails. Check the build succeeds first
- [x] 4.8 `nix flake check` passes, coverage floors included

## 5. Notes

- [x] 5.1 Remove the `l` key from the README's key tables
- [x] 5.2 Create a follow-up bean for giving the three real error producers a home: a watcher that failed to start, an unreadable scan directory, and an unreadable store registry. They were as good as invisible in the panel, but auto-rescan dying in silence deserves better than nothing
