package pollers

import (
	"context"
	"log"
	"time"
)

type runner struct {
	logger   *log.Logger
	poolRate time.Duration
}

func NewRunner(logger *log.Logger, poolRate time.Duration) *runner {
	return &runner{
		logger:   logger,
		poolRate: poolRate,
	}
}

func (r *runner) Run(ctx context.Context, p Poller) error {
	ticker := time.NewTicker(r.poolRate)

	for {
		select {
		case <-ctx.Done():
			r.logger.Println("poller stopped")
			return ctx.Err()
		case <-ticker.C:
			runnerCtx := context.Background()

			err := p.Poll(runnerCtx)
			if err != nil {
				r.logger.Printf("error polling: %v\n", err)
			}
		}
	}
}
