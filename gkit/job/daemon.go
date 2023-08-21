package job

import (
	"context"
	"time"
)

type Daemon struct {
	Work Work
	Rate time.Duration
}

func (job *Daemon) Run(ctx context.Context) error {
	throttle := time.NewTicker(job.Rate)
	handle := RecoverInterceptor(job.Work)
	for {
		select {
		case <-throttle.C:
			_ = handle(ctx)
		}
	}
}
