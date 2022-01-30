package run

import (
	"context"
	"time"

	"github.com/dihedron/rafter/logging"
)

func FollowerRoutine(ctx context.Context, logger logging.Logger, done chan<- bool) {
	ticker := time.NewTicker(LoopInterval)
	logger.Info("FOLLOWER: background checker started ticking every %+v ms", LoopInterval)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			logger.Info("FOLLOWER: done, exiting")
			done <- true
			break loop
		case <-ticker.C:
			logger.Info("FOLLOWER: woken up")
		}
	}
}
