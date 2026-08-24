# Bug Reproduction

Fail the first flush of `[old]`, append `fresh`, then retry. The retry returns the old error or sends `[fresh]` instead of the original snapshot.

```sh
go test ./internal/flushengine -run '^TestFailedFlushRetriesOriginalSnapshotBeforeNewWindow$' -race -count=1
```

Expected: two sink calls, with the second batch equal to `[old]`.
