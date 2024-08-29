package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"ml_facade/internal/models/postgres_models"
)

func TestParseInputs(t *testing.T) {
	// Create an instance of MlService
	m := MlService{}

	// Test Case 1: Valid AMQP Delivery
	t.Run("Valid AMQP Delivery", func(t *testing.T) {
		// Arrange
		sensor := postgres_models.Sensor{MachineID: 123}
		sensorBytes, _ := json.Marshal(sensor)
		msgs := []amqp.Delivery{
			{Body: sensorBytes},
		}

		// Act
		inputs, err := m.parseInputs(msgs)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(inputs) != 1 {
			t.Fatalf("expected 1 input, got %d", len(inputs))
		}
		if inputs[0].MachineID != 123 {
			t.Errorf("expected MachineID to be '123', got '%d'", inputs[0].MachineID)
		}
	})

	// Test Case 2: Invalid AMQP Delivery (malformed JSON)
	t.Run("Invalid AMQP Delivery", func(t *testing.T) {
		// Arrange
		msgs := []amqp.Delivery{
			{Body: []byte(`{invalid json}`)},
		}

		// Act
		inputs, err := m.parseInputs(msgs)

		// Assert
		if err == nil {
			t.Fatal("expected an error, got none")
		}
		if inputs != nil {
			t.Fatalf("expected no inputs, got %v", inputs)
		}
	})

	// Test Case 3: Valid io.Reader
	t.Run("Valid io.Reader", func(t *testing.T) {
		// Arrange
		sensors := []postgres_models.Sensor{
			{MachineID: 123},
			{MachineID: 456},
		}
		sensorBytes, _ := json.Marshal(sensors)
		reader := bytes.NewReader(sensorBytes)

		// Act
		inputs, err := m.parseInputs(reader)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(inputs) != 2 {
			t.Fatalf("expected 2 inputs, got %d", len(inputs))
		}
		if inputs[0].MachineID != 123 {
			t.Errorf("expected MachineID to be 123, got '%d'", inputs[0].MachineID)
		}
		if inputs[1].MachineID != 456 {
			t.Errorf("expected MachineID to be 456, got '%d'", inputs[1].MachineID)
		}
	})

	// Test Case 4: Invalid io.Reader (malformed JSON)
	t.Run("Invalid io.Reader", func(t *testing.T) {
		// Arrange
		reader := bytes.NewReader([]byte(`{invalid json}`))

		// Act
		inputs, err := m.parseInputs(reader)

		// Assert
		if err == nil {
			t.Fatal("expected an error, got none")
		}
		if inputs != nil {
			t.Fatalf("expected no inputs, got %v", inputs)
		}
	})

	// Test Case 5: Unsupported body type
	t.Run("Unsupported Body Type", func(t *testing.T) {
		// Act
		inputs, err := m.parseInputs(123) // Unsupported type

		// Assert
		if err == nil {
			t.Fatal("expected an error, got none")
		}
		expectedError := errors.New("unsupported body type: int")
		if err.Error() != expectedError.Error() {
			t.Errorf("expected error '%v', got '%v'", expectedError, err)
		}
		if inputs != nil {
			t.Fatalf("expected no inputs, got %v", inputs)
		}
	})
}
