## 1. Implementation

- [x] 1.1 Add a single method on the model that returns the content width of the main panel, and give it the comment that currently sits on `docRegion`: this number and the width the panel is drawn at must agree, or the lipgloss box re-wraps the rows and the viewport's reported position stops matching the screen
- [x] 1.2 Replace the three derivations with calls to it: `renderFrame` at ui.go:1219, `docRegion` at docview.go:23, and the log viewport in `recalcLayout` at ui.go:878
- [x] 1.3 Leave `specsSplit` as it is. It already takes the content width as an argument and already carries the same warning; it is a consumer of the new method, not a fourth copy

## 2. Verification

- [x] 2.1 A test asserts the agreement directly rather than by implication: the width `docRegion` reports equals the width `renderFrame` passes to `renderDetailPanel`, for several terminal widths including the 60-column minimum
- [x] 2.2 Every rendered view is byte-for-byte what it was before. A refactor that changes a single column is not a refactor, and the baseline capture in the scratchpad is the reference
- [x] 2.3 Reintroduce one of the three literals and confirm 2.1 fails. Check the build succeeds first: a revert that does not compile makes `go test` print nothing, which reads like a pass
- [x] 2.4 `nix flake check` passes, coverage floors included

## 3. Notes

- [x] 3.1 This is the one part of the padding work worth keeping whatever is decided about how the padding looks. It removes a class of bug that has already been described twice in comments, on `docRegion` and on `specsSplit`, without being fixed either time
