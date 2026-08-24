# Bug Reproduction

Cancel a subscription for `alerts-critical`, subscribe to the same topic again, and publish `rule-fired`. The second subscription inherits a closed channel instead of receiving the event.

```sh
go test ./internal/subscription -run '^TestCancellingOneSubscriberDoesNotCloseNextSubscription$' -race -count=1
```

The failure is deterministic: the second read reports `second subscription inherited a closed channel`.
