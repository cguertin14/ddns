package config

import (
	"github.com/spf13/viper"
)

const (
	githubToken     string = "GITHUB_TOKEN"
	cloudflareToken string = "CLOUDFLARE_TOKEN"
	recordName      string = "RECORD_NAME"
	zoneName        string = "ZONE_NAME"
)

type Config struct {
	GithubToken     string `mapstructure:"GITHUB_TOKEN"`
	CloudflareToken string `mapstructure:"CLOUDFLARE_TOKEN"`
	RecordName      string `mapstructure:"RECORD_NAME"`
	ZoneName        string `mapstructure:"ZONE_NAME"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	if err := v.BindEnv(
		githubToken,
		cloudflareToken,
		recordName,
		zoneName,
	); err != nil {
		return nil, err
	}

	cfg := &Config{}
	err := v.Unmarshal(cfg)
	return cfg, err
}
