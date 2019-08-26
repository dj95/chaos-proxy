package roundtrip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`OK`))
	}))

	tests := []struct {
		description string

		bandwidth      int
		latencyMin     int
		latencyMax     int
		lossRate       float64
		headerOverhead int

		method string
		url    string

		expectedError bool
	}{
		{
			description: "do a request",

			bandwidth:      8388608,
			latencyMin:     10,
			latencyMax:     100,
			lossRate:       0.1,
			headerOverhead: 40,

			method: "GET",
			url:    s.URL,

			expectedError: false,
		},
	}

	for _, test := range tests {
		viper.Set("conn.bandwidth", test.bandwidth)
		viper.Set("conn.latency.min", test.latencyMin)
		viper.Set("conn.latency.max", test.latencyMax)
		viper.Set("conn.loss_rate", test.lossRate)
		viper.Set("conn.overhead", test.headerOverhead)

		lossy := NewLossy()

		req, _ := http.NewRequest(
			test.method,
			test.url,
			nil,
		)

		_, err := lossy.RoundTrip(req)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		viper.Set("conn.bandwidth", "")
		viper.Set("conn.latency.min", "")
		viper.Set("conn.latency.max", "")
		viper.Set("conn.loss_rate", "")
		viper.Set("conn.overhead", "")
	}
}
