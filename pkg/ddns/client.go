package ddns

import (
	"github.com/cguertin14/ddns/pkg/cloudflare"
)

type Client struct {
	cloudflare cloudflare.Interface
}

type Dependencies struct {
	Cloudflare cloudflare.Interface
}

func NewClient(deps Dependencies) *Client {
	return &Client{deps.Cloudflare}
}
