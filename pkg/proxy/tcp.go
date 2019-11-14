package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/cevatbarisyilmaz/lossy"
	log "github.com/sirupsen/logrus"

	"github.com/dj95/deception-proxy/pkg/config"
)

// TCPProxy Proxy tcp requests
type TCPProxy struct {
	Target   *config.Target
	listener net.Listener
	closed   bool
}

// Config Return the current confiuration.
func (p *TCPProxy) Config() *config.Target {
	return p.Target
}

// StartListener Start a listener and handle incoming connections.
func (p *TCPProxy) StartListener() error {
	// intialize a new listener
	listener, err := net.Listen(
		p.Target.Protocol,
		fmt.Sprintf(":%d", p.Target.ListenPort),
	)

	// error handling
	if err != nil {
		// log the accepted connection
		log.WithFields(log.Fields{
			"protocol": "tcp",
			"target":   p.Target.Target,
		}).Errorf("error creating listener: %s", err.Error())

		return err
	}

	// save the listener
	p.listener = listener

	// iniialize the close indicator
	p.closed = false

	// run in an infinite loop
	go func() {
		for {
			// accept a new connection
			conn, err := listener.Accept()

			// log the accepted connection
			log.WithFields(log.Fields{
				"protocol": "tcp",
				"target":   p.Target.Target,
			}).Infof("accepting connection")

			// error handling
			if err != nil {
				log.WithFields(log.Fields{
					"protocol": "tcp",
					"target":   p.Target.Target,
				}).Warnf("accepting the connection failed: %s", err.Error())

				// if the proxy is closed...
				if p.closed {
					// ...return
					return
				}

				continue
			}

			// handle the requests concurrent
			go p.handleRequest(conn)
		}
	}()

	return nil
}

// Shutdown Close the listener to shutdown the proxy
func (p *TCPProxy) Shutdown() error {
	log.WithFields(log.Fields{
		"target": p.Target.Target,
	}).Info("shutting down proxy")

	p.closed = true

	return p.listener.Close()
}

func (p *TCPProxy) handleRequest(clientConn net.Conn) {
	// connect to the remote address
	conn, err := net.Dial(
		p.Target.Protocol,
		p.Target.Target,
	)

	// error handling
	if err != nil {
		return
	}

	// wrap it in a lossy connection
	lossyConn := lossy.Conn(
		conn,
		p.Target.Bandwidth,
		time.Duration(p.Target.Latency.Min),
		time.Duration(p.Target.Latency.Max),
		p.Target.LossRate,
		40,
	)

	// close the lossy connection on exit
	defer lossyConn.Close()

	// forward the data between the client and the target
	// with the configured lossy connection
	p.forwardData(lossyConn, clientConn)
}

// forwardData Forward data between the target and the client
func (p *TCPProxy) forwardData(lossyConn net.Conn, conn net.Conn) {
	// initialize a new channel for channelDone
	waitGroup := new(sync.WaitGroup)

	waitGroup.Add(1)
	// forward packets between target and client
	go func() {
		// unblock the waitgroup on return
		defer waitGroup.Done()

		// forward the request
		_, err := io.Copy(lossyConn, conn)

		// if an error occurred
		if err != nil {
			log.Infof("cannot forward data from client to target")
		}
	}()

	waitGroup.Add(1)
	// forward packets between client and target
	go func() {
		// unblock the waitgroup on return
		defer waitGroup.Done()

		// forward the request
		_, err := io.Copy(conn, lossyConn)

		// if an error occurred
		if err != nil {
			log.Infof("cannot forward data from target to client")
		}
	}()

	// wait until all data is sent
	waitGroup.Wait()

	// close the target connection
	lossyConn.Close()

	// close the client connection
	conn.Close()
}
