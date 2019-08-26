package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthzHandler(t *testing.T) {
	tests := []struct {
		description  string
		expectedCode int
		expectedBody string
	}{
		{
			description:  "success",
			expectedCode: 200,
			expectedBody: "OK",
		},
		{
			description:  "non initializable deception middleware",
			expectedCode: 200,
			expectedBody: "OK",
		},
	}

	s, _ := Setup()

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			"/healthz",
			nil,
		)

		s.ServeHTTP(w, req)

		assert.Equalf(t, test.expectedCode, w.Code, test.description)
		assert.Equalf(t, test.expectedBody, w.Body.String(), test.description)
	}
}
