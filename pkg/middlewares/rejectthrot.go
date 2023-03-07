package middlewares

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/atomic"
)

// NewRejectionThrottling создаёт миддлвейер, который при превышении переданного
// RPD (rate per duration, продолжительность кулдауна задаётся) отклоняет
// входящие запросы.
func NewRejectionThrottling(
	ctx context.Context,
	maxRatePerDuration int,
	duration time.Duration,
) func(root http.Handler) http.Handler {
	var (
		ticker  = time.NewTicker(duration)
		counter atomic.Int64
		maxRPD  = int64(maxRatePerDuration)
	)

	go func() {
		for {
			select {
			case _, ok := <-ticker.C:
				if !ok {
					return
				}
				counter.Store(0)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCounterVal := counter.Inc()
			if newCounterVal > maxRPD {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
