/*
Package config contains the cli configuration
*/
package config

import (
	"fmt"
	"strings"

	"github.com/controlplaneio/netassertv2-l4-client/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config contains the global configuration.
type Config struct {
	LogLevel    log.Level
	LogEncoding log.Encoding
	Protocol    string
	TargetHost  string
	TargetPort  uint16
	Message     string
	Timeout     uint
	Attempts    uint
	Period      uint
	SuccThrPec  uint
}

// Init config
func (c *Config) Init(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	var lasterror error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name

		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)
			if err := cmd.Flags().Set(configName, fmt.Sprintf("%v", val)); err != nil {
				lasterror = err
			}
		}
	})
	return lasterror
}
