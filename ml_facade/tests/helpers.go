package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"testing"
)

func dClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

func checkStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func checkJSON(t *testing.T, got, want any) {
	if !equal(want, got) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			got, want)
	}
}

func equal(a, b any) bool {
	return bytes.Equal(mustJSONMarshal(a), mustJSONMarshal(b))
}

func mustJSONMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
