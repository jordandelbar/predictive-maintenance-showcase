package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *Server) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/health", a.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/predict", a.predictHandler)
	router.HandlerFunc(http.MethodPost, "/v1/threshold", a.thresholdHandler)

	return a.recoverPanic(a.rateLimit(router))
}
