package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"ml_facade/internal/models/postgres_models"
	"net/http"
	"strconv"
	"time"
)

// HandleMlServiceRequest processes the request to the ML service, performs anomaly detection,
// and records the results in the database.
//
// This function takes a request body and its origin, reads and parses the request body,
// forwards the data to an ML service, retrieves a threshold, determines if an anomaly
// has occurred, and records the entire transaction.
func (m *MlService) HandleMlServiceRequest(body any, origin string) (postgres_models.MlServiceResponse, int, error) {
	defer m.wg.Done()

	// Parse the inputs based on the body type
	inputs, err := m.parseInputs(body)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	mlRequestBody, err := m.createMLRequestBody(inputs)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	modelResponse, err := m.forwardRequestToMLService(mlRequestBody)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	// Check that our input length is the same as the reconstruction errors length
	if len(modelResponse.ReconstructionErrors) != len(inputs) {
		return postgres_models.MlServiceResponse{}, 0, fmt.Errorf("mismatch between number of inputs and reconstruction errors")
	}

	anomalies, anomalyCounter, err := m.processAnomalies(inputs, modelResponse)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	err = m.insertRecord(inputs, modelResponse, anomalies, anomalyCounter, origin)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	return modelResponse, anomalyCounter, nil
}

// parseInputs handles the input parsing based on the body type
func (m *MlService) parseInputs(body any) ([]postgres_models.Sensor, error) {
	var inputs []postgres_models.Sensor

	switch v := body.(type) {
	case []amqp.Delivery:
		for _, msg := range v {
			var input postgres_models.Sensor
			err := json.Unmarshal(msg.Body, &input)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, input)
		}
	case io.Reader:
		data, err := io.ReadAll(v)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &inputs)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported body type: %T", body)
	}

	return inputs, nil
}

// createMLRequestBody converts the Sensor struct into the format required by the ML service
func (m *MlService) createMLRequestBody(input []postgres_models.Sensor) ([]byte, error) {
	// Initialize a 2D slice
	inputValues := make([][]float64, len(input))

	for i, sensor := range input {
		inputValues[i] = []float64{
			sensor.Sensor00, sensor.Sensor01, sensor.Sensor02, sensor.Sensor03, sensor.Sensor04,
			sensor.Sensor05, sensor.Sensor06, sensor.Sensor07, sensor.Sensor08, sensor.Sensor09,
			sensor.Sensor10, sensor.Sensor11, sensor.Sensor12, sensor.Sensor13, sensor.Sensor14,
			sensor.Sensor15, sensor.Sensor16, sensor.Sensor17, sensor.Sensor18, sensor.Sensor19,
			sensor.Sensor20, sensor.Sensor21, sensor.Sensor22, sensor.Sensor23, sensor.Sensor24,
			sensor.Sensor25, sensor.Sensor26, sensor.Sensor27, sensor.Sensor28, sensor.Sensor29,
			sensor.Sensor30, sensor.Sensor31, sensor.Sensor32, sensor.Sensor33, sensor.Sensor34,
			sensor.Sensor35, sensor.Sensor36, sensor.Sensor37, sensor.Sensor38, sensor.Sensor39,
			sensor.Sensor40, sensor.Sensor41, sensor.Sensor42, sensor.Sensor43, sensor.Sensor44,
			sensor.Sensor45, sensor.Sensor46, sensor.Sensor47, sensor.Sensor48, sensor.Sensor49,
			sensor.Sensor50, sensor.Sensor51,
		}
	}

	mlRequest := map[string]interface{}{
		"input_values": inputValues,
	}

	return json.Marshal(mlRequest)
}

// forwardRequestToMLService forwards the given request body to the ML service,
// expecting a JSON response containing the model response. It returns the decoded
// model response and any errors encountered during the request or response handling.
func (m *MlService) forwardRequestToMLService(body []byte) (postgres_models.MlServiceResponse, error) {
	var postError = errors.New("error making POST request to model service")
	var encodingError = errors.New("error decoding response body from model service")
	var modelResponse postgres_models.MlServiceResponse

	reqBody := bytes.NewReader(body)
	resp, err := m.client.Post(m.config.MlServiceUri()+"/predict", "application/json", reqBody)
	if err != nil {
		m.logger.Error(err.Error())
		return modelResponse, postError
	}
	if resp.StatusCode != http.StatusOK {
		err := errors.New(resp.Status)
		m.logger.Error(err.Error())
		return modelResponse, postError
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&modelResponse)
	if err != nil {
		m.logger.Error(err.Error())
		return modelResponse, encodingError
	}

	return modelResponse, nil
}

// processAnomalies determines anomalies and counts them
func (m *MlService) processAnomalies(
	inputs []postgres_models.Sensor,
	modelResponse postgres_models.MlServiceResponse,
) ([]bool, int, error) {
	var anomalies []bool
	var anomalyCounter int

	for i, input := range inputs {
		threshold, err := m.fetchOrCacheThreshold(input.MachineID)
		if err != nil {
			return nil, 0, err
		}

		reconstructionError := modelResponse.ReconstructionErrors[i]
		anomaly, counter, err := m.determineAnomaly(input.MachineID, reconstructionError, threshold)
		if err != nil {
			return nil, 0, err
		}

		// We only take the last counter
		anomalyCounter = counter
		anomalies = append(anomalies, anomaly)
	}

	return anomalies, anomalyCounter, nil
}

// fetchOrCacheThreshold retrieves the threshold for the given machineID. It first
// checks the cache. If a valid cached entry is found, it returns the cached value.
// Otherwise, it fetches the threshold from the database, caches it, and returns it.
func (m *MlService) fetchOrCacheThreshold(machineID int) (float64, error) {
	cacheKey := fmt.Sprintf("threshold_%s", strconv.Itoa(machineID))
	entry, found := m.thresholdCache.Load(cacheKey)

	if found {
		cachedEntry := entry.(cacheEntry)
		if time.Now().Before(cachedEntry.expiration) {
			return cachedEntry.value, nil
		}
		m.thresholdCache.Delete(cacheKey)
	}

	threshold, err := m.thresholdModel.Get(machineID)
	if err != nil {
		return 0, err
	}

	m.thresholdCache.Store(cacheKey, cacheEntry{value: threshold, expiration: time.Now().Add(1 * time.Minute)})
	return threshold, nil
}

// determineAnomaly calculates if a given reconstruction error is an anomaly based on the provided threshold.
// It increments or decrements the anomaly counter for the given machineID accordingly.
func (m *MlService) determineAnomaly(machineID int, reconstructionError, threshold float64) (bool, int, error) {
	var anomaly bool
	var anomalyCounter int
	var err error

	if reconstructionError > threshold {
		anomaly = true
		anomalyCounter, err = m.thresholdModel.Increment(machineID)
		if err != nil {
			m.logger.Error(err.Error())
			return false, 0, err
		}
	} else {
		anomaly = false
		anomalyCounter, err = m.thresholdModel.Decrement(machineID)
		if err != nil {
			m.logger.Error(err.Error())
			return false, 0, err
		}
	}

	return anomaly, anomalyCounter, nil
}

// insertRecord inserts a new record into the database, containing sensor data, model response, anomaly flag, and anomaly counter.
func (m *MlService) insertRecord(
	inputs []postgres_models.Sensor,
	modelResponse postgres_models.MlServiceResponse,
	anomalies []bool,
	anomalyCounter int,
	origin string) error {

	if len(inputs) != len(anomalies) {
		return errors.New("mismatch between number of inputs and anomalies")
	}

	// Prepare the records for bulk insert
	records := make([]postgres_models.Record, len(inputs))
	for i, input := range inputs {
		records[i] = postgres_models.Record{
			SensorData:          input,
			ReconstructionError: modelResponse.ReconstructionErrors[i],
			Anomaly:             anomalies[i],
			AnomalyCounter:      anomalyCounter,
			Origin:              origin,
		}
	}

	err := m.sensorModel.Insert(records)
	if err != nil {
		return err
	}

	return nil
}
