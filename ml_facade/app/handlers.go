package app

import (
	"bytes"
	"encoding/json"
	"io"
	"ml_facade/internal/data"
	"net/http"
)

func (app *application) predictHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Sensor
	var modelResponse data.ModelResponse

	// Read request body into a buffer
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.logger.Error("error reading request body: %v", err)
		app.errorResponse(w, r, http.StatusBadRequest, "error reading request body")
		return
	}
	defer r.Body.Close()

	// Decode JSON directly from the buffer
	err = json.Unmarshal(body, &input)
	if err != nil {
		app.logger.Error("error decoding JSON from request body: %v", err)
		app.errorResponse(w, r, http.StatusBadRequest, "invalid JSON in request body")
		return
	}

	// Forward request to ML model service
	reqBody := bytes.NewReader(body)
	// TODO: model service in config
	resp, err := app.mlModelClient.Post("http://localhost:3000/predict", "application/json", reqBody)
	if err != nil {
		app.logger.Error("error making POST request to model service: %v", err)
		app.errorResponse(w, r, http.StatusInternalServerError, "error making POST request to model service")
		return
	}
	if resp.StatusCode != http.StatusOK {
		app.logger.Error("wrong status code from model service: %v", resp.Status)
		app.errorResponse(w, r, http.StatusInternalServerError, "error making POST request to model service")
		return
	}
	defer resp.Body.Close()

	// Decode response from model service
	err = json.NewDecoder(resp.Body).Decode(&modelResponse)
	if err != nil {
		app.logger.Error("error decoding response body from model service: %v", err)
		app.errorResponse(w, r, http.StatusInternalServerError, "error decoding response body from model service")
		return
	}

	threshold, err := app.models.Threshold.Get(input.MachineID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var anomaly bool

	if modelResponse.ReconstructionError > threshold {
		anomaly = true
	} else {
		anomaly = false
	}

	// Merge responses
	record := data.Record{
		SensorData:    input,
		ModelResponse: modelResponse,
		Anomaly:       anomaly,
	}

	err = app.models.Sensor.Insert(&record)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write response
	headers := make(http.Header)
	headers.Set("test", "/v1/sensors/")
	response := envelope{"reconstruction_error": record.ModelResponse.ReconstructionError}
	err = app.writeJSON(w, http.StatusCreated, response)
	if err != nil {
		app.logger.Error("error writing JSON response: %v", err)
	}
}

func (app *application) thresholdHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Threshold

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logger.Error("error decoding JSON from request body: %v", err)
		app.errorResponse(w, r, http.StatusBadRequest, "invalid JSON in request body")
		return
	}

	err = app.models.Threshold.Insert(input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"threshold": input.Threshold})
	if err != nil {
		app.logger.Error("error writing JSON response: %v", err)
	}
}
