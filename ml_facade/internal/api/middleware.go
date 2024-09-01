package api

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

func (a *Server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *Server) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(a.config.ApiServer.Limiter.Rps), a.config.ApiServer.Limiter.Burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.config.ApiServer.Limiter.Enabled {
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}
			if !limiter.Allow() {
				a.rateLimitExceededResponse(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
