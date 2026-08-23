package ingest

import "context"

type Job struct { ID string; Fail bool }
type Result struct { ID string }
type Pipeline struct { jobs []Job }

func New(jobs []Job) *Pipeline { return &Pipeline{jobs: append([]Job(nil), jobs...)} }
func (p *Pipeline) Run(ctx context.Context) ([]Result, error) { return collect(ctx, p.jobs) }
