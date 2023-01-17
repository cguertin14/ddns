package main

import (
	"context"
	"log"

	"github.com/cguertin14/ddns/pkg/cloudflare"
	"github.com/cguertin14/ddns/pkg/config"
	"github.com/cguertin14/ddns/pkg/ddns"
	gh "github.com/cguertin14/ddns/pkg/github"
	"github.com/cguertin14/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s\n", err)
	}

	// initialize logger
	appLogger := logger.Initialize(logger.Config{
		Level: cfg.LogLevel,
	})
	ctx := context.WithValue(context.Background(), logger.CtxKey, appLogger)

	// cloudflare client
	cfClient, err := cloudflare.NewClient(cfg.CloudflareToken)
	if err != nil {
		appLogger.Fatalln(err)
	}

	// github client
	ghClient := gh.NewClient(ctx, cfg.GithubToken)

	// initialize app
	app := ddns.NewClient(ddns.Dependencies{
		Cloudflare: cfClient,
		Github:     ghClient,
	})

	// run app
	report, err := app.Run(ctx, *cfg)
	if err != nil {
		appLogger.Fatalf("failed to run app: %s\n", err)
	}

	// success message
	if report.DnsChanged {
		appLogger.Infof("Successfully updated DNS Record %q to %q.\n", cfg.RecordName, report.NewIP)
	} else {
		appLogger.Infoln("DNS record unchanged.")
	}
}
