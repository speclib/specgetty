## 1. Implementation

- [ ] 1.1 Split the specs tab's single border into one per pane, keeping the existing width split and the column of gap between them
- [ ] 1.2 Colour every border by containment: lit when the focused region lies inside it, dim otherwise. With `unify-focus-state` in, each border asks one comparison rather than joining two with `&&`
- [ ] 1.3 Leave the panel border's rule alone in effect while restating it in these terms. It is lit when the keyboard is in the detail area, which is what containment already says about it, so no line of its behaviour changes
- [ ] 1.4 Keep the selected spec's highlight dimming as well. The border says which pane; the highlight says which row, and losing it would make the list look inert when the content has the keys

## 2. Verification

- [ ] 2.1 Count the lit borders in a rendered frame for each stop on the ring and assert the exact set, not just the count: list focused gives panel plus list, content focused gives panel plus content, log focused gives the log alone
- [ ] 2.2 The two specs borders never light together, and never both go dim while the keyboard is on that tab
- [ ] 2.3 The panel border's colour is unchanged from before this change, at every stop on the ring. It is the one thing here that must not move
- [ ] 2.4 The specs tab is still exactly `m.height` lines and `m.width` columns with two borders where there was one, at three widths including the 60-column minimum
- [ ] 2.5 Reintroduce the single border and confirm 2.1 fails. Check the build succeeds first
- [ ] 2.6 `nix flake check` passes, coverage floors included

## 3. Notes

- [ ] 3.1 Independent of the mockup. `box-the-tab-content` stands on its own and this can be reverted without it, which is why the two are not one change
