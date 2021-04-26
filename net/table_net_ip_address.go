package net

import (
	"context"
	"net"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableNetIPAddress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "net_ip_address",
		List: &plugin.ListConfig{
			Hydrate: tableNetIPAddressList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("ip_address"),
			ItemFromKey: tableNetIPAddressFromKey,
			Hydrate:     tableNetIPAddressGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "ip_address", Type: proto.ColumnType_STRING, Transform: transform.FromField("IPAddress")},

			// Other columns
			{Name: "hosts", Type: proto.ColumnType_JSON, Hydrate: tableNetIPAddressGetHosts, Transform: transform.FromField("Hosts")},
		},
	}
}

type tableNetIPAddressRow struct {
	IPAddress string
}

func tableNetIPAddressFromKey(_ context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	ip_address := quals["ip_address"].GetStringValue()
	item := &tableNetIPAddressRow{
		IPAddress: ip_address,
	}
	return item, nil
}

func tableNetIPAddressList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO - Needs to support quals so that joins, select...in, etc can work
	return nil, nil
}

func tableNetIPAddressGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	i := h.Item.(*tableNetIPAddressRow)
	return i, nil
}

func tableNetIPAddressGetHosts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*tableNetIPAddressRow)
	hosts, err := net.LookupAddr(item.IPAddress)
	if err != nil {
		if e, ok := err.(*net.DNSError); ok {
			if e.IsNotFound {
				return struct{ Hosts []string }{}, nil
			}
		}
		return nil, err
	}
	return struct{ Hosts []string }{Hosts: hosts}, nil
}
