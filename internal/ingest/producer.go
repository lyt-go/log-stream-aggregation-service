package ingest

import (
 "context"
 "errors"
)

func produce(ctx context.Context, jobs []Job, out chan<- Result, errs chan<- error) {
 for _, job := range jobs {
  if job.Fail { errs <- errors.New("collector rejected job " + job.ID); return }
  select { case out <- Result{ID: job.ID}: case <-ctx.Done(): return }
 }
 close(out)
}
