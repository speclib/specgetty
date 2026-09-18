## 1. Implementation

- [x] 1.1 Replace `activeView` and `specsFocus` with one `focus` value over four states: the detail content, the spec list, the spec content, and the log panel. There is one keyboard, and `tab` moves it around a ring; two fields for one position is a product type standing in for a sum type
- [x] 1.2 Delete the paired assignments that keep them in step. `case m.activeView == viewLog: m.activeView = viewDetail; m.specsFocus = specsFocusList` becomes a single assignment, and the reset on leaving the specs tab stops being something to remember
- [x] 1.3 Rewrite the reads. `m.activeView == viewDetail && m.detailTab == tabSpecs && m.specsFocus == specsFocusContent` becomes one comparison, which is also the shape the border colouring will want in `light-the-focused-pane`
- [x] 1.4 Keep `viewDetail` and `viewLog` where `renderPanel` uses them to choose a title. That is a different question from where the keyboard is, and collapsing the two would trade one conflation for another

## 2. Verification

- [x] 2.1 The ring is asserted directly, at every length it has: spec list to spec content to log and round on the specs tab with the log open; spec list to spec content and back with it closed; detail to log and back elsewhere; and `tab` doing nothing on a single-pane tab with the log closed
- [x] 2.2 Leaving the specs tab and returning puts the keyboard back on the spec list, which `openspec/specs/specs-tab/spec.md` already requires and which the paired assignments were there to deliver
- [x] 2.3 Every rendered view is byte-for-byte what it was before, captured across three terminal widths. A refactor that changes a column is not a refactor
- [x] 2.4 Reintroduce a state the old pairing could reach and the new value cannot, and confirm it no longer compiles rather than no longer happens. That is the point of the collapse
- [x] 2.5 `nix flake check` passes, coverage floors included

## 3. Notes

- [x] 3.1 Same shape as `unify-content-width`: a fact written in two places, kept in agreement by hand, with the agreement enforced nowhere. That one was three copies of a width; this one is two fields for one position
