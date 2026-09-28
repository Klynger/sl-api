package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"

	errorsCommon "sl-api/api/routes/common/errors"
)

const (
	loginRateLimit  = 10
	loginRateWindow = time.Minute
)

// LoginRateLimiter throttles login attempts per client IP. It caps brute-force
// guessing and, just as importantly, sheds floods before they reach the bcrypt
// comparison, which is deliberately expensive and would otherwise be a cheap way
// to exhaust CPU. Counters are kept in memory, which is correct for a single
// instance; running more than one instance needs a shared store.
//
// Keying is by connection remote address, which is correct when nothing sits in
// front of the app. Behind a proxy or load balancer this must switch to a
// trusted X-Forwarded-For / X-Real-IP value (see the rate-limiting plan).
func LoginRateLimiter() func(http.Handler) http.Handler {
	return httprate.Limit(
		loginRateLimit,
		loginRateWindow,
		httprate.WithKeyByIP(),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			errorsCommon.TooManyRequests(w, errorsCommon.NewError(
				errorsCommon.CodeRateLimited,
				"too many requests, please try again later",
			))
		}),
	)
}
