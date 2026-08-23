package ingest

import (
 "context"
 "fmt"
)

func collect(ctx context.Context, jobs []Job) ([]Result, error) {
 out, errs := start(ctx, jobs)
 results := make([]Result, 0, len(jobs))
 for {
  select {
  case item, ok := <-out:
   if !ok { return results, nil }
   results = append(results, item)
  case <-ctx.Done():
   return results, fmt.Errorf("pipeline stopped: %w", ctx.Err())
  case <-errs:
   continue
  }
 }
}
