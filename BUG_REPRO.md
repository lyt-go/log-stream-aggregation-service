# Cancelled request state poisons a later dispatch

## Bug

A request state returned to the pool keeps its cancelled context. The worker also blocks the route and the dispatcher caches the cancellation, so a later request with a fresh context returns the old error and no receipt.

## Trigger

Dispatch `old payload` on `edge-west` with an already cancelled context. Then dispatch `fresh payload` on the same route with a new background context.

## Error

```text
fresh dispatch after cancellation returned context canceled
fresh dispatch receipt = "", want accepted:fresh payload
```
