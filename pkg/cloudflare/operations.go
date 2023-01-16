package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

func (c Client) ZoneIDByName(name string) (string, error) {
	return c.api.ZoneIDByName(name)
}

func (c Client) ListDNSRecords(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.ListDNSRecordsParams) ([]cloudflare.DNSRecord, *cloudflare.ResultInfo, error) {
	return c.api.ListDNSRecords(ctx, rc, params)
}

func (c Client) UpdateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.UpdateDNSRecordParams) error {
	return c.api.UpdateDNSRecord(ctx, rc, params)
}
