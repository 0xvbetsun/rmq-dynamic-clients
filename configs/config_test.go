// Package configs provides project configuration structure
package configs

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{name: "merge configs from env vars and conf file", want: &Config{
			AMQP: &AMQP{Url: "amqp://guest:guest@localhost:5672", Queue: "queue_name", Durable: true, AutoAck: true, AutoDelete: true},
			Docs: &Docs{Port: 3000},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Chdir("..")
			os.Setenv("AMQP_SERVER_URL", "amqp://guest:guest@localhost:5672")
			os.Setenv("AMQP_QUEUE_NAME", "queue_name")
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
