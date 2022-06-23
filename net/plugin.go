package net

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-net",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"net_certificate":    tableNetCertificate(ctx),
			"net_connection":     tableNetConnection(ctx),
			"net_dns_record":     tableNetDNSRecord(ctx),
			"net_dns_reverse":    tableNetDNSReverse(ctx),
			"net_http_request":   tableNetHTTPRequest(),
			"net_tls_connection": tableNetTLSConnection(ctx),
		},
	}
	return p
}
