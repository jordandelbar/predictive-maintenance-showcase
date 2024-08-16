package test

import (
	"io"
	"log"
	"testing"
)

func closeOrLog(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

func checkStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}
