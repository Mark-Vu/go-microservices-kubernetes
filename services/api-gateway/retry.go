package main

import (
	"context"
	"log"
	"time"
)

// InitgRPCServiceWithRetry keeps calling initFn until it succeeds or ctx is cancelled.
// - name is only for logs
// - interval controls how often to retry (simple fixed-interval retry)
func InitgRPCServiceWithRetry[T any](
	ctx context.Context,
	name string,
	interval time.Duration,
	initFn func(context.Context) (T, error),
) (T, error) {
	var zero T

	for {
		v, err := initFn(ctx)
		if err == nil {
			log.Printf("%s client initialized", name)
			return v, nil
		}

		log.Printf("%s init failed (will retry): %v", name, err)

		select {
		case <-ctx.Done():
			log.Printf("%s init cancelled", name)
			return zero, ctx.Err()
		case <-time.After(interval):
		}
	}
}
