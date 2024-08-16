package api

import (
	"net/http"
)

func (a *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.Env,
			"version":     a.version,
		},
	}

	err := a.writeJSON(w, http.StatusOK, env)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
