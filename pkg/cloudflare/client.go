package cloudflare

import (
	"context"

	cloudflare "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

// Interface is an interface for the
// cloudflare SDK.
type Interface interface {
	ZoneIDByName(name string) (string, error)
	ListDNSRecords(ctx context.Context, zoneID string, params dns.RecordListParams) ([]dns.RecordResponse, error)
	UpdateDNSRecord(ctx context.Context, zoneID string, recordID string, params dns.RecordUpdateParams) (dns.RecordResponse, error)
}

type Client struct {
	api *cloudflare.Client
}

// Make sure Client struct
// implements Interface interface
var _ Interface = &Client{}

func NewClient(token string) *Client {
	client := cloudflare.NewClient(
		option.WithAPIToken(token),
	)
	return &Client{api: client}
}
