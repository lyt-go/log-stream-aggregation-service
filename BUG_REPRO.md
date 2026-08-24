# Recovered decoder panic poisons its route

## Bug

With telemetry disabled, recovering a decoder panic triggers another panic through a typed-nil reporter. The route remains quarantined, so its next healthy frame is rejected without a decoded result.

## Trigger

Create a decoding runner with reporting disabled. Process a corrupt frame on a route, catch the escaped panic, and then process a healthy frame on the same route.

## Error

```text
corrupt frame escaped recovery with panic: runtime error: invalid memory address or nil pointer dereference
healthy frame after recovered failure returned decode route quarantined
healthy frame entries = [], want one decoded entry
```
