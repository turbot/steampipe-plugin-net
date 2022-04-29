package net

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/genkiroid/cert"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableNetCertificate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_certificate",
		Description: "Certificate details for a domain.",
		List: &plugin.ListConfig{
			Hydrate:    tableNetCertificateList,
			KeyColumns: plugin.SingleColumn("domain"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Domain name the certificate represents."},
			{Name: "common_name", Type: proto.ColumnType_STRING, Description: "Common name for the certificate."},
			{Name: "not_after", Type: proto.ColumnType_DATETIME, Description: "Time when the certificate expires. Also see not_before."},
			{Name: "is_revoked", Type: proto.ColumnType_BOOL, Description: "Indicates whether the certificate was revoked, or not."},
			// Other columns
			{Name: "chain", Type: proto.ColumnType_JSON, Description: "Certificate chain."},
			{Name: "country", Type: proto.ColumnType_STRING, Description: "Country for the certificate."},
			{Name: "dns_names", Type: proto.ColumnType_JSON, Transform: transform.FromField("DNSNames"), Description: "DNS names for the certificate."},
			{Name: "email_addresses", Type: proto.ColumnType_JSON, Description: "Email addresses for the certificate."},
			{Name: "ip_address", Type: proto.ColumnType_IPADDR, Transform: transform.FromField("IPAddress"), Description: "IP address associated with the domain."},
			{Name: "ip_addresses", Type: proto.ColumnType_JSON, Transform: transform.FromField("IPAddresses"), Description: "Array of IP addresses associated with the domain."},
			{Name: "is_ca", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsCertificateAuthority"), Description: "True if the certificate represents a certificate authority."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "Issuer of the certificate."},
			{Name: "issuer_name", Type: proto.ColumnType_STRING, Description: "Issuer of the certificate."},
			{Name: "issuing_certificate_url", Type: proto.ColumnType_JSON, Transform: transform.FromField("IssuingCertificateURL"), Description: "List of URLs of the issuing certificates."},
			{Name: "locality", Type: proto.ColumnType_STRING, Description: "Locality of the certificate."},
			{Name: "not_before", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the certificate is valid from. Also see not_after."},
			{Name: "organization", Type: proto.ColumnType_STRING, Description: "Organization of the certificate."},
			{Name: "ou", Type: proto.ColumnType_JSON, Transform: transform.FromField("OU"), Description: "Organizational Unit of the certificate."},
			{Name: "public_key_algorithm", Type: proto.ColumnType_STRING, Description: "Public key algorithm used by the certificate."},
			{Name: "public_key_length", Type: proto.ColumnType_INT, Description: ""},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Serial number of the certificate."},
			{Name: "signature_algorithm", Type: proto.ColumnType_STRING, Description: "Signature algorithm of the certificate."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the certificate."},
			{Name: "subject", Type: proto.ColumnType_STRING, Description: "Subject of the certificate."},
			{Name: "crl_distribution_points", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRLDistributionPoints"), Description: "A CRL distribution point (CDP) is a location on an LDAP directory server or Web server where a CA publishes CRLs."},
			{Name: "ocsp_server", Type: proto.ColumnType_JSON, Transform: transform.FromField("OCSPServer"), Description: "The Online Certificate Status Protocol (OCSP) is a protocol for determining the status of a digital certificate without requiring Certificate Revocation Lists (CRLs. The revocation check is by an online protocol is timely and does not require fetching large lists of revoked certificate on the client side. This test suite can be used to test OCSP Responder implementations."},
			{Name: "protocol", Type: proto.ColumnType_STRING, Hydrate: getProtocolDetails, Transform: transform.FromField("Protocol"), Description: "The TLS version used by the connection."},
			{Name: "cipher_suite", Type: proto.ColumnType_STRING, Hydrate: getProtocolDetails, Transform: transform.FromField("CipherSuite"), Description: "The cipher suite negotiated for the connection."},
		},
	}
}

// Define our own structure for certificate information since the cert
// package has multiple partial structures
type tableNetCertificateRow struct {
	// Common
	Domain     string    `json:"domain,omitempty"`
	CommonName string    `json:"common_name,omitempty"`
	NotAfter   time.Time `json:"not_after,omitempty"`
	IsRevoked  bool      `json:"is_revoked,omitempty"`
	// Other
	Chain                  []tableNetCertificateRow `json:"chain,omitempty"`
	Country                string                   `json:"country,omitempty"`
	DNSNames               []string                 `json:"dns_names,omitempty"`
	EmailAddresses         []string                 `json:"email_addresses,omitempty"`
	IPAddress              string                   `json:"ip_address,omitempty"`
	IPAddresses            []net.IP                 `json:"ip_addresses,omitempty"`
	IsCertificateAuthority bool                     `json:"is_certificate_authority,omitempty"`
	Issuer                 string                   `json:"issuer,omitempty"`
	IssuerName             string                   `json:"issuer_name,omitempty"`
	IssuingCertificateURL  []string                 `json:"issuing_certificate_url,omitempty"`
	Locality               string                   `json:"locality,omitempty"`
	NotBefore              time.Time                `json:"not_before,omitempty"`
	Organization           string                   `json:"organization,omitempty"`
	OU                     []string                 `json:"ou,omitempty"`
	PublicKeyAlgorithm     string                   `json:"public_key_algorithm,omitempty"`
	PublicKeyLength        int                      `json:"public_key_length,omitempty"`
	SignatureAlgorithm     string                   `json:"signature_algorithm,omitempty"`
	SerialNumber           string                   `json:"serial_number,omitempty"`
	State                  string                   `json:"state,omitempty"`
	Subject                string                   `json:"subject,omitempty"`
	CRLDistributionPoints  []string                 `json:"crl_distribution_points,omitempty"`
	OCSPServer             []string                 `json:"ocsp_server,omitempty"`
}

func tableNetCertificateList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Options for certificate retrieval - https://github.com/genkiroid/cert/blob/master/cert.go
	// TODO - should we support env vars here?
	cert.SkipVerify = false
	cert.UTC = true
	cert.TimeoutSeconds = 3 // short, certificates should be fast
	cert.CipherSuite = ""   // all cipher suites are supported by default

	quals := d.KeyColumnQuals
	dn := quals["domain"].GetStringValue()

	// Get the full certificate chain, so we can include it in results
	items, err := cert.NewCerts([]string{dn})
	if err != nil {
		return nil, err
	}

	// Should not happen. If it does, then assume the cert was not found.
	if len(items) <= 0 {
		return nil, nil
	}

	// Errors getting the certificate are returned in the Error field. For
	// example, a DNS timeout looking up the certificate.
	if items[0].Error != "" {
		return nil, err
	}

	// PRE: Certificate was found without error

	chain := items[0].CertChain()
	if len(chain) <= 0 {
		return nil, errors.New("Certificate chain should never be empty: " + dn)
	}

	certRows := []tableNetCertificateRow{}
	for _, i := range chain {
		c := tableNetCertificateRow{}

		// Multiple Subject fields are commonly used, so are elevated to
		// top level columns.
		//
		// In some cases (e.g. Country) multiple items are possible, but very very
		// rare, so we pull out the first item to the top level for convenience.
		// The full data is always available in the Subject field that these are
		// extracted from if needed. We considered making them into a comma separated
		// string, but decided on the simpler first item model.
		c.CommonName = i.Subject.CommonName
		if len(i.Subject.Country) > 0 {
			c.Country = i.Subject.Country[0]
		}
		if len(i.Subject.Province) > 0 {
			c.State = i.Subject.Province[0]
		}
		if len(i.Subject.Locality) > 0 {
			c.Locality = i.Subject.Locality[0]
		}
		if len(i.Subject.Organization) > 0 {
			c.Organization = i.Subject.Organization[0]
		}
		// OU is an array. Naming here is tricky, but ultimately ou feels simple
		// and common enough to be best. Also considered ous and organizational_unit(s).
		c.OU = i.Subject.OrganizationalUnit

		c.DNSNames = i.DNSNames
		c.EmailAddresses = i.EmailAddresses
		c.IPAddresses = i.IPAddresses
		c.IsCertificateAuthority = i.IsCA
		c.IssuerName = i.Issuer.CommonName
		c.Issuer = i.Issuer.String()
		c.IssuingCertificateURL = i.IssuingCertificateURL
		c.NotAfter = i.NotAfter
		c.NotBefore = i.NotBefore
		c.PublicKeyAlgorithm = i.PublicKeyAlgorithm.String()
		// Represent the serial number as 32 hex characters, with leading zeros.
		// This appears to be consistent with the Qualys SSL display.
		c.SerialNumber = fmt.Sprintf("%032x", i.SerialNumber)
		c.SignatureAlgorithm = i.SignatureAlgorithm.String()
		c.Subject = i.Subject.String()
		c.CRLDistributionPoints = i.CRLDistributionPoints
		c.OCSPServer = i.OCSPServer

		isRevoked, err := isRevokedCertificate(ctx, i.CRLDistributionPoints, c.SerialNumber)
		if err != nil {
			return nil, err
		}
		c.IsRevoked = *isRevoked

		// validDomain := true
		// verifyErr := i.VerifyHostname(dn)
		// if verifyErr != nil {
		// 	validDomain = false
		// }
		// panic(validDomain)

		var bitLen int
		switch publicKey := i.PublicKey.(type) {
		case *rsa.PublicKey:
			bitLen = publicKey.N.BitLen()
		case *ecdsa.PublicKey:
			bitLen = publicKey.Curve.Params().BitSize
		default:
		}
		c.PublicKeyLength = bitLen

		certRows = append(certRows, c)
	}

	// The first certificate in the chain is always the one we've requested.
	item := certRows[0]
	// Add the other dependent (e.g. certificate authority) certificates as a
	// single JSON array for reference. They have exactly the same format as
	// the table, so could possibly be returned as rows instead. It seemed
	// better to keep it focused on one row per domain, which is the main point
	// of certificate interaction.
	item.Chain = certRows[1:]

	// The primary certificate in the request has extra details we can pull
	// out from the request. Add those now.
	item.Domain = dn
	item.IPAddress = items[0].IP

	d.StreamListItem(ctx, item)

	return nil, nil
}

func getProtocolDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getProtocolDetails")

	data := h.Item.(tableNetCertificateRow)

	cfg := tls.Config{}
	addr := net.JoinHostPort(data.Domain, "443")
	conn, err := tls.Dial("tcp", addr, &cfg)
	if err != nil {
		return nil, errors.New("TLS connection failed: " + err.Error())
	}

	var tlsVersion string
	switch conn.ConnectionState().Version {
	case 769:
		tlsVersion = "TLS v1.0"
	case 770:
		tlsVersion = "TLS v1.1"
	case 771:
		tlsVersion = "TLS v1.2"
	case 772:
		tlsVersion = "TLS v1.3"
	}

	return map[string]string{
		"Protocol":    tlsVersion,
		"CipherSuite": tls.CipherSuiteName(conn.ConnectionState().CipherSuite),
	}, nil
}

// Checks if the certificate was revoked
func isRevokedCertificate(ctx context.Context, crlDistributionPoints []string, serialNumber string) (*bool, error) {
	plugin.Logger(ctx).Trace("isRevokedCertificate")

	isRevoked := false

	for _, crlDistributionPoint := range crlDistributionPoints {
		crlInfo, err := fetchCRL(crlDistributionPoint)
		if err != nil {
			return nil, err
		}

		// Check CRL is not outdated
		if crlInfo.TBSCertList.NextUpdate.Before(time.Now()) {
			return nil, errors.New("CRL is outdated")
		}

		// Check if the certificate is listed in Certificate Revocation List (CRL)
		for _, i := range crlInfo.TBSCertList.RevokedCertificates {
			if fmt.Sprintf("%032x", i.SerialNumber) == serialNumber {
				isRevoked = true
				return &isRevoked, nil
			}
		}
	}
	return &isRevoked, nil
}

// Fetch CRL list
func fetchCRL(url string) (*pkix.CertificateList, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, errors.New("failed to retrieve CRL")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return x509.ParseCRL(body)
}
