package lossytransport

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cevatbarisyilmaz/lossy"
)

// LossOptions Bundles the loss options for lossy
type LossOptions struct {
	Bandwidth      int
	MinLatency     time.Duration
	MaxLatency     time.Duration
	LossRate       float64
	HeaderOverhead int
}

// NewLossyTransport Creates a new http transport with loss connection
func NewLossyTransport(options *LossOptions) *http.Transport {
	return &http.Transport{
		Dial:                  dial(options),
		DialTLS:               dialTLS(options),
		DisableCompression:    true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
}

func dial(options *LossOptions) func(string, string) (net.Conn, error) {
	return func(network, address string) (net.Conn, error) {
		// connect to the remote address
		conn, err := net.Dial(network, address)

		// error handling
		if err != nil {
			return nil, err
		}

		// wrap it in a lossy connection
		lossyConn := lossy.Conn(
			conn,
			options.Bandwidth,
			options.MinLatency,
			options.MaxLatency,
			options.LossRate,
			options.HeaderOverhead,
		)

		return lossyConn, nil
	}
}

func dialTLS(options *LossOptions) func(string, string) (net.Conn, error) {
	return func(network, address string) (net.Conn, error) {
		// connect to the remote address
		conn, err := net.Dial(network, address)

		// error handling
		if err != nil {
			return nil, err
		}

		// wrap it in a lossy connection
		lossyConn := lossy.Conn(
			conn,
			options.Bandwidth,
			options.MinLatency,
			options.MaxLatency,
			options.LossRate,
			40,
		)

		// create a new tls config
		cfg := new(tls.Config)

		// get the url as url object
		u, err := url.Parse(fmt.Sprintf("https://%s", address))

		// error handling
		if err != nil {
			return nil, err
		}

		// set the server name for sni
		cfg.ServerName = u.Host[:strings.LastIndex(u.Host, ":")]

		// create the new tls client
		tlsConn := tls.Client(lossyConn, cfg)

		// initialize an error channel
		errorChannel := make(chan error, 2)

		// timeout for the tls handshake
		timer := time.AfterFunc(
			10*time.Second,
			func() {
				errorChannel <- fmt.Errorf("TLS handshake timeout")
			},
		)

		// run the handshake
		go func() {
			err := tlsConn.Handshake()
			timer.Stop()
			errorChannel <- err
		}()

		// verify if the handshake was successful
		if err := <-errorChannel; err != nil {
			lossyConn.Close()

			return nil, err
		}

		return tlsConn, nil
	}
}
