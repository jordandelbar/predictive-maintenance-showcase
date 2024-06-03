package app

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.Env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
