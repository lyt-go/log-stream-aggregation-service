# Recovered parser panic leaves ingestion state active

## Bug

After a parser panic is recovered, the in-flight lease remains active, the supervisor remains stopped, and the engine caches the recovered error. The next healthy batch is rejected with stale state and no parsed entry.

## Trigger

Process a `malformed` batch, inspect the active batch count after the recovered error, and immediately process `healthy batch` through the same engine.

## Error

```text
active batches after recovered panic = 1, want 0
healthy batch after recovered panic returned parser panic recovered: malformed log batch
healthy batch entries = [], want one parsed entry
```
