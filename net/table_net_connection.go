package net

import (
	"context"
	"errors"
	"net"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableNetConnection(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_connection",
		Description: "Test network connectivity to an address.",
		List: &plugin.ListConfig{
			Hydrate: tableNetConnectionList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "address"},
				{Name: "protocol", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "protocol", Type: proto.ColumnType_STRING, Description: "Protocol type: tcp, tcp4 (IPv4-only), tcp6 (IPv6-only), udp, udp4 (IPv4-only), udp6 (IPv6-only), ip, ip4 (IPv4-only), ip6 (IPv6-only), unix, unixgram or unixpacket."},
			{Name: "address", Type: proto.ColumnType_STRING, Description: "Address to connect to, as specified in https://golang.org/pkg/net/#Dial."},
			// Other columns
			{Name: "connected", Type: proto.ColumnType_BOOL, Description: "True if the connection was successful."},
			{Name: "error", Type: proto.ColumnType_STRING, Description: "Error message if the connection failed."},
			{Name: "local_address", Type: proto.ColumnType_STRING, Description: "Local address (ip:port) for the successful connection."},
			{Name: "remote_address", Type: proto.ColumnType_STRING, Description: "Remote address (ip:port) for the successful connection."},
		},
	}
}

type connectionRow struct {
	Protocol      string `json:"protocol"`
	Address       string `json:"address"`
	Connected     bool   `json:"connected"`
	Error         string `json:"error"`
	LocalAddress  string `json:"local_address"`
	RemoteAddress string `json:"remote_address"`
}

func tableNetConnectionList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	var protocol, address string
	if quals["protocol"] != nil {
		protocol = quals["protocol"].GetStringValue()
	} else {
		// Default to TCP
		protocol = "tcp"
	}
	if quals["address"] != nil {
		address = quals["address"].GetStringValue()
	} else {
		return nil, errors.New("address must be specified")
	}
	connectionResult, err := net.DialTimeout(protocol, address, GetConfigTimeout(ctx, d))
	r := connectionRow{
		Protocol: protocol,
		Address:  address,
	}
	if err == nil {
		r.Connected = true
		r.LocalAddress = connectionResult.LocalAddr().String()
		r.RemoteAddress = connectionResult.RemoteAddr().String()
	} else {
		r.Error = err.Error()
	}
	d.StreamListItem(ctx, r)
	return nil, nil
}
