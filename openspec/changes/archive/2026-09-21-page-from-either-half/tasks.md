## 1. Implementation

- [x] 1.1 Add a predicate that reports whether a document is displayed, beside the one that reports whether it holds the keyboard. The difference is the focus check on a split tab, and nothing else
- [x] 1.2 Point `pgdown`, `ctrl+f`, `pgup`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` at it
- [x] 1.3 Leave `j` and `k` on the focus-gated predicate, so the list is still walked rather than the document scrolled
- [x] 1.4 Leave the reading position on the focus-gated predicate. `specs-tab` reports it only while the content holds the keyboard, which is deliberate and was not what broke

## 2. Verification

- [x] 2.1 Page and jump on the properties tab with the row list holding the keyboard, and assert the content moved. This is the reported defect
- [x] 2.2 The same on the specs tab, which has under-delivered against `document-viewer` since it shipped
- [x] 2.3 `j` and `k` move the list and leave the document where it is, on both tabs
- [x] 2.4 The log panel still takes these keys when it holds the keyboard
- [x] 2.5 An open change is unaffected, since it was never a split
- [x] 2.6 The panel title still reports no position while a list holds the keyboard
- [x] 2.7 Revert task 1.2 and confirm the properties test fails. Confirmed twice: the defect was measured before the fix (pgdown left the offset at 0 on both split tabs, while an open change scrolled), and the reverted build failed `TestPagingWorksWhileTheListHoldsTheKeyboard` on both subtests
- [x] 2.8 `nix flake check` passes, coverage floors included
