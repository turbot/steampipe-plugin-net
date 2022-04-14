package net

import (
	"context"
	"fmt"
	"strings"

	"github.com/miekg/dns"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableNetDNSRecord(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_dns_record",
		Description: "DNS records associated with a given domain.",
		List: &plugin.ListConfig{
			Hydrate: tableDNSRecordList,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "domain", Require: plugin.Required, Operators: []string{"="}},
				{Name: "type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dns_server", Require: plugin.Optional, Operators: []string{"="}, CacheMatch: "exact"},
			},
		},
		Columns: []*plugin.Column{
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Domain name for the record."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the DNS record: A, CNAME, MX, etc"},
			{Name: "dns_server", Type: proto.ColumnType_STRING, Description: "Type of the DNS record: A, CNAME, MX, etc", Transform: transform.FromQual("dns_server")},
			{Name: "ip", Transform: transform.FromField("IP"), Type: proto.ColumnType_IPADDR, Description: "IP address for the record, such as for A records."},
			{Name: "target", Type: proto.ColumnType_STRING, Description: "Target of the record, such as the target address for CNAME records."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "Priority of the record, such as for MX records."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Value of the record, such as the text of a TXT record."},
			{Name: "ttl", Transform: transform.FromField("TTL"), Type: proto.ColumnType_INT, Description: "Time To Live in seconds for the record in DNS cache."},
			{Name: "serial", Type: proto.ColumnType_INT, Description: "Serial number of the record, such as for SOA records."},
			{Name: "mbox", Type: proto.ColumnType_STRING, Description: "Specifies the hostmaster email address."},
			{Name: "min_ttl", Type: proto.ColumnType_INT, Transform: transform.FromField("MinTTL"), Description: "Specifies the hostmaster email address."},
			{Name: "refresh", Type: proto.ColumnType_INT, Description: "Specifies the SOA refresh interval. The value configures how often a name server should check it's primary server to see if there has been any updates to the zone which it does by comparing Serial numbers."},
			{Name: "retry", Type: proto.ColumnType_INT, Description: "Specifies SOA retry value. The value indicates how long a name server should wait to retry an attempt to get fresh zone data from the primary name server if the first attempt should fail."},
			{Name: "expire", Type: proto.ColumnType_INT, Description: "Specifies SOA expire value. A name server will no longer consider itself Authoritative if it hasn't been able to refresh the zone data in the time limit declared in this value."},
		},
	}
}

type tableDNSRecordRow struct {
	Domain    string
	Type      string
	DNSServer string
	IP        string
	Target    string
	TTL       uint32
	Priority  uint16
	Value     string
	Serial    uint32
	Mbox      string
	MinTTL    uint32
	Refresh   uint32
	Retry     uint32
	Expire    uint32
}

func getTypeQuals(typeQualsWrapper *proto.Quals) []string {
	if typeQualsWrapper == nil {
		var allTypes []string
		return append(allTypes, "A", "AAAA", "CERT", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT")
	}
	var types []string
	typeQuals := typeQualsWrapper.Quals[0].Value
	if qualList := typeQuals.GetListValue(); qualList != nil {
		for _, q := range qualList.Values {
			types = append(types, q.GetStringValue())
		}
	} else {
		types = append(types, typeQuals.GetStringValue())
	}
	return types
}

func dnsTypeToDNSLibTypeEnum(recordType string) (uint16, error) {
	switch recordType {
	case "A":
		return dns.TypeA, nil
	case "AAAA":
		return dns.TypeAAAA, nil
	case "CERT":
		return dns.TypeCERT, nil
	case "CNAME":
		return dns.TypeCNAME, nil
	case "MX":
		return dns.TypeMX, nil
	case "NS":
		return dns.TypeNS, nil
	case "PTR":
		return dns.TypePTR, nil
	case "SOA":
		return dns.TypeSOA, nil
	case "SRV":
		return dns.TypeSRV, nil
	case "TXT":
		return dns.TypeTXT, nil
	}
	return dns.TypeANY, fmt.Errorf("Unsupported DNS record type: %s", recordType)
}

func getRecords(domain string, dnsType string, answer dns.RR) []tableDNSRecordRow {
	var records []tableDNSRecordRow
	switch typedRecord := answer.(type) {
	case *dns.A:
		records = append(records, tableDNSRecordRow{
			Domain: domain,
			Type:   dnsType,
			IP:     typedRecord.A.String(),
			TTL:    typedRecord.Hdr.Ttl,
		})
	case *dns.AAAA:
		records = append(records, tableDNSRecordRow{
			Domain: domain,
			Type:   dnsType,
			IP:     typedRecord.AAAA.String(),
			TTL:    typedRecord.Hdr.Ttl,
		})
	case *dns.CERT:
		records = append(records, tableDNSRecordRow{
			Domain: domain,
			Type:   dnsType,
			TTL:    typedRecord.Hdr.Ttl,
			Value:  typedRecord.String(),
		})
	case *dns.CNAME:
		records = append(records, tableDNSRecordRow{
			Domain: domain,
			Type:   dnsType,
			Target: typedRecord.Target,
			TTL:    typedRecord.Hdr.Ttl,
		})
	case *dns.MX:
		records = append(records, tableDNSRecordRow{
			Domain:   domain,
			Type:     dnsType,
			Priority: typedRecord.Preference,
			Target:   typedRecord.Mx,
			TTL:      typedRecord.Hdr.Ttl,
		})
	case *dns.NS:
		records = append(records, tableDNSRecordRow{
			Domain: domain,
			Type:   dnsType,
			Target: typedRecord.Ns,
			TTL:    typedRecord.Hdr.Ttl,
		})
	case *dns.PTR:
		records = append(records, tableDNSRecordRow{
			Domain: domain,
			Type:   dnsType,
			Target: typedRecord.Ptr,
			TTL:    typedRecord.Hdr.Ttl,
		})
	case *dns.SOA:
		records = append(records, tableDNSRecordRow{
			Domain:  domain,
			Type:    dnsType,
			Target:  typedRecord.Ns,
			TTL:     typedRecord.Hdr.Ttl,
			Serial:  typedRecord.Serial,
			Mbox:    typedRecord.Mbox,
			MinTTL:  typedRecord.Minttl,
			Refresh: typedRecord.Refresh,
			Retry:   typedRecord.Retry,
			Expire:  typedRecord.Expire,
		})
	case *dns.SRV:
		records = append(records, tableDNSRecordRow{
			Domain:   domain,
			Type:     dnsType,
			Priority: typedRecord.Priority,
			Target:   typedRecord.Target,
			TTL:      typedRecord.Hdr.Ttl,
		})
	case *dns.TXT:
		for _, txt := range typedRecord.Txt {
			records = append(records, tableDNSRecordRow{
				Domain: domain,
				Type:   dnsType,
				TTL:    typedRecord.Hdr.Ttl,
				Value:  txt,
			})
		}
	}
	return records
}

func tableDNSRecordList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	queryCols := d.QueryContext.Columns

	// You must pass 1 or more domain quals to the query
	if d.KeyColumnQuals["domain"] == nil {
		logger.Trace("tableDNSRecordList", "No domain quals provided")
		return nil, nil
	}
	domain := d.KeyColumnQualString("domain")

	typeQualsWrapper := d.QueryContext.UnsafeQuals["type"]
	types := getTypeQuals(typeQualsWrapper)

	c := new(dns.Client)
	// Ensure a single request of the same question, type and class at a time.
	c.SingleInflight = true
	// Use our configuration for the timeout
	c.Timeout = GetConfigTimeout(ctx, d)

	var dnsServer string
	if d.KeyColumnQuals["dns_server"] != nil {
		dnsServer = d.KeyColumnQualString("dns_server")
		if !strings.HasSuffix(dnsServer, ":53") {
			dnsServer = fmt.Sprintf("%s:53", dnsServer)
		}
	} else {
		dnsServer = GetConfigDNSServerAndPort(ctx, d)
	}

	logger.Trace("tableDNSRecordList", "Cols", queryCols)
	logger.Info("tableDNSRecordList", "Domain", domain)
	logger.Trace("tableDNSRecordList", "Types", types)
	logger.Info("tableDNSRecordList", "DNS Server", dnsServer)

	for _, dnsType := range types {
		dnsTypeEnumVal, err := dnsTypeToDNSLibTypeEnum(dnsType)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(domain), dnsTypeEnumVal)
		m.RecursionDesired = true
		r, _, err := c.Exchange(m, dnsServer)
		if err != nil {
			return nil, err
		}
		if r.Rcode != dns.RcodeSuccess {
			return nil, err
		}

		logger.Trace("tableDNSRecordList", "Question", r.Question)
		logger.Trace("tableDNSRecordList", "Answer", r.Answer)
		logger.Trace("tableDNSRecordList", "Extra", r.Extra)
		logger.Trace("tableDNSRecordList", "NS", r.Ns)

		for _, answer := range r.Answer {
			for _, record := range getRecords(domain, dnsType, answer) {
				logger.Trace("tableDNSRecordList", "Record", record)
				d.StreamListItem(ctx, record)
			}
		}
	}

	return nil, nil
}
