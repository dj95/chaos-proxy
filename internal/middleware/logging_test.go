package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`OK`))
	}))

	tests := []struct {
		description  string
		method       string
		route        string
		expectedCode int
	}{
		{
			description:  "call the health route to trigger the logging middleware",
			method:       "GET",
			route:        "/healthz",
			expectedCode: 200,
		},
	}

	viper.Set("conn.target", s.URL)

	middleware := Logging()

	for _, test := range tests {
		req, _ := http.NewRequest(
			test.method,
			test.route,
			nil,
		)

		// create a new response writer
		responseWriter := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(responseWriter)
		c.Request = req

		// test the middleware function
		middleware(c)

		assert.Equalf(
			t,
			test.expectedCode,
			responseWriter.Code,
			test.description,
		)
	}
}
