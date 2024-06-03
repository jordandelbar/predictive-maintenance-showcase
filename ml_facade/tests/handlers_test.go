package test

import (
	"fmt"
	"net/http"
	"testing"
)

// Check if the healthcheck route works normally
func TestHealthcheck(t *testing.T) {
	// Arrange
	url := fmt.Sprintf("http://localhost:%d/v1/healthcheck", testCfg.Port)

	// Act
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer dClose(resp.Body)

	// Assert
	checkStatus(t, resp.StatusCode, http.StatusOK)
}
