package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ml_facade/cmd/app"
	"ml_facade/config"
	"ml_facade/internal/models/postgres_models"
	"ml_facade/internal/models/redis_models"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	machineID        = 1234
	postgresUsername = "monitor"
	postgresPassword = "test"
	postgresDatabase = "monitoring"
)

type MockResponse struct {
	Message string `json:"message"`
}

var testCfg config.Config

func terminateOrLog(container testcontainers.Container, ctx context.Context) {
	err := container.Terminate(ctx)
	if err != nil {
		log.Println("Failed to terminate:", err)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup PostgreSQL container
	postgresContainer, err := setupPostgresContainer(ctx)
	if err != nil {
		fmt.Println("Failed to start PostgreSQL container:", err)
		os.Exit(1)
	}
	defer terminateOrLog(postgresContainer, ctx)

	postgresHost, err := postgresContainer.Host(ctx)
	if err != nil {
		fmt.Println("Failed to get PostgreSQL container host:", err)
		os.Exit(1)
	}
	postgresPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Println("Failed to get PostgreSQL container port:", err)
		os.Exit(1)
	}

	// Setup Redis container
	redisContainer, err := setupRedisContainer(ctx)
	if err != nil {
		fmt.Println("Failed to start Redis container:", err)
		os.Exit(1)
	}
	defer terminateOrLog(redisContainer, ctx)

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		fmt.Println("Failed to get Redis container host:", err)
		os.Exit(1)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		fmt.Println("Failed to get Redis container port:", err)
		os.Exit(1)
	}

	mockMlService, mlServiceHost, mlServicePort := startMockMLService()
	defer mockMlService.Close()

	testCfg.ApiServer.Port, err = findAvailablePort()
	if err != nil {
		fmt.Println(err)
	}
	testCfg.Env = "test"
	testCfg.PostgresDB.Host = postgresHost
	testCfg.PostgresDB.Port = postgresPort.Port()
	testCfg.PostgresDB.Username = postgresUsername
	testCfg.PostgresDB.Password = postgresPassword
	testCfg.PostgresDB.DatabaseName = postgresDatabase
	testCfg.RedisDB.Host = redisHost
	testCfg.RedisDB.Port = redisPort.Port()
	testCfg.MlService.Host = mlServiceHost
	testCfg.MlService.Port = mlServicePort
	testCfg.PostgresDB.MaxOpenConns = 25
	testCfg.PostgresDB.MaxIdleConns = 25
	testCfg.PostgresDB.MaxIdleTime = 5 * time.Minute

	go app.StartApp(testCfg)

	if !waitForAPI(fmt.Sprintf("http://localhost:%d/health", testCfg.ApiServer.Port), 30, 1*time.Second) {
		fmt.Println("API did not start in time")
		os.Exit(1)
	}

	InitializeRedisThreshold(testCfg)
	ApplyMigration(testCfg)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func startMockMLService() (*httptest.Server, string, string) {
	handler := http.NewServeMux()
	handler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := MockResponse{Message: "Hello, I am a mocked server!"}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Fatal(err)
		}
	})
	handler.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		response := postgres_models.MlServiceResponse{ReconstructionErrors: []float64{0.2}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Fatal(err)
		}
	})

	server := httptest.NewServer(handler)
	// Parse server URL
	parsedURL, err := url.Parse(server.URL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse server URL: %v", err))
	}

	// Extract host and port
	host, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		panic(fmt.Sprintf("failed to split host and port: %v", err))
	}

	return server, host, port
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer closeOrLog(listener)

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func waitForAPI(url string, maxRetries int, retryInterval time.Duration) bool {
	time.Sleep(retryInterval)
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		fmt.Println(resp)
		if err == nil && resp.StatusCode == http.StatusOK {
			closeOrLog(resp.Body)
			return true
		}
		time.Sleep(retryInterval)
	}
	return false
}

func setupPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     postgresUsername,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDatabase,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func setupRedisContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func InitializeRedisThreshold(config config.Config) {
	threshold := redis_models.Threshold{
		MachineID: machineID,
		Threshold: 1.0,
	}
	jsonData, err := json.Marshal(threshold)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal threshold: %v", err))
	}
	_, _ = http.Post(fmt.Sprintf("http://localhost:%d/v1/threshold", config.ApiServer.Port),
		"application/json", bytes.NewBuffer(jsonData))
}

func ApplyMigration(config config.Config) {
	m, err := migrate.New(
		"file://../migrations",
		fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			config.PostgresDB.Username, config.PostgresDB.Password, config.PostgresDB.Port, config.PostgresDB.DatabaseName))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
