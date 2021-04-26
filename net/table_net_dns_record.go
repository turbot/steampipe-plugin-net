package net

import (
	"context"
	"fmt"

	"github.com/miekg/dns"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableNetDNSRecord(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_dns_record",
		Description: "DNS records associated with a given domain.",
		List: &plugin.ListConfig{
			Hydrate:    tableDNSRecordList,
			KeyColumns: plugin.SingleColumn("domain"),
		},
		Columns: []*plugin.Column{
			{Name: "domain", Type: proto.ColumnType_STRING},
			{Name: "type", Type: proto.ColumnType_STRING},
			{Name: "ip", Transform: transform.FromField("IP"), Type: proto.ColumnType_IPADDR},
			{Name: "target", Type: proto.ColumnType_STRING},
			{Name: "priority", Type: proto.ColumnType_INT},
			{Name: "value", Type: proto.ColumnType_STRING},
			{Name: "ttl", Transform: transform.FromField("TTL"), Type: proto.ColumnType_INT},
		},
	}
}

type tableDNSRecordRow struct {
	Domain   string
	Type     string
	IP       string
	Target   string
	TTL      uint32
	Priority uint16
	Value    string
}

func getDomainQuals(domainQualsWrapper *proto.Quals) []string {
	var domains []string
	domainQuals := domainQualsWrapper.Quals[0].Value
	if qualList := domainQuals.GetListValue(); qualList != nil {
		for _, q := range qualList.Values {
			domains = append(domains, q.GetStringValue())
		}
	} else {
		domains = append(domains, domainQuals.GetStringValue())
	}
	return domains
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
	return dns.TypeANY, fmt.Errorf("Unsupported DNS record type: %g", recordType)
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
			Domain: domain,
			Type:   dnsType,
			Target: typedRecord.Ns,
			TTL:    typedRecord.Hdr.Ttl,
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
	domainQualsWrapper := d.QueryContext.Quals["domain"]

	// You must pass 1 or more domain quals to the query
	if domainQualsWrapper == nil {
		logger.Trace("tableDNSRecordList", "No domain quals provided")
		return nil, nil
	}

	domains := getDomainQuals(domainQualsWrapper)

	typeQualsWrapper := d.QueryContext.Quals["type"]
	types := getTypeQuals(typeQualsWrapper)

	logger.Trace("tableDNSRecordList", "Cols", queryCols)
	logger.Trace("tableDNSRecordList", "Domains", domains)
	logger.Trace("tableDNSRecordList", "Types", types)

	c := new(dns.Client)
	// Ensure a single request of the same question, type and class at a time.
	c.SingleInflight = true
	// Use our configuration for the timeout
	c.Timeout = GetConfigTimeout(ctx, d)

	dnsServer := GetConfigDNSServerAndPort(ctx, d)
	logger.Trace("tableDNSRecordList", "DNS Server", dnsServer)

	for _, domain := range domains {
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
	}

	return nil, nil
}
