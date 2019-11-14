// Package config Parse and create config objects.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Target Represent a target in the configuration.
type Target struct {
	ID         string   `json:"id"`
	Protocol   string   `json:"protocol" mapstructure:"protocol"`
	Target     string   `json:"target" mapstructure:"target"`
	Bandwidth  int      `json:"bandwidth" mapstructure:"bandwidth"`
	Overhead   string   `json:"overhead" mapstructure:"overhead"`
	Latency    *Latency `json:"latency" mapstructure:"latency"`
	LossRate   float64  `json:"loss_rate" mapstructure:"loss_rate"`
	ListenPort int      `json:"listen_port" mapstructure:"listen_port"`
}

// Latency Hold the maximal and minimal latency.
type Latency struct {
	Min int `json:"min" mapstructure:"min"`
	Max int `json:"max" mapstructure:"max"`
}

// Parse Parse the configuration into a target map.
func Parse(input map[string]interface{}) map[string]*Target {
	// initialize the output
	output := map[string]*Target{}

	// iterate through the configured connections
	for key := range input {
		// create a new target object and get its pointer
		target := new(Target)

		// unmarshal the object
		err := viper.UnmarshalKey(
			"conn."+key,
			target,
		)

		// error handling
		if err != nil {
			fmt.Printf("%s\n", err.Error())

			continue
		}

		// save the key as id
		target.ID = key

		// save the target if no error occurred
		output[key] = target
	}

	return output
}
