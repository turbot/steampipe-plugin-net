package net

import (
	"context"
	"crypto/tls"

	"golang.org/x/exp/slices"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
)

func getQualListValues(ctx context.Context, quals map[string]*proto.QualValue, qualName string) []string {
	values := make([]string, 0)
	if quals[qualName].GetStringValue() != "" {
		values = append(values, quals[qualName].GetStringValue())
	} else if quals[qualName].GetListValue() != nil {
		for _, value := range quals[qualName].GetListValue().Values {
			str := value.GetStringValue()
			values = append(values, str)
		}
	}
	return values
}

// List all cipher suites supported by `crypto/tls` package
func cipherSuites() []*tls.CipherSuite {
	allCiphers := tls.CipherSuites()

	// also append insecure ciphers
	allCiphers = append(allCiphers, tls.InsecureCipherSuites()...)
	return allCiphers
}

// List all cipher suites supported by TLS v1.3
func cipherSuitesTLS13() []string {
	allCiphers := cipherSuites()
	var ciphersTLS13 []string

	for _, i := range allCiphers {
		if slices.Contains(i.SupportedVersions, tls.VersionTLS13) {
			ciphersTLS13 = append(ciphersTLS13, i.Name)
		}
	}
	return ciphersTLS13
}

// List all cipher suites supported by TLS v1.2
func cipherSuitesTLS12() []string {
	allCiphers := cipherSuites()
	var ciphersTLS12 []string

	for _, i := range allCiphers {
		if slices.Contains(i.SupportedVersions, tls.VersionTLS12) {
			ciphersTLS12 = append(ciphersTLS12, i.Name)
		}
	}
	return ciphersTLS12
}

// List all cipher suites supported by TLS v1.0 - TLS v1.1
func cipherSuitesUptoTLS11() []string {
	allCiphers := cipherSuites()
	var ciphersUptoTLS11 []string

	for _, i := range allCiphers {
		if slices.Contains(i.SupportedVersions, tls.VersionTLS12) || slices.Contains(i.SupportedVersions, tls.VersionTLS11) || slices.Contains(i.SupportedVersions, tls.VersionTLS10) {
			ciphersUptoTLS11 = append(ciphersUptoTLS11, i.Name)
		}
	}
	return ciphersUptoTLS11
}

// Check if given cipher suite is supported by the given protocol version
func isSupported(protocol string, cipher string) bool {
	switch protocol {
	case "TLS v1.0":
	case "TLS v1.1":
		ciphers := cipherSuitesUptoTLS11()
		return slices.Contains(ciphers, cipher)
	case "TLS v1.2":
		ciphers := cipherSuitesTLS12()
		return slices.Contains(ciphers, cipher)
	case "TLS v1.3":
		ciphers := cipherSuitesTLS13()
		return slices.Contains(ciphers, cipher)
	}
	return false
}
