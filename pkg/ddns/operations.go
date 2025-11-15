package ddns

import (
	"context"
	"fmt"
	"io"
	"net/http"

	_ "embed"

	"github.com/cguertin14/ddns/pkg/config"
	"github.com/cguertin14/logger"
	cloudflare "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
)

type RunReport struct {
	DnsChanged bool
	NewIP      string
}

type PRReport struct {
	NewIP, OldIP         string
	ZoneName, RecordName string
}

// returns wether or not DNS changed
func (c Client) Run(ctx context.Context, cfg config.Config) (RunReport, error) {
	// fetch zone ID from name
	zoneID, err := c.cloudflare.ZoneIDByName(ctx, cfg.ZoneName)
	if err != nil {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to fetch zone ID: %s", err)
	}

	// list dns records and find the one needed
	fullRecordName := fmt.Sprintf("%s.%s", cfg.RecordName, cfg.ZoneName)
	records, err := c.cloudflare.ListDNSRecords(ctx, zoneID, dns.RecordListParams{
		Name: cloudflare.F(dns.RecordListParamsName{
			Exact: cloudflare.F(fullRecordName),
		}),
	})
	if err != nil {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to fetch dns records: %s", err)
	}
	if len(records) == 0 {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to find dns record on cloudflare")
	}

	// fetch public IP
	newIP, err := getPublicIP()
	if err != nil {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to fetch public IP: %s", err)
	}
	logs := logger.NewFromContextOrDefault(ctx)
	logs.Infof("Current IP is %s", newIP)

	// check if IP changed and act if it did
	record := records[0]
	if record.Content != newIP {
		// step 1: update dns record
		if _, err := c.cloudflare.UpdateDNSRecord(ctx, zoneID, record.ID, dns.RecordUpdateParams{
			Body: dns.ARecordParam{
				Name:    cloudflare.F(fullRecordName),
				Content: cloudflare.F(newIP),
				Type:    cloudflare.F(dns.ARecordTypeA),
				TTL:     cloudflare.F(dns.TTL(record.TTL)),
			},
		}); err != nil {
			return RunReport{DnsChanged: false}, fmt.Errorf("failed to update dns record: %s", err)
		}

		return RunReport{
			DnsChanged: true,
			NewIP:      newIP,
		}, nil
	}

	return RunReport{DnsChanged: false}, nil
}

func getPublicIP() (string, error) {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
