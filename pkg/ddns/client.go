package ddns

import (
	"github.com/cguertin14/ddns/pkg/cloudflare"
	gh "github.com/cguertin14/ddns/pkg/github"
)

type Client struct {
	cloudflare cloudflare.Interface
	github     gh.Interface
}

type Dependencies struct {
	Cloudflare cloudflare.Interface
	Github     gh.Interface
}

func NewClient(deps Dependencies) *Client {
	return &Client{deps.Cloudflare, deps.Github}
}
