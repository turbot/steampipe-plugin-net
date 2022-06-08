package net

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/turbot/steampipe-plugin-net/constants"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableNetTLSConnection(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_tls_connection",
		Description: "Check server TLS connectivity to an address.",
		List: &plugin.ListConfig{
			Hydrate: tableNetTLSConnectionList,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "address", Require: plugin.Required, Operators: []string{"="}},
				{Name: "version", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "cipher_suite_name", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "address", Type: proto.ColumnType_STRING, Description: "Address to connect to, as specified in https://golang.org/pkg/net/#Dial.", Transform: transform.FromQual("address")},
			{Name: "server_name", Type: proto.ColumnType_STRING, Description: "The server name indication extension sent by the client."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "The TLS version used by the connection."},
			{Name: "cipher_suite_name", Type: proto.ColumnType_STRING, Description: "The cipher suite negotiated for the connection."},
			{Name: "cipher_suite_id", Type: proto.ColumnType_STRING, Description: "The ID of the cipher suite."},
			{Name: "handshake_completed", Type: proto.ColumnType_BOOL, Description: "True if the handshake has concluded."},
			{Name: "error", Type: proto.ColumnType_STRING, Description: "Error message if the connection failed."},
			{Name: "fallback_scsv_supported", Type: proto.ColumnType_BOOL, Description: "True if the TLS fallback SCSV is enabled to prevent protocol downgrade attacks.", Hydrate: checkFallbackSCSVSupport, Transform: transform.FromValue()},
			{Name: "alpn_supported", Type: proto.ColumnType_BOOL, Description: "True if the ALPN is supported.", Hydrate: checkAPLNSupport, Transform: transform.FromValue()},
			{Name: "local_address", Type: proto.ColumnType_STRING, Description: "Local address (ip:port) for the successful connection."},
			{Name: "remote_address", Type: proto.ColumnType_STRING, Description: "Remote address (ip:port) for the successful connection."},
		},
	}
}

type tlsConnectionRow struct {
	Version            string `json:"version"`
	CipherSuiteName    string `json:"cipher_suite_name"`
	CipherSuiteID      string `json:"cipher_suite_id"`
	ServerName         string `json:"server_name"`
	HandshakeCompleted bool   `json:"handshake_completed"`
	Error              string `json:"error"`
	LocalAddress       string `json:"local_address"`
	RemoteAddress      string `json:"remote_address"`
}

//// LIST FUNCTION

func tableNetTLSConnectionList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Trace("tableNetTLSConnectionList")

	// You must pass 1 or more domain quals to the query
	quals := d.KeyColumnQuals
	address := d.KeyColumnQualString("address")

	// By default, consider all available protocols and ciphers
	var ciphers []string
	for _, c := range cipherSuites() {
		ciphers = append(ciphers, c.Name)
	}
	protocols := []string{"TLS v1.3", "TLS v1.2", "TLS v1.1", "TLS v1.0"}

	// Check for additional quals
	if d.KeyColumnQuals["version"] != nil {
		protocols = getQualListValues(ctx, quals, "version")
	}
	if d.KeyColumnQuals["cipher_suite_name"] != nil {
		ciphers = getQualListValues(ctx, quals, "cipher_suite_name")
	}

	var wg sync.WaitGroup
	for _, protocol := range protocols {
		for _, cipher := range ciphers {
			wg.Add(1)
			go func(p string, c string) {
				row := getTLSConnectionRowData(ctx, address, p, c)
				d.StreamListItem(ctx, row)

				wg.Done()
			}(protocol, cipher)
		}
	}
	wg.Wait()

	return nil, nil
}

func getTLSConnectionRowData(ctx context.Context, address string, protocol string, cipher string) tlsConnectionRow {
	r := tlsConnectionRow{
		Version:         protocol,
		CipherSuiteName: cipher,
		CipherSuiteID:   fmt.Sprintf("0x%04x", constants.CipherSuites[cipher]),
	}
	if isSupported(protocol, cipher) {
		conn, err := getTLSConnection(ctx, address, protocol, cipher)
		if err == nil {
			if conn != nil {
				r.ServerName = conn.ConnectionState().ServerName
				r.HandshakeCompleted = conn.ConnectionState().HandshakeComplete
				r.LocalAddress = conn.LocalAddr().String()
				r.RemoteAddress = conn.RemoteAddr().String()
			}
		} else {
			r.Error = err.Error()
		}
	} else {
		r.Error = "unsupported protocol-cipher combination"
	}

	return r
}

// Initiate a TLS handshake and return TLS connection
func getTLSConnection(ctx context.Context, address string, protocol string, cipher string) (*tls.Conn, error) {
	cfg := tls.Config{
		Rand:               rand.Reader,
		InsecureSkipVerify: true,
	}

	if protocol != "" {
		if _, ok := constants.TLSVersions[protocol]; !ok {
			return nil, fmt.Errorf("%s is not a valid protocol version. Possible values are: TLS v1.0, TLS v1.1, TLS v1.2 and TLS v1.3", protocol)
		}
		cfg.MaxVersion = constants.TLSVersions[protocol]
		cfg.MinVersion = constants.TLSVersions[protocol]
	}

	if cipher != "" {
		if _, ok := constants.CipherSuites[cipher]; !ok {
			return nil, fmt.Errorf("%s is not a valid cipher suite", cipher)
		}
		cfg.CipherSuites = []uint16{constants.CipherSuites[cipher]}
	}

	conn, err := tls.DialWithDialer(&net.Dialer{}, "tcp", address, &cfg)
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.tableNetCertificateList", "TLS connection failed: ", err)
		return nil, errors.New(err.Error())
	}

	return conn, nil
}

// Check if TLS Fallback Signaling Cipher Suite Value supported
func checkFallbackSCSVSupport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(tlsConnectionRow)

	// Return nil, if connection is failed
	if data.Error != "" {
		return nil, nil
	}

	cfg := tls.Config{
		Rand:               rand.Reader,
		InsecureSkipVerify: true,
		CipherSuites:       []uint16{constants.CipherSuites["TLS_FALLBACK_SCSV"]},
	}

	addr := d.KeyColumnQualString("address")

	conn, err := tls.DialWithDialer(&net.Dialer{}, "tcp", addr, &cfg)
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.checkFallbackSCSVSupport", "check_fallback_scsv_support", err)
		return false, nil
	}

	if conn == nil {
		return false, nil
	}

	return true, nil
}

// Check if Application-Layer Protocol Negotiation (ALPN) supported
func checkAPLNSupport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(tlsConnectionRow)

	// Return nil, if connection is failed
	if data.Error != "" {
		return nil, nil
	}

	cfg := tls.Config{
		Rand:               rand.Reader,
		InsecureSkipVerify: true,
		NextProtos:         []string{"http/0.9", "http/1.0", "http/1.1", "spdy/1", "spdy/2", "spdy/3", "stun.turn", "stun.nat-discovery", "h2", "h2c", "webrtc", "c-webrtc", "ftp", "imap", "pop3", "managesieve", "coap", "xmpp-client", "xmpp-server", "acme-tls/1", "mqtt", "dot", "ntske/1", "sunrpc", "h3", "smb", "irc", "nntp", "nnsp", "doq"}, // A list of all available TLS ALPN protocol. Please refer: https://www.iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml#alpn-protocol-ids
	}

	addr := d.KeyColumnQualString("address")

	conn, err := tls.DialWithDialer(&net.Dialer{}, "tcp", addr, &cfg)
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.checkAPLNSupport", "check_tls_alpn_support", err)
		return nil, err
	}

	if conn.ConnectionState().HandshakeComplete && conn.ConnectionState().NegotiatedProtocol != "" {
		return true, nil
	}

	return nil, nil
}

// Parse TLS version to human-readable format
func parseTLSVersion(p uint16) string {
	switch p {
	case tls.VersionTLS10:
		return "TLS v1.0"
	case tls.VersionTLS11:
		return "TLS v1.1"
	case tls.VersionTLS12:
		return "TLS v1.2"
	case tls.VersionTLS13:
		return "TLS v1.3"
	case tls.VersionSSL30:
		return "SSL v3"
	default:
		return "unknown"
	}
}
