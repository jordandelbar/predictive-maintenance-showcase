package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ml_facade/internal/models/postgres_models"
	"net/http"
	"testing"
)

// Check if the healthcheck route works normally
func TestHealthcheck(t *testing.T) {
	// Arrange
	url := fmt.Sprintf("http://localhost:%d/health", testCfg.Port)

	// Act
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer closeOrLog(resp.Body)

	// Assert
	checkStatus(t, resp.StatusCode, http.StatusOK)
}

// Check if the routes return a MethodNotAllowed if wrong request is sent
func TestWrongMethod(t *testing.T) {
	// Arrange
	url1 := fmt.Sprintf("http://localhost:%d/v1/predict", testCfg.Port)
	url2 := fmt.Sprintf("http://localhost:%d/health", testCfg.Port)

	// Act
	resp1, err := http.Get(url1)
	if err != nil {
		t.Fatal(err)
	}
	resp2, err := http.Post(url2, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer closeOrLog(resp1.Body)
	defer closeOrLog(resp2.Body)

	// Assert
	checkStatus(t, resp1.StatusCode, http.StatusMethodNotAllowed)
	checkStatus(t, resp2.StatusCode, http.StatusMethodNotAllowed)
}

// Check if request to a wrong route send the StatusNotFound response
func TestWrongRoute(t *testing.T) {
	// Arrange
	url := fmt.Sprintf("http://localhost:%d/v1/wrongroute", testCfg.Port)

	// Act
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer closeOrLog(resp.Body)

	// Assert
	checkStatus(t, resp.StatusCode, http.StatusNotFound)
}

func TestPredictRoute(t *testing.T) {
	// Arrange
	url := fmt.Sprintf("http://localhost:%d/v1/predict", testCfg.Port)
	sensorData := postgres_models.Sensor{
		MachineID: machineID,
		Sensor00:  0.1,
		Sensor01:  0.2,
		Sensor02:  0.3,
		Sensor03:  0.4,
		Sensor04:  0.5,
		Sensor05:  0.6,
		Sensor06:  0.7,
		Sensor07:  0.8,
		Sensor08:  0.9,
		Sensor09:  1.0,
		Sensor10:  1.1,
		Sensor11:  1.2,
		Sensor12:  1.3,
		Sensor13:  1.4,
		Sensor14:  1.5,
		Sensor15:  1.6,
		Sensor16:  1.7,
		Sensor17:  1.8,
		Sensor18:  1.9,
		Sensor19:  2.0,
		Sensor20:  2.1,
		Sensor21:  2.2,
		Sensor22:  2.3,
		Sensor23:  2.4,
		Sensor24:  2.5,
		Sensor25:  2.6,
		Sensor26:  2.7,
		Sensor27:  2.8,
		Sensor28:  2.9,
		Sensor29:  3.0,
		Sensor30:  3.1,
		Sensor31:  3.2,
		Sensor32:  3.3,
		Sensor33:  3.4,
		Sensor34:  3.5,
		Sensor35:  3.6,
		Sensor36:  3.7,
		Sensor37:  3.8,
		Sensor38:  3.9,
		Sensor39:  4.0,
		Sensor40:  4.1,
		Sensor41:  4.2,
		Sensor42:  4.3,
		Sensor43:  4.4,
		Sensor44:  4.5,
		Sensor45:  4.6,
		Sensor46:  4.7,
		Sensor47:  4.8,
		Sensor48:  4.9,
		Sensor49:  5.0,
		Sensor50:  5.1,
		Sensor51:  5.2,
	}
	jsonData, err := json.Marshal([]postgres_models.Sensor{sensorData})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Act
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	defer closeOrLog(resp.Body)

	// Assert
	checkStatus(t, resp.StatusCode, http.StatusCreated)
}
