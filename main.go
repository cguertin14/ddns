package main

import (
	"context"
	"log"

	"github.com/cguertin14/ddns/pkg/cloudflare"
	"github.com/cguertin14/ddns/pkg/config"
	"github.com/cguertin14/ddns/pkg/ddns"
	"github.com/cguertin14/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	// initialize logger
	appLogger := logger.Initialize(logger.Config{
		Level:     cfg.LogLevel,
		Formatter: logger.ServiceFormatter,
	})
	ctx := context.WithValue(context.Background(), logger.CtxKey, appLogger)

	// cloudflare client
	cfClient := cloudflare.NewClient(cfg.CloudflareToken)

	// initialize app
	app := ddns.NewClient(ddns.Dependencies{
		Cloudflare: cfClient,
	})

	// run app
	report, err := app.Run(ctx, *cfg)
	if err != nil {
		appLogger.Fatalf("failed to run app: %s", err)
	}

	// success message
	if report.DnsChanged {
		appLogger.Infof("successfully updated dns Record %v to %v", cfg.RecordName, report.NewIP)
	} else {
		appLogger.Infoln("dns record unchanged")
	}
}
