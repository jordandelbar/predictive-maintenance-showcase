package api

import (
	"fmt"
	"net/http"
)

func (a *Server) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	a.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (a *Server) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}
	err := a.writeJSON(w, status, env)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

func (a *Server) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	a.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (a *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	a.errorResponse(w, r, http.StatusNotFound, message)
}

func (a *Server) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this ressource", r.Method)
	a.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (a *Server) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	a.errorResponse(w, r, http.StatusTooManyRequests, message)
}
