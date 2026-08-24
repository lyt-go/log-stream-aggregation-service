# Temporary delivery failure leaves a phantom commit

## Bug

A temporary sink failure advances the stream cursor and leaves the batch reserved. Retrying the same batch then reports success without storing its log entries.

## Trigger

Create a delivery coordinator, configure the sink to fail the next attempt for a batch, deliver that batch once, and immediately deliver the same batch again. Inspect the checkpoint after the first call and the sink contents after the retry.

## Error

```text
cursor after failed delivery = 42, want 0
stored entries after successful retry = [], want both log entries
```
