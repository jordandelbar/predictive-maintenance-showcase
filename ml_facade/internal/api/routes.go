package api

import (
	"expvar"
	"net/http"
)

func (a *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", a.methodCheck(a.healthcheckHandler, http.MethodGet))
	mux.HandleFunc("/v1/predict", a.methodCheck(a.predictHandler, http.MethodPost))
	mux.HandleFunc("/v1/threshold", a.methodCheck(a.thresholdHandler, http.MethodPost))
	mux.HandleFunc("/metrics", a.methodCheck(expvar.Handler().ServeHTTP, http.MethodGet))

	mux.HandleFunc("/", a.notFoundResponse)

	return a.recoverPanic(a.rateLimit(mux))
}
