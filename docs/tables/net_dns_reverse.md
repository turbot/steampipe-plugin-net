---
title: "Steampipe Table: net_dns_reverse - Query OCI DNS Reverse Zones using SQL"
description: "Allows users to query DNS Reverse Zones in OCI, specifically the reverse DNS lookup information, providing insights into IP address mapping."
---

# Table: net_dns_reverse - Query OCI DNS Reverse Zones using SQL

Oracle Cloud Infrastructure's DNS service is a scalable, reliable, and managed Domain Name System (DNS) solution. It enables developers and businesses to route end users to internet applications by translating human-readable domain names (like www.example.com) into numeric IP addresses (like 192.0.2.1) that computers use to connect to each other. A Reverse DNS (rDNS) is the determination of a domain name associated with an IP address via querying DNS.

## Table Usage Guide

The `net_dns_reverse` table provides insights into Reverse DNS Zones within Oracle Cloud Infrastructure's DNS service. As a network administrator, explore reverse DNS lookup details through this table, including IP address mapping and associated metadata. Utilize it to uncover information about IP addresses, such as their associated domain names, aiding in network troubleshooting and security investigations.

## Examples

### Find host names for an IP address
Discover the host names associated with a specific IP address to better understand network connections and potential security risks.

```sql+postgres
select
  *
from
  net_dns_reverse
where
  ip_address = '1.1.1.1';
```

```sql+sqlite
select
  *
from
  net_dns_reverse
where
  ip_address = '1.1.1.1';
```