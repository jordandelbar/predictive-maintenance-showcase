package api

import (
	"net/http"
)

func (a *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", a.healthcheckHandler)
	mux.HandleFunc("/v1/predict", a.predictHandler)
	mux.HandleFunc("/v1/threshold", a.thresholdHandler)

	return a.recoverPanic(a.rateLimit(mux))
}
