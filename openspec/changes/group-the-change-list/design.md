## Context

The change list is one table with a three-state filter in front of it. The
filter was the least visible state in the application:

```
  q quit  jk/↑↓ navigate  ←→/1-3 tabs  ⏎ view  / search  f mode:archived
                                                          ^^^^^^^^^^^^^^
                                              the whole indication of state
```

Grouping says the same thing by position, permanently, and for both groups at
once. That removes the filter rather than improving it.

## Decisions

### The groups replace the modes rather than joining them

A fourth mode would have left the three that made the state invisible. Removing
them takes a shipped key, a shipped flag and a shipped config setting with it,
which is the cost:

```
  gone                              what says it now
  ────                              ────────────────
  f cycles active/archived/both     the group a row sits in
  change_mode: in config.yml        nothing to configure
  --change-mode on the command line nothing to pass
  emptyListMessage(mode)            a group header counting zero
```

The price is that a project with a large archive always renders it. Active is
the first group, so nothing is ever scrolled past to reach the work in hand, and
`/` still filters. That is a real cost, accepted rather than mitigated.

### A dead configuration key is reported

`yaml.Unmarshal` ignores keys it does not know, which was verified rather than
assumed: a config carrying `change_mode` parses without error and the key
vanishes. So removing the field silently converts a working setting into a line
that does nothing.

The requirement being deleted here is the one that says the opposite:

> An unrecognised mode ... the application SHALL report the unknown value
> together with the valid ones, rather than falling back silently

Honouring that on the way out costs one line. The command-line side needs
nothing: urfave/cli already rejects an unknown flag, which is the same treatment
`--zoom` got when it was removed, and it gets the same scenario.

### Two fixed orders instead of a sort

What sorting was wanted for was seeing recent work first. The archive is
currently ordered oldest first, by accident: `os.ReadDir` sorts by the
date-prefixed directory name and nothing re-sorts it.

```
  today                                   after
  row  1  2026-03-31-fix-openspec...      row  1  2026-09-21-list-the-repos...
  row 35  2026-09-21-list-the-repos...    row 35  2026-03-31-fix-openspec...
```

Reversing that group is not a sorting feature; it is choosing the other fixed
order. A real sort would need a sort key per field, because the rendered value
is wrong to sort on: `tasks` produces `9/12` and `10/12`, and string order puts
ten before nine. That belongs with `specgetty-vru8`, where the columns that
would want sorting are already planned.

### The columns do not change between groups

The date column already renders blank for an active change, so one header row
and one column set serve both groups:

```go
"date": { value: func(r changeRow) string {
    if r.ci.ArchiveDate.IsZero() { return "" }
    return r.ci.ArchiveDate.Format("2006-01-02")
}}
```

That is what makes this one table rather than two stacked. The `archived` state
column becomes redundant against the group header and leaves the defaults; it
stays in the registry for anyone who wants it back.

### The cursor stays on changes, and only the arithmetic learns about headers

`renderTable` is shared with the project picker and assumes one entry is one
line:

```go
offset := 0
if cursor >= bodyHeight {
    offset = cursor - bodyHeight + 1      // cursor indexes rows
}
```

Two ways to break that assumption. Either the cursor indexes display entries and
learns to skip headers, or it keeps indexing changes and the offset is computed
in display lines. The second is the smaller change: the selection key, the search
filter, `g` and `G` and all three actions already work on changes and none of
them should have to know a header exists.

The grouped renderer is the change list's own, reusing `layoutFields` and
`fitCell` rather than teaching `renderTable` about groups. The picker was
rewritten two changes ago and has no use for headers.

## Risks

### The scroll arithmetic is where this goes wrong

Every off-by-one in this change lives in one place: the mapping from a change's
index to the line it is drawn on. It is invisible until the list is longer than
the panel and the cursor is near a group boundary, which is exactly the case a
test has to construct deliberately rather than stumble into.

```
  index  line   what is drawn
    -     0     ACTIVE (2)
    0     1     add-thing
    1     2     fix-other
    -     3     ARCHIVED (35)
    2     4     list-the-repos-not-the-stores
```

### A zero group looks like a bug until you read it

`ACTIVE (0)` with nothing under it is deliberate: "nothing in flight" is an
answer, and an absent header would be indistinguishable from a filter having
hidden it. It will still read as odd the first time. The count is what makes it
legible, which is why the header carries one at all.

### Removing a setting someone may be using

`change_mode` and `--change-mode` are in a released version and documented in
`src/config.yml`. Anyone who set the config key gets one report and then a list
that ignores it; anyone who scripted the flag gets a hard failure from the flag
parser. The flag failing loudly is correct. The config being reported rather
than swallowed is the part that needed deciding, and it was.

### Rollback

Its own commit. The modes, the key, the flag and the config setting return
together because they are one feature. Nothing new is written to a user's
configuration, so a revert cannot strand one.
