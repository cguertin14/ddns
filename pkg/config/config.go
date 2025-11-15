package config

import (
	"github.com/spf13/viper"
)

var (
	mustHaveEnvVars = []string{
		"CLOUDFLARE_TOKEN",
		"RECORD_NAME",
		"ZONE_NAME",
	}

	otherEnvVars = []string{
		"LOG_LEVEL",
	}
)

type Config struct {
	CloudflareToken string `mapstructure:"CLOUDFLARE_TOKEN"`
	RecordName      string `mapstructure:"RECORD_NAME"`
	ZoneName        string `mapstructure:"ZONE_NAME"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	// set fields from env vars
	for _, env := range mustHaveEnvVars {
		v.MustBindEnv(env)
	}
	for _, env := range otherEnvVars {
		if err := v.BindEnv(env); err != nil {
			return nil, err
		}
	}

	cfg := &Config{}
	err := v.Unmarshal(cfg)
	return cfg, err
}
