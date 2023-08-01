package net

import (
	"context"
	"net"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableNetDNSReverse(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_dns_reverse",
		Description: "Reverse DNS lookup from an IP address.",
		List: &plugin.ListConfig{
			Hydrate:    tableNetDNSReverseList,
			KeyColumns: plugin.SingleColumn("ip_address"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "ip_address", Type: proto.ColumnType_IPADDR, Transform: transform.FromField("IPAddress"), Description: "IP address to lookup."},
			// Other columns
			{Name: "domains", Type: proto.ColumnType_JSON, Description: "Domain names associated with the IP address."},
		},
	}
}

type tableNetDNSReverseRow struct {
	IPAddress string   `json:"ip_address"`
	Domains   []string `json:"domains"`
}

func tableNetDNSReverseList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	ip := quals["ip_address"].GetInetValue().GetAddr()
	result := tableNetDNSReverseRow{
		IPAddress: ip,
		Domains:   []string{},
	}
	domains, err := net.LookupAddr(ip)
	if err != nil {
		if e, ok := err.(*net.DNSError); ok {
			if e.IsNotFound {
				return result, nil
			}
		}
		return nil, err
	}
	result.Domains = domains
	d.StreamListItem(ctx, result)
	return nil, nil
}
