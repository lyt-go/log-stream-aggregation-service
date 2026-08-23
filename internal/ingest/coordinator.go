package ingest

import "context"

func start(ctx context.Context, jobs []Job) (<-chan Result, <-chan error) {
 out := make(chan Result)
 errs := make(chan error)
 go func() { produce(ctx, jobs, out, errs) }()
 return out, errs
}
