package middlewares

import (
	"context"
	"net/http"
	"time"
)

// NewWaitThrottling создаёт миддлвейер, который при превышении переданного
// RPD (rate per duration, продолжительность кулдауна задаётся) заставляет
// входящие запросы ожидать появления места.
func NewWaitThrottling(
	ctx context.Context,
	maxRatePerDuration int,
	duration time.Duration,
) func(root http.Handler) http.Handler {
	var (
		ticker = time.NewTicker(duration)
		passes = make(chan struct{})
		stop   = func() {
			ticker.Stop()
			close(passes)
		}
	)

	go func() {
		for {
			select {
			case _, ok := <-ticker.C:
				if !ok {
					return
				}

				for i := 0; i < maxRatePerDuration; i++ {
					select {
					case passes <- struct{}{}:
					case <-ctx.Done():
						stop()
						return
					}
				}

			case <-ctx.Done():
				stop()
				return
			}
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-ctx.Done():
				return
			case <-r.Context().Done():
				return
			case _, ok := <-passes:
				if !ok {
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
