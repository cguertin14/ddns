package config

import (
	"github.com/spf13/viper"
)

var (
	envVars = []string{
		"UPDATE_GITHUB_TERRAFORM",
		"GITHUB_BASE_BRANCH",
		"GITHUB_FILE_PATH",
		"GITHUB_REPO_OWNER",
		"GITHUB_REPO_NAME",
		"GITHUB_TOKEN",
		"CLOUDFLARE_TOKEN",
		"RECORD_NAME",
		"ZONE_NAME",
		"LOG_LEVEL",
	}
)

type Config struct {
	UpdateGithubTerraform bool   `mapstructure:"UPDATE_GITHUB_TERRAFORM"`
	GithubBaseBranch      string `mapstructure:"GITHUB_BASE_BRANCH"`
	GithubFilePath        string `mapstructure:"GITHUB_FILE_PATH"`
	GithubRepoOwner       string `mapstructure:"GITHUB_REPO_OWNER"`
	GithubRepoName        string `mapstructure:"GITHUB_REPO_NAME"`
	GithubToken           string `mapstructure:"GITHUB_TOKEN"`
	CloudflareToken       string `mapstructure:"CLOUDFLARE_TOKEN"`
	RecordName            string `mapstructure:"RECORD_NAME"`
	ZoneName              string `mapstructure:"ZONE_NAME"`
	LogLevel              string `mapstructure:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	// set fields from env vars
	for _, env := range envVars {
		if err := v.BindEnv(env); err != nil {
			return nil, err
		}
	}

	cfg := &Config{}
	err := v.Unmarshal(cfg)
	return cfg, err
}
