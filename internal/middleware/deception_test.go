package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type RecorderNotifier struct {
	http.ResponseWriter

	Recorder *httptest.ResponseRecorder
}

func NewRecorderNotifier(recorder *httptest.ResponseRecorder) *RecorderNotifier {
	return &RecorderNotifier{
		Recorder: recorder,
	}
}

func (r *RecorderNotifier) Header() http.Header {
	return r.Recorder.Header()
}

func (r *RecorderNotifier) Write(input []byte) (int, error) {
	return r.Recorder.Write(input)
}

func (r *RecorderNotifier) WriteHeader(status int) {
	r.Recorder.WriteHeader(status)
}

func (r *RecorderNotifier) CloseNotify() <-chan bool {
	return make(<-chan bool, 1)
}

func TestDeception(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`OK`))
	}))

	tests := []struct {
		description   string
		method        string
		url           string
		route         string
		expectedCode  int
		expectedError bool
	}{
		{
			description:   "index route",
			method:        "GET",
			url:           s.URL,
			route:         "/",
			expectedCode:  200,
			expectedError: false,
		},
		{
			description:   "wrong configured url",
			method:        "GET",
			url:           "non\tparsable\turl",
			route:         "/",
			expectedCode:  200,
			expectedError: true,
		},
		{
			description:   "skip to healthz route",
			method:        "GET",
			url:           s.URL,
			route:         "/healthz",
			expectedCode:  200,
			expectedError: false,
		},
	}

	for _, test := range tests {
		fmt.Printf("----- %s\n", test.description)
		// configure the middleware
		viper.Set("conn.target", test.url)

		// initialize the middleware
		mw, err := Deception()

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if err != nil {
			continue
		}

		// create a new response writer and request
		responseWriter := NewRecorderNotifier(httptest.NewRecorder())
		req, _ := http.NewRequest(
			test.method,
			test.route,
			nil,
		)

		c, _ := gin.CreateTestContext(responseWriter)
		c.Request = req

		// test the middleware function
		mw(c)

		assert.Equalf(
			t,
			test.expectedCode,
			responseWriter.Recorder.Code,
			test.description,
		)
		fmt.Printf("----- %v\n", responseWriter.Recorder.Body.String())

		viper.Set("conn.target", "")
		fmt.Printf("----- test end\n")
	}
}
