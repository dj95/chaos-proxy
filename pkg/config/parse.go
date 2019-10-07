package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Target Represent a target in the configuration.
type Target struct {
	Protocol   string   `mapstructure:"protocol"`
	Target     string   `mapstructure:"target"`
	Bandwidth  int      `mapstructure:"bandwidth"`
	Overhead   string   `mapstructure:"overhead"`
	Latency    *Latency `mapstructure:"latency"`
	LossRate   float64  `mapstructure:"loss_rate"`
	ListenPort int      `mapstructure:"listen_port"`
}

// Latency Hold the maximal and minimal latency.
type Latency struct {
	Min int `mapstructure:"min"`
	Max int `mapstructure:"max"`
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

		// save the target if no error occured
		output[key] = target
	}

	return output
}
