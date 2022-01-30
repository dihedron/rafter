package run

import (
	"context"
	"time"

	"github.com/dihedron/rafter/logging"
)

const LoopInterval time.Duration = 2 * time.Second

func LeaderRoutine(ctx context.Context, logger logging.Logger, done chan<- bool) {
	ticker := time.NewTicker(LoopInterval)
	logger.Info("LEADER: background checker started ticking every %+v ms", LoopInterval)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			logger.Info("LEADER: done, exiting")
			done <- true
			break loop
		case <-ticker.C:
			logger.Info("LEADER: woken up")
		}
	}
}
