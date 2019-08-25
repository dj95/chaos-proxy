package roundtrip

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/dj95/deception-proxy/pkg/lossytransport"
)

// Lossy Do http requests in a proxy with packet loss, latency and bandwidth limitation
type Lossy struct {
	http.RoundTripper

	LossOptions   *lossytransport.LossOptions
	LossTransport *http.Transport
}

// NewLossy Initializes a new LossyRoundTrip that can be used for http proxies
func NewLossy() *Lossy {
	// create the loss options to use
	lossOptions := &lossytransport.LossOptions{
		Bandwidth:      int(viper.GetInt("conn.bandwidth") / 8),
		MinLatency:     10 * time.Millisecond,
		MaxLatency:     100 * time.Millisecond,
		LossRate:       viper.GetFloat64("conn.loss_rate"),
		HeaderOverhead: 40,
	}

	return &Lossy{
		LossOptions: lossOptions,
		LossTransport: lossytransport.NewLossyTransport(
			lossOptions,
		),
	}
}

// RoundTrip Do the given request and return the response
func (ltr *Lossy) RoundTrip(req *http.Request) (*http.Response, error) {
	// debug message
	log.Debugf("Round trip")

	// create a new http client with the lossy transport
	client := &http.Client{
		Transport: ltr.LossTransport,
	}

	// clone the request
	r, err := http.NewRequest(
		req.Method,
		req.URL.String(),
		req.Body,
	)

	// debug message
	log.Debugf("cloned request: %v - err: %v", r, err)

	// do the request
	resp, err := client.Do(r)

	// debug errors
	log.Debugf("Did request with error: %v", err)

	return resp, err
}
