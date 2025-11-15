package ddns

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "embed"

	"github.com/cguertin14/ddns/pkg/config"
	"github.com/cguertin14/logger"
	legacy_cf "github.com/cloudflare/cloudflare-go"
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
	zoneID, err := c.cloudflare.ZoneIDByName(cfg.ZoneName)
	if err != nil {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to fetch zone ID: %s", err)
	}
	identifier := legacy_cf.ZoneIdentifier(zoneID)

	// list dns records and find the one needed
	records, _, err := c.cloudflare.ListDNSRecords(ctx, identifier, legacy_cf.ListDNSRecordsParams{
		Name: fmt.Sprintf("%s.%s", cfg.RecordName, cfg.ZoneName),
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
		if _, err := c.cloudflare.UpdateDNSRecord(ctx, identifier, legacy_cf.UpdateDNSRecordParams{
			Name:    cfg.RecordName,
			ID:      record.ID,
			Content: newIP,
			Type:    "A",
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
