package proxy

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dj95/chaos-proxy/pkg/config"
)

func TestTCPStartListener(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`OK`))
	}))

	randomPort := rand.Intn(10000) + 30000

	tests := []struct {
		description     string
		proxy           Proxy
		closeConnection bool
		expectedError   bool
	}{
		{
			description: "port out of range",
			proxy: &TCPProxy{
				Target: &config.Target{
					ListenPort: 666666,
				},
			},
			expectedError: true,
		},
		{
			description: "http test",
			proxy: &TCPProxy{
				Target: &config.Target{
					Protocol:   "tcp",
					Target:     strings.TrimPrefix(s.URL, "http://"),
					ListenPort: randomPort,
					Latency: &config.Latency{
						Min: 10,
						Max: 100,
					},
				},
			},
			expectedError: false,
		},
		{
			description:     "close connection before http request",
			closeConnection: true,
			proxy: &TCPProxy{
				Target: &config.Target{
					Protocol:   "tcp",
					Target:     strings.TrimPrefix(s.URL, "http://"),
					ListenPort: randomPort + 1,
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
		err := test.proxy.StartListener()

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if test.expectedError {
			continue
		}

		if test.closeConnection {
			test.proxy.Shutdown()
		}

		time.Sleep(1.0 * time.Second)

		req, _ := http.NewRequest(
			"GET",
			fmt.Sprintf("http://127.0.0.1:%d/", randomPort),
			nil,
		)

		res, err := http.DefaultClient.Do(req)

		fmt.Printf("%v\n", err)
		fmt.Printf("%d\n", res.StatusCode)
	}
}
