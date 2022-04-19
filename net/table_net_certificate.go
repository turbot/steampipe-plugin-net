package net

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/genkiroid/cert"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
			{Name: "not_after", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the certificate expires. Also see not_before."},
			// Other columns
			{Name: "chain", Type: proto.ColumnType_JSON, Description: "Certificate chain."},
			{Name: "country", Type: proto.ColumnType_STRING, Description: "Country for the certificate."},
			{Name: "dns_names", Type: proto.ColumnType_JSON, Transform: transform.FromField("DNSNames"), Description: "DNS names for the certificate."},
			{Name: "email_addresses", Type: proto.ColumnType_JSON, Description: "Email addresses for the certificate."},
			{Name: "ip_address", Type: proto.ColumnType_IPADDR, Transform: transform.FromField("IPAddress"), Description: "IP address associated with the domain."},
			{Name: "ip_addresses", Type: proto.ColumnType_JSON, Transform: transform.FromField("IPAddresses"), Description: "Array of IP addresses associated with the domain."},
			{Name: "is_ca", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsCA"), Description: "True if the certificate represents a certificate authority."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "Issuer of the certificate."},
			{Name: "issuing_certificate_url", Type: proto.ColumnType_JSON, Transform: transform.FromField("IssuingCertificateURL"), Description: "List of URLs of the issuing certificates."},
			{Name: "locality", Type: proto.ColumnType_STRING, Description: "Locality of the certificate."},
			{Name: "not_before", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the certificate is valid from. Also see not_after."},
			{Name: "organization", Type: proto.ColumnType_STRING, Description: "Organization of the certificate."},
			{Name: "ou", Type: proto.ColumnType_JSON, Transform: transform.FromField("OU"), Description: "Organizational Unit of the certificate."},
			{Name: "public_key_algorithm", Type: proto.ColumnType_STRING, Description: "Public key algorithm used by the certificate."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Serial number of the certificate."},
			{Name: "signature_algorithm", Type: proto.ColumnType_STRING, Description: "Signature algorithm of the certificate."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the certificate."},
			{Name: "subject", Type: proto.ColumnType_STRING, Description: "Subject of the certificate."},
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
	// Other
	Chain                  []tableNetCertificateRow `json:"chain,omitempty"`
	Country                string                   `json:"country,omitempty"`
	DNSNames               []string                 `json:"dns_names,omitempty"`
	EmailAddresses         []string                 `json:"email_addresses,omitempty"`
	IPAddress              string                   `json:"ip_address,omitempty"`
	IPAddresses            []net.IP                 `json:"ip_addresses,omitempty"`
	IsCertificateAuthority bool                     `json:"is_certificate_authority,omitempty"`
	Issuer                 string                   `json:"issuer,omitempty"`
	IssuingCertificateURL  []string                 `json:"issuing_certificate_url,omitempty"`
	Locality               string                   `json:"locality,omitempty"`
	NotBefore              time.Time                `json:"not_before,omitempty"`
	Organization           string                   `json:"organization,omitempty"`
	OU                     []string                 `json:"ou,omitempty"`
	PublicKeyAlgorithm     string                   `json:"public_key_algorithm,omitempty"`
	SignatureAlgorithm     string                   `json:"signature_algorithm,omitempty"`
	SerialNumber           string                   `json:"serial_number,omitempty"`
	State                  string                   `json:"state,omitempty"`
	Subject                string                   `json:"subject,omitempty"`
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
