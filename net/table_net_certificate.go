package net

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"golang.org/x/crypto/ocsp"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type OCSP struct {
	StatusString           string     `json:"status"`
	RevokedAt              *time.Time `json:"revoked_at,omitempty"`
	RevocationReasonString string     `json:"revocation_reason,omitempty"`
}

//// TABLE DEFINITION

func tableNetCertificate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_certificate",
		Description: "Certificate details for a domain.",
		List: &plugin.ListConfig{
			Hydrate: tableNetCertificateList,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "domain", Require: plugin.Required, Operators: []string{"="}},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Domain name the certificate represents."},
			{Name: "common_name", Type: proto.ColumnType_STRING, Description: "Common name for the certificate."},
			{Name: "not_after", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the certificate expires. Also see not_before."},
			{Name: "is_revoked", Type: proto.ColumnType_BOOL, Hydrate: getRevocationInformation, Description: "Indicates whether the certificate was revoked, or not."},
			{Name: "transparent", Type: proto.ColumnType_BOOL, Hydrate: getCertificateTransparencyLogs, Transform: transform.FromValue(), Description: "Indicates whether certificate is visible in certificate transparency logs."},
			{Name: "is_ca", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsCertificateAuthority"), Description: "True if the certificate represents a certificate authority."},
			// Other columns
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Serial number of the certificate."},
			{Name: "subject", Type: proto.ColumnType_STRING, Description: "Subject of the certificate."},
			{Name: "public_key_algorithm", Type: proto.ColumnType_STRING, Description: "Public key algorithm used by the certificate."},
			{Name: "public_key_length", Type: proto.ColumnType_INT, Description: "Specifies the size of the key."},
			{Name: "signature_algorithm", Type: proto.ColumnType_STRING, Description: "Signature algorithm of the certificate."},
			{Name: "ip_address", Type: proto.ColumnType_IPADDR, Transform: transform.FromField("IPAddress"), Description: "IP address associated with the domain."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "Issuer of the certificate."},
			{Name: "issuer_name", Type: proto.ColumnType_STRING, Description: "Issuer of the certificate."},
			{Name: "chain", Type: proto.ColumnType_JSON, Description: "Certificate chain."},
			{Name: "country", Type: proto.ColumnType_STRING, Description: "Country for the certificate."},
			{Name: "dns_names", Type: proto.ColumnType_JSON, Transform: transform.FromField("DNSNames"), Description: "DNS names for the certificate."},
			{Name: "crl_distribution_points", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRLDistributionPoints"), Description: "A CRL distribution point (CDP) is a location on an LDAP directory server or Web server where a CA publishes CRLs."},
			{Name: "ocsp_server", Type: proto.ColumnType_JSON, Transform: transform.FromField("OCSPServer"), Description: "The Online Certificate Status Protocol (OCSP) is a protocol for determining the status of a digital certificate without requiring Certificate Revocation Lists (CRLs. The revocation check is by an online protocol is timely and does not require fetching large lists of revoked certificate on the client side. This test suite can be used to test OCSP Responder implementations."},
			{Name: "ocsp", Type: proto.ColumnType_JSON, Hydrate: getRevocationInformation, Transform: transform.FromField("OCSP"), Description: "OCSP server details about the certificate."},
			{Name: "email_addresses", Type: proto.ColumnType_JSON, Description: "Email addresses for the certificate."},
			{Name: "ip_addresses", Type: proto.ColumnType_JSON, Transform: transform.FromField("IPAddresses"), Description: "Array of IP addresses associated with the domain."},
			{Name: "issuing_certificate_url", Type: proto.ColumnType_JSON, Transform: transform.FromField("IssuingCertificateURL"), Description: "List of URLs of the issuing certificates."},
			{Name: "locality", Type: proto.ColumnType_STRING, Description: "Locality of the certificate."},
			{Name: "not_before", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the certificate is valid from. Also see not_after."},
			{Name: "organization", Type: proto.ColumnType_STRING, Description: "Organization of the certificate."},
			{Name: "ou", Type: proto.ColumnType_JSON, Transform: transform.FromField("OU"), Description: "Organizational Unit of the certificate."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the certificate."},
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

	rawCert *x509.Certificate `json:"-"`
}

type Cert struct {
	IssuerCaID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
}

//// LIST FUNCTION

func tableNetCertificateList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Trace("tableNetCertificateList")

	// You must pass 1 or more domain quals to the query
	if d.KeyColumnQuals["domain"] == nil {
		plugin.Logger(ctx).Trace("tableDNSRecordList", "No domain quals provided")
		return nil, nil
	}
	dn := d.KeyColumnQualString("domain")

	// Create TLS config
	cfg := tls.Config{
		Rand:               rand.Reader,
		InsecureSkipVerify: true,
	}

	addr := net.JoinHostPort(dn, "443")
	dialer := &net.Dialer{
		Timeout: time.Duration(3) * time.Second, // short, certificates should be fast
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &cfg)
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.tableNetCertificateList", "TLS connection failed: ", err)
		return nil, errors.New("TLS connection failed: " + err.Error())
	}
	items := conn.ConnectionState().PeerCertificates

	// Should not happen. If it does, then assume the cert was not found.
	if len(items) <= 0 {
		return nil, nil
	}

	chain := items
	if len(chain) <= 0 {
		return nil, errors.New("Certificate chain can not be empty: " + dn)
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
		if i.Issuer.CommonName != "" {
			c.IssuerName = i.Issuer.CommonName
		} else {
			if len(i.Issuer.Organization) > 0 && len(i.Issuer.OrganizationalUnit) > 0 {
				c.IssuerName = fmt.Sprintf("%s / %s", i.Issuer.Organization[0], i.Issuer.OrganizationalUnit[0])
			}
		}
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

		var bitLen int
		switch publicKey := i.PublicKey.(type) {
		case *rsa.PublicKey:
			bitLen = publicKey.N.BitLen()
		case *ecdsa.PublicKey:
			bitLen = publicKey.Curve.Params().BitSize
		default:
		}
		c.PublicKeyLength = bitLen
		c.rawCert = i

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
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.tableNetCertificateList", "error retrieving host from network address", err)
		return nil, fmt.Errorf("failed to extract host from network address: %v", err)
	}
	item.IPAddress = host

	d.StreamListItem(ctx, item)

	return nil, nil
}

//// HYDRATE FUNCTIONS

// Check if certificate is transparent
func getCertificateTransparencyLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(tableNetCertificateRow)
	domainName := data.CommonName
	serialNumber := data.SerialNumber

	// crt.sh is a web interface to a distributed database called the certificate transparency logs.
	// To validate if domain certificate is transparent, check your certificate in certificate transparency logs
	var certs []Cert
	baseURL := "https://crt.sh/"
	url := fmt.Sprintf("%s?q=%s&match==&output=json", baseURL, domainName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve certificate transparency log: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate transparency log: %v", err)
	}

	err = json.Unmarshal(body, &certs)
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.getCertificateTransparencyLogs", "unmarshal_error", err)
		return nil, nil
	}

	// If certificate record found in certificate transparency logs, return transparent as true
	isTransparent := false
	for _, c := range certs {
		if c.SerialNumber == serialNumber {
			isTransparent = true
			break
		}
	}
	return isTransparent, nil
}

// Check certificate revocation information
// This function checks both CRL and OCSP server to check for certificate revocation status
func getRevocationInformation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("net_certificate.getRevocationInformation")

	certRevocationInfo := map[string]interface{}{}
	data := h.Item.(tableNetCertificateRow)

	// Check Certificate Revocation List (CRL) to verify certificate revocation status
	isRevokedByCAOrOwner, err := isCertificateRevokedByCA(ctx, data.CRLDistributionPoints, data.SerialNumber)
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.getCertificateTransparencyLogs", "error getting revocation information from CRL", err)
	}

	// Check Online Certificate Status Protocol (OCSP) to verify certificate revocation status
	ocspCertificateRevocationInfo, err := fetchOCSPDetails(ctx, data)
	if err != nil {
		plugin.Logger(ctx).Error("net_certificate.getCertificateTransparencyLogs", "error getting revocation information from OCSP server", err)
	}

	if ocspCertificateRevocationInfo != nil {
		certRevocationInfo["OCSP"] = ocspCertificateRevocationInfo
	}

	if isRevokedByCAOrOwner == nil && ocspCertificateRevocationInfo == nil {
		return nil, errors.New("unable to retrieve certificate revocation information")
	}

	isRevoked := false
	if *isRevokedByCAOrOwner || ocspCertificateRevocationInfo.StatusString == "revoked" {
		isRevoked = true
	}
	certRevocationInfo["IsRevoked"] = isRevoked

	return certRevocationInfo, nil
}

// getOCSPDetails queries the ocsp_server as given in the certificate and fetches the ocsp status
// adapted from https://github.com/crtsh/ocsp_monitor/blob/e5a2a490acb05dafb0d46f4d0f32c89b1e91a1b5/ocsp_monitor.go#L233
func fetchOCSPDetails(ctx context.Context, data tableNetCertificateRow) (*OCSP, error) {

	plugin.Logger(ctx).Trace("net_certificate.fetchOCSPDetails")

	ocspData := OCSP{}

	if len(data.Chain) == 0 {
		plugin.Logger(ctx).Trace("could not find a certificate chain")
		return nil, nil
	}

	cert := data.rawCert
	// the first element of the chain is the certificate of the issuer
	issuerCert := data.Chain[0].rawCert

	if len(cert.OCSPServer) == 0 {
		plugin.Logger(ctx).Trace("could not find OCSP verification server")
		return nil, nil
	}
	requestUrl := cert.OCSPServer[0]

	ocspBytes, err := ocsp.CreateRequest(cert, issuerCert, &ocsp.RequestOptions{})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(ocspBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/ocsp-request")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "Steampipe")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OCSP response: %v", err)
	}

	ocspResponse, err := ocsp.ParseResponseForCert(body, cert, issuerCert)
	if err != nil {
		return nil, err
	}

	ocspData = OCSP{}
	switch ocspResponse.Status {
	case ocsp.Good:
		ocspData.StatusString = "good"
	case ocsp.Unknown:
		ocspData.StatusString = "unknown"
	case ocsp.Revoked:
		ocspData.RevokedAt = &ocspResponse.RevokedAt
		ocspData.RevocationReasonString = getOCSPRevocationReasonString(ocspResponse.RevocationReason)
		ocspData.StatusString = "revoked"
	default:
		ocspData.StatusString = "unexpected"
	}

	return &ocspData, nil
}

// Checks if the certificate was revoked
func isCertificateRevokedByCA(ctx context.Context, crlDistributionPoints []string, serialNumber string) (*bool, error) {
	plugin.Logger(ctx).Trace("isCertificateRevokedByCA")

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

// Parse OCSP revocation status to a human-readable format
func getOCSPRevocationReasonString(reasonCode int) string {
	switch reasonCode {
	case ocsp.Unspecified:
		return "unspecified"
	case ocsp.KeyCompromise:
		return "key-compromise"
	case ocsp.CACompromise:
		return "ca-compromise"
	case ocsp.AffiliationChanged:
		return "affiliation-changed"
	case ocsp.Superseded:
		return "superseded"
	case ocsp.CessationOfOperation:
		return "cessation-of-operation"
	case ocsp.CertificateHold:
		return "certificate-hold"
	case ocsp.RemoveFromCRL:
		return "remove-from-crl"
	case ocsp.PrivilegeWithdrawn:
		return "privilefe-withdrawn"
	case ocsp.AACompromise:
		return "aa-compromise"
	}
	return "unknown"
}
