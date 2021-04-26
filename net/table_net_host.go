package net

import (
	"context"
	"net"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableNetHost(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "net_host",
		List: &plugin.ListConfig{
			Hydrate: tableNetHostList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("host"),
			ItemFromKey: tableNetHostFromKey,
			Hydrate:     tableNetHostGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "host", Type: proto.ColumnType_STRING},
			// Other columns
			{Name: "cname", Type: proto.ColumnType_STRING, Hydrate: tableNetHostGetCNAME, Transform: transform.FromField("CNAME")},
			{Name: "ip_addresses", Type: proto.ColumnType_JSON, Hydrate: tableNetHostGetIPAddresses, Transform: transform.FromField("IPAddresses")},
			{Name: "ns", Type: proto.ColumnType_JSON, Hydrate: tableNetHostGetNSRecords, Transform: transform.FromField("Hosts")},
			{Name: "txt", Type: proto.ColumnType_JSON, Hydrate: tableNetHostGetTXTRecords, Transform: transform.FromField("TXT")},
		},
	}
}

type tableNetHostRow struct {
	Host string
}

func tableNetHostFromKey(_ context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	host := quals["host"].GetStringValue()
	item := &tableNetHostRow{
		Host: host,
	}
	return item, nil
}

func tableNetHostList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO - Needs to support quals so that joins, select...in, etc can work
	return nil, nil
}

func tableNetHostGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	i := h.Item.(*tableNetHostRow)
	return i, nil
}

func tableNetHostGetIPAddresses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*tableNetHostRow)
	ipAddresses, err := net.LookupHost(item.Host)
	if err != nil {
		if e, ok := err.(*net.DNSError); ok {
			if e.IsNotFound {
				return struct{ IpAddresses []string }{}, nil
			}
		}
		return nil, err
	}
	return struct{ IPAddresses []string }{IPAddresses: ipAddresses}, nil
}

func tableNetHostGetTXTRecords(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*tableNetHostRow)
	txts, err := net.LookupTXT(item.Host)
	if err != nil {
		if e, ok := err.(*net.DNSError); ok {
			if e.IsNotFound {
				return struct{ TXT []string }{}, nil
			}
		}
		return nil, err
	}
	return struct{ TXT []string }{TXT: txts}, nil
}

func tableNetHostGetNSRecords(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*tableNetHostRow)
	records, err := net.LookupNS(item.Host)
	if err != nil {
		if e, ok := err.(*net.DNSError); ok {
			if e.IsNotFound {
				return struct{ Hosts []string }{}, nil
			}
		}
		return nil, err
	}
	hosts := []string{}
	for _, i := range records {
		hosts = append(hosts, i.Host)
	}
	return struct{ Hosts []string }{Hosts: hosts}, nil
}

func tableNetHostGetCNAME(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*tableNetHostRow)
	cname, err := net.LookupCNAME(item.Host)
	if err != nil {
		if e, ok := err.(*net.DNSError); ok {
			if e.IsNotFound {
				return struct{ CNAME string }{}, nil
			}
		}
		return nil, err
	}
	return struct{ CNAME string }{CNAME: cname}, nil
}
