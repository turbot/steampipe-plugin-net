package net

import (
	"context"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type netConfig struct {
	Timeout   *int    `hcl:"timeout"`
	DNSServer *string `hcl:"dns_server"`
}

func ConfigInstance() interface{} {
	return &netConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) netConfig {
	if connection == nil || connection.Config == nil {
		return netConfig{}
	}
	config, _ := connection.Config.(netConfig)
	return config
}

func GetConfigTimeout(ctx context.Context, d *plugin.QueryData) time.Duration {
	// default to 2000ms
	ts := 2000
	config := GetConfig(d.Connection)
	if config.Timeout != nil {
		ts = *config.Timeout
	}
	return time.Millisecond * time.Duration(ts)
}

func GetConfigDNSServerAndPort(ctx context.Context, d *plugin.QueryData) string {
	// default to Google
	s := "8.8.8.8:53"
	config := GetConfig(d.Connection)
	if config.DNSServer != nil {
		s = *config.DNSServer
	}
	return s
}
