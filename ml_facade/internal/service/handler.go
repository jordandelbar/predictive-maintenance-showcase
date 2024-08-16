package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"io"
	"ml_facade/internal/models/postgres_models"
	"net/http"
	"reflect"
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

	var input postgres_models.Sensor
	var reader io.Reader

	reader, err := getReaderFromBody(body)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	err = json.Unmarshal(data, &input)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	mlRequestBody, err := createMLRequestBody(input)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	modelResponse, err := m.forwardRequestToMLService(mlRequestBody)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	threshold, err := m.fetchOrCacheThreshold(input.MachineID)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	anomaly, anomalyCounter, err := m.determineAnomaly(input.MachineID, modelResponse.ReconstructionError, threshold)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	err = m.insertRecord(input, modelResponse, anomaly, anomalyCounter, origin)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	return modelResponse, anomalyCounter, nil
}

func getReaderFromBody(body any) (io.Reader, error) {
	switch v := body.(type) {
	case amqp091.Delivery:
		return bytes.NewReader(v.Body), nil
	case io.Reader:
		return v, nil
	default:
		return nil, fmt.Errorf("body type not supported: %v", reflect.TypeOf(body))
	}
}

// createMLRequestBody converts the Sensor struct into the format required by the ML service
func createMLRequestBody(input postgres_models.Sensor) ([]byte, error) {
	// Collect all the sensor values into a slice
	inputValues := []float64{
		input.Sensor00, input.Sensor01, input.Sensor02, input.Sensor03, input.Sensor04,
		input.Sensor05, input.Sensor06, input.Sensor07, input.Sensor08, input.Sensor09,
		input.Sensor10, input.Sensor11, input.Sensor12, input.Sensor13, input.Sensor14,
		input.Sensor15, input.Sensor16, input.Sensor17, input.Sensor18, input.Sensor19,
		input.Sensor20, input.Sensor21, input.Sensor22, input.Sensor23, input.Sensor24,
		input.Sensor25, input.Sensor26, input.Sensor27, input.Sensor28, input.Sensor29,
		input.Sensor30, input.Sensor31, input.Sensor32, input.Sensor33, input.Sensor34,
		input.Sensor35, input.Sensor36, input.Sensor37, input.Sensor38, input.Sensor39,
		input.Sensor40, input.Sensor41, input.Sensor42, input.Sensor43, input.Sensor44,
		input.Sensor45, input.Sensor46, input.Sensor47, input.Sensor48, input.Sensor49,
		input.Sensor50, input.Sensor51,
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
	input postgres_models.Sensor,
	modelResponse postgres_models.MlServiceResponse,
	anomaly bool,
	anomalyCounter int,
	origin string) error {

	record := postgres_models.Record{
		SensorData:     input,
		ModelResponse:  modelResponse,
		Anomaly:        anomaly,
		AnomalyCounter: anomalyCounter,
		Origin:         origin,
	}

	err := m.sensorModel.Insert(&record)
	if err != nil {
		return err
	}

	return nil
}
