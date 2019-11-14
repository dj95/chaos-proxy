// Package proxy Implement proxy components for different protocols.
package proxy

import (
	"fmt"

	"github.com/dj95/chaos-proxy/pkg/config"
)

// Proxy Interface for the proxy server of one target.
type Proxy interface {
	// status functions
	Config() *config.Target

	// management functions
	StartListener() error
	Shutdown() error
}

// New Create and initialize a new proxy connection
func New(target *config.Target) (Proxy, error) {
	// initialize a new proxy variable
	var proxy Proxy

	// create the proxy based on the protocol type
	switch target.Protocol {
	case "tcp":
		proxy = &TCPProxy{
			Target: target,
		}
	default:
		return nil, fmt.Errorf("wrong protocol given for target")
	}

	// return the proxy
	return proxy, nil
}
