package ddns

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cguertin14/ddns/pkg/config"
	legacy_cf "github.com/cloudflare/cloudflare-go"
)

func (c Client) Run(ctx context.Context, cfg config.Config) error {
	// fetch zone ID from name
	zoneID, err := c.cloudflare.ZoneIDByName(cfg.ZoneName)
	if err != nil {
		return fmt.Errorf("failed to fetch zone ID: %s\n", err)
	}
	identifier := legacy_cf.ZoneIdentifier(zoneID)

	// list dns records and find the one needed
	records, _, err := c.cloudflare.ListDNSRecords(ctx, identifier, legacy_cf.ListDNSRecordsParams{
		Name: fmt.Sprintf("%s.%s", cfg.RecordName, cfg.ZoneName),
	})
	if err != nil {
		return fmt.Errorf("failed to fetch dns records: %s\n", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("failed to find dns record on cloudflare\n")
	}

	// fetch public IP
	ipAddress, err := getPublicIP()
	if err != nil {
		return fmt.Errorf("failed to fetch public IP: %s", err)
	}

	// check if IP changed and act if it did
	record := records[0]
	if record.Content != ipAddress {
		// step 1: update dns record
		if err := c.cloudflare.UpdateDNSRecord(context.TODO(), identifier, legacy_cf.UpdateDNSRecordParams{
			ZoneID:  zoneID,
			Name:    cfg.RecordName,
			ID:      record.ID,
			Content: ipAddress,
			Type:    "A",
		}); err != nil {
			return fmt.Errorf("failed to update dns record: %s\n", err)
		}

		// step 2: TODO: open PR on homelab repo
		// - create branch
		// -
	}

	return nil
}

func getPublicIP() (string, error) {
	res, err := http.Get("https://ifconfig.me")
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
