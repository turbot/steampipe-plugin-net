package net

import (
	"context"
	"errors"
	"net"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableNetConnection(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "net_connection",
		List: &plugin.ListConfig{
			Hydrate:    tableNetConnectionList,
			KeyColumns: plugin.AllColumns([]string{"network", "address"}),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "network", Type: proto.ColumnType_STRING, Description: "Network type: tcp, tcp4 (IPv4-only), tcp6 (IPv6-only), udp, udp4 (IPv4-only), udp6 (IPv6-only), ip, ip4 (IPv4-only), ip6 (IPv6-only), unix, unixgram or unixpacket."},
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
	Network       string `json:"network"`
	Address       string `json:"address"`
	Connected     bool   `json:"connected"`
	Error         string `json:"error"`
	LocalAddress  string `json:"local_address"`
	RemoteAddress string `json:"remote_address"`
}

func tableNetConnectionList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	var network, address string
	// NOTE: This is designed for "AnyColumn" instead of "AllColumns" to make usage easier,
	// but AnyColumn currently stops after it has a single value rather than collecting all
	// that are specified.
	if quals["network"] != nil {
		network = quals["network"].GetStringValue()
	} else {
		// Default to TCP
		network = "tcp"
	}
	if quals["address"] != nil {
		address = quals["address"].GetStringValue()
	} else {
		return nil, errors.New("address must be specified")
	}
	connectionResult, err := net.DialTimeout(network, address, GetConfigTimeout(ctx, d))
	r := connectionRow{
		Network: network,
		Address: address,
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
