# Bug Reproduction

Resolve a typed-nil decoder, recover its panic, register a healthy replacement on the same route, and decode again. The replacement remains poisoned by the old failure.

```sh
go test ./internal/decodepipeline -run '^TestTypedNilDecoderPanicDoesNotPoisonReplacement$' -race -count=1
```

Expected: `decoded:fresh`. Actual: the old decoder panic is returned.
