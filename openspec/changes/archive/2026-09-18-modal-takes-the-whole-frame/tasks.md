## 1. Implementation

- [x] 1.1 Settle the question the bean left open: the full-frame modal is the intended design. It is how specgetty has looked since the first commit, and compositing would change how every modal reads. Record it in the spec rather than in a commit message
- [x] 1.2 Drop the `background` parameter from `placeOverlay` and rename it to say what it does, since it builds a frame rather than laying one thing over another
- [x] 1.3 Fix the comment in `View()` that describes confirmation modals sitting on top of the picker. There is no on top: drawing the picker and then a modal discards the picker
- [x] 1.4 Fix the same claim where the tests repeat it, so the next reader is not told twice that overlays composite

## 2. Verification

- [x] 2.1 A test pins the requirement directly: with a confirmation up, the frame holds the modal and none of the view behind it
- [x] 2.2 The existing precedence and frame-size tests still pass, with their comments now matching what they assert
- [x] 2.3 `nix flake check` passes, coverage floors included
