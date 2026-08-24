# Buffered frame ownership leaks across pipeline stages

## Bug

Queued frames retain memory returned to a buffer pool, queue views expose pending frames, and archive snapshots expose stored frames. Reuse or caller mutation changes data owned by later pipeline stages.

## Trigger

Submit `alpha` and `bravo` before flushing. Inspect and mutate the pending view, flush it, then mutate an archive snapshot and read the archive again.

## Error

```text
pending frames = ["bravo" "bravo"], want [alpha bravo]
archived frames after pending view mutation = ["Xravo" "Xravo"], want [alpha bravo]
archived frames after snapshot mutation = ["Zravo" "Zravo"], want [alpha bravo]
```
