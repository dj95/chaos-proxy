package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		description    string
		input          map[string]interface{}
		viperConfig    map[string]interface{}
		expectedResult map[string]*Target
	}{
		{
			description: "successful parsing",
			input: map[string]interface{}{
				"testconn": ``,
			},
			viperConfig: map[string]interface{}{
				"protocol":  "tcp",
				"target":    "127.0.0.1:80",
				"bandwidth": 8388608,
				"overhead":  "v4_tcp_max",
				"latency": map[string]int{
					"min": 10,
					"max": 100,
				},
				"loss_rate":   0.1,
				"listen_port": 8080,
			},
			expectedResult: map[string]*Target{
				"testconn": &Target{
					Protocol:  "tcp",
					Target:    "127.0.0.1:80",
					Bandwidth: 8388608,
					Overhead:  "v4_tcp_max",
					Latency: &Latency{
						Min: 10,
						Max: 100,
					},
					LossRate:   0.1,
					ListenPort: 8080,
				},
			},
		},
		{
			description: "cannot unmarshal",
			input: map[string]interface{}{
				"testconn": ``,
			},
			viperConfig: map[string]interface{}{
				"protocol":  "tcp",
				"target":    "127.0.0.1:80",
				"bandwidth": "this is false",
				"overhead":  "v4_tcp_max",
				"latency": map[string]int{
					"min": 10,
					"max": 100,
				},
				"loss_rate":   0.1,
				"listen_port": 8080,
			},
			expectedResult: map[string]*Target{},
		},
	}

	for _, test := range tests {
		for key := range test.input {
			viper.Set("conn."+key, test.viperConfig)
		}

		result := Parse(test.input)

		assert.Truef(t, reflect.DeepEqual(test.expectedResult, result), test.description)
	}
}
