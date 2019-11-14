package router

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/dj95/deception-proxy/pkg/proxy"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		description   string
		targetURL     string
		expectedError bool
	}{
		{
			description:   "success",
			targetURL:     "",
			expectedError: false,
		},
	}

	for _, test := range tests {
		viper.Set("conn.target", test.targetURL)

		_, err := Setup(
			[]proxy.Proxy{},
		)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		viper.Set("conn.target", "")
	}
}
