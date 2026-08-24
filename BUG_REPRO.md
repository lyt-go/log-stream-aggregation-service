# Bug Reproduction

Dispatch a cancelled envelope to `edge-east`, then dispatch a fresh envelope to the same route. The second request inherits the old cancellation instead of returning its own receipt.

```sh
go test ./internal/fanout -run '^TestCancelledDeliveryDoesNotPoisonNextEnvelope$' -race -count=1
```

Expected: `delivered:msg-new:fresh`. Actual: the second call returns `context canceled`.
