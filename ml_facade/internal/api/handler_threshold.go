package api

import (
	"encoding/json"
	"ml_facade/internal/models/redis_models"
	"net/http"
)

// thresholdHandler handles incoming threshold data, inserts it into the database, and returns a success response.
func (a *Server) thresholdHandler(w http.ResponseWriter, r *http.Request) {
	var input redis_models.Threshold

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		a.logger.Error(err.Error())
		a.errorResponse(w, r, http.StatusBadRequest, "invalid JSON in request body")
		return
	}

	err = a.redisModel.Insert(input)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	err = a.writeJSON(w, http.StatusCreated, envelope{"threshold": input.Threshold})
	if err != nil {
		a.logger.Error(err.Error())
	}
}
