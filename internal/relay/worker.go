package relay

import (
	"context"
	"errors"
)

func deliver(ctx context.Context, jobs []Job, cancel context.CancelFunc) ([]Receipt, error) {
	receipts := make([]Receipt, 0, len(jobs))
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return receipts, ctx.Err()
		default:
		}
		if job.Fail {
			cancel()
			return receipts, errors.New("relay rejected job " + job.ID)
		}
		receipts = append(receipts, Receipt{ID: job.ID})
	}
	return receipts, nil
}
