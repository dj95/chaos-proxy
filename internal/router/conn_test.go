package router

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dj95/chaos-proxy/pkg/config"
	"github.com/dj95/chaos-proxy/pkg/proxy"
)

type MockProxy struct {
	id        string
	failStart bool
	failClose bool
}

func (m *MockProxy) Config() *config.Target {
	return &config.Target{
		ID: m.id,
	}
}

func (m *MockProxy) StartListener() error {
	if m.failStart {
		return fmt.Errorf("")
	}

	return nil
}

func (m *MockProxy) Shutdown() error {
	if m.failClose {
		return fmt.Errorf("")
	}

	return nil
}

func TestConnHandler(t *testing.T) {
	tests := []struct {
		description  string
		expectedCode int
		expectedBody string
	}{
		{
			description:  "get all clients",
			expectedCode: http.StatusOK,
			expectedBody: "[{\"id\":\"\",\"protocol\":\"\",\"target\":\"\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":0}]",
		},
	}

	s, _ := Setup(
		[]proxy.Proxy{
			&proxy.TCPProxy{
				Target: &config.Target{},
			},
		},
	)

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			"/api/conn",
			nil,
		)

		s.ServeHTTP(w, req)

		assert.Equalf(t, test.expectedCode, w.Code, test.description)
		assert.Equalf(t, test.expectedBody, w.Body.String(), test.description)
	}
}

func TestConnIDHandler(t *testing.T) {
	tests := []struct {
		description  string
		id           string
		expectedCode int
		expectedBody string
	}{
		{
			description:  "id exists",
			id:           "foobar",
			expectedCode: http.StatusOK,
			expectedBody: "{\"id\":\"foobar\",\"protocol\":\"\",\"target\":\"\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":0}",
		},
		{
			description:  "id does not exist",
			id:           "nonexistant",
			expectedCode: http.StatusNotFound,
			expectedBody: "{}",
		},
	}

	s, _ := Setup(
		[]proxy.Proxy{
			&proxy.TCPProxy{
				Target: &config.Target{
					ID: "foobar",
				},
			},
		},
	)

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			"/api/conn/"+test.id,
			nil,
		)

		s.ServeHTTP(w, req)

		assert.Equalf(t, test.expectedCode, w.Code, test.description)
		assert.Equalf(t, test.expectedBody, w.Body.String(), test.description)
	}
}

func TestConnUpdateHandler(t *testing.T) {
	randomPort := rand.Intn(10000) + 40000

	tests := []struct {
		description  string
		id           string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			description:  "id exists",
			id:           "foobar",
			body:         fmt.Sprintf("{\"id\":\"foobar\",\"protocol\":\"tcp\",\"target\":\"127.0.0.1:8080\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":%d}", randomPort),
			expectedCode: http.StatusOK,
			expectedBody: "{\"msg\":\"\",\"status\":\"success\"}",
		},
		{
			description:  "invalid request body",
			id:           "nonexistant",
			body:         "this is no valid request body",
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"msg\":\"invalid request body\",\"status\":\"error\"}",
		},
		{
			description:  "create new proxy",
			id:           "foo",
			body:         fmt.Sprintf("{\"id\":\"foobar\",\"protocol\":\"tcp\",\"target\":\"127.0.0.1:8080\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":%d}", randomPort+1),
			expectedCode: http.StatusOK,
			expectedBody: "{\"msg\":\"\",\"status\":\"success\"}",
		},
		{
			description:  "invalid new config",
			id:           "foo",
			body:         fmt.Sprintf("{\"id\":\"foo\",\"protocol\":\"\",\"target\":\"127.0.0.1:8080\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":%d}", randomPort+1),
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"msg\":\"cannot create updated proxy: wrong protocol given for target\",\"status\":\"error\"}",
		},
		{
			description:  "create new proxy with same port so it cannot start",
			id:           "bar",
			body:         fmt.Sprintf("{\"id\":\"bar\",\"protocol\":\"tcp\",\"target\":\"127.0.0.1:8080\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":%d}", randomPort+1),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"msg\":\"cannot start updated proxy\",\"status\":\"error\"}",
		},
		{
			description:  "saved proxy cannot be started",
			id:           "duh",
			body:         fmt.Sprintf("{\"id\":\"bar\",\"protocol\":\"tcp\",\"target\":\"127.0.0.1:8080\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":%d}", randomPort+1),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"msg\":\"udated proxy and restored saved proxy cannot start\",\"status\":\"error\"}",
		},
		{
			description:  "fail to shutdown a running proxy",
			id:           "failtoclose",
			body:         fmt.Sprintf("{\"id\":\"failtoclose\",\"protocol\":\"tcp\",\"target\":\"127.0.0.1:8080\",\"bandwidth\":0,\"overhead\":\"\",\"latency\":null,\"loss_rate\":0,\"listen_port\":%d}", randomPort+1),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"msg\":\"cannot shutdown proxy\",\"status\":\"error\"}",
		},
	}

	p := proxy.TCPProxy{
		Target: &config.Target{
			ID:         "foobar",
			Protocol:   "tcp",
			Target:     "127.0.0.1:8080",
			ListenPort: randomPort,
		},
	}

	p.StartListener()

	p2 := proxy.TCPProxy{
		Target: &config.Target{
			ID:         "duh",
			Protocol:   "tcp",
			Target:     "127.0.0.1:8080",
			ListenPort: randomPort + 2,
		},
	}

	p2.StartListener()

	p2.Config().Protocol = ""

	time.Sleep(500 * time.Millisecond)

	s, _ := Setup(
		[]proxy.Proxy{
			&p,
			&p2,
			&proxy.TCPProxy{
				Target: &config.Target{
					ID:         "invalidshutdown",
					Protocol:   "tcp",
					Target:     "127.0.0.1:8080",
					ListenPort: randomPort + 3,
				},
			},
			&MockProxy{
				id:        "failtoclose",
				failClose: true,
			},
		},
	)

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			"/api/conn/"+test.id,
			bytes.NewBufferString(test.body),
		)

		s.ServeHTTP(w, req)

		assert.Equalf(t, test.expectedCode, w.Code, test.description)
		assert.Equalf(t, test.expectedBody, w.Body.String(), test.description)
	}
}
