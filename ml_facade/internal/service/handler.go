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

	var inputs []postgres_models.Sensor
	//var reader io.Reader

	msgs, ok := body.([]amqp.Delivery)
	if !ok {
		return postgres_models.MlServiceResponse{}, 0, fmt.Errorf("expected []amqp.Delivery, got %T", body)
	}

	for _, msg := range msgs {
		var input postgres_models.Sensor

		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			return postgres_models.MlServiceResponse{}, 0, err
		}

		inputs = append(inputs, input)
	}

	//reader, err := getReaderFromBody(body)
	//if err != nil {
	//	return postgres_models.MlServiceResponse{}, 0, err
	//}
	//
	//data, err := io.ReadAll(reader)
	//if err != nil {
	//	return postgres_models.MlServiceResponse{}, 0, err
	//}
	//
	//err = json.Unmarshal(data, &inputs)
	//if err != nil {
	//	return postgres_models.MlServiceResponse{}, 0, err
	//}

	mlRequestBody, err := createMLRequestBody(inputs)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	modelResponse, err := m.forwardRequestToMLService(mlRequestBody)
	if err != nil {
		return postgres_models.MlServiceResponse{}, 0, err
	}

	var anomalyCounter int
	for _, input := range inputs {
		threshold, err := m.fetchOrCacheThreshold(input.MachineID)
		if err != nil {
			return postgres_models.MlServiceResponse{}, 0, err
		}

		anomaly, counter, err := m.determineAnomaly(input.MachineID, modelResponse.ReconstructionError, threshold)
		if err != nil {
			return postgres_models.MlServiceResponse{}, 0, err
		}
		anomalyCounter += counter

		err = m.insertRecord(input, modelResponse, anomaly, anomalyCounter, origin)
		if err != nil {
			return postgres_models.MlServiceResponse{}, 0, err
		}
	}

	return modelResponse, anomalyCounter, nil
}

func getReaderFromBody(body any) (io.Reader, error) {
	switch v := body.(type) {
	case []amqp.Delivery:
		var combined []byte
		for _, msg := range v {
			combined = append(combined, msg.Body...)
		}
		return bytes.NewReader(combined), nil
	case amqp.Delivery:
		return bytes.NewReader(v.Body), nil
	case io.Reader:
		return v, nil
	default:
		return nil, fmt.Errorf("body type not supported: %v", reflect.TypeOf(body))
	}
}

// createMLRequestBody converts the Sensor struct into the format required by the ML service
func createMLRequestBody(input []postgres_models.Sensor) ([]byte, error) {
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
