package proxy

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dj95/deception-proxy/pkg/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		description    string
		target         *config.Target
		expectedError  bool
		expectedResult Proxy
	}{
		{
			description: "success",
			target: &config.Target{
				Protocol:  "tcp",
				Target:    "127.0.0.1:80",
				Bandwidth: 8388608,
				Overhead:  "v4_tcp_max",
				Latency: &config.Latency{
					Min: 10,
					Max: 100,
				},
				LossRate:   0.1,
				ListenPort: 8080,
			},
			expectedError: false,
			expectedResult: &TCPProxy{
				Target: &config.Target{
					Protocol:  "tcp",
					Target:    "127.0.0.1:80",
					Bandwidth: 8388608,
					Overhead:  "v4_tcp_max",
					Latency: &config.Latency{
						Min: 10,
						Max: 100,
					},
					LossRate:   0.1,
					ListenPort: 8080,
				},
			},
		},
		{
			description: "wrong protocol",
			target: &config.Target{
				Protocol:  "foobar",
				Target:    "127.0.0.1:80",
				Bandwidth: 8388608,
				Overhead:  "v4_tcp_max",
				Latency: &config.Latency{
					Min: 10,
					Max: 100,
				},
				LossRate:   0.1,
				ListenPort: 8080,
			},
			expectedError:  true,
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		result, err := New(test.target)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if err != nil {
			continue
		}

		assert.Truef(t, reflect.DeepEqual(test.expectedResult, result), test.description)
	}
}
