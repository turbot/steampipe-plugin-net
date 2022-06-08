package constants

import "crypto/tls"

// A map of TLS versions, along with their IDs
var TLSVersions = map[string]uint16{
	"TLS v1.0": tls.VersionTLS10,
	"TLS v1.1": tls.VersionTLS11,
	"TLS v1.2": tls.VersionTLS12,
	"TLS v1.3": tls.VersionTLS13,
}
