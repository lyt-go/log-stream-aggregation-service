# Bug Reproduction

Recover a decoder panic for `batch-old`, then process a healthy batch on the same route. The healthy batch inherits the recovered panic.

```sh
go test ./internal/batchpipeline -run '^TestRecoveredBatchPanicDoesNotPoisonNextCommit$' -race -count=1
```

Expected: `committed:batch-new:decoded:healthy`. Actual: the old worker panic is returned.
