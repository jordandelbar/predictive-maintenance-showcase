package api

import (
	"ml_facade/internal/models/postgres_models"
	"net/http"
)

// predictHandler handles incoming sensor data, processes it through the ML model,
// determines anomalies based on reconstruction error and threshold, and stores
// the results in the database. It writes the reconstruction error and anomaly
// counter as a JSON response.
func (a *Server) predictHandler(w http.ResponseWriter, r *http.Request) {
	var modelResponse postgres_models.MlServiceResponse

	a.wg.Add(2)

	modelResponse, anomalyCounter, err := a.service.HandleMlServiceRequest(r.Body, "api")

	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	a.writeResponse(w, modelResponse.ReconstructionErrors, anomalyCounter)
}

// writeResponse writes the reconstruction error and anomaly counter as a JSON response with status code 201 Created.
func (a *Server) writeResponse(w http.ResponseWriter, reconstructionErrors []float64, anomalyCounter int) {
	defer a.wg.Done()

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	response := envelope{"reconstruction_errors": reconstructionErrors, "anomaly_counter": anomalyCounter}
	err := a.writeJSON(w, http.StatusCreated, response)
	if err != nil {
		a.logger.Error(err.Error())
	}
}
