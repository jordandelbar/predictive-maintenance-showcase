package service

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

// TestGetReaderFromBody tests the getReaderFromBody function
func TestGetReaderFromBody(t *testing.T) {
	tests := []struct {
		name    string
		body    any
		wantErr bool
	}{
		{
			name:    "AMQP Delivery",
			body:    amqp091.Delivery{Body: []byte("test message")},
			wantErr: false,
		},
		{
			name:    "io.Reader",
			body:    strings.NewReader("test reader"),
			wantErr: false,
		},
		{
			name:    "Unsupported type",
			body:    123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReader, err := getReaderFromBody(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("getReaderFromBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				buf := new(bytes.Buffer)
				buf.ReadFrom(gotReader)
				switch v := tt.body.(type) {
				case amqp091.Delivery:
					if got := buf.String(); got != string(v.Body) {
						t.Errorf("getReaderFromBody() got = %v, want %v", got, string(v.Body))
					}
				case *strings.Reader:
					if got := buf.String(); got != "test reader" {
						t.Errorf("getReaderFromBody() got = %v, want %v", got, "test reader")
					}
				default:
					t.Errorf("getReaderFromBody() unexpected type = %v", reflect.TypeOf(tt.body))
				}
			}
		})
	}
}
