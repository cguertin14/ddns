package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

// Interface is an interface for the
// cloudflare SDK.
type Interface interface {
	ZoneIDByName(name string) (string, error)
	ListDNSRecords(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.ListDNSRecordsParams) ([]cloudflare.DNSRecord, *cloudflare.ResultInfo, error)
	UpdateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.UpdateDNSRecordParams) error
}

type Client struct {
	api Interface
}

// Make sure Client struct
// implements Interface interface
var _ Interface = &Client{}

func NewClient(token string) (*Client, error) {
	client, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, err
	}
	return &Client{api: client}, nil
}
