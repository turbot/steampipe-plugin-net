package constants

// A map of TLS versions, along with their IDs
var TLSVersions = map[string]uint16{
	"TLS v1.0": 0x0301,
	"TLS v1.1": 0x0302,
	"TLS v1.2": 0x0303,
	"TLS v1.3": 0x0304,
	"SSL v3":   0x0300,
}
