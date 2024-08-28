package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var useRabbitmq bool

const apiUrl = "http://localhost:4000/v1/predict"

var apiHeaders = map[string]string{
	"Content-Type": "application/json",
}

// toFloat converts a string to a float64, returning 0.0 in case of an error
func toFloat(value string) float64 {
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}
	return 0.0
}

// sendPrediction sends the data to the prediction endpoint
func sendPredictionAPI(data SensorData, counter *uint64) {
	listData := make([]SensorDataPayload, 1)
	listData = append(listData, data.SensorDataPayload)
	machineStatus := data.MachineStatus
	jsonData, err := json.Marshal(listData)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for key, value := range apiHeaders {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer closeOrLog(resp.Body)
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime).Milliseconds()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	count := atomic.AddUint64(counter, 1)
	fmt.Printf("%s %d ms %d rows processed, machine status: %s\n", body, elapsedTime, count, machineStatus)
}

func sendPredictionRabbit(ch *amqp.Channel, data SensorData, counter *uint64) {
	rabbitmqQueue := "test"

	// Convert data to JSON
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	// Publish the message
	startTime := time.Now()
	err = ch.Publish(
		"",
		rabbitmqQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Milliseconds()

	count := atomic.AddUint64(counter, 1)
	fmt.Printf("%d ms %d rows processed\n", elapsedTime, count)
}

func worker(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, dataCh <-chan SensorData, counter *uint64, useRabbitmq bool, ch *amqp.Channel) {
	defer wg.Done()
	for {
		select {
		case data, ok := <-dataCh:
			if !ok {
				return
			}

			if err := limiter.Wait(ctx); err != nil {
				log.Println("Error waiting for limiter: ", err)
				continue
			}

			if useRabbitmq {
				sendPredictionRabbit(ch, data, counter)
			} else {
				sendPredictionAPI(data, counter)
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	var rps = RequestPerSecond{}

	flag.BoolVar(&useRabbitmq, "rabbitmq", false, "Use RabbitMQ for sending data")
	flag.IntVar(&rps.rate, "requests", 10, "Requests per second")
	flag.IntVar(&rps.rateBurst, "requests-burst", 2, "Requests per second burst")
	flag.Parse()

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Navigate up the data directory
	basePath := filepath.Join(cwd, "..", "..")
	filePath := filepath.Join(basePath, "data", "sensor.csv")

	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer closeOrLog(csvFile)

	reader := csv.NewReader(csvFile)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Error reading CSV headers: %v", err)
	}

	limiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(rps.rate)), rps.rateBurst)

	var wg sync.WaitGroup
	dataCh := make(chan SensorData, 100)
	var counter uint64

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const numWorkers = 30
	var ch *amqp.Channel
	var conn *amqp.Connection

	if useRabbitmq {
		rabbitmqURI := "amqp://guest:guest@localhost:5672/"
		conn, err = amqp.Dial(rabbitmqURI)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
		defer conn.Close()

		ch, err = conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
		}
		defer ch.Close()

		_, err = ch.QueueDeclare(
			"test", // name
			true,   // durable
			false,  // delete when unused
			false,  // exclusive
			true,   // no-wait
			nil,    // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
		}
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, limiter, dataCh, &counter, useRabbitmq, ch)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV row: %v", err)
		}

		data := SensorData{
			SensorDataPayload: SensorDataPayload{
				MachineID: 7,
				Sensor00:  toFloat(row[2]),
				Sensor01:  toFloat(row[3]),
				Sensor02:  toFloat(row[4]),
				Sensor03:  toFloat(row[5]),
				Sensor04:  toFloat(row[6]),
				Sensor05:  toFloat(row[7]),
				Sensor06:  toFloat(row[8]),
				Sensor07:  toFloat(row[9]),
				Sensor08:  toFloat(row[10]),
				Sensor09:  toFloat(row[11]),
				Sensor10:  toFloat(row[12]),
				Sensor11:  toFloat(row[13]),
				Sensor12:  toFloat(row[14]),
				Sensor13:  toFloat(row[15]),
				Sensor14:  toFloat(row[16]),
				Sensor15:  toFloat(row[17]),
				Sensor16:  toFloat(row[18]),
				Sensor17:  toFloat(row[19]),
				Sensor18:  toFloat(row[20]),
				Sensor19:  toFloat(row[21]),
				Sensor20:  toFloat(row[22]),
				Sensor21:  toFloat(row[23]),
				Sensor22:  toFloat(row[24]),
				Sensor23:  toFloat(row[25]),
				Sensor24:  toFloat(row[26]),
				Sensor25:  toFloat(row[27]),
				Sensor26:  toFloat(row[28]),
				Sensor27:  toFloat(row[29]),
				Sensor28:  toFloat(row[30]),
				Sensor29:  toFloat(row[31]),
				Sensor30:  toFloat(row[32]),
				Sensor31:  toFloat(row[33]),
				Sensor32:  toFloat(row[34]),
				Sensor33:  toFloat(row[35]),
				Sensor34:  toFloat(row[36]),
				Sensor35:  toFloat(row[37]),
				Sensor36:  toFloat(row[38]),
				Sensor37:  toFloat(row[39]),
				Sensor38:  toFloat(row[40]),
				Sensor39:  toFloat(row[41]),
				Sensor40:  toFloat(row[42]),
				Sensor41:  toFloat(row[43]),
				Sensor42:  toFloat(row[44]),
				Sensor43:  toFloat(row[45]),
				Sensor44:  toFloat(row[46]),
				Sensor45:  toFloat(row[47]),
				Sensor46:  toFloat(row[48]),
				Sensor47:  toFloat(row[49]),
				Sensor48:  toFloat(row[50]),
				Sensor49:  toFloat(row[51]),
				Sensor50:  toFloat(row[52]),
				Sensor51:  toFloat(row[53]),
			},
			MachineStatus: row[54],
		}

		dataCh <- data
	}
	close(dataCh)
	wg.Wait()
}
