package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dj95/deception-proxy/pkg/config"
)

func TestTCPStartListener(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`OK`))
	}))

	tests := []struct {
		description   string
		proxy         Proxy
		expectedError bool
	}{
		{
			description: "port out of range",
			proxy: &TCPProxy{
				&config.Target{
					ListenPort: 666666,
				},
			},
			expectedError: true,
		},
		{
			description: "http test",
			proxy: &TCPProxy{
				&config.Target{
					Protocol:   "tcp",
					Target:     strings.TrimLeft(s.URL, "http://"),
					ListenPort: 8080,
					Latency: &config.Latency{
						Min: 10,
						Max: 100,
					},
				},
			},
			expectedError: false,
		},
	}

	for _, test := range tests {
		go func(test *struct {
			description   string
			proxy         Proxy
			expectedError bool
		}) {
			err := test.proxy.StartListener()

			assert.Equalf(t, test.expectedError, err != nil, test.description)
		}(&test)

		time.Sleep(1.0 * time.Second)

		req, _ := http.NewRequest("GET", "http://127.0.0.1:8080", nil)

		_, _ = http.DefaultClient.Do(req)
	}
}
