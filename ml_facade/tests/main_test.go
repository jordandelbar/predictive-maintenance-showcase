package test

import (
	"fmt"
	"ml_facade/app"
	"ml_facade/config"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

var testCfg config.Config

func TestMain(m *testing.M) {
	var err error

	testCfg.Port, err = findAvailablePort()
	if err != nil {
		fmt.Println(err)
	}
	testCfg.Env = "test"
	testCfg.Db.Host = "localhost"
	testCfg.Db.Port = "5432"
	testCfg.Db.Username = "monitor"
	testCfg.Db.Password = "test"
	testCfg.Db.DatabaseName = "monitoring"
	testCfg.Rdb.Host = "localhost"
	testCfg.Rdb.Port = "6379"
	testCfg.MlService.Host = "localhost"
	testCfg.MlService.Port = "3000"
	testCfg.Db.MaxOpenConns = 25
	testCfg.Db.MaxIdleConns = 25
	testCfg.Db.MaxIdleTime = 5 * time.Minute

	go app.StartApp(testCfg)

	if !waitForAPI(fmt.Sprintf("http://localhost:%d/v1/healthcheck", testCfg.Port), 30, 10*time.Microsecond) {
		fmt.Println("API did not start in time")
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func waitForAPI(url string, maxRetries int, retryInterval time.Duration) bool {
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		time.Sleep(retryInterval)
	}
	return false
}
