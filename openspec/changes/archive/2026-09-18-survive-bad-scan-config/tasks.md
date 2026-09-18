## 1. Implementation

- [x] 1.1 Replace the `log.Fatal(err)` in `Walk`'s glob expansion with the handling the walk loop below it already uses: log and skip when `ignore_dir_errors` is set, return the error when it is not
- [x] 1.2 Close `results` on the new early-return path, so a caller ranging over the channel cannot deadlock. `close(results)` currently sits after `errors.Wait()` and is reached on the happy path only; a `defer` at the top of `Walk` covers every return
- [x] 1.3 Make `skip` ignore empty haystack entries, and test the leading slash with `strings.HasPrefix` rather than `f[0:1]`, so the comparison cannot go out of range whatever the entry holds
- [x] 1.4 Guard the glob test in `Walk` the same way: `globPath[len(globPath)-1:]` panics on an empty include, which is the same bug as 1.3 on the other half of the config. Use `strings.HasSuffix` and skip empty includes

## 2. Verification

- [x] 2.1 `Walk` with a glob include into a missing parent, errors ignored: returns nil, still reports the projects under the other includes, and does not exit
- [x] 2.2 The same call with errors not ignored: returns the error, and the results channel is closed so the consumer finishes
- [x] 2.3 `skip` with an empty entry in the haystack returns false instead of panicking
- [x] 2.4 `Walk` with an empty string in the include list completes over the remaining includes
- [x] 2.5 Reintroduce each of the four defects in turn and confirm the matching test fails. Check that the build succeeds first: a revert that does not compile makes `go test` print nothing, which reads like a pass
- [x] 2.6 `nix flake check` passes, coverage floors included
