package proxy

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/dj95/deception-proxy/pkg/lossytransport"
)

type LossyRoundTrip struct {
	http.RoundTripper

	LossOptions   *lossytransport.LossOptions
	LossTransport *http.Transport
}

func NewLossyRoundTrip() *LossyRoundTrip {
	lossOptions := &lossytransport.LossOptions{
		Bandwidth:  int(viper.GetInt("conn.bandwidth") / 8),
		MinLatency: 10 * time.Millisecond,
		MaxLatency: 100 * time.Millisecond,
		LossRate:   viper.GetFloat64("conn.loss_rate"),
	}

	return &LossyRoundTrip{
		LossOptions: lossOptions,
		LossTransport: lossytransport.NewLossyTransport(
			lossOptions,
		),
	}
}

func (ltr *LossyRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
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
