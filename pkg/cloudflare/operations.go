package cloudflare

import (
	"context"
	"fmt"

	cloudflare "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/zones"
)

func (c Client) ZoneIDByName(ctx context.Context, name string) (string, error) {
	// List zones with a filter on the name
	page, err := c.api.Zones.List(ctx, zones.ZoneListParams{
		Name: cloudflare.F(name),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	// Check if we got any results
	if len(page.Result) == 0 {
		return "", fmt.Errorf("zone not found: %s", name)
	}

	return page.Result[0].ID, nil
}

func (c Client) ListDNSRecords(ctx context.Context, zoneID string, params dns.RecordListParams) ([]dns.RecordResponse, error) {
	// Set the ZoneID in params
	params.ZoneID = cloudflare.F(zoneID)

	page, err := c.api.DNS.Records.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS records: %w", err)
	}

	// Return the results from the page
	return page.Result, nil
}

func (c Client) UpdateDNSRecord(ctx context.Context, zoneID string, recordID string, params dns.RecordUpdateParams) (dns.RecordResponse, error) {
	// Set the ZoneID in params
	params.ZoneID = cloudflare.F(zoneID)

	record, err := c.api.DNS.Records.Update(ctx, recordID, params)
	if err != nil {
		return dns.RecordResponse{}, fmt.Errorf("failed to update DNS record: %w", err)
	}

	return *record, nil
}
